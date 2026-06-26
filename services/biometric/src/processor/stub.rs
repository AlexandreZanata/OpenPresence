use std::collections::HashMap;
use std::sync::{Arc, RwLock};

use image::{DynamicImage, GenericImageView};

use super::{
    decode_jpeg, embedding_from_seed, liveness_from_image, passes_enroll_liveness,
    passes_enroll_quality, passes_punch_liveness, BiometricProcessor, EnrollResult, FraudFlag,
    ProcessorError, VerifyResult,
};
use crate::config::RECOGNITION_THRESHOLD;
use crate::pipeline::{cosine_similarity, preprocess_for_recognition};

/// Deterministic processor for dev/test when ONNX models are not on disk.
pub struct StubProcessor {
    profiles: Arc<RwLock<HashMap<String, Vec<f32>>>>,
}

impl StubProcessor {
    pub fn new() -> Self {
        Self {
            profiles: Arc::new(RwLock::new(HashMap::new())),
        }
    }

    fn profile_key(tenant_id: &str, employee_id: &str) -> String {
        format!("{tenant_id}:{employee_id}")
    }

    fn quality_score(img: &DynamicImage) -> f32 {
        let (w, h) = img.dimensions();
        if w < 64 || h < 64 {
            return 0.3;
        }
        0.85
    }
}

impl Default for StubProcessor {
    fn default() -> Self {
        Self::new()
    }
}

impl BiometricProcessor for StubProcessor {
    fn verify_punch(
        &self,
        frame_jpeg: &[u8],
        employee_id: &str,
        tenant_id: &str,
    ) -> Result<VerifyResult, ProcessorError> {
        let img = decode_jpeg(frame_jpeg)?;
        let liveness_score = liveness_from_image(&img);
        let is_live = passes_punch_liveness(liveness_score);
        let embedding = embedding_from_seed(&format!("{tenant_id}:{employee_id}:frame"));
        let key = Self::profile_key(tenant_id, employee_id);

        let (is_recognized, recognition_confidence) = {
            let guard = self.profiles.read().unwrap();
            match guard.get(&key) {
                Some(stored) => {
                    let sim = cosine_similarity(&embedding, stored);
                    (sim >= RECOGNITION_THRESHOLD, sim)
                }
                None => (true, 0.95_f32),
            }
        };

        let mut fraud_flags = Vec::new();
        if !is_live {
            fraud_flags.push(FraudFlag {
                fraud_type: "LIVENESS_FAILED".into(),
                severity: "HIGH".into(),
                metadata_json: format!(r#"{{"liveness_score":{liveness_score}}}"#),
            });
        }

        Ok(VerifyResult {
            is_live,
            liveness_score,
            is_recognized,
            recognition_confidence,
            matched_employee_id: employee_id.to_string(),
            embedding,
            fraud_flags,
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
        let liveness_score = liveness_from_image(&img);
        let liveness_ok = passes_enroll_liveness(liveness_score);
        let quality_score = Self::quality_score(&img);
        let quality_ok = passes_enroll_quality(quality_score);
        let accepted = liveness_ok && quality_ok;
        let landmarks = [[0.0_f32; 2]; 5];
        let _ = preprocess_for_recognition(&img, &landmarks);
        let embedding = embedding_from_seed(&format!("{tenant_id}:{employee_id}:enroll"));

        if accepted {
            let key = Self::profile_key(tenant_id, employee_id);
            self.profiles.write().unwrap().insert(key, embedding.clone());
        }

        let mut fraud_flags = Vec::new();
        if !liveness_ok {
            fraud_flags.push(FraudFlag {
                fraud_type: "LIVENESS_FAILED".into(),
                severity: "HIGH".into(),
                metadata_json: format!(r#"{{"liveness_score":{liveness_score}}}"#),
            });
        }
        if !quality_ok {
            fraud_flags.push(FraudFlag {
                fraud_type: "LOW_QUALITY".into(),
                severity: "MEDIUM".into(),
                metadata_json: format!(r#"{{"quality_score":{quality_score}}}"#),
            });
        }

        Ok(EnrollResult {
            is_live: liveness_ok,
            liveness_score,
            quality_score,
            embedding: if accepted { embedding } else { vec![] },
            fraud_flags,
        })
    }

    fn delete_profile(&self, employee_id: &str, tenant_id: &str) -> Result<(), ProcessorError> {
        let key = Self::profile_key(tenant_id, employee_id);
        self.profiles.write().unwrap().remove(&key);
        Ok(())
    }

    fn is_ready(&self) -> bool {
        true
    }
}
