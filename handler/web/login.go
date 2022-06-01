package web

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/myOmikron/echotools/utilitymodels"
	"regexp"
	"strconv"
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
}

type LoginArg struct {
	Provider string `query:"provider"`
}

var ReLdapProvider = regexp.MustCompile(`ldap(\d+)`)

func (w *Wrapper) Login(c echo.Context) error {
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

	return c.Render(200, "login", &LoginData{
		PageTitle:               "Login - " + w.Config.General.SiteName,
		LoginProvider:           lp,
		ForgotPasswordLink:      forgotPasswordLink,
		ForgotPasswordAvailable: forgotPasswordAvailable,
		SelectedProvider:        loginArg.Provider,
		RegistrationEnabled:     registrationEnabled,
	})
}
