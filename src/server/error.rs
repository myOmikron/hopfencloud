use std::io;

#[derive(Debug)]
pub(crate) enum ServerError {
    IO(io::Error),
}

impl From<io::Error> for ServerError {
    fn from(value: io::Error) -> Self {
        Self::IO(value)
    }
}

impl From<ServerError> for String {
    fn from(value: ServerError) -> Self {
        match value {
            ServerError::IO(err) => format!("IOError: {err}"),
        }
    }
}
