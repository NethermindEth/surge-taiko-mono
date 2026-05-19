import {
  useMutation,
  useQuery,
  useQueryClient,
  type UseQueryResult,
} from "@tanstack/react-query";
import { request } from "../lib/apiClient";
import type {
  CreateLambdaRequest,
  LambdaGroup,
  LambdaView,
} from "../types/api";

const KEYS = {
  list: ["registry", "lambdas"] as const,
  detail: (id: number) => ["registry", "lambdas", id] as const,
};

export function useLambdas(): UseQueryResult<LambdaGroup[]> {
  return useQuery({
    queryKey: KEYS.list,
    queryFn: () => request<LambdaGroup[]>("/admin/registry/lambdas"),
    staleTime: 5_000,
  });
}

export function useGetLambda(id: number | undefined): UseQueryResult<LambdaView> {
  return useQuery({
    enabled: id !== undefined,
    queryKey:
      id !== undefined ? KEYS.detail(id) : ["registry", "lambdas", "none"],
    queryFn: () => request<LambdaView>(`/admin/registry/lambdas/${id}`),
  });
}

export function useCreateLambda() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (body: CreateLambdaRequest) =>
      request<LambdaView>("/admin/registry/lambdas", { method: "POST", body }),
    onSuccess: (data) => {
      qc.invalidateQueries({ queryKey: KEYS.list });
      qc.setQueryData(KEYS.detail(data.id), data);
    },
  });
}

export function useDeleteLambda() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: number) =>
      request<void>(`/admin/registry/lambdas/${id}`, { method: "DELETE" }),
    onSuccess: (_, id) => {
      qc.invalidateQueries({ queryKey: KEYS.list });
      qc.removeQueries({ queryKey: KEYS.detail(id) });
    },
  });
}
