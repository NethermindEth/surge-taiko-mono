export type RoleName = "admin" | "user" | (string & {});

export interface ApiErrorBody {
  error: { code: string; message: string };
}

export interface ChallengeResponse {
  message: string;
  expires_at: number;
}

export interface VerifyResponse {
  token: string;
  expires_at: number;
}

export interface Role {
  id: number;
  name: RoleName;
}

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

export type UpsertMemberRequest =
  | { role: "admin" }
  | {
      role: "user";
      attributes?: { kyc?: boolean; blacklisted?: boolean };
    };

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
  lambda_id: number | null;
  lambda_name: string | null;
}

export interface EntryInput {
  role: RoleName;
  lambda_id?: number | null;
}

export interface CreateRuleRequest {
  contract_address: string;
  function_selector: string;
  mode: "allow" | "deny";
  entries: EntryInput[];
}

export interface ReplaceRuleRequest {
  mode: "allow" | "deny";
  entries: EntryInput[];
}

export interface UpdateEntryRequest {
  lambda_id: number | null;
}

export type LhsKind = "calldata" | "attribute";
export type RhsKind = "tx_origin" | "msg_sender" | "literal";
export type Condition = "eq" | "neq" | "gt" | "lt" | "gte" | "lte";

export interface LambdaRuleView {
  id: number;
  selector: string;
  lhs_kind: LhsKind;
  lhs_offset: number | null;
  lhs_attribute: string | null;
  condition: Condition;
  rhs_kind: RhsKind;
  rhs_value: string | null;
}

export interface LambdaView {
  id: number;
  name: string;
  role: RoleName;
  description: string | null;
  rules: LambdaRuleView[];
  in_use: boolean;
}

export interface LambdaGroup {
  role: RoleName;
  lambdas: LambdaView[];
}

export interface LambdaRuleInput {
  selector: string;
  lhs_kind: LhsKind;
  lhs_offset?: number | null;
  lhs_attribute?: string | null;
  condition: Condition;
  rhs_kind: RhsKind;
  rhs_value?: string | null;
}

export interface CreateLambdaRequest {
  name: string;
  role: RoleName;
  description?: string | null;
  rules: LambdaRuleInput[];
}

export interface RoleAttribute {
  name: string;
  type: string;
}

export interface RoleAttributesGroup {
  role: RoleName;
  attributes: RoleAttribute[];
}

export interface SyntheticSelector {
  method: string;
  selector: string;
}

export interface RevokeTokensResponse {
  revoked: number;
}
