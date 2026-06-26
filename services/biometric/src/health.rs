use std::sync::Arc;

use axum::http::StatusCode;
use axum::response::IntoResponse;
use axum::routing::get;
use axum::{Json, Router};
use serde::Serialize;

use crate::processor::BiometricProcessor;

#[derive(Serialize)]
struct HealthBody {
    status: &'static str,
}

pub fn router(processor: Arc<dyn BiometricProcessor>) -> Router {
    let ready = processor.clone();
    Router::new()
        .route(
            "/health/live",
            get(|| async { (StatusCode::OK, Json(HealthBody { status: "ok" })) }),
        )
        .route(
            "/health/ready",
            get(move || {
                let proc = ready.clone();
                async move {
                    if proc.is_ready() {
                        (StatusCode::OK, Json(HealthBody { status: "ready" })).into_response()
                    } else {
                        (
                            StatusCode::SERVICE_UNAVAILABLE,
                            Json(HealthBody { status: "not_ready" }),
                        )
                            .into_response()
                    }
                }
            }),
        )
}
