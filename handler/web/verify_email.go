package web

import (
	"github.com/myOmikron/hopfencloud/models/db"
	"github.com/myOmikron/hopfencloud/modules/logger"

	"github.com/labstack/echo/v4"
	"github.com/myOmikron/echotools/utilitymodels"
)

type VerifyEmailArg struct {
	Token string `query:"token"`
}

type VerifyEmailData struct {
	PageTitle string
	Token     string
}

func (w *Wrapper) VerifyEmailGet(c echo.Context) error {
	var arg VerifyEmailArg
	if err := c.Bind(&arg); err != nil {
		//TODO: Display error message
		logger.Error(err.Error())
		return c.String(500, "Internal server error")
	}

	return c.Render(200, "auth/verify_email", &VerifyEmailData{
		PageTitle: "Verify email - " + w.Settings.SiteName,
		Token:     arg.Token,
	})
}

type VerifyEmailRequest struct {
	Token string `form:"token"`
}

func (w *Wrapper) VerifyEmailPost(c echo.Context) error {
	var arg VerifyEmailRequest

	if err := c.Bind(&arg); err != nil {
		//TODO: Display error message
		logger.Error(err.Error())
		return c.String(500, "Internal server error")
	}

	if arg.Token == "" {
		//TODO: Display error message
		return c.String(400, "Invalid token")
	}

	var count int64
	userMailConfirmation := db.UserMailConfirmation{}
	w.DB.Preload("User").Find(&userMailConfirmation, "token = ?", arg.Token).Count(&count)
	if count != 1 {
		//TODO: Display error message
		return c.String(400, "Invalid token")
	}

	if userMailConfirmation.User.AuthKey != "local" {
		//TODO: Invalid user for mail change
		logger.Info("Non local user tried mail change")
		return c.String(500, "Internal server error")
	}

	var user utilitymodels.LocalUser
	w.DB.Find(&user, "id = ?", userMailConfirmation.User.AuthID).Count(&count)
	if count != 1 {
		//TODO: Invalid internal user model
		logger.Info("No local user was found.")
		return c.String(500, "Internal server error")
	}

	user.Email = &userMailConfirmation.Mail
	w.DB.Save(&user)
	w.DB.Delete(&userMailConfirmation)

	return c.Redirect(302, "/login")
}
