import {
  useMutation,
  useQuery,
  useQueryClient,
  type UseQueryResult,
} from "@tanstack/react-query";
import { request } from "../../lib/apiClient";
import type {
  CreateRuleRequest,
  EntryInput,
  EntryView,
  ReplaceRuleRequest,
  RuleView,
  UpdateEntryRequest,
} from "../../types/api";

const KEYS = {
  list: (contract?: string) => ["rules", "list", contract ?? "all"] as const,
  detail: (id: number) => ["rules", "detail", id] as const,
};

export function useListRules(contract?: string): UseQueryResult<RuleView[]> {
  return useQuery({
    queryKey: KEYS.list(contract),
    queryFn: () => {
      const params = new URLSearchParams();
      if (contract) params.set("contract", contract);
      params.set("limit", "1000");
      return request<RuleView[]>(
        `/admin/registry/rules${params.toString() ? `?${params.toString()}` : ""}`,
      );
    },
    staleTime: 10_000,
  });
}

export function useGetRule(id: number | undefined): UseQueryResult<RuleView> {
  return useQuery({
    enabled: id !== undefined,
    queryKey: id !== undefined ? KEYS.detail(id) : ["rules", "detail", "none"],
    queryFn: () => request<RuleView>(`/admin/registry/rules/${id}`),
  });
}

export function useCreateRule() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (body: CreateRuleRequest) =>
      request<RuleView>("/admin/registry/rules", { method: "POST", body }),
    onSuccess: (data) => {
      qc.invalidateQueries({ queryKey: ["rules", "list"] });
      qc.setQueryData(KEYS.detail(data.id), data);
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
      qc.invalidateQueries({ queryKey: ["rules", "list"] });
      qc.setQueryData(KEYS.detail(data.id), data);
    },
  });
}

export function useDeleteRule() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: number) =>
      request<void>(`/admin/registry/rules/${id}`, { method: "DELETE" }),
    onSuccess: (_, id) => {
      qc.invalidateQueries({ queryKey: ["rules", "list"] });
      qc.removeQueries({ queryKey: KEYS.detail(id) });
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
      qc.invalidateQueries({ queryKey: KEYS.detail(vars.ruleId) });
      qc.invalidateQueries({ queryKey: ["rules", "list"] });
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
      qc.invalidateQueries({ queryKey: KEYS.detail(vars.ruleId) });
      qc.invalidateQueries({ queryKey: ["rules", "list"] });
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
      qc.invalidateQueries({ queryKey: KEYS.detail(vars.ruleId) });
      qc.invalidateQueries({ queryKey: ["rules", "list"] });
    },
  });
}
