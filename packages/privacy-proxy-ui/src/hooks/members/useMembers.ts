import {
  useMutation,
  useQuery,
  useQueryClient,
  type UseQueryResult,
} from "@tanstack/react-query";
import { request } from "../../lib/apiClient";
import type {
  MemberView,
  RevokeTokensResponse,
  UpsertMemberRequest,
} from "../../types/api";
import { normalizeAddress } from "../../lib/format";

const KEYS = {
  list: (role?: string) => ["members", "list", role ?? "all"] as const,
  detail: (eoa: string) =>
    ["members", "detail", normalizeAddress(eoa)] as const,
};

export function useListMembers(role?: string): UseQueryResult<MemberView[]> {
  return useQuery({
    queryKey: KEYS.list(role),
    queryFn: () => {
      const params = new URLSearchParams();
      if (role) params.set("role", role);
      params.set("limit", "1000");
      return request<MemberView[]>(`/admin/members?${params.toString()}`);
    },
    staleTime: 10_000,
  });
}

export function useGetMember(eoa: string | undefined): UseQueryResult<MemberView> {
  return useQuery({
    enabled: !!eoa,
    queryKey: eoa ? KEYS.detail(eoa) : ["members", "detail", "none"],
    queryFn: () =>
      request<MemberView>(`/admin/members/${normalizeAddress(eoa as string)}`),
  });
}

export function useUpsertMember() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({
      eoa,
      body,
    }: {
      eoa: string;
      body: UpsertMemberRequest;
    }) =>
      request<MemberView>(`/admin/members/${normalizeAddress(eoa)}`, {
        method: "PUT",
        body,
      }),
    onSuccess: (data) => {
      qc.invalidateQueries({ queryKey: ["members", "list"] });
      qc.setQueryData(KEYS.detail(data.eoa_address), data);
    },
  });
}

export function useDeleteMember() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (eoa: string) =>
      request<void>(`/admin/members/${normalizeAddress(eoa)}`, {
        method: "DELETE",
      }),
    onSuccess: (_, eoa) => {
      qc.invalidateQueries({ queryKey: ["members", "list"] });
      qc.removeQueries({ queryKey: KEYS.detail(eoa) });
    },
  });
}

export function useRevokeMemberTokens() {
  return useMutation({
    mutationFn: (eoa: string) =>
      request<RevokeTokensResponse>(
        `/admin/members/${normalizeAddress(eoa)}/tokens`,
        { method: "DELETE" },
      ),
  });
}
