import { useEffect, useMemo, useState } from "react";
import toast from "react-hot-toast";
import { Drawer } from "../common/Drawer";
import { Badge, ModeBadge } from "../common/Badge";
import { SelectorPicker } from "./SelectorPicker";
import { RuleDrawer } from "./RuleDrawer";
import {
  useCreateBinding,
  useListRules,
} from "../../hooks/rules/useRules";
import { isAddress, normalizeAddress } from "../../lib/format";
import { findCommonSelector } from "../../lib/selectors";
import { useSyntheticSelectors } from "../../hooks/useSyntheticSelectors";
import { AdminApiError } from "../../lib/apiClient";
import type { RuleView } from "../../types/api";

interface BindingDrawerProps {
  open: boolean;
  onClose: () => void;
  /** Pre-filled contract address when launched from a contract row. */
  presetContract?: string;
  /** Pre-filled selector when launched from a contract-selector row. */
  presetSelector?: string;
}

export function BindingDrawer({
  open,
  onClose,
  presetContract,
  presetSelector,
}: BindingDrawerProps) {
  const rules = useListRules();
  const synthetic = useSyntheticSelectors();
  const create = useCreateBinding();

  const [contract, setContract] = useState("");
  const [contractError, setContractError] = useState<string | null>(null);
  const [selector, setSelector] = useState("");
  const [selectedRuleId, setSelectedRuleId] = useState<number | null>(null);
  const [ruleDrawerOpen, setRuleDrawerOpen] = useState(false);

  useEffect(() => {
    if (!open) return;
    setContract(presetContract ?? "");
    setSelector(presetSelector ?? "");
    setSelectedRuleId(null);
    setContractError(null);
  }, [open, presetContract, presetSelector]);

  const hexSelector = useMemo(() => {
    if (!selector) return "";
    if (selector.startsWith("0x")) return selector.toLowerCase();
    const found = (synthetic.data ?? []).find((s) => s.method === selector);
    return found ? found.selector.toLowerCase() : "";
  }, [selector, synthetic.data]);

  const compatibleRules = useMemo<RuleView[]>(() => {
    if (!rules.data) return [];
    if (!hexSelector) return rules.data;
    return rules.data.filter(
      (r) => r.selector.toLowerCase() === hexSelector,
    );
  }, [rules.data, hexSelector]);

  const selectedRule = compatibleRules.find((r) => r.id === selectedRuleId);

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!isAddress(contract)) {
      setContractError("Enter a 20-byte hex address.");
      return;
    }
    if (selectedRuleId === null) {
      toast.error("Pick a rule to apply, or create a new one.");
      return;
    }
    try {
      await create.mutateAsync({
        contract_address: normalizeAddress(contract),
        rule_id: selectedRuleId,
      });
      toast.success("Rule applied to contract");
      onClose();
    } catch (err) {
      toast.error(
        err instanceof AdminApiError ? err.message : (err as Error).message,
      );
    }
  };

  return (
    <Drawer
      open={open}
      onClose={onClose}
      title="Apply rule to contract"
      widthClass="max-w-xl"
      footer={
        <div className="flex items-center justify-end gap-2">
          <button
            type="button"
            onClick={onClose}
            className="rounded-lg px-3 py-2 text-sm font-medium text-surge-muted hover:bg-surge-card-hover hover:text-surge-text"
          >
            Cancel
          </button>
          <button
            type="submit"
            form="binding-form"
            disabled={create.isPending || selectedRuleId === null}
            className="rounded-lg bg-surge-primary px-4 py-2 text-sm font-semibold text-white shadow-sm hover:bg-surge-primary/90 disabled:cursor-not-allowed disabled:opacity-50"
          >
            {create.isPending ? "Applying…" : "Apply rule"}
          </button>
        </div>
      }
    >
      <form id="binding-form" onSubmit={onSubmit} className="space-y-5">
        <div>
          <label className="block text-xs font-medium uppercase tracking-wide text-surge-muted">
            Contract address
          </label>
          <input
            type="text"
            value={contract}
            onChange={(e) => {
              setContract(e.target.value);
              setContractError(null);
            }}
            placeholder="0x..."
            className="mt-1 w-full rounded-lg border border-surge-border bg-surge-card-hover/40 px-3 py-2 font-mono text-sm outline-none focus:border-surge-secondary"
          />
          {contractError ? (
            <p className="mt-1 text-xs text-red-600">{contractError}</p>
          ) : null}
        </div>

        <div>
          <label className="block text-xs font-medium uppercase tracking-wide text-surge-muted">
            Selector to gate
          </label>
          <p className="mt-1 text-[11px] text-surge-muted">
            Filters the rule list below to rules whose selector matches.
            Leave blank to see all rules.
          </p>
          <div className="mt-2">
            <SelectorPicker
              value={selector}
              onChange={(v) => {
                setSelector(v);
                setSelectedRuleId(null);
              }}
            />
          </div>
        </div>

        <div>
          <div className="flex items-center justify-between">
            <label className="text-xs font-medium uppercase tracking-wide text-surge-muted">
              Rule
            </label>
            <button
              type="button"
              onClick={() => setRuleDrawerOpen(true)}
              className="rounded-md border border-surge-secondary px-2 py-1 text-xs font-semibold text-surge-primary hover:bg-surge-secondary/10"
            >
              + Create new rule
            </button>
          </div>
          {rules.isLoading ? (
            <p className="mt-2 text-xs text-surge-muted">Loading rules…</p>
          ) : compatibleRules.length === 0 ? (
            <div className="mt-2 rounded-xl border border-dashed border-surge-border bg-surge-card-hover/30 p-3 text-sm text-surge-muted">
              {hexSelector
                ? "No existing rules match this selector. Create a new one above."
                : "No rules defined yet. Create one to bind it to a contract selector."}
            </div>
          ) : (
            <div className="mt-2 max-h-72 space-y-2 overflow-auto rounded-xl border border-surge-border bg-surge-card p-2">
              {compatibleRules.map((r) => {
                const checked = selectedRuleId === r.id;
                const sigHint = findCommonSelector(r.selector)?.signature;
                return (
                  <label
                    key={r.id}
                    className={`flex cursor-pointer items-start gap-3 rounded-lg border px-3 py-2 transition ${
                      checked
                        ? "border-surge-primary bg-surge-primary/5"
                        : "border-surge-border bg-surge-card hover:bg-surge-card-hover"
                    }`}
                  >
                    <input
                      type="radio"
                      name="rule"
                      checked={checked}
                      onChange={() => setSelectedRuleId(r.id)}
                      className="mt-1"
                    />
                    <div className="min-w-0 flex-1">
                      <div className="flex flex-wrap items-center gap-2">
                        <span className="font-mono text-sm font-semibold text-surge-text">
                          {r.name}
                        </span>
                        <ModeBadge mode={r.mode} />
                        <Badge tone="aqua" className="font-mono text-[10px]">
                          {r.selector}
                        </Badge>
                        {sigHint ? (
                          <span className="text-[10px] text-surge-muted">
                            {sigHint}
                          </span>
                        ) : null}
                      </div>
                      {r.description ? (
                        <p className="mt-1 text-xs text-surge-muted">
                          {r.description}
                        </p>
                      ) : null}
                      <p className="mt-1 text-[11px] text-surge-muted">
                        {r.entries.length} role entr
                        {r.entries.length === 1 ? "y" : "ies"} · used by{" "}
                        {r.binding_count} contract(s)
                      </p>
                    </div>
                  </label>
                );
              })}
            </div>
          )}
          {selectedRule ? (
            <p className="mt-1 text-[11px] text-surge-muted">
              Will gate{" "}
              <span className="font-mono text-surge-text">
                {selectedRule.selector}
              </span>{" "}
              on this contract with the{" "}
              <span className="font-semibold text-surge-text">
                {selectedRule.mode}
              </span>{" "}
              policy.
            </p>
          ) : null}
        </div>
      </form>

      <RuleDrawer
        open={ruleDrawerOpen}
        onClose={() => setRuleDrawerOpen(false)}
        presetSelector={selector}
        onCreated={(rule) => setSelectedRuleId(rule.id)}
      />
    </Drawer>
  );
}
