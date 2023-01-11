use actix_toolbox::tb_middleware::{setup_logging_mw, LoggingMiddlewareConfig};
use actix_web::{App, HttpServer};
use log::info;

use crate::models::config::Config;
use crate::server::error::ServerError;

mod error;

/**
Starts the server
*/
pub(crate) async fn start_server(config: &Config) -> Result<(), ServerError> {
    info!(
        "Starting to listen on http://{}:{}",
        &config.server.listen_address, config.server.listen_port
    );

    HttpServer::new(|| App::new().wrap(setup_logging_mw(LoggingMiddlewareConfig::default())))
        .bind((
            config.server.listen_address.as_str(),
            config.server.listen_port,
        ))?
        .run()
        .await?;

    Ok(())
}
