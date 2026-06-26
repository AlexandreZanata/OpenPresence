//! Thresholds aligned with docs/BUSINESS-RULES.md (BR-002, BR-010).

pub const LIVENESS_THRESHOLD_PUNCH: f32 = 0.80;
pub const LIVENESS_THRESHOLD_ENROLL: f32 = 0.85;
pub const QUALITY_THRESHOLD_ENROLL: f32 = 0.70;
pub const RECOGNITION_THRESHOLD: f32 = 0.75;
pub const EMBEDDING_DIM: usize = 512;

pub const ENV_MODELS_PATH: &str = "ONNX_MODELS_PATH";
pub const ENV_GRPC_ADDR: &str = "BIOMETRIC_GRPC_ADDR";
pub const ENV_HTTP_ADDR: &str = "BIOMETRIC_HTTP_ADDR";
pub const ENV_USE_STUB: &str = "BIOMETRIC_USE_STUB";

pub fn grpc_addr() -> String {
    std::env::var(ENV_GRPC_ADDR).unwrap_or_else(|_| "0.0.0.0:9090".into())
}

pub fn http_addr() -> String {
    std::env::var(ENV_HTTP_ADDR).unwrap_or_else(|_| "0.0.0.0:9091".into())
}

pub fn use_stub_processor() -> bool {
    std::env::var(ENV_USE_STUB)
        .map(|v| v == "1" || v.eq_ignore_ascii_case("true"))
        .unwrap_or_else(|_| std::env::var(ENV_MODELS_PATH).is_err())
}
