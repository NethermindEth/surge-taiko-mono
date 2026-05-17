// Mirrors of the privacy-proxy response shapes. Kept hand-written and small
// so changes on the server are loud here too. See:
// packages/privacy-proxy/docs/admin-api.md
// packages/privacy-proxy/src/admin/{members,registry,lambdas,roles}.rs

export type RoleName = "admin" | "user" | (string & {});

/** Server error body for any 4xx/5xx. */
export interface ApiErrorBody {
  error: { code: string; message: string };
}

/** GET /auth/challenge */
export interface ChallengeResponse {
  message: string;
  expires_at: number;
}

/** POST /auth/verify */
export interface VerifyResponse {
  token: string;
  expires_at: number;
}

/** Capability 11 — GET /admin/roles */
export interface Role {
  id: number;
  name: RoleName;
}

/** Capability 12, 13, 14 — member shape. `attributes` is null for admin. */
export interface MemberView {
  eoa_address: string;
  role: RoleName;
  attributes: UserAttributes | null;
  created_at: number;
}

export interface UserAttributes {
  kyc: boolean;
  blacklisted: boolean;
}

/** Capability 14 body — PUT /admin/members/:eoa */
export type UpsertMemberRequest =
  | { role: "admin" }
  | {
      role: "user";
      // Both fields are individually optional. Omitted fields preserve the
      // row's current value; defaults to false on first insert.
      attributes?: { kyc?: boolean; blacklisted?: boolean };
    };

/** Capability 1, 2 — rule shape. */
export interface RuleView {
  id: number;
  contract_address: string;
  function_selector: string;
  mode: "allow" | "deny";
  entries: EntryView[];
}

export interface EntryView {
  id: number;
  role: RoleName;
  lambda_name: string | null;
}

export interface EntryInput {
  role: RoleName;
  lambda_name?: string | null;
}

/** Capability 3 — POST /admin/registry/rules */
export interface CreateRuleRequest {
  contract_address: string;
  function_selector: string;
  mode: "allow" | "deny";
  entries: EntryInput[];
}

/** Capability 4 — PUT /admin/registry/rules/:id */
export interface ReplaceRuleRequest {
  mode: "allow" | "deny";
  entries: EntryInput[];
}

/** Capability 7 — PUT /admin/registry/rules/:id/entries/:entry_id */
export interface UpdateEntryRequest {
  lambda_name: string | null;
}

/** Capability 9 — GET /admin/registry/lambdas */
export interface LambdaGroup {
  role: RoleName;
  lambdas: LambdaDescriptor[];
}

export interface LambdaDescriptor {
  name: string;
  description: string;
  /** 4-byte selectors the lambda is built to evaluate. Empty means selector-agnostic. */
  expected_selectors: string[];
}

/** Capability 10 — GET /admin/registry/synthetic-selectors */
export interface SyntheticSelector {
  method: string;
  selector: string;
}

/** Capability 16 — DELETE /admin/members/:eoa/tokens */
export interface RevokeTokensResponse {
  revoked: number;
}
