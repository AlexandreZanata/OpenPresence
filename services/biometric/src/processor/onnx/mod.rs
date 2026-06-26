//! ONNX Runtime face pipeline — detection, liveness ensemble, embedding.

mod face_processor;
mod processor_impl;

pub use processor_impl::OnnxProcessor;
