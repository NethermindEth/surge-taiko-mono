import { useEffect, useMemo, useState } from "react";
import toast from "react-hot-toast";
import { useCreateLambda } from "../../hooks/useLambdas";
import { useRoleAttributes } from "../../hooks/useRoleAttributes";
import { useSyntheticSelectors } from "../../hooks/useSyntheticSelectors";
import { AdminApiError } from "../../lib/apiClient";
import { findCommonSelector } from "../../lib/selectors";
import type { CommonSelectorParam } from "../../lib/selectors";
import {
  hasArrayParam,
  isParseError,
  parseSignature,
  type ParsedParam,
} from "../../lib/signature";
import { categorize, encodeLiteral } from "../../lib/abiTypes";
import type {
  Condition,
  LambdaRuleInput,
  LambdaView,
  LhsKind,
  RhsKind,
  RoleAttribute,
  RoleName,
} from "../../types/api";

interface LambdaBuilderProps {
  role: RoleName;
  selector: string;
  onCancel: () => void;
  onCreated: (lambda: LambdaView) => void;
}

interface DraftRule {
  key: string;
  lhs_kind: LhsKind;
  lhs_param_index: number | null;
  lhs_attribute: string | null;
  condition: Condition;
  rhs_kind: RhsKind;
  rhs_value: string;
}

interface ResolvedParams {
  source: "common" | "synthetic" | "signature";
  signature?: string;
  params: ParsedParam[];
}

const CONDITION_LABELS: Record<Condition, string> = {
  eq: "=",
  neq: "≠",
  gt: ">",
  lt: "<",
  gte: "≥",
  lte: "≤",
};

const ALL_CONDITIONS: Condition[] = ["eq", "neq", "gt", "lt", "gte", "lte"];
const EQUALITY_ONLY_CONDITIONS: Condition[] = ["eq", "neq"];

const RHS_LABEL: Record<RhsKind, string> = {
  tx_origin: "tx.origin",
  msg_sender: "msg.sender",
  literal: "Literal value",
};

function conditionsFor(ty: string | null): Condition[] {
  const cat = categorize(ty);
  if (cat.kind === "bool" || cat.kind === "address") return EQUALITY_ONLY_CONDITIONS;
  return ALL_CONDITIONS;
}

function rhsKindsFor(ty: string | null): RhsKind[] {
  // tx.origin and msg.sender are address-valued, so they only make sense to
  // compare against an address LHS. Everything else can only check against a
  // typed literal.
  const cat = categorize(ty);
  if (cat.kind === "address") return ["tx_origin", "msg_sender", "literal"];
  return ["literal"];
}

function newDraftRule(): DraftRule {
  return {
    key: Math.random().toString(36).slice(2),
    lhs_kind: "calldata",
    lhs_param_index: 0,
    lhs_attribute: null,
    condition: "eq",
    rhs_kind: "tx_origin",
    rhs_value: "",
  };
}

function syntheticParams(method: string): CommonSelectorParam[] {
  if (method === "eth_getStorageAt") {
    return [
      { name: "target", type: "address" },
      { name: "slot", type: "bytes32" },
    ];
  }
  return [{ name: "target", type: "address" }];
}

function toParams(params: CommonSelectorParam[]): ParsedParam[] {
  return params.map((p) => ({ name: p.name, type: p.type }));
}

function lhsType(
  draft: DraftRule,
  resolved: ResolvedParams | null,
  attrs: RoleAttribute[],
): string | null {
  if (draft.lhs_kind === "calldata") {
    if (!resolved || draft.lhs_param_index === null) return null;
    return resolved.params[draft.lhs_param_index]?.type ?? null;
  }
  if (!draft.lhs_attribute) return null;
  return attrs.find((a) => a.name === draft.lhs_attribute)?.type ?? null;
}

