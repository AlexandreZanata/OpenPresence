//! Pure math for liveness ensemble and embedding comparison.

/// Cosine similarity between two equal-length vectors. Returns 0.0 if either norm is zero.
pub fn cosine_similarity(a: &[f32], b: &[f32]) -> f32 {
    if a.len() != b.len() || a.is_empty() {
        return 0.0;
    }
    let mut dot = 0.0_f32;
    let mut norm_a = 0.0_f32;
    let mut norm_b = 0.0_f32;
    for i in 0..a.len() {
        dot += a[i] * b[i];
        norm_a += a[i] * a[i];
        norm_b += b[i] * b[i];
    }
    if norm_a == 0.0 || norm_b == 0.0 {
        return 0.0;
    }
    dot / (norm_a.sqrt() * norm_b.sqrt())
}

fn softmax(logits: &[f32; 3]) -> [f32; 3] {
    let max = logits[0].max(logits[1]).max(logits[2]);
    let exps = [
        (logits[0] - max).exp(),
        (logits[1] - max).exp(),
        (logits[2] - max).exp(),
    ];
    let sum = exps[0] + exps[1] + exps[2];
    [exps[0] / sum, exps[1] / sum, exps[2] / sum]
}

/// Average MiniFASNet V2 + V1SE softmax outputs. Index 0 = live class.
pub fn ensemble_liveness(v2_logits: &[f32; 3], v1se_logits: &[f32; 3]) -> (bool, f32) {
    let v2 = softmax(v2_logits);
    let v1se = softmax(v1se_logits);
    let live_score = (v2[0] + v1se[0]) / 2.0;
    let is_live = live_score >= crate::config::LIVENESS_THRESHOLD_PUNCH;
    (is_live, live_score)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_cosine_similarity_identical_vectors() {
        let v = vec![1.0_f32; 512];
        let sim = cosine_similarity(&v, &v);
        assert!((sim - 1.0).abs() < 1e-5);
    }

    #[test]
    fn test_face_recognition_same_person() {
        let mut a = vec![0.0_f32; 512];
        let mut b = vec![0.0_f32; 512];
        for i in 0..512 {
            a[i] = (i as f32 * 0.01).sin();
            b[i] = (i as f32 * 0.01).sin() * 0.98 + 0.01;
        }
        assert!(cosine_similarity(&a, &b) >= 0.75);
    }

    #[test]
    fn test_face_recognition_different_persons() {
        let mut a = vec![0.0_f32; 512];
        let mut b = vec![0.0_f32; 512];
        for i in 0..512 {
            a[i] = (i as f32).sin();
            b[i] = (i as f32).cos();
        }
        assert!(cosine_similarity(&a, &b) < 0.65);
    }

    #[test]
    fn test_liveness_real_face_above_threshold() {
        let live_v2 = [3.0_f32, 0.1, 0.1];
        let live_v1se = [2.8_f32, 0.2, 0.1];
        let (is_live, score) = ensemble_liveness(&live_v2, &live_v1se);
        assert!(is_live);
        assert!(score >= 0.85);
    }

    #[test]
    fn test_liveness_printed_photo_rejected() {
        let spoof_v2 = [0.1_f32, 3.0, 0.2];
        let spoof_v1se = [0.2_f32, 2.5, 0.3];
        let (is_live, score) = ensemble_liveness(&spoof_v2, &spoof_v1se);
        assert!(!is_live);
        assert!(score < 0.80);
    }
}
