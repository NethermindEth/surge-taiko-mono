import {
  useMutation,
  useQuery,
  useQueryClient,
  type UseQueryResult,
} from "@tanstack/react-query";
import { request } from "../../lib/apiClient";
import type {
  BindingView,
  CreateBindingRequest,
  CreateRuleRequest,
  EntryInput,
  EntryView,
  ReplaceRuleRequest,
  RuleView,
  UpdateEntryRequest,
} from "../../types/api";

const KEYS = {
  rules: () => ["rules", "list"] as const,
  rule: (id: number) => ["rules", "detail", id] as const,
  bindings: (contract?: string, ruleId?: number) =>
    ["bindings", "list", contract ?? "all", ruleId ?? "all"] as const,
};

export function useListRules(): UseQueryResult<RuleView[]> {
  return useQuery({
    queryKey: KEYS.rules(),
    queryFn: () => request<RuleView[]>("/admin/registry/rules?limit=1000"),
    staleTime: 5_000,
  });
}

export function useGetRule(id: number | undefined): UseQueryResult<RuleView> {
  return useQuery({
    enabled: id !== undefined,
    queryKey: id !== undefined ? KEYS.rule(id) : ["rules", "detail", "none"],
    queryFn: () => request<RuleView>(`/admin/registry/rules/${id}`),
  });
}

export function useCreateRule() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (body: CreateRuleRequest) =>
      request<RuleView>("/admin/registry/rules", { method: "POST", body }),
    onSuccess: (data) => {
      qc.invalidateQueries({ queryKey: KEYS.rules() });
      qc.setQueryData(KEYS.rule(data.id), data);
    },
  });
}

export function useReplaceRule() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ id, body }: { id: number; body: ReplaceRuleRequest }) =>
      request<RuleView>(`/admin/registry/rules/${id}`, {
        method: "PUT",
        body,
      }),
    onSuccess: (data) => {
      qc.invalidateQueries({ queryKey: KEYS.rules() });
      qc.invalidateQueries({ queryKey: ["bindings", "list"] });
      qc.setQueryData(KEYS.rule(data.id), data);
    },
  });
}

export function useDeleteRule() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: number) =>
      request<void>(`/admin/registry/rules/${id}`, { method: "DELETE" }),
    onSuccess: (_, id) => {
      qc.invalidateQueries({ queryKey: KEYS.rules() });
      qc.removeQueries({ queryKey: KEYS.rule(id) });
    },
  });
}

export function useAddEntry() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ ruleId, body }: { ruleId: number; body: EntryInput }) =>
      request<EntryView>(`/admin/registry/rules/${ruleId}/entries`, {
        method: "POST",
        body,
      }),
    onSuccess: (_, vars) => {
      qc.invalidateQueries({ queryKey: KEYS.rule(vars.ruleId) });
      qc.invalidateQueries({ queryKey: KEYS.rules() });
    },
  });
}

export function useUpdateEntry() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({
      ruleId,
      entryId,
      body,
    }: {
      ruleId: number;
      entryId: number;
      body: UpdateEntryRequest;
    }) =>
      request<EntryView>(
        `/admin/registry/rules/${ruleId}/entries/${entryId}`,
        { method: "PUT", body },
      ),
    onSuccess: (_, vars) => {
      qc.invalidateQueries({ queryKey: KEYS.rule(vars.ruleId) });
      qc.invalidateQueries({ queryKey: KEYS.rules() });
    },
  });
}

export function useDeleteEntry() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ ruleId, entryId }: { ruleId: number; entryId: number }) =>
      request<void>(
        `/admin/registry/rules/${ruleId}/entries/${entryId}`,
        { method: "DELETE" },
      ),
    onSuccess: (_, vars) => {
      qc.invalidateQueries({ queryKey: KEYS.rule(vars.ruleId) });
      qc.invalidateQueries({ queryKey: KEYS.rules() });
    },
  });
}

export function useListBindings(opts?: {
  contract?: string;
  ruleId?: number;
}): UseQueryResult<BindingView[]> {
  return useQuery({
    queryKey: KEYS.bindings(opts?.contract, opts?.ruleId),
    queryFn: () => {
      const p = new URLSearchParams();
      if (opts?.contract) p.set("contract", opts.contract);
      if (opts?.ruleId !== undefined) p.set("rule_id", String(opts.ruleId));
      p.set("limit", "1000");
      return request<BindingView[]>(
        `/admin/registry/bindings?${p.toString()}`,
      );
    },
    staleTime: 5_000,
  });
}

export function useCreateBinding() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (body: CreateBindingRequest) =>
      request<BindingView>("/admin/registry/bindings", {
        method: "POST",
        body,
      }),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["bindings", "list"] });
      qc.invalidateQueries({ queryKey: KEYS.rules() });
    },
  });
}

export function useDeleteBinding() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: number) =>
      request<void>(`/admin/registry/bindings/${id}`, { method: "DELETE" }),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["bindings", "list"] });
      qc.invalidateQueries({ queryKey: KEYS.rules() });
    },
  });
}
