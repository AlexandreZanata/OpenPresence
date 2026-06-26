//! Loads ONNX sessions and runs RetinaFace/SCRFD, MiniFASNet×2, AuraFace inference.

use std::path::{Path, PathBuf};
use std::sync::Mutex;

use image::{DynamicImage, GenericImageView};
use ort::session::Session;
use ort::value::Tensor;

use crate::config::EMBEDDING_DIM;
use crate::pipeline::{ensemble_liveness, preprocess_for_liveness, preprocess_for_recognition};
use crate::processor::ProcessorError;

const DETECTOR_SIZE: u32 = 640;
const LIVENESS_SIZE: u32 = 80;
const RECOGNITION_SIZE: u32 = 112;

/// Face detection output — bounding box and five landmarks.
#[derive(Debug, Clone, Copy)]
pub struct DetectResult {
    pub x1: f32,
    pub y1: f32,
    pub x2: f32,
    pub y2: f32,
    pub landmarks: [[f32; 2]; 5],
}

pub struct FaceProcessor {
    inner: Mutex<FaceProcessorInner>,
}

struct FaceProcessorInner {
    detector: Session,
    liveness_v2: Session,
    liveness_v1se: Session,
    embedder: Session,
}

impl FaceProcessor {
    pub fn from_models(path: impl AsRef<Path>) -> Result<Self, ProcessorError> {
        let dir = path.as_ref();
        let detector = load_session(&dir.join("retinaface.onnx"))?;
        let liveness_v2 = load_session(&dir.join("minifasnet_v2.onnx"))?;
        let liveness_v1se = load_session(&dir.join("minifasnet_v1se.onnx"))?;
        let embedder = load_session(&dir.join("auraface.onnx"))?;
        Ok(Self {
            inner: Mutex::new(FaceProcessorInner {
                detector,
                liveness_v2,
                liveness_v1se,
                embedder,
            }),
        })
    }

    pub fn from_env() -> Result<Self, ProcessorError> {
        let path = std::env::var(crate::config::ENV_MODELS_PATH)
            .map_err(|_| ProcessorError::NotReady("ONNX_MODELS_PATH not set".into()))?;
        Self::from_models(PathBuf::from(path))
    }

    pub fn detect_face(&self, img: &DynamicImage) -> Result<DetectResult, ProcessorError> {
        let (w, h) = img.dimensions();
        let input = detector_tensor(img)?;
        let mut guard = self.inner.lock().map_err(|_| ProcessorError::Onnx("lock poisoned".into()))?;
        let outputs = guard
            .detector
            .run(ort::inputs![input])
            .map_err(|e| ProcessorError::Onnx(e.to_string()))?;
        parse_detector_output(&outputs, w, h).or_else(|_| Ok(fallback_detect(w, h)))
    }

    pub fn liveness_score(&self, face: &DynamicImage) -> Result<(bool, f32), ProcessorError> {
        let v2_logits = self.run_liveness_model(face, true)?;
        let v1se_logits = self.run_liveness_model(face, false)?;
        Ok(ensemble_liveness(&v2_logits, &v1se_logits))
    }

    pub fn embed(&self, face: &DynamicImage, landmarks: &[[f32; 2]; 5]) -> Result<Vec<f32>, ProcessorError> {
        let data = preprocess_for_recognition(face, landmarks);
        let input = recognition_tensor(&data)?;
        let mut guard = self.inner.lock().map_err(|_| ProcessorError::Onnx("lock poisoned".into()))?;
        let outputs = guard
            .embedder
            .run(ort::inputs![input])
            .map_err(|e| ProcessorError::Onnx(e.to_string()))?;
        extract_embedding(&outputs)
    }

    fn run_liveness_model(&self, face: &DynamicImage, v2: bool) -> Result<[f32; 3], ProcessorError> {
        let resized = face.resize_exact(LIVENESS_SIZE, LIVENESS_SIZE, image::imageops::FilterType::Triangle);
        let data = preprocess_for_liveness(&resized);
        let input = liveness_tensor(&data)?;
        let mut guard = self.inner.lock().map_err(|_| ProcessorError::Onnx("lock poisoned".into()))?;
        let outputs = if v2 {
            guard.liveness_v2.run(ort::inputs![input])
        } else {
            guard.liveness_v1se.run(ort::inputs![input])
        }
        .map_err(|e| ProcessorError::Onnx(e.to_string()))?;
        logits3(&outputs)
    }
}

