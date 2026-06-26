use std::sync::Arc;

use super::{BiometricProcessor, ProcessorError, StubProcessor};

/// Builds the active processor: stub when configured or ONNX when models are present.
pub fn build_processor() -> Result<Arc<dyn BiometricProcessor>, ProcessorError> {
    if crate::config::use_stub_processor() {
        tracing::warn!("BIOMETRIC_USE_STUB active — using deterministic stub processor");
        return Ok(Arc::new(StubProcessor::new()));
    }

    #[cfg(feature = "onnx")]
    {
        use crate::processor::onnx::OnnxProcessor;
        tracing::info!("Loading ONNX biometric models from ONNX_MODELS_PATH");
        let processor: Arc<dyn BiometricProcessor> = OnnxProcessor::from_env()?;
        Ok(processor)
    }

    #[cfg(not(feature = "onnx"))]
    Err(ProcessorError::NotReady(
        "ONNX models requested but binary built without `onnx` feature".into(),
    ))
}
