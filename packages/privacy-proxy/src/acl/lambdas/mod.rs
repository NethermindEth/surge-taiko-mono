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
pub struct LambdaSpec<C: 'static> {
    pub name: &'static str,
    pub description: &'static str,
    pub expected_selector: Option<[u8; 4]>,
    pub run: fn(&LambdaCtx<C>) -> bool,
}
