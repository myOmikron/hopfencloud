package tasks

import (
	"bytes"
	"fmt"
	"github.com/myOmikron/hopfencloud/modules/logger"
	"net"
	"net/smtp"
	"net/url"
	"strconv"
	"text/template"
	"time"

	"github.com/myOmikron/hopfencloud/models/db"

	"github.com/jordan-wright/email"
	"github.com/myOmikron/echotools/worker"
)

func SendRegistrationMail(
	receiver string,
	accountName string,
	token string,
	settings *db.Settings,
	mailTemplates *template.Template,
) worker.Task {
	return worker.NewTask(func() error {
		var body bytes.Buffer

		urlLink, _ := url.Parse(settings.PublicAddress)
		urlLink.Path = "/verify_email"
		q := urlLink.Query()
		q.Set("token", token)
		urlLink.RawQuery = q.Encode()

		data := make(map[string]string)
		data["SiteName"] = settings.SiteName
		data["RegistrationLink"] = urlLink.String()
		data["Username"] = accountName
		if err := mailTemplates.ExecuteTemplate(&body, "registration", data); err != nil {
			return err
		}

		// Authentication.
		authentication := smtp.PlainAuth("", settings.SMTPUser, settings.SMTPPassword, settings.SMTPHost)

		// Sending email.
		e := email.NewEmail()
		e.From = fmt.Sprintf("%s <%s>", settings.SMTPFrom, settings.SMTPUser)
		e.To = []string{receiver}
		e.Subject = "Hopfencloud Account Registration"
		e.Text = bytes.TrimSpace(body.Bytes())
		e.Headers = map[string][]string{
			"Date": {time.Now().Format(time.RFC1123Z)},
		}

		err := e.Send(net.JoinHostPort(settings.SMTPHost, strconv.Itoa(int(settings.SMTPPort))), authentication)
		if err != nil {
			logger.Error(err.Error())
		}

		return err
	})
}
