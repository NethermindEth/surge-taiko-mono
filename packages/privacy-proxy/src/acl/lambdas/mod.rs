pub mod attributes;
pub mod eval;
pub mod loader;

use serde::{Deserialize, Serialize};

#[derive(Clone, Copy, Debug, Eq, PartialEq, Serialize, Deserialize)]
#[serde(rename_all = "lowercase")]
pub enum Condition {
    Eq,
    Neq,
    Gt,
    Lt,
    Gte,
    Lte,
}

impl Condition {
    pub fn from_db(s: &str) -> Option<Self> {
        match s {
            "eq" => Some(Self::Eq),
            "neq" => Some(Self::Neq),
            "gt" => Some(Self::Gt),
            "lt" => Some(Self::Lt),
            "gte" => Some(Self::Gte),
            "lte" => Some(Self::Lte),
            _ => None,
        }
    }

    pub fn as_db(&self) -> &'static str {
        match self {
            Self::Eq => "eq",
            Self::Neq => "neq",
            Self::Gt => "gt",
            Self::Lt => "lt",
            Self::Gte => "gte",
            Self::Lte => "lte",
        }
    }
}

#[derive(Clone, Debug, Eq, PartialEq, Serialize, Deserialize)]
#[serde(tag = "kind", rename_all = "lowercase")]
pub enum Lhs {
    Calldata { offset: u32 },
    Attribute { name: String },
}

#[derive(Clone, Debug, Eq, PartialEq, Serialize, Deserialize)]
#[serde(tag = "kind", rename_all = "snake_case")]
pub enum Rhs {
    TxOrigin,
    MsgSender,
    Literal {
        #[serde(rename = "value")]
        value_hex: String,
    },
}

#[derive(Clone, Debug, Serialize)]
pub struct LambdaRule {
    pub id: i64,
    pub selector: [u8; 4],
    pub lhs: Lhs,
    pub condition: Condition,
    pub rhs: Rhs,
}

#[derive(Clone, Debug, Serialize)]
pub struct Lambda {
    pub id: i64,
    pub name: String,
    pub role_id: i64,
    pub role: String,
    pub description: Option<String>,
    pub rules: Vec<LambdaRule>,
}
