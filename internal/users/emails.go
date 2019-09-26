package users

import (
	"bytes"
	"context"
	fmt "fmt"
	"html/template"

	emails "github.com/tcfw/evntsrc/internal/emails/protos"
	users "github.com/tcfw/evntsrc/internal/users/protos"
)

const (
	emailTemplatesRecreation = `
<p>Hi,</p>
<p><br/></p>
<p>We received a request to create an account using this email, however an account
already exists with the matching email on the {{.CreatedAt.String}}.</p>
<p><br/></p>
<p>If you believe this is incorrect or a secuity concern, please contact us 
at hello@evntsrc.io.</p>
<p><br/></p>
<p>If you have forgotten your password for this account, please make a password
reset request on the forgot password page.</p>
<p><br/></p>
<p>Warmest regards,<br/>
Evntsrc.io</p>
`
)

func (s *server) sendRecreationEmail(user *users.User) error {
	emailBody, err := renderTemplate(emailTemplatesRecreation, user)
	if err != nil {
		return err
	}

	_, err = s.emails.Send(context.Background(), &emails.Email{
		To: []*emails.Recipient{
			&emails.Recipient{Email: user.Email, Name: user.Name},
		},
		Subject: "You've already created an account",
		Html:    emailBody,
		Headers: map[string]string{
			"X-Priority": "High",
		},
	})

	return err
}

func renderTemplate(templateString string, data interface{}) (string, error) {
	temp, err := template.New("emailTemplate").Parse(templateString)
	if err != nil {
		return "", fmt.Errorf("Failed to init template: %s", err)
	}

	buff := &bytes.Buffer{}

	if err := temp.Execute(buff, data); err != nil {
		return "", fmt.Errorf("Failed to execute template: %s", err)
	}

	return buff.String(), nil
}
