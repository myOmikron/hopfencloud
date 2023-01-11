use actix_toolbox::logging::LoggingConfig;
use rorm::config::DatabaseConfig;
use serde::Deserialize;

#[derive(Deserialize)]
#[serde(rename_all = "PascalCase")]
pub(crate) struct Server {
    pub(crate) listen_address: String,
    pub(crate) listen_port: u16,
    pub(crate) secret_key: String,
}

#[derive(Deserialize)]
#[serde(rename_all = "PascalCase")]
pub(crate) struct Config {
    pub(crate) server: Server,
    pub(crate) logging: LoggingConfig,
    pub(crate) database: DatabaseConfig,
}
