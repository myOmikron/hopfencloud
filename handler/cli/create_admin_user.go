package cli

import (
	"errors"
	"net/mail"

	"github.com/myOmikron/hopfencloud/models/db"
	"github.com/myOmikron/hopfencloud/modules/logger"

	"github.com/myOmikron/echotools/database"
	"github.com/myOmikron/echotools/utilitymodels"
)

var (
	ErrUsernameEmpty         = errors.New("username must not be empty")
	ErrEmailEmpty            = errors.New("email must not be empty")
	ErrPasswordEmpty         = errors.New("password must not be empty")
	ErrEmailInvalid          = errors.New("email is not valid")
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrEmailAlreadyExists    = errors.New("email already exists")
)

type CreateAdminUserRequest struct {
	Username string
	Password string
	Email    string
}

type CreateAdminUserResponse struct {
	ErrorMessage *string
}

func (c *CLI) CreateAdminUser(req CreateAdminUserRequest, res *CreateAdminUserResponse) error {

	if req.Username == "" {
		return ErrUsernameEmpty
	}

	if req.Password == "" {
		return ErrPasswordEmpty
	}

	if req.Email == "" {
		return ErrEmailEmpty
	}

	address, err := mail.ParseAddress(req.Email)
	if err != nil {
		return ErrEmailInvalid
	}

	selected := make([]utilitymodels.LocalUser, 0)
	c.DB.Find(&selected, "username = ? OR email = ?", req.Username, address.Address)

	if len(selected) != 0 {
		for _, user := range selected {
			if user.Username == req.Username {
				return ErrUsernameAlreadyExists
			}
		}
		return ErrEmailAlreadyExists
	}

	var count int64
	c.DB.Where(&db.AccountEmailVerification{}, "email = ?", address.Address).Count(&count)
	if count != 0 {
		return ErrEmailAlreadyExists
	}

	localUser, err := database.CreateLocalUser(c.DB, req.Username, req.Password, &address.Address)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	account := db.Account{
		AuthID:  localUser.ID,
		AuthKey: "local",
	}
	c.DB.Create(&account)

	return nil
}
