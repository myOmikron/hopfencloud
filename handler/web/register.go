package web

import (
	"net/mail"

	"github.com/myOmikron/hopfencloud/models/db"
	"github.com/myOmikron/hopfencloud/modules/crypt"
	"github.com/myOmikron/hopfencloud/modules/logger"
	"github.com/myOmikron/hopfencloud/modules/tasks"

	"github.com/labstack/echo/v4"
	"github.com/myOmikron/echotools/database"
	"github.com/myOmikron/echotools/utilitymodels"
)

type RegisterData struct {
	PageTitle string
}

func (w *Wrapper) RegisterGet(c echo.Context) error {
	return c.Render(200, "auth/register", &RegisterData{
		PageTitle: "Registration - " + w.Settings.SiteName,
	})
}

type RegisterRequest struct {
	Username string `form:"username"`
	Email    string `form:"email"`
	Password string `form:"password"`
}

func (w *Wrapper) RegisterPost(c echo.Context) error {
	var req RegisterRequest

	if err := c.Bind(&req); err != nil {
		//TODO: Display error message
		logger.Error(err.Error())
		return err
	}

	if req.Username == "" {
		//TODO: Display error message
		return c.String(400, "Username must not be empty")
	}

	if req.Password == "" {
		//TODO: Display error message
		return c.String(400, "Password must not be empty")
	}

	if req.Email == "" {
		//TODO: Display error message
		return c.String(400, "Email must not be empty")
	}

	address, err := mail.ParseAddress(req.Email)
	if err != nil {
		//TODO: Display error message
		return c.String(400, "Email is invalid")
	}

	selected := make([]utilitymodels.LocalUser, 0)
	w.DB.Find(&selected, "username = ? OR email = ?", req.Username, address.Address)

	if len(selected) != 0 {
		for _, user := range selected {
			if user.Username == req.Username {
				//TODO: Display error message
				return c.String(400, "Username already exists")
			}
		}
		//TODO: Display error message
		return c.String(400, "Email already exists")
	}

	var count int64
	w.DB.Where(&db.AccountEmailVerification{}, "email = ?", address.Address).Count(&count)
	if count != 0 {
		//TODO: Display error message
		return c.String(400, "Email already exists")
	}

	localUser, err := database.CreateLocalUser(w.DB, req.Username, req.Password, nil)
	if err != nil {
		logger.Error(err.Error())
		//TODO: Display error message
		return c.String(500, "Account creation failed")
	}

	account := db.Account{
		AuthID:  localUser.ID,
		AuthKey: "local",
	}
	w.DB.Create(&account)

	var token string
	for {
		token, err = crypt.GetToken()
		if err != nil {
			//TODO: Display error message
			return c.String(500, "Internal server error")
		}

		w.DB.Find(&db.AccountEmailVerification{}, "token = ?", token).Count(&count)
		if count == 0 {
			break
		}
	}

	w.DB.Create(&db.AccountEmailVerification{
		Account: account,
		Email:   address.Address,
		Token:   token,
	})

	w.WorkerPool.AddTask(tasks.SendRegistrationMail(address.Address, req.Username, token, w.Settings, w.MailTemplates))

	//TODO: Render prettier template
	return c.String(200, "Account created, email must be confirmed before you can login")
}
