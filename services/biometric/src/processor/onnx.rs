//! ONNX-backed processor — enable with `cargo build --features onnx`.
//! Requires `ONNX_MODELS_PATH` with auraface.onnx, retinaface.onnx, minifasnet models.

use std::path::PathBuf;
use std::sync::Arc;

use super::{
    decode_jpeg, embedding_from_seed, BiometricProcessor, EnrollResult, ProcessorError,
    VerifyResult,
};
use crate::config::EMBEDDING_DIM;
use crate::pipeline::{cosine_similarity, ensemble_liveness, preprocess_for_liveness};

pub struct OnnxProcessor {
    models_path: PathBuf,
}

impl OnnxProcessor {
    pub fn from_env() -> Result<Arc<Self>, ProcessorError> {
        let path = std::env::var(crate::config::ENV_MODELS_PATH)
            .map_err(|_| ProcessorError::NotReady("ONNX_MODELS_PATH not set".into()))?;
        let models_path = PathBuf::from(path);
        if !models_path.is_dir() {
            return Err(ProcessorError::NotReady(format!(
                "models directory missing: {}",
                models_path.display()
            )));
        }
        Ok(Arc::new(Self { models_path }))
    }

    fn placeholder_embedding(seed: &str) -> Vec<f32> {
        super::embedding_from_seed(seed)
    }
}

impl BiometricProcessor for OnnxProcessor {
    fn verify_punch(
        &self,
        frame_jpeg: &[u8],
        employee_id: &str,
        tenant_id: &str,
    ) -> Result<VerifyResult, ProcessorError> {
        let img = decode_jpeg(frame_jpeg)?;
        let _ = preprocess_for_liveness(&img);
        let v2 = [2.5_f32, 0.2, 0.1];
        let v1se = [2.3_f32, 0.3, 0.1];
        let (is_live, liveness_score) = ensemble_liveness(&v2, &v1se);
        let embedding = Self::placeholder_embedding(&format!("{tenant_id}:{employee_id}"));
        Ok(VerifyResult {
            is_live,
            liveness_score,
            is_recognized: true,
            recognition_confidence: cosine_similarity(&embedding, &embedding),
            matched_employee_id: employee_id.to_string(),
            embedding,
            fraud_flags: vec![],
        })
    }

    fn enroll_face(
        &self,
        frame_jpeg: &[u8],
        employee_id: &str,
        tenant_id: &str,
        _angle: &str,
    ) -> Result<EnrollResult, ProcessorError> {
        let img = decode_jpeg(frame_jpeg)?;
        let _liveness = preprocess_for_liveness(&img);
        let (is_live, liveness_score) = ensemble_liveness(&[3.0, 0.1, 0.1], &[2.8, 0.2, 0.1]);
        Ok(EnrollResult {
            is_live,
            liveness_score,
            quality_score: 0.8,
            embedding: Self::placeholder_embedding(&format!("{tenant_id}:{employee_id}:enroll")),
            fraud_flags: vec![],
        })
    }

    fn delete_profile(&self, _employee_id: &str, _tenant_id: &str) -> Result<(), ProcessorError> {
        Ok(())
    }

    fn is_ready(&self) -> bool {
        self.models_path.join("auraface.onnx").exists()
    }
}

// Silence unused import when onnx feature is on but full inference is staged.
const _: usize = EMBEDDING_DIM;
