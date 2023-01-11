use std::io;

use actix_web::cookie::KeyError;

#[derive(Debug)]
pub(crate) enum ServerError {
    IO(io::Error),
    Base64DecodingFailed(base64::DecodeError),
    KeyError(KeyError),
}

impl From<io::Error> for ServerError {
    fn from(value: io::Error) -> Self {
        Self::IO(value)
    }
}

impl From<base64::DecodeError> for ServerError {
    fn from(value: base64::DecodeError) -> Self {
        Self::Base64DecodingFailed(value)
    }
}

impl From<KeyError> for ServerError {
    fn from(value: KeyError) -> Self {
        Self::KeyError(value)
    }
}

impl From<ServerError> for String {
    fn from(value: ServerError) -> Self {
        match value {
            ServerError::IO(err) => format!("IOError: {err}"),
            ServerError::Base64DecodingFailed(err) => format!("Base64 decode error: {err}"),
            ServerError::KeyError(err) => format!("Invalid SecretKey: {err}"),
        }
    }
}
