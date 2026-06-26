use std::sync::Arc;
use std::time::Duration;

use biometric_service::grpc::proto::biometric_service_client::BiometricServiceClient;
use biometric_service::grpc::proto::{EnrollFaceRequest, VerifyPunchRequest};
use biometric_service::{grpc, StubProcessor};
use image::RgbImage;
use tokio_stream::wrappers::TcpListenerStream;
use tonic::transport::{Channel, Server};

fn tiny_jpeg() -> Vec<u8> {
    let img =
        image::DynamicImage::ImageRgb8(RgbImage::from_pixel(128, 128, image::Rgb([90, 120, 200])));
    let mut buf = Vec::new();
    let mut cursor = std::io::Cursor::new(&mut buf);
    img.write_to(&mut cursor, image::ImageFormat::Jpeg).unwrap();
    buf
}

#[tokio::test]
async fn grpc_verify_punch_and_enroll_roundtrip() {
    let processor = Arc::new(StubProcessor::new());
    let grpc_service = grpc::BiometricGrpcService::new(processor);
    let listener = tokio::net::TcpListener::bind("127.0.0.1:0").await.unwrap();
    let grpc_addr = listener.local_addr().unwrap();

    tokio::spawn(async move {
        Server::builder()
            .add_service(grpc::proto::biometric_service_server::BiometricServiceServer::new(
                grpc_service,
            ))
            .serve_with_incoming(TcpListenerStream::new(listener))
            .await
            .unwrap();
    });

    tokio::time::sleep(Duration::from_millis(150)).await;

    let channel = Channel::from_shared(format!("http://{grpc_addr}"))
        .unwrap()
        .connect()
        .await
        .unwrap();
    let mut client = BiometricServiceClient::new(channel);
    let jpeg = tiny_jpeg();

    let enroll = client
        .enroll_face(EnrollFaceRequest {
            frame_jpeg: jpeg.clone(),
            employee_id: "emp-1".into(),
            tenant_id: "tenant-1".into(),
            angle: "FRONTAL".into(),
        })
        .await
        .unwrap()
        .into_inner();
    assert!(enroll.liveness_score > 0.0);
    assert!(!enroll.embedding.is_empty());

    let verify = client
        .verify_punch(VerifyPunchRequest {
            frame_jpeg: jpeg,
            employee_id: "emp-1".into(),
            tenant_id: "tenant-1".into(),
        })
        .await
        .unwrap()
        .into_inner();
    assert_eq!(verify.matched_employee_id, "emp-1");
    assert!(!verify.embedding.is_empty());
}
