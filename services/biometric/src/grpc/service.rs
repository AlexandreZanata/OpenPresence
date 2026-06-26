use std::sync::Arc;

use tonic::{Request, Response, Status};

use super::proto::biometric_service_server::BiometricService;
use super::proto::{
    DeleteProfileRequest, DeleteProfileResponse, EnrollFaceRequest, EnrollFaceResponse,
    FraudFlag as ProtoFraudFlag, VerifyPunchRequest, VerifyPunchResponse,
};
use crate::processor::{BiometricProcessor, ProcessorError};

pub struct BiometricGrpcService<P: BiometricProcessor> {
    processor: Arc<P>,
}

impl<P: BiometricProcessor> BiometricGrpcService<P> {
    pub fn new(processor: Arc<P>) -> Self {
        Self { processor }
    }
}

fn map_fraud_flags(flags: &[crate::processor::FraudFlag]) -> Vec<ProtoFraudFlag> {
    flags
        .iter()
        .map(|f| ProtoFraudFlag {
            fraud_type: f.fraud_type.clone(),
            severity: f.severity.clone(),
            metadata_json: f.metadata_json.clone(),
        })
        .collect()
}

fn embedding_to_bytes(embedding: &[f32]) -> Vec<u8> {
    let mut bytes = Vec::with_capacity(embedding.len() * 4);
    for v in embedding {
        bytes.extend_from_slice(&v.to_le_bytes());
    }
    bytes
}

fn map_processor_error(err: ProcessorError) -> Status {
    match err {
        ProcessorError::InvalidFrame(msg) => Status::invalid_argument(msg),
        ProcessorError::Onnx(msg) => Status::internal(msg),
        ProcessorError::NotReady(msg) => Status::failed_precondition(msg),
    }
}

#[tonic::async_trait]
impl<P: BiometricProcessor + 'static> BiometricService for BiometricGrpcService<P> {
    async fn verify_punch(
        &self,
        request: Request<VerifyPunchRequest>,
    ) -> Result<Response<VerifyPunchResponse>, Status> {
        let req = request.into_inner();
        let result = self
            .processor
            .verify_punch(&req.frame_jpeg, &req.employee_id, &req.tenant_id)
            .map_err(map_processor_error)?;

        Ok(Response::new(VerifyPunchResponse {
            is_live: result.is_live,
            liveness_score: result.liveness_score,
            is_recognized: result.is_recognized,
            recognition_confidence: result.recognition_confidence,
            matched_employee_id: result.matched_employee_id,
            fraud_flags: map_fraud_flags(&result.fraud_flags),
            embedding: embedding_to_bytes(&result.embedding),
        }))
    }

    async fn enroll_face(
        &self,
        request: Request<EnrollFaceRequest>,
    ) -> Result<Response<EnrollFaceResponse>, Status> {
        let req = request.into_inner();
        let result = self
            .processor
            .enroll_face(&req.frame_jpeg, &req.employee_id, &req.tenant_id, &req.angle)
            .map_err(map_processor_error)?;

        Ok(Response::new(EnrollFaceResponse {
            is_live: result.is_live,
            liveness_score: result.liveness_score,
            quality_score: result.quality_score,
            embedding: embedding_to_bytes(&result.embedding),
            fraud_flags: map_fraud_flags(&result.fraud_flags),
        }))
    }

    async fn delete_profile(
        &self,
        request: Request<DeleteProfileRequest>,
    ) -> Result<Response<DeleteProfileResponse>, Status> {
        let req = request.into_inner();
        self.processor
            .delete_profile(&req.employee_id, &req.tenant_id)
            .map_err(map_processor_error)?;
        Ok(Response::new(DeleteProfileResponse { success: true }))
    }
}
