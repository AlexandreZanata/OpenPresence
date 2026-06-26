//! OpenPresence biometric service — liveness ensemble and face embedding extraction.

pub mod config;
pub mod grpc;
pub mod health;
pub mod pipeline;
pub mod processor;

pub use processor::{BiometricProcessor, ProcessorError, StubProcessor};
