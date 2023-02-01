use std::fmt::{Display, Formatter};
use std::io;

use actix_web::cookie::KeyError;

#[derive(Debug)]
pub(crate) enum ServerStartError {
    IO(io::Error),
    Base64DecodingFailed(base64::DecodeError),
    InvalidSecretKey(KeyError),
}

impl Display for ServerStartError {
    fn fmt(&self, f: &mut Formatter<'_>) -> std::fmt::Result {
        match self {
            ServerStartError::IO(err) => write!(f, "IOError: {err}"),
            ServerStartError::Base64DecodingFailed(err) => write!(f, "Base64 decode error: {err}"),
            ServerStartError::InvalidSecretKey(err) => write!(f, "Invalid SecretKey: {err}"),
        }
    }
}

impl From<io::Error> for ServerStartError {
    fn from(value: io::Error) -> Self {
        Self::IO(value)
    }
}

impl From<base64::DecodeError> for ServerStartError {
    fn from(value: base64::DecodeError) -> Self {
        Self::Base64DecodingFailed(value)
    }
}

impl From<KeyError> for ServerStartError {
    fn from(value: KeyError) -> Self {
        Self::InvalidSecretKey(value)
    }
}