export function LambdaBuilder({
  role,
  selector,
  onCancel,
  onCreated,
}: LambdaBuilderProps) {
  const create = useCreateLambda();
  const attrsQ = useRoleAttributes();
  const syntheticQ = useSyntheticSelectors();

  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [signatureInput, setSignatureInput] = useState("");
  const [signatureError, setSignatureError] = useState<string | null>(null);
  const [resolved, setResolved] = useState<ResolvedParams | null>(null);
  const [rules, setRules] = useState<DraftRule[]>([newDraftRule()]);

  const roleAttributes = useMemo<RoleAttribute[]>(
    () => attrsQ.data?.find((g) => g.role === role)?.attributes ?? [],
    [attrsQ.data, role],
  );

  useEffect(() => {
    if (!selector) {
      setResolved(null);
      return;
    }
    const common = findCommonSelector(selector);
    if (common) {
      setResolved({
        source: "common",
        signature: common.signature,
        params: toParams(common.params),
      });
      setSignatureError(null);
      return;
    }
    const synthetic = (syntheticQ.data ?? []).find(
      (s) => s.selector.toLowerCase() === selector.toLowerCase(),
    );
    if (synthetic) {
      setResolved({
        source: "synthetic",
        signature: synthetic.method,
        params: toParams(syntheticParams(synthetic.method)),
      });
      setSignatureError(null);
      return;
    }
    setResolved(null);
  }, [selector, syntheticQ.data]);

  const onSignatureBlur = () => {
    if (!signatureInput.trim()) {
      setSignatureError(null);
      return;
    }
    const parsed = parseSignature(signatureInput);
    if (isParseError(parsed)) {
      setSignatureError(parsed.error);
      setResolved(null);
      return;
    }
    if (parsed.selector.toLowerCase() !== selector.toLowerCase()) {
      setSignatureError(
        `Signature resolves to selector ${parsed.selector}; doesn't match the rule selector ${selector}.`,
      );
      setResolved(null);
      return;
    }
    if (hasArrayParam(parsed.params)) {
      setSignatureError("Array parameters are not supported in lambda rules.");
      setResolved(null);
      return;
    }
    setSignatureError(null);
    setResolved({
      source: "signature",
      signature: `${parsed.name}(${parsed.params.map((p) => p.type).join(",")})`,
      params: parsed.params,
    });
  };

  const updateRule = (idx: number, patch: Partial<DraftRule>) => {
    setRules((rs) =>
      rs.map((r, i) => {
        if (i !== idx) return r;
        const next = { ...r, ...patch };
        // Reset rhs_value when the LHS or rhs_kind changes — a typed input from
        // the prior context will rarely round-trip into the new one.
        const lhsChanged =
          patch.lhs_kind !== undefined ||
          patch.lhs_param_index !== undefined ||
          patch.lhs_attribute !== undefined;
        const rhsKindChanged =
          patch.rhs_kind !== undefined && patch.rhs_kind !== r.rhs_kind;
        if (lhsChanged || rhsKindChanged) next.rhs_value = "";
        // Snap the condition and rhs_kind back into the allowed sets for the
        // new LHS type — e.g. switching to a uint LHS hides tx.origin/msg.sender.
        if (lhsChanged) {
          const nextTy = lhsType(next, resolved, roleAttributes);
          const allowedConds = conditionsFor(nextTy);
          if (!allowedConds.includes(next.condition)) next.condition = "eq";
          const allowedRhs = rhsKindsFor(nextTy);
          if (!allowedRhs.includes(next.rhs_kind)) {
            next.rhs_kind = "literal";
            next.rhs_value = "";
          }
        }
        return next;
      }),
    );
  };

  const addRule = () => setRules((rs) => [...rs, newDraftRule()]);
  const removeRule = (idx: number) =>
    setRules((rs) => rs.filter((_, i) => i !== idx));

  const canSubmit = useMemo(() => {
    if (!name.trim()) return false;
    if (!resolved) return false;
    if (rules.length === 0) return false;
    for (const r of rules) {
      if (r.lhs_kind === "calldata") {
        if (
          r.lhs_param_index === null ||
          r.lhs_param_index < 0 ||
          r.lhs_param_index >= resolved.params.length
        ) {
          return false;
        }
      } else {
        if (
          !r.lhs_attribute ||
          !roleAttributes.some((a) => a.name === r.lhs_attribute)
        ) {
          return false;
        }
      }
      if (r.rhs_kind === "literal") {
        const ty = lhsType(r, resolved, roleAttributes);
        if (!encodeLiteral(r.rhs_value, ty).ok) return false;
      }
    }
    return true;
  }, [name, resolved, rules, roleAttributes]);

  const onSubmit = async () => {
    if (!resolved) return;
    const built: LambdaRuleInput[] = rules.map((r) => {
      const base = { selector, condition: r.condition };
      let lhs: Pick<
        LambdaRuleInput,
        "lhs_kind" | "lhs_offset" | "lhs_attribute"
      >;
      if (r.lhs_kind === "calldata") {
        const idx = r.lhs_param_index!;
        lhs = {
          lhs_kind: "calldata",
          lhs_offset: 4 + 32 * idx,
          lhs_attribute: null,
        };
      } else {
        lhs = {
          lhs_kind: "attribute",
          lhs_offset: null,
          lhs_attribute: r.lhs_attribute,
        };
      }
      let rhs: Pick<LambdaRuleInput, "rhs_kind" | "rhs_value">;
      if (r.rhs_kind === "literal") {
        const ty = lhsType(r, resolved, roleAttributes);
        const encoded = encodeLiteral(r.rhs_value, ty);
        rhs = {
          rhs_kind: "literal",
          rhs_value: encoded.ok ? encoded.hex : null,
        };
      } else {
        rhs = { rhs_kind: r.rhs_kind, rhs_value: null };
      }
      return { ...base, ...lhs, ...rhs };
    });

    try {
      const created = await create.mutateAsync({
        name: name.trim(),
        role,
        description: description.trim() || undefined,
        rules: built,
      });
      toast.success(`Lambda "${created.name}" created`);
      onCreated(created);
    } catch (e) {
      toast.error(
        e instanceof AdminApiError ? e.message : (e as Error).message,
      );
    }
  };

  return (
    <div className="space-y-4 rounded-2xl border border-surge-border bg-surge-card-hover/40 p-4">
      <header className="flex items-start justify-between gap-3">
        <div>
          <h3 className="text-sm font-semibold text-surge-text">
            Create new lambda
          </h3>
          <p className="text-xs text-surge-muted">
            Comparison rules on selector{" "}
            <span className="font-mono text-surge-text">{selector || "—"}</span>
            , for role{" "}
            <span className="font-semibold text-surge-text">{role}</span>.
          </p>
        </div>
        <button
          type="button"
          onClick={onCancel}
          className="rounded-md px-2 py-1 text-xs font-medium text-surge-muted hover:bg-surge-card hover:text-surge-text"
        >
          Cancel
        </button>
      </header>

      <div className="grid grid-cols-1 gap-3 md:grid-cols-2">
        <div>
          <label className="block text-[10px] font-semibold uppercase tracking-wide text-surge-muted">
            Name
          </label>
          <input
            type="text"
            value={name}
            onChange={(e) => setName(e.target.value)}
            placeholder="e.g. erc20_self_only"
            className="mt-1 w-full rounded-md border border-surge-border bg-surge-card px-2 py-1.5 text-sm outline-none focus:border-surge-secondary"
          />
        </div>
        <div>
          <label className="block text-[10px] font-semibold uppercase tracking-wide text-surge-muted">
            Description (optional)
          </label>
          <input
            type="text"
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            placeholder="What does this lambda do?"
            className="mt-1 w-full rounded-md border border-surge-border bg-surge-card px-2 py-1.5 text-sm outline-none focus:border-surge-secondary"
          />
        </div>
      </div>

      {!resolved ? (
        <div className="rounded-xl border border-dashed border-surge-border bg-surge-card p-3">
          <label className="block text-[10px] font-semibold uppercase tracking-wide text-surge-muted">
            Function signature
          </label>
          <p className="mt-1 text-xs text-surge-muted">
            This selector isn't in our common list. Paste the full function
            signature so we can resolve its parameters.
          </p>
          <input
            type="text"
            value={signatureInput}
            onChange={(e) => setSignatureInput(e.target.value)}
            onBlur={onSignatureBlur}
            placeholder="transfer(address to,uint256 amount)"
            className="mt-2 w-full rounded-md border border-surge-border bg-surge-card-hover/40 px-2 py-1.5 font-mono text-sm outline-none focus:border-surge-secondary"
          />
          {signatureError ? (
            <p className="mt-1 text-xs text-red-600">{signatureError}</p>
          ) : null}
        </div>
      ) : (
        <div className="rounded-xl border border-surge-border bg-surge-card p-3 text-xs text-surge-muted">
          <span className="font-medium text-surge-text">Resolved:</span>{" "}
          <code className="font-mono">{resolved.signature}</code>
          {resolved.source === "signature" ? (
            <button
              type="button"
              onClick={() => {
                setResolved(null);
                setSignatureInput("");
              }}
              className="ml-2 text-surge-secondary hover:underline"
            >
              change
            </button>
          ) : null}
        </div>
      )}

      <div className="space-y-3">
        <div className="flex items-center justify-between">
          <h4 className="text-xs font-semibold uppercase tracking-wide text-surge-muted">
            Comparison rules ({rules.length})
          </h4>
          <button
            type="button"
            onClick={addRule}
            disabled={!resolved}
            className="rounded-md border border-dashed border-surge-border px-2 py-1 text-xs font-medium text-surge-muted hover:border-surge-secondary hover:text-surge-text disabled:cursor-not-allowed disabled:opacity-40"
          >
            + Add rule
          </button>
        </div>

        {rules.map((r, idx) => {
          const ty = lhsType(r, resolved, roleAttributes);
          return (
            <div
              key={r.key}
              className="space-y-2 rounded-xl border border-surge-border bg-surge-card p-3"
            >
              <div className="flex items-center justify-between">
                <span className="text-[10px] font-semibold uppercase tracking-wide text-surge-muted">
                  Rule #{idx + 1}
                </span>
                {rules.length > 1 ? (
                  <button
                    type="button"
                    onClick={() => removeRule(idx)}
                    aria-label="Remove rule"
                    className="rounded-md p-1 text-surge-muted hover:bg-surge-card-hover hover:text-red-600"
                  >
                    <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.6" strokeLinecap="round" strokeLinejoin="round" aria-hidden>
                      <path d="M3 6h18M8 6V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2M6 6l1 14a2 2 0 0 0 2 2h6a2 2 0 0 0 2-2l1-14" />
                    </svg>
                  </button>
                ) : null}
              </div>

              <div className="grid grid-cols-1 gap-2 md:grid-cols-3">
                <div>
                  <label className="block text-[10px] font-semibold uppercase tracking-wide text-surge-muted">
                    LHS source
                  </label>
                  <div className="mt-1 inline-flex w-full rounded-md border border-surge-border bg-surge-card-hover/40 p-0.5 text-xs">
                    {(["calldata", "attribute"] as const).map((k) => (
                      <button
                        key={k}
                        type="button"
                        onClick={() =>
                          updateRule(idx, {
                            lhs_kind: k,
                            lhs_param_index: k === "calldata" ? 0 : null,
                            lhs_attribute:
                              k === "attribute"
                                ? roleAttributes[0]?.name ?? null
                                : null,
                          })
                        }
                        className={`flex-1 rounded px-2 py-1 transition ${
                          r.lhs_kind === k
                            ? "bg-surge-card text-surge-text shadow-sm"
                            : "text-surge-muted hover:text-surge-text"
                        }`}
                      >
                        {k === "calldata" ? "Calldata arg" : "Role attribute"}
                      </button>
                    ))}
                  </div>
                </div>

                <div>
                  <label className="block text-[10px] font-semibold uppercase tracking-wide text-surge-muted">
                    {r.lhs_kind === "calldata" ? "Argument" : "Attribute"}
                  </label>
                  {r.lhs_kind === "calldata" ? (
                    <select
                      value={r.lhs_param_index ?? 0}
                      onChange={(e) =>
                        updateRule(idx, {
                          lhs_param_index: Number(e.target.value),
                        })
                      }
                      disabled={!resolved || resolved.params.length === 0}
                      className="mt-1 w-full rounded-md border border-surge-border bg-surge-card-hover/40 px-2 py-1.5 text-sm outline-none focus:border-surge-secondary disabled:opacity-50"
                    >
                      {resolved && resolved.params.length > 0 ? (
                        resolved.params.map((p, i) => (
                          <option key={i} value={i}>
                            arg{i}
                            {p.name ? ` (${p.name})` : ""}: {p.type}
                          </option>
                        ))
                      ) : (
                        <option value={0}>No parameters available</option>
                      )}
                    </select>
                  ) : (
                    <select
                      value={r.lhs_attribute ?? ""}
                      onChange={(e) =>
                        updateRule(idx, { lhs_attribute: e.target.value })
                      }
                      disabled={roleAttributes.length === 0}
                      className="mt-1 w-full rounded-md border border-surge-border bg-surge-card-hover/40 px-2 py-1.5 text-sm outline-none focus:border-surge-secondary disabled:opacity-50"
                    >
                      {roleAttributes.length === 0 ? (
                        <option value="">No attributes for this role</option>
                      ) : null}
                      {roleAttributes.map((a) => (
                        <option key={a.name} value={a.name}>
                          {a.name}: {a.type}
                        </option>
                      ))}
                    </select>
                  )}
                </div>

                <div>
                  <label className="block text-[10px] font-semibold uppercase tracking-wide text-surge-muted">
                    Condition
                  </label>
                  <div
                    className={`mt-1 grid gap-1 ${
                      conditionsFor(ty).length <= 2
                        ? "grid-cols-2"
                        : "grid-cols-6"
                    }`}
                  >
                    {conditionsFor(ty).map((c) => (
                      <button
                        key={c}
                        type="button"
                        onClick={() => updateRule(idx, { condition: c })}
                        className={`rounded-md px-1.5 py-1.5 text-sm font-semibold transition ${
                          r.condition === c
                            ? "bg-surge-primary text-white"
                            : "border border-surge-border bg-surge-card text-surge-text hover:bg-surge-card-hover"
                        }`}
                      >
                        {CONDITION_LABELS[c]}
                      </button>
                    ))}
                  </div>
                </div>
              </div>

              <div className="grid grid-cols-1 gap-2 md:grid-cols-3">
                <div className="md:col-span-1">
                  <label className="block text-[10px] font-semibold uppercase tracking-wide text-surge-muted">
                    RHS source
                  </label>
                  <select
                    value={r.rhs_kind}
                    onChange={(e) =>
                      updateRule(idx, { rhs_kind: e.target.value as RhsKind })
                    }
                    className="mt-1 w-full rounded-md border border-surge-border bg-surge-card-hover/40 px-2 py-1.5 text-sm outline-none focus:border-surge-secondary"
                  >
                    {rhsKindsFor(ty).map((k) => (
                      <option key={k} value={k}>
                        {RHS_LABEL[k]}
                      </option>
                    ))}
                  </select>
                </div>
                {r.rhs_kind === "literal" ? (
                  <div className="md:col-span-2">
                    <LiteralInput
                      type={ty}
                      value={r.rhs_value}
                      onChange={(v) => updateRule(idx, { rhs_value: v })}
                    />
                  </div>
                ) : null}
              </div>
            </div>
          );
        })}
      </div>

      <div className="flex items-center justify-end gap-2">
        <button
          type="button"
          onClick={onCancel}
          className="rounded-lg px-3 py-2 text-sm font-medium text-surge-muted hover:bg-surge-card-hover hover:text-surge-text"
        >
          Cancel
        </button>
        <button
          type="button"
          onClick={onSubmit}
          disabled={!canSubmit || create.isPending}
          className="rounded-lg bg-surge-primary px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-surge-primary/90 disabled:cursor-not-allowed disabled:opacity-50"
        >
          {create.isPending ? "Saving…" : "Save lambda"}
        </button>
      </div>
    </div>
  );
}

