package web

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/myOmikron/echotools/auth"
	"github.com/myOmikron/echotools/middleware"
	"github.com/myOmikron/echotools/utilitymodels"
)

type LoginProvider struct {
	Name                 string
	RegistrationDisabled bool
	Selected             bool
	Identifier           string
}

type LoginData struct {
	PageTitle               string
	LoginProvider           []LoginProvider
	ForgotPasswordAvailable bool
	ForgotPasswordLink      string
	SelectedProvider        string
	RegistrationEnabled     bool
	RedirectTo              string
}

type LoginArg struct {
	Provider   string `query:"provider"`    // Will be empty for local auth
	RedirectTo string `query:"redirect_to"` // If empty, login will redirect to "/"
}

var ReLdapProvider = regexp.MustCompile(`ldap(\d+)`)

func (w *Wrapper) LoginGet(c echo.Context) error {
	localProvider := LoginProvider{
		Name:                 "Local",
		RegistrationDisabled: w.Config.General.RegistrationDisabled,
		Identifier:           "",
	}

	lp := make([]LoginProvider, 0)
	lp = append(lp, localProvider)

	forgotPasswordAvailable := false
	forgotPasswordLink := ""
	registrationEnabled := !w.Config.General.RegistrationDisabled

	ldapProviders := make([]utilitymodels.LDAPProvider, 0)
	w.DB.Find(&ldapProviders)

	loginArg := LoginArg{}
	if err := c.Bind(&loginArg); err != nil {
		//TODO: Render template with error message
		return err
	}
	switch {
	case loginArg.Provider == "":
		lp[0].Selected = true
		forgotPasswordAvailable = true
		forgotPasswordLink = "/forgotPassword"

		for _, ldapProvider := range ldapProviders {
			tmpProvider := LoginProvider{
				Name:       "LDAP - " + ldapProvider.Name,
				Identifier: fmt.Sprintf("ldap%d", ldapProvider.ID),
			}
			lp = append(lp, tmpProvider)
		}
	case ReLdapProvider.MatchString(loginArg.Provider):
		submatch := ReLdapProvider.FindStringSubmatch(loginArg.Provider)[1]
		number, err := strconv.Atoi(submatch)
		if err != nil {
			//TODO: Render template with error message
			return err
		}

		for _, ldapProvider := range ldapProviders {
			tmpProvider := LoginProvider{
				Name:       "LDAP - " + ldapProvider.Name,
				Identifier: fmt.Sprintf("ldap%d", ldapProvider.ID),
			}
			if int(ldapProvider.ID) == number {
				tmpProvider.Selected = true
				registrationEnabled = false
			}
			lp = append(lp, tmpProvider)
		}
	default:
		//TODO: Display "invalid login provider selected"
	}

	if loginArg.RedirectTo == "" {
		loginArg.RedirectTo = "/"
	}

	return c.Render(200, "login", &LoginData{
		PageTitle:               "Login - " + w.Config.General.SiteName,
		LoginProvider:           lp,
		ForgotPasswordLink:      forgotPasswordLink,
		ForgotPasswordAvailable: forgotPasswordAvailable,
		SelectedProvider:        loginArg.Provider,
		RegistrationEnabled:     registrationEnabled,
		RedirectTo:              loginArg.RedirectTo,
	})
}

type LoginRequest struct {
	Username   string `form:"username"`
	Password   string `form:"password"`
	Provider   string `form:"provider"`    // Will be an empty string for local auth
	RememberMe string `form:"remember_me"` // Will be "on" if checkbox was set
	RedirectTo string `form:"redirect_to"` // If empty, the user will be redirected to this page
}

func (w *Wrapper) LoginPost(c echo.Context) error {
	req := LoginRequest{}
	if err := c.Bind(&req); err != nil {
		//TODO: Display error page
		return err
	}

	if req.Username == "" || req.Password == "" {
		//TODO: Display error page
		return c.String(400, "Username or password was empty")
	}
	if req.RememberMe != "" && req.RememberMe != "on" {
		//TODO: Display error page
		return c.String(400, "Invalid value for remember_me")
	}

	switch {
	case req.Provider == "":
		user, err := auth.AuthenticateLocalUser(w.DB, req.Username, req.Password)
		if err != nil {
			//TODO: Display error message
			return c.String(400, err.Error())
		}

		if user == nil {
			//TODO: Display error message
			return c.String(500, "Internal server error")
		}

		if err := middleware.Login(w.DB, user, c, req.RememberMe == ""); err != nil {
			//TODO: Display error message
			return err
		}

	case ReLdapProvider.MatchString(req.Provider):
		//TODO: Let user login
		return c.String(500, "Not implemented yet")
	default:
		//TODO: Display error: Unknown Login Provider
		return c.String(500, "Not implemented yet")
	}

	return c.Redirect(302, req.RedirectTo)
}
