use image::DynamicImage;
use thiserror::Error;

use crate::config::{EMBEDDING_DIM, LIVENESS_THRESHOLD_ENROLL, LIVENESS_THRESHOLD_PUNCH};

#[derive(Debug, Error)]
pub enum ProcessorError {
    #[error("invalid frame: {0}")]
    InvalidFrame(String),
    #[error("onnx runtime: {0}")]
    Onnx(String),
    #[error("processor not ready: {0}")]
    NotReady(String),
}

#[derive(Debug, Clone)]
pub struct VerifyResult {
    pub is_live: bool,
    pub liveness_score: f32,
    pub is_recognized: bool,
    pub recognition_confidence: f32,
    pub matched_employee_id: String,
    pub embedding: Vec<f32>,
    pub fraud_flags: Vec<FraudFlag>,
}

#[derive(Debug, Clone)]
pub struct EnrollResult {
    pub is_live: bool,
    pub liveness_score: f32,
    pub quality_score: f32,
    pub embedding: Vec<f32>,
    pub fraud_flags: Vec<FraudFlag>,
}

#[derive(Debug, Clone)]
pub struct FraudFlag {
    pub fraud_type: String,
    pub severity: String,
    pub metadata_json: String,
}

pub trait BiometricProcessor: Send + Sync {
    fn verify_punch(
        &self,
        frame_jpeg: &[u8],
        employee_id: &str,
        tenant_id: &str,
    ) -> Result<VerifyResult, ProcessorError>;

    fn enroll_face(
        &self,
        frame_jpeg: &[u8],
        employee_id: &str,
        tenant_id: &str,
        angle: &str,
    ) -> Result<EnrollResult, ProcessorError>;

    fn delete_profile(&self, employee_id: &str, tenant_id: &str) -> Result<(), ProcessorError>;

    fn is_ready(&self) -> bool;
}

pub fn decode_jpeg(frame_jpeg: &[u8]) -> Result<DynamicImage, ProcessorError> {
    image::ImageReader::new(std::io::Cursor::new(frame_jpeg))
        .with_guessed_format()
        .map_err(|e| ProcessorError::InvalidFrame(e.to_string()))?
        .decode()
        .map_err(|e| ProcessorError::InvalidFrame(e.to_string()))
}

pub fn embedding_from_seed(seed: &str) -> Vec<f32> {
    let mut out = vec![0.0_f32; EMBEDDING_DIM];
    let bytes = seed.as_bytes();
    for (i, v) in out.iter_mut().enumerate() {
        let b = bytes[i % bytes.len()] as f32;
        *v = ((i as f32 * 0.13 + b).sin() + 1.0) / 2.0;
    }
    out
}

pub fn liveness_from_image(img: &DynamicImage) -> f32 {
    let data = crate::pipeline::preprocess_for_liveness(img);
    let mean: f32 = data.iter().sum::<f32>() / data.len().max(1) as f32;
    mean.clamp(0.0, 1.0)
}

pub fn passes_punch_liveness(score: f32) -> bool {
    score >= LIVENESS_THRESHOLD_PUNCH
}

pub fn passes_enroll_liveness(score: f32) -> bool {
    score >= LIVENESS_THRESHOLD_ENROLL
}

pub fn passes_enroll_quality(score: f32) -> bool {
    score >= crate::config::QUALITY_THRESHOLD_ENROLL
}

mod factory;
mod stub;

#[cfg(feature = "onnx")]
mod onnx;

pub use factory::build_processor;
pub use stub::StubProcessor;

#[cfg(feature = "onnx")]
pub use onnx::OnnxProcessor;