interface LiteralInputProps {
  type: string | null;
  value: string;
  onChange: (v: string) => void;
}

function LiteralInput({ type, value, onChange }: LiteralInputProps) {
  const cat = categorize(type);
  const result = value ? encodeLiteral(value, type) : null;
  const showError = !!value && result && !result.ok;

  const placeholder = ((): string => {
    switch (cat.kind) {
      case "bool":
        return "true / false";
      case "address":
        return "0x… (40 hex chars)";
      case "uint":
        return `unsigned integer (0 .. 2^${cat.bits}-1)`;
      case "int":
        return `signed integer (-2^${cat.bits - 1} .. 2^${cat.bits - 1}-1)`;
      case "bytes":
        return `0x… (${cat.size * 2} hex chars)`;
      case "unknown":
        return "0x.. (32 bytes) | address | decimal integer";
    }
  })();

  return (
    <div>
      <label className="block text-[10px] font-semibold uppercase tracking-wide text-surge-muted">
        Literal value
        {type ? (
          <span className="ml-1 text-surge-muted/80">({type})</span>
        ) : null}
      </label>
      {cat.kind === "bool" ? (
        <select
          value={value}
          onChange={(e) => onChange(e.target.value)}
          className="mt-1 w-full rounded-md border border-surge-border bg-surge-card-hover/40 px-2 py-1.5 text-sm outline-none focus:border-surge-secondary"
        >
          <option value="">Select…</option>
          <option value="true">true</option>
          <option value="false">false</option>
        </select>
      ) : (
        <input
          type="text"
          value={value}
          onChange={(e) => onChange(e.target.value)}
          placeholder={placeholder}
          className="mt-1 w-full rounded-md border border-surge-border bg-surge-card-hover/40 px-2 py-1.5 font-mono text-sm outline-none focus:border-surge-secondary"
        />
      )}
      {showError && result && !result.ok ? (
        <p className="mt-1 text-[11px] text-red-600">{result.error}</p>
      ) : null}
    </div>
  );
}
