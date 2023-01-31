use std::fmt::{Display, Formatter};

use actix_toolbox::tb_middleware::actix_session::SessionInsertError;
use actix_web::body::BoxBody;
use actix_web::HttpResponse;
use log::{error, trace};
use serde::Serialize;
use serde_repr::Serialize_repr;

pub(crate) use crate::handler::api::auth::{login, test};

mod auth;

pub(crate) type ApiResult<T> = Result<T, ApiError>;

#[derive(Serialize)]
pub(crate) struct ErrorResponse {
    status_code: ApiStatusCode,
    message: String,
}

impl ErrorResponse {
    fn new(status_code: ApiStatusCode, message: String) -> Self {
        Self {
            status_code,
            message,
        }
    }
}

#[derive(Serialize_repr)]
#[repr(u16)]
pub(crate) enum ApiStatusCode {
    LoginFailed = 1000,
    InternalServerError = 2000,
    DatabaseError = 2001,
    SessionError = 2002,
}

#[derive(Debug)]
pub(crate) enum ApiError {
    LoginFailed,
    DatabaseError(rorm::Error),
    InvalidHash(argon2::password_hash::Error),
    SessionInsert(SessionInsertError),
}

impl From<rorm::Error> for ApiError {
    fn from(value: rorm::Error) -> Self {
        Self::DatabaseError(value)
    }
}

impl From<argon2::password_hash::Error> for ApiError {
    fn from(value: argon2::password_hash::Error) -> Self {
        Self::InvalidHash(value)
    }
}

impl From<SessionInsertError> for ApiError {
    fn from(value: SessionInsertError) -> Self {
        Self::SessionInsert(value)
    }
}

impl Display for ApiError {
    fn fmt(&self, f: &mut Formatter<'_>) -> std::fmt::Result {
        match self {
            ApiError::LoginFailed => write!(f, "Login failed"),
            ApiError::DatabaseError(_) => write!(f, "Database error occurred"),
            ApiError::InvalidHash(_) => write!(f, "Internal server error"),
            ApiError::SessionInsert(_) => write!(f, "Session error occurred"),
        }
    }
}

impl actix_web::ResponseError for ApiError {
    fn error_response(&self) -> HttpResponse<BoxBody> {
        match self {
            ApiError::LoginFailed => {
                trace!("Login failed");
                HttpResponse::Ok().json(ErrorResponse::new(
                    ApiStatusCode::LoginFailed,
                    self.to_string(),
                ))
            }
            ApiError::DatabaseError(err) => {
                error!("Database error: {err}");
                HttpResponse::Ok().json(ErrorResponse::new(
                    ApiStatusCode::DatabaseError,
                    self.to_string(),
                ))
            }
            ApiError::InvalidHash(err) => {
                error!("Hashing error: {err}");
                HttpResponse::Ok().json(ErrorResponse::new(
                    ApiStatusCode::InternalServerError,
                    self.to_string(),
                ))
            }
            ApiError::SessionInsert(err) => {
                error!("Session insert: {err}");
                HttpResponse::Ok().json(ErrorResponse::new(
                    ApiStatusCode::SessionError,
                    self.to_string(),
                ))
            }
        }
    }
}
