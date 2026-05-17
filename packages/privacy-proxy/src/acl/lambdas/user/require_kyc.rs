use super::UserCallerInfo;
use crate::acl::lambdas::LambdaCtx;

pub fn run(ctx: &LambdaCtx<UserCallerInfo>) -> bool {
    ctx.caller_info.kyc
}

#[cfg(test)]
mod tests {
    use super::*;
    use alloy::primitives::Address;

    fn info(kyc: bool) -> UserCallerInfo {
        UserCallerInfo {
            eoa: Address::ZERO,
            kyc,
            blacklisted: false,
        }
    }

    #[test]
    fn allows_when_kyc_true() {
        let i = info(true);
        assert!(run(&LambdaCtx {
            caller_info: &i,
            selector: [0; 4],
            call_data: &[],
        }));
    }

    #[test]
    fn denies_when_kyc_false() {
        let i = info(false);
        assert!(!run(&LambdaCtx {
            caller_info: &i,
            selector: [0; 4],
            call_data: &[],
        }));
    }
}
