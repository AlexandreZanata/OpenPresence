use std::net::SocketAddr;

use biometric_service::config::{grpc_addr, http_addr};
use biometric_service::grpc::proto::biometric_service_server::BiometricServiceServer;
use biometric_service::grpc::BiometricGrpcService;
use biometric_service::health;
use biometric_service::processor::build_processor;
use tonic::transport::Server;
use tracing_subscriber::EnvFilter;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    tracing_subscriber::fmt()
        .with_env_filter(EnvFilter::from_default_env().add_directive("info".parse()?))
        .init();

    let processor = build_processor()?;
    let grpc_service = BiometricGrpcService::new(processor.clone());

    let grpc_socket: SocketAddr = grpc_addr().parse()?;
    let http_socket: SocketAddr = http_addr().parse()?;

    let grpc_server = Server::builder()
        .add_service(BiometricServiceServer::new(grpc_service))
        .serve(grpc_socket);

    let app = health::router(processor);
    let http_server = async move {
        let listener = tokio::net::TcpListener::bind(http_socket).await?;
        tracing::info!("HTTP health listening on {http_socket}");
        axum::serve(listener, app).await
    };

    tracing::info!("gRPC listening on {grpc_socket}");
    let (grpc_res, http_res) = tokio::join!(grpc_server, http_server);
    grpc_res?;
    http_res?;
    Ok(())
}
