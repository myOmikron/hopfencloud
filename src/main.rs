//!
//! # Hopfencloud
//!
//! A file sharing and storage platform with multiuser support.
//!
#![deny(missing_docs)]
// Disable unused_imports and dead_code while generating models
// as the generation overwrites the main function and performs an early exit
#![cfg_attr(feature = "rorm-main", allow(unused_imports, dead_code))]

use std::fs;
use std::io;
use std::io::Write;
use std::path::Path;
use std::process::exit;
use std::time::Duration;

use actix_toolbox::logging::setup_logging;
use actix_web::cookie::Key;
use argon2::password_hash::rand_core::OsRng;
use argon2::password_hash::SaltString;
use argon2::{Argon2, PasswordHasher};
use base64::prelude::BASE64_STANDARD;
use base64::Engine;
use clap::{Parser, Subcommand};
use log::error;
use rorm::{insert, query, Model};
use rorm::{Database, DatabaseConfiguration};
use tokio::time::sleep;

use crate::config::Config;
use crate::models::{User, UserInsert};
use crate::server::start_server;

mod config;
mod handler;
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
    CreateAdminUser {
        #[clap(long = "config-path")]
        #[clap(default_value_t = String::from("/etc/hopfencloud/config.toml"))]
        #[clap(help = "Path to the configuration file")]
        config_path: String,
    },
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
            let db = init_db(&config).await?;

            start_server(&config, db).await?;
        }
        Command::GenKey => {
            let key = Key::generate();
            let encoded = BASE64_STANDARD.encode(key.master());
            println!("{}", &encoded);
        }
        Command::CreateAdminUser { config_path } => {
            let config = read_config(&config_path)?;
            setup_logging(&config.logging)?;
            let db = init_db(&config).await?;

            let stdin = io::stdin();
            let mut stdout = io::stdout();

            let mut username = String::new();

            print!("Enter a username: ");
            stdout.flush().unwrap();
            stdin.read_line(&mut username).unwrap();
            let username = username.trim();

            if query!(&db, (User::F.username,))
                .condition(User::F.username.equals(username))
                .optional()
                .await
                .unwrap()
                .is_some()
            {
                eprintln!("There is already a user with that name");
                exit(1);
            }

            let password = rpassword::prompt_password("Enter password: ").unwrap();

            let salt = SaltString::generate(&mut OsRng);
            let hashed_password = Argon2::default()
                .hash_password(password.as_bytes(), &salt)
                .unwrap()
                .to_string();

            insert!(&db, UserInsert)
                .single(&UserInsert {
                    username: username.to_owned(),
                    password: hashed_password,
                    is_admin: true,
                    last_login: None,
                })
                .await
                .map_err(|e| format!("Failed to create user: {e}"))?;

            sleep(Duration::from_millis(500)).await;

            println!("Created user {username}");
        }
    }

    Ok(())
}

/**
Initialize the connection to the database
*/
async fn init_db(config: &Config) -> Result<Database, String> {
    let mut db_connect_options = DatabaseConfiguration::new(config.database.driver.clone());

    // Disables the logging from sqlx
    db_connect_options.disable_logging = Some(true);

    Database::connect(db_connect_options).await.map_err(|e| {
        error!("Error while initializing the database: {e}");

        e.to_string()
    })
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
