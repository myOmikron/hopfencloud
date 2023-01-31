//! This module holds the definitions for the configurations

use actix_toolbox::logging::LoggingConfig;
use rorm::config::DatabaseConfig;
use serde::Deserialize;

/// Configuration for server specifics
#[derive(Deserialize)]
#[serde(rename_all = "PascalCase")]
pub struct ServerConfig {
    /// Address the server binds to
    pub listen_address: String,
    /// Port the server binds to
    pub listen_port: u16,
    /// Secret key, used for cookie generation
    pub secret_key: String,
    /// The uri this server is accessible from the outside
    pub origin_uri: String,
    /// Name of the origin
    ///
    /// This name is only used by webauthn to display a user friendly name
    /// instead of the server uri
    pub origin_name: String,
}

/// Configuration
#[derive(Deserialize)]
#[serde(rename_all = "PascalCase")]
pub struct Config {
    /// Server configuration
    pub server: ServerConfig,
    /// Logging configuration
    pub logging: LoggingConfig,
    /// Database configuration
    pub database: DatabaseConfig,
}
