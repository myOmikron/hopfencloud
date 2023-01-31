use rorm::{ForeignModel, Model, Patch};

use crate::models::User;

/// Representation of a file
#[derive(Model)]
pub struct File {
    /// Unique identifier of a file
    #[rorm(id)]
    pub id: i64,

    /// Owner of the file
    #[rorm(on_delete = "Cascade", on_update = "Cascade")]
    pub owner: ForeignModel<User>,

    /// Filename
    #[rorm(max_length = 1024)]
    pub file_name: String,

    /// Parent of the file
    #[rorm(on_delete = "Cascade", on_update = "Cascade")]
    pub parent: Option<ForeignModel<File>>,
}

#[derive(Patch)]
#[rorm(model = "File")]
pub(crate) struct FileInsert {
    pub(crate) file_name: String,
    pub(crate) parent: Option<ForeignModel<File>>,
    pub(crate) owner: ForeignModel<User>,
}
