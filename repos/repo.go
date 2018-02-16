package repo

import "errors"

type Repo struct {
	Name                  string
	URL                   string
	Protocol              string
	Branch                string
	WebhookSecret         string
	Username              string
	PasswordToken         string
	SSHKey                string
	WebhookSecretRequired bool
}

func (r Repo) populate() error {
	r.Name = "TODO"
	r.Protocol = "TODO"

	return errors.New("not implemented")
}
