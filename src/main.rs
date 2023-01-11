use std::fs;
use std::path::Path;

use actix_toolbox::logging::setup_logging;
use actix_web::cookie::Key;
use base64::prelude::BASE64_STANDARD;
use base64::Engine;
use clap::{Parser, Subcommand};

use crate::models::config::Config;
use crate::server::start_server;

mod models;
mod server;

#[derive(Subcommand)]
#[clap(version)]
pub(crate) enum Command {
    Start {
        #[clap(long = "config-path")]
        #[clap(default_value_t = String::from("/etc/hopfencloud/config.toml"))]
        #[clap(help = "Path to the configuration file")]
        config_path: String,
    },
    GenKey,
}

#[derive(Parser)]
pub(crate) struct Cli {
    #[clap(subcommand)]
    command: Command,
}

#[rorm::rorm_main]
#[tokio::main]
async fn main() -> Result<(), String> {
    let cli = Cli::parse();

    match cli.command {
        Command::Start { config_path } => {
            let config = read_config(&config_path)?;

            setup_logging(&config.logging)?;

            start_server(&config).await?;
        }
        Command::GenKey => {
            let key = Key::generate();
            let encoded = BASE64_STANDARD.encode(key.master());
            println!("{}", &encoded);
        }
    }

    Ok(())
}

/**
Reads a [Config] from the given path
*/
fn read_config(config_path: &str) -> Result<Config, String> {
    let path = Path::new(config_path);
    if !path.exists() {
        return Err(format!(
            "Did not found configuration file at {config_path}."
        ));
    }

    let content = fs::read_to_string(path)
        .map_err(|e| format!("Error while reading configuration file: {e}"))?;

    toml::from_str::<Config>(&content)
        .map_err(|e| format!("Error while parsing configuration file: {e}"))
}
