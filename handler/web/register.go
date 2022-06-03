package web

import (
	"github.com/myOmikron/hopfencloud/models/db"
	"github.com/myOmikron/hopfencloud/modules/logger"
	"github.com/myOmikron/hopfencloud/modules/tasks"

	"github.com/labstack/echo/v4"
	"github.com/myOmikron/echotools/database"
	"github.com/myOmikron/echotools/utilitymodels"
	"github.com/myOmikron/echotools/worker"
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
		return c.String(200, "Username must not be empty")
	}

	if req.Password == "" {
		//TODO: Display error message
		return c.String(200, "Password must not be empty")
	}

	if req.Email == "" {
		//TODO: Display error message
		return c.String(200, "Email must not be empty")
	}

	selected := make([]utilitymodels.LocalUser, 0)
	w.DB.Where("username = ? OR email = ?", req.Username, req.Email).Find(&utilitymodels.LocalUser{})
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

	localUser, err := database.CreateLocalUser(w.DB, req.Username, req.Password, &req.Email)
	if err != nil {
		logger.Error(err.Error())
		//TODO: Display error message
		return c.String(500, "User creation failed")
	}

	w.DB.Create(&db.User{
		AuthID:  localUser.ID,
		AuthKey: "local",
	})

	w.WorkerPool.AddTask(worker.NewTask(tasks.SendRegistrationMail))

	//TODO: Render prettier template
	return c.String(200, "Account created, email must be confirmed before you can login")
}
