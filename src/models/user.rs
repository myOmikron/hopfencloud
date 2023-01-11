use rorm::{BackRef, ForeignModel, Model, Patch};

use crate::models::File;

/**
User model
*/
#[derive(Model)]
pub(crate) struct User {
    #[rorm(id)]
    pub(crate) id: i64,

    #[rorm(max_length = 255)]
    pub(crate) username: String,
    #[rorm(max_length = 1024)]
    pub(crate) password: String,
    pub(crate) is_admin: bool,

    #[rorm(field = "File::F.owner")]
    pub(crate) files: BackRef<File>,

    #[rorm(auto_create_time)]
    pub(crate) created_at: chrono::NaiveDateTime,
    pub(crate) last_login: Option<chrono::NaiveDateTime>,
}

#[derive(Patch)]
#[rorm(model = "User")]
pub(crate) struct UserInsert {
    pub(crate) username: String,
    pub(crate) password: String,
    pub(crate) is_admin: bool,
    pub(crate) last_login: Option<chrono::NaiveDateTime>,
}

/**
Security key of a user.
*/
#[derive(Model)]
pub(crate) struct UserSecurityKey {
    #[rorm(max_length = 38, primary_key)]
    pub(crate) uuid: String,

    #[rorm(on_delete = "Cascade", on_update = "Cascade")]
    pub(crate) user: ForeignModel<User>,
    pub(crate) key: Vec<u8>,

    #[rorm(auto_create_time)]
    pub(crate) created_at: chrono::NaiveDateTime,
}

#[derive(Patch)]
#[rorm(model = "UserSecurityKey")]
pub(crate) struct UserSecurityKeyInsert {
    pub(crate) uuid: String,
    pub(crate) user: ForeignModel<User>,
    pub(crate) key: Vec<u8>,
}
