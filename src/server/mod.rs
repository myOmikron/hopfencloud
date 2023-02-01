use actix_toolbox::tb_middleware::{
    setup_logging_mw, DBSessionStore, LoggingMiddlewareConfig, SessionMiddleware,
};
use actix_web::cookie::Key;
use actix_web::middleware::Compress;
use actix_web::web::{get, post, scope, Data, JsonConfig, PayloadConfig};
use actix_web::{App, HttpServer};
use base64::prelude::BASE64_STANDARD;
use base64::Engine;
use rorm::Database;

use crate::config::Config;
use crate::handler::api;
use crate::server::error::ServerStartError;
use crate::server::middleware::authentication_required::AuthenticationRequired;

mod error;
mod middleware;

/**
Starts the server

**Parameter**:
- `config`: Reference to a [Config].
- `database`: [Database] instance
*/
pub(crate) async fn start_server(
    config: &Config,
    database: Database,
) -> Result<(), ServerStartError> {
    let key = Key::try_from(
        BASE64_STANDARD
            .decode(&config.server.secret_key)?
            .as_slice(),
    )?;

    HttpServer::new(move || {
        App::new()
            .wrap(setup_logging_mw(LoggingMiddlewareConfig::default()))
            .wrap(
                SessionMiddleware::builder(DBSessionStore::new(database.clone()), key.clone())
                    .cookie_name("session".to_string())
                    .build(),
            )
            .wrap(Compress::default())
            .app_data(Data::new(database.clone()))
            .app_data(PayloadConfig::default())
            .app_data(JsonConfig::default())
            .route("/api/v1/auth/login", post().to(api::login))
            .service(
                scope("/api/v1")
                    .wrap(AuthenticationRequired)
                    .route("test", get().to(api::test)),
            )
    })
    .bind((
        config.server.listen_address.as_str(),
        config.server.listen_port,
    ))?
    .run()
    .await?;

    Ok(())
}
