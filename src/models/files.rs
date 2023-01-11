use rorm::{ForeignModel, Model, Patch};

use crate::models::User;

/**
Representation of a file
*/
#[derive(Model)]
pub(crate) struct File {
    #[rorm(id)]
    pub(crate) id: i64,

    #[rorm(on_delete = "Cascade", on_update = "Cascade")]
    pub(crate) owner: ForeignModel<User>,

    #[rorm(max_length = 1024)]
    pub(crate) file_name: String,
    #[rorm(on_delete = "Cascade", on_update = "Cascade")]
    pub(crate) parent: Option<ForeignModel<File>>,
}

#[derive(Patch)]
#[rorm(model = "File")]
pub(crate) struct FileInsert {
    pub(crate) file_name: String,
    pub(crate) parent: Option<ForeignModel<File>>,
    pub(crate) owner: ForeignModel<User>,
}
