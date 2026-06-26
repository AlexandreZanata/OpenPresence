pub mod service;

pub mod proto {
    tonic::include_proto!("openpresence.biometric.v1");
}

pub use service::BiometricGrpcService;
