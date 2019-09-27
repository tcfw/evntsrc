package users

import (
	"bytes"
	"context"
	"encoding/base64"
	fmt "fmt"
	"html/template"
	"net/url"

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
	emailTemplatesValidation = `
<p>Welcome!</p>
<p>Please click the button below to verify your email</p>
<div style="text-align: center">
<a href="https://staging.evntsrc.io/verify/{{.}}" style="color: white; background: #5F539B; display: inline-block; padding: 10px 30px; border-radius: 4px; text-decoration: none; font-weight: bold;">Verify</a>
</div>
<p><br/></p>
<p>If you did not create an account, please disregard this email.</p>
<p>If you're having trouble using the button, copy and paste the following link into your browser</p>
<p><a href="https://staging.evntsrc.io/verify/{{.}}">https://staging.evntsrc.io/verify/{{.}}</a></p>
<p>Warmest regards,<br/>
Evntsrc.io</p>
`
)

func (s *server) sendRecreationEmail(ctx context.Context, user *users.User) error {
	emailBody, err := renderTemplate(emailTemplatesRecreation, user)
	if err != nil {
		return err
	}

	_, err = s.emails.Send(ctx, &emails.Email{
		To: []*emails.Recipient{
			{Email: user.Email, Name: user.Name},
		},
		Subject: "You've already created an account",
		Html:    emailBody,
		Headers: map[string]string{
			"X-Priority": "High",
		},
	})

	return err
}

func (s *server) sendValidationEmail(ctx context.Context, user *users.User) error {
	token := fmt.Sprintf("%s?e=%s", string(user.Metadata[mdValidationtoken]), url.QueryEscape(base64.StdEncoding.EncodeToString([]byte(user.Email))))

	emailBody, err := renderTemplate(emailTemplatesValidation, token)
	if err != nil {
		return err
	}

	_, err = s.emails.Send(ctx, &emails.Email{
		To: []*emails.Recipient{
			{Email: user.Email, Name: user.Name},
		},
		Subject: "Please verify your email",
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
