use rorm::{BackRef, ForeignModel, Model, Patch};

use crate::models::File;

/**
User model
*/
#[derive(Model)]
pub struct User {
    /// Unique identifier of an user
    #[rorm(primary_key)]
    pub uuid: Vec<u8>,

    /// The username is used for log in as well identification by admins
    #[rorm(max_length = 255, unique)]
    pub username: String,

    /// This name is displayed for other users
    pub display_name: String,

    /// Hashed password
    #[rorm(max_length = 1024)]
    pub password_hash: String,

    /// Flag whether the user is an administrative user
    pub is_admin: bool,

    /// List of files the user owns
    #[rorm(field = "File::F.owner")]
    pub files: BackRef<File>,

    /// Datetime when the user was created
    #[rorm(auto_create_time)]
    pub created_at: chrono::NaiveDateTime,
    /// Datetime when the users' last login was
    pub last_login: Option<chrono::NaiveDateTime>,
}

#[derive(Patch)]
#[rorm(model = "User")]
pub(crate) struct UserInsert {
    pub(crate) username: String,
    pub(crate) password_hash: String,
    pub(crate) is_admin: bool,
    pub(crate) last_login: Option<chrono::NaiveDateTime>,
}

/**
Security key of a user.
*/
#[derive(Model)]
pub struct UserSecurityKey {
    /// Unique identifier of a security key
    #[rorm(id)]
    pub id: i64,

    /// Owner of the key
    #[rorm(on_delete = "Cascade", on_update = "Cascade")]
    pub user: ForeignModel<User>,

    /// Content of the key
    pub key: Vec<u8>,

    /// Time when the security key was added
    #[rorm(auto_create_time)]
    pub created_at: chrono::NaiveDateTime,
}

#[derive(Patch)]
#[rorm(model = "UserSecurityKey")]
pub(crate) struct UserSecurityKeyInsert {
    pub(crate) id: i64,
    pub(crate) user: ForeignModel<User>,
    pub(crate) key: Vec<u8>,
}
