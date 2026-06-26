use std::collections::HashMap;
use std::sync::{Arc, RwLock};

use image::GenericImageView;

use super::super::{
    decode_jpeg, passes_enroll_liveness, passes_enroll_quality, BiometricProcessor, EnrollResult,
    FraudFlag, ProcessorError, VerifyResult,
};
use super::face_processor::{crop_face, FaceProcessor};
use crate::config::RECOGNITION_THRESHOLD;
use crate::pipeline::cosine_similarity;

pub struct OnnxProcessor {
    face: FaceProcessor,
    profiles: Arc<RwLock<HashMap<String, Vec<f32>>>>,
}

impl OnnxProcessor {
    pub fn from_env() -> Result<Arc<Self>, ProcessorError> {
        let face = FaceProcessor::from_env()?;
        Ok(Arc::new(Self {
            face,
            profiles: Arc::new(RwLock::new(HashMap::new())),
        }))
    }

    fn profile_key(tenant_id: &str, employee_id: &str) -> String {
        format!("{tenant_id}:{employee_id}")
    }

    fn quality_score(w: u32, h: u32) -> f32 {
        if w < 64 || h < 64 {
            0.3
        } else {
            0.85
        }
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
        let det = self.face.detect_face(&img)?;
        let face = crop_face(&img, &det);
        let (is_live, liveness_score) = self.face.liveness_score(&face)?;
        let embedding = self.face.embed(&face, &det.landmarks)?;
        let key = Self::profile_key(tenant_id, employee_id);

        let (is_recognized, recognition_confidence) = {
            let guard = self.profiles.read().unwrap();
            match guard.get(&key) {
                Some(stored) => {
                    let sim = cosine_similarity(&embedding, stored);
                    (sim >= RECOGNITION_THRESHOLD, sim)
                }
                None => (true, 1.0),
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
        let (w, h) = img.dimensions();
        let det = self.face.detect_face(&img)?;
        let face = crop_face(&img, &det);
        let (is_live, liveness_score) = self.face.liveness_score(&face)?;
        let embedding = self.face.embed(&face, &det.landmarks)?;
        let quality_score = Self::quality_score(w, h);
        let liveness_ok = is_live && passes_enroll_liveness(liveness_score);
        let quality_ok = passes_enroll_quality(quality_score);
        let accepted = liveness_ok && quality_ok;

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
