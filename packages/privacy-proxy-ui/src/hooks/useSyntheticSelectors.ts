import { useQuery, type UseQueryResult } from "@tanstack/react-query";
import { request } from "../lib/apiClient";
import type { SyntheticSelector } from "../types/api";

export function useSyntheticSelectors(): UseQueryResult<SyntheticSelector[]> {
  return useQuery({
    queryKey: ["registry", "synthetic-selectors"],
    queryFn: () =>
      request<SyntheticSelector[]>("/admin/registry/synthetic-selectors"),
    staleTime: Infinity,
  });
}

/**
 * Look up "0xff010001" → "eth_getBalance". Returns undefined if not a known
 * synthetic selector.
 */
export function useSelectorName(
  selector: string | undefined,
): string | undefined {
  const { data } = useSyntheticSelectors();
  if (!selector || !data) return undefined;
  return data.find((s) => s.selector.toLowerCase() === selector.toLowerCase())
    ?.method;
}
