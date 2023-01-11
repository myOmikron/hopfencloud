use actix_toolbox::tb_middleware::Session;
use actix_web::web::{Data, Json};
use argon2::{Argon2, PasswordHash, PasswordVerifier};
use rorm::{query, Database, Model};
use serde::{Deserialize, Serialize};

use crate::handler::api::{ApiError, ApiResult};
use crate::models::User;

#[derive(Deserialize)]
pub(crate) struct LoginRequest {
    username: String,
    password: String,
}

#[derive(Serialize)]
pub(crate) struct LoginResponse {
    success: bool,
}

pub(crate) async fn login(
    req: Json<LoginRequest>,
    db: Data<Database>,
    session: Session,
) -> ApiResult<Json<LoginResponse>> {
    if let Some(user) = query!(&db, User)
        .condition(User::F.username.equals(&req.username))
        .optional()
        .await?
    {
        Argon2::default()
            .verify_password(req.password.as_bytes(), &PasswordHash::new(&user.password)?)
            .map_err(|_| ApiError::LoginFailed)?;

        Ok(Json(LoginResponse { success: true }))
    } else {
        Err(ApiError::LoginFailed)
    }
}