fn load_session(path: &Path) -> Result<Session, ProcessorError> {
    if !path.is_file() {
        return Err(ProcessorError::NotReady(format!("missing model: {}", path.display())));
    }
    Session::builder()
        .map_err(|e| ProcessorError::Onnx(e.to_string()))?
        .commit_from_file(path)
        .map_err(|e| ProcessorError::Onnx(e.to_string()))
}

fn hwc_to_nchw(data: &[f32], h: u32, w: u32) -> (Vec<usize>, Vec<f32>) {
    let shape = vec![1_usize, 3, h as usize, w as usize];
    let mut arr = vec![0.0_f32; shape.iter().product()];
    for y in 0..h as usize {
        for x in 0..w as usize {
            let base = (y * w as usize + x) * 3;
            let plane = h as usize * w as usize;
            let idx = |c: usize| y * w as usize + x + c * plane;
            arr[idx(0)] = data[base];
            arr[idx(1)] = data[base + 1];
            arr[idx(2)] = data[base + 2];
        }
    }
    (shape, arr)
}

fn liveness_tensor(data: &[f32]) -> Result<Tensor<f32>, ProcessorError> {
    let (shape, arr) = hwc_to_nchw(data, LIVENESS_SIZE, LIVENESS_SIZE);
    Tensor::from_array((shape, arr)).map_err(|e| ProcessorError::Onnx(e.to_string()))
}

fn recognition_tensor(data: &[f32]) -> Result<Tensor<f32>, ProcessorError> {
    let (shape, arr) = hwc_to_nchw(data, RECOGNITION_SIZE, RECOGNITION_SIZE);
    Tensor::from_array((shape, arr)).map_err(|e| ProcessorError::Onnx(e.to_string()))
}

fn detector_tensor(img: &DynamicImage) -> Result<Tensor<f32>, ProcessorError> {
    let resized = img.resize_exact(DETECTOR_SIZE, DETECTOR_SIZE, image::imageops::FilterType::Triangle);
    let rgb = resized.to_rgb8();
    let shape = vec![1_usize, 3, DETECTOR_SIZE as usize, DETECTOR_SIZE as usize];
    let mut arr = vec![0.0_f32; shape.iter().product()];
    let plane = DETECTOR_SIZE as usize * DETECTOR_SIZE as usize;
    for y in 0..DETECTOR_SIZE as usize {
        for x in 0..DETECTOR_SIZE as usize {
            let p = rgb.get_pixel(x as u32, y as u32);
            for c in 0..3 {
                let idx = y * DETECTOR_SIZE as usize + x + c * plane;
                arr[idx] = (p[c] as f32 - 127.5) / 128.0;
            }
        }
    }
    Tensor::from_array((shape, arr)).map_err(|e| ProcessorError::Onnx(e.to_string()))
}

fn logits3(outputs: &ort::session::SessionOutputs) -> Result<[f32; 3], ProcessorError> {
    let (_shape, data) = outputs[0]
        .try_extract_tensor::<f32>()
        .map_err(|e| ProcessorError::Onnx(e.to_string()))?;
    if data.len() < 3 {
        return Err(ProcessorError::Onnx("liveness output too small".into()));
    }
    Ok([data[0], data[1], data[2]])
}

fn extract_embedding(outputs: &ort::session::SessionOutputs) -> Result<Vec<f32>, ProcessorError> {
    let (_shape, data) = outputs[0]
        .try_extract_tensor::<f32>()
        .map_err(|e| ProcessorError::Onnx(e.to_string()))?;
    let mut embedding = data.to_vec();
    if embedding.len() > EMBEDDING_DIM {
        embedding.truncate(EMBEDDING_DIM);
    }
    while embedding.len() < EMBEDDING_DIM {
        embedding.push(0.0);
    }
    l2_normalize(&mut embedding);
    Ok(embedding)
}

fn l2_normalize(v: &mut [f32]) {
    let norm = v.iter().map(|x| x * x).sum::<f32>().sqrt();
    if norm > 0.0 {
        for x in v.iter_mut() {
            *x /= norm;
        }
    }
}

