//! gRPC enrollment E2E — BR-001 (angles), BR-002 (liveness), BR-003 (quality).

use std::sync::Arc;
use std::time::Duration;

use biometric_service::config::{LIVENESS_THRESHOLD_ENROLL, QUALITY_THRESHOLD_ENROLL};
use biometric_service::grpc::proto::biometric_service_client::BiometricServiceClient;
use biometric_service::grpc::proto::EnrollFaceRequest;
use biometric_service::processor::BiometricProcessor;
use biometric_service::{grpc, StubProcessor};
use image::{Rgb, RgbImage};
use tokio_stream::wrappers::TcpListenerStream;
use tonic::transport::{Channel, Server};

const ANGLES: [&str; 3] = ["FRONTAL", "LEFT_15", "RIGHT_15"];

fn jpeg_from_rgb(w: u32, h: u32, color: Rgb<u8>) -> Vec<u8> {
    let img = image::DynamicImage::ImageRgb8(RgbImage::from_pixel(w, h, color));
    let mut buf = Vec::new();
    let mut cursor = std::io::Cursor::new(&mut buf);
    img.write_to(&mut cursor, image::ImageFormat::Jpeg).unwrap();
    buf
}

fn valid_enroll_jpeg() -> Vec<u8> {
    jpeg_from_rgb(128, 128, Rgb([240, 240, 240]))
}

fn low_liveness_jpeg() -> Vec<u8> {
    jpeg_from_rgb(128, 128, Rgb([0, 0, 0]))
}

fn low_quality_jpeg() -> Vec<u8> {
    jpeg_from_rgb(32, 32, Rgb([240, 240, 240]))
}

async fn spawn_client() -> (BiometricServiceClient<Channel>, tokio::task::JoinHandle<()>) {
    let processor: Arc<dyn BiometricProcessor> = Arc::new(StubProcessor::new());
    let grpc_service = grpc::BiometricGrpcService::new(processor);
    let listener = tokio::net::TcpListener::bind("127.0.0.1:0").await.unwrap();
    let grpc_addr = listener.local_addr().unwrap();

    let server = tokio::spawn(async move {
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
    (BiometricServiceClient::new(channel), server)
}

async fn enroll(
    client: &mut BiometricServiceClient<Channel>,
    jpeg: Vec<u8>,
    angle: &str,
    employee_id: &str,
) -> biometric_service::grpc::proto::EnrollFaceResponse {
    client
        .enroll_face(EnrollFaceRequest {
            frame_jpeg: jpeg,
            employee_id: employee_id.into(),
            tenant_id: "tenant-e2e".into(),
            angle: angle.into(),
        })
        .await
        .unwrap()
        .into_inner()
}

#[tokio::test]
async fn enrollment_e2e_br001_three_angles_success() {
    let (mut client, _server) = spawn_client().await;
    let jpeg = valid_enroll_jpeg();

    for angle in ANGLES {
        let resp = enroll(&mut client, jpeg.clone(), angle, "emp-br001").await;
        assert!(resp.is_live, "angle {angle} should pass liveness");
        assert!(resp.liveness_score >= LIVENESS_THRESHOLD_ENROLL);
        assert!(resp.quality_score >= QUALITY_THRESHOLD_ENROLL);
        assert!(!resp.embedding.is_empty(), "angle {angle} should store embedding");
        assert!(resp.fraud_flags.is_empty(), "angle {angle} should have no fraud flags");
    }
}

#[tokio::test]
async fn enrollment_e2e_br002_liveness_fail_rejected() {
    let (mut client, _server) = spawn_client().await;
    let resp = enroll(
        &mut client,
        low_liveness_jpeg(),
        "FRONTAL",
        "emp-br002",
    )
    .await;

    assert!(!resp.is_live);
    assert!(resp.liveness_score < LIVENESS_THRESHOLD_ENROLL);
    assert!(resp.embedding.is_empty());
    assert!(
        resp.fraud_flags
            .iter()
            .any(|f| f.fraud_type == "LIVENESS_FAILED")
    );
}

#[tokio::test]
async fn enrollment_e2e_br003_low_quality_rejected() {
    let (mut client, _server) = spawn_client().await;
    let resp = enroll(
        &mut client,
        low_quality_jpeg(),
        "FRONTAL",
        "emp-br003",
    )
    .await;

    assert!(resp.is_live, "liveness may pass on small bright frame");
    assert!(resp.quality_score < QUALITY_THRESHOLD_ENROLL);
    assert!(resp.embedding.is_empty());
    assert!(
        resp.fraud_flags
            .iter()
            .any(|f| f.fraud_type == "LOW_QUALITY")
    );
}

/// Writes JPEG fixtures to `tests/fixtures/` for grpcurl manual runs (`--ignored`).
#[tokio::test]
#[ignore = "run once to refresh tests/fixtures/*.jpg"]
async fn write_enrollment_fixture_jpegs() {
    use std::fs;
    use std::path::PathBuf;

    let dir = PathBuf::from(env!("CARGO_MANIFEST_DIR")).join("tests/fixtures");
    fs::create_dir_all(&dir).unwrap();
    fs::write(dir.join("valid_128.jpg"), valid_enroll_jpeg()).unwrap();
    fs::write(dir.join("low_liveness_128.jpg"), low_liveness_jpeg()).unwrap();
    fs::write(dir.join("low_quality_32.jpg"), low_quality_jpeg()).unwrap();
}
