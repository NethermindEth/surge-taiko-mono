import { useQuery, type UseQueryResult } from "@tanstack/react-query";
import { request } from "../lib/apiClient";
import type { RoleAttribute, RoleAttributesGroup, RoleName } from "../types/api";

interface RawGroup {
  role: RoleName;
  attributes: Array<string | { name: string; type?: string }>;
}

export function useRoleAttributes(): UseQueryResult<RoleAttributesGroup[]> {
  return useQuery({
    queryKey: ["registry", "role-attributes"],
    queryFn: async () => {
      const raw = await request<RawGroup[]>("/admin/registry/role-attributes");
      return raw.map((g) => ({
        role: g.role,
        attributes: g.attributes.map<RoleAttribute>((a) =>
          typeof a === "string"
            ? { name: a, type: "unknown" }
            : { name: a.name, type: a.type ?? "unknown" },
        ),
      }));
    },
    staleTime: Infinity,
  });
}