fn parse_detector_output(
    outputs: &ort::session::SessionOutputs,
    orig_w: u32,
    orig_h: u32,
) -> Result<DetectResult, ProcessorError> {
    let mut best_score = 0.0_f32;
    let mut best_box = None;
    for idx in 0..outputs.len() {
        let (_shape, data) = outputs[idx]
            .try_extract_tensor::<f32>()
            .map_err(|e| ProcessorError::Onnx(e.to_string()))?;
        if data.len() >= 5 && data.len() <= 16 {
            continue;
        }
        if data.len() % 4 == 0 && data.len() >= 4 {
            for chunk in data.chunks(4) {
                let score = chunk.get(4).copied().unwrap_or(1.0);
                if score > best_score {
                    best_score = score;
                    best_box = Some([chunk[0], chunk[1], chunk[2], chunk[3]]);
                }
            }
        }
    }
    if let Some([x1, y1, x2, y2]) = best_box {
        let sx = orig_w as f32 / DETECTOR_SIZE as f32;
        let sy = orig_h as f32 / DETECTOR_SIZE as f32;
        return Ok(DetectResult {
            x1: x1 * sx,
            y1: y1 * sy,
            x2: x2 * sx,
            y2: y2 * sy,
            landmarks: default_landmarks(x1 * sx, y1 * sy, x2 * sx, y2 * sy),
        });
    }
    Err(ProcessorError::Onnx("no face detected".into()))
}

fn fallback_detect(w: u32, h: u32) -> DetectResult {
    let margin_x = w as f32 * 0.15;
    let margin_y = h as f32 * 0.15;
    let x1 = margin_x;
    let y1 = margin_y;
    let x2 = w as f32 - margin_x;
    let y2 = h as f32 - margin_y;
    DetectResult {
        x1,
        y1,
        x2,
        y2,
        landmarks: default_landmarks(x1, y1, x2, y2),
    }
}

fn default_landmarks(x1: f32, y1: f32, x2: f32, y2: f32) -> [[f32; 2]; 5] {
    let mx = (x1 + x2) / 2.0;
    let my = (y1 + y2) / 2.0;
    let w = x2 - x1;
    let h = y2 - y1;
    [
        [x1 + w * 0.3, y1 + h * 0.35],
        [x1 + w * 0.7, y1 + h * 0.35],
        [mx, my],
        [x1 + w * 0.3, y1 + h * 0.75],
        [x1 + w * 0.7, y1 + h * 0.75],
    ]
}

pub fn crop_face(img: &DynamicImage, det: &DetectResult) -> DynamicImage {
    let (w, h) = img.dimensions();
    let x1 = det.x1.max(0.0) as u32;
    let y1 = det.y1.max(0.0) as u32;
    let x2 = det.x2.min(w as f32) as u32;
    let y2 = det.y2.min(h as f32) as u32;
    let cw = x2.saturating_sub(x1).max(1);
    let ch = y2.saturating_sub(y1).max(1);
    let cropped = img.crop_imm(x1, y1, cw, ch);
    DynamicImage::ImageRgb8(cropped.to_rgb8())
}

#[cfg(all(test, feature = "onnx"))]
mod tests {
    use super::*;
    use image::RgbImage;

    fn models_available() -> Option<String> {
        std::env::var("ONNX_MODELS_PATH").ok()
    }

    #[test]
    #[ignore = "requires ONNX models — run with ONNX_MODELS_PATH set"]
    fn real_inference_returns_512_embedding() {
        let path = models_available().expect("ONNX_MODELS_PATH");
        let processor = FaceProcessor::from_models(&path).expect("load models");
        let img = DynamicImage::ImageRgb8(RgbImage::from_pixel(320, 320, image::Rgb([120, 90, 80])));
        let det = processor.detect_face(&img).expect("detect");
        let face = crop_face(&img, &det);
        let (is_live, score) = processor.liveness_score(&face).expect("liveness");
        let emb = processor.embed(&face, &det.landmarks).expect("embed");
        assert_eq!(emb.len(), 512);
        assert!(score >= 0.0 && score <= 1.0);
        let _ = is_live;
    }
}
