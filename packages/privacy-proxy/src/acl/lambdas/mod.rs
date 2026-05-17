pub mod user;

/// Per-call context handed to a lambda. The attribute type `C` is the
/// typed `*CallerInfo` struct for the caller's role; lambdas under each
/// role's directory pick a concrete `C` (e.g. `UserCallerInfo`) and read
/// fields directly instead of probing JSON.
pub struct LambdaCtx<'a, C> {
    pub caller_info: &'a C,
    pub selector: [u8; 4],
    pub call_data: &'a [u8],
}

/// In-build metadata for a single lambda. Each role's `registry()` builds
/// a `HashMap<&str, &LambdaSpec<RoleAttrs>>`; entries are referenced by
/// `name` from `access_rule_entries.lambda_name`.
///
/// `expected_selectors` is the set of 4-byte selectors the lambda is built
/// to evaluate. An empty slice means the lambda is selector-agnostic
/// (e.g. attribute-only predicates like `require_kyc`). Admins authoring
/// rules use this to confirm a lambda is paired with a compatible
/// `function_selector`; the proxy itself does not enforce the match.
pub struct LambdaSpec<C: 'static> {
    pub name: &'static str,
    pub description: &'static str,
    pub expected_selectors: &'static [[u8; 4]],
    pub run: fn(&LambdaCtx<C>) -> bool,
}
