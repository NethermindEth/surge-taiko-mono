import { useQuery, type UseQueryResult } from "@tanstack/react-query";
import { request } from "../lib/apiClient";
import type { Role } from "../types/api";

export function useRoles(): UseQueryResult<Role[]> {
  return useQuery({
    queryKey: ["registry", "roles"],
    queryFn: () => request<Role[]>("/admin/roles"),
    staleTime: Infinity,
  });
}
