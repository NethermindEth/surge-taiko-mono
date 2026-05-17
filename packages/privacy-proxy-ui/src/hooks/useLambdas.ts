import { useQuery, type UseQueryResult } from "@tanstack/react-query";
import { request } from "../lib/apiClient";
import type { LambdaGroup } from "../types/api";

export function useLambdas(): UseQueryResult<LambdaGroup[]> {
  return useQuery({
    queryKey: ["registry", "lambdas"],
    queryFn: () => request<LambdaGroup[]>("/admin/registry/lambdas"),
    staleTime: Infinity, // static per build
  });
}
