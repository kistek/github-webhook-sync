package repotrack

import (
	"errors"
	"fmt"
)

// RepoTrack represents git(hub) repository and sync metadata
type RepoTrack struct {
	Name                  string `yaml:"name"`
	URL                   string `yaml:"url"`
	Protocol              string
	Branch                string `yaml:"branch"`
	WebhookSecret         string `yaml:"webhook_secret"`
	Username              string `yaml:"username"`
	PasswordToken         string `yaml:"password_token"`
	SSHKey                string `yaml:"ssh_key"`
	WebhookSecretRequired bool   `yaml:"webhook_secret_required"`
}

// NewRepoTrack returns a RepoTrack struct with default values
func NewRepoTrack() RepoTrack {

	var rt RepoTrack
	rt.WebhookSecretRequired = true
	return rt
}

// Populate RepoTrack struct with metadata from other fields
func (r RepoTrack) Populate() error {
	r.Name = "TODO"
	r.Protocol = "TODO"

	return errors.New("not implemented")
}

func (r RepoTrack) String() string {
	return fmt.Sprintf("\"%s\" <%s> %t", r.Name, r.URL, r.WebhookSecretRequired)
}
