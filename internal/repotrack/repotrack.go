package repotrack

import (
	"fmt"
	"strings"
)

// RepoTrack represents git(hub) repository and sync metadata
// fields missing yaml tags are not imported from yaml but are added later
type RepoTrack struct {
	Name                  string
	URL                   string `yaml:"url"`
	Protocol              string
	Branch                string `yaml:"branch"`
	WebhookSecret         string `yaml:"webhook_secret"`
	Username              string `yaml:"username"`
	PasswordToken         string `yaml:"password_token"`
	SSHKey                string `yaml:"ssh_key"`
	CommitID              string
	WebhookSecretRequired bool `yaml:"webhook_secret_required"`
}

// NewRepoTrack returns a RepoTrack struct with default values
func NewRepoTrack() *RepoTrack {

	var rt RepoTrack
	rt.WebhookSecretRequired = true
	return &rt
}

// Populate RepoTrack struct with metadata from other fields
// TODO should this validate and also return error?
func Populate(r *RepoTrack) {

	pathParts := strings.Split(r.URL, "/")
	tail := pathParts[len(pathParts)-1]
	nameParts := strings.Split(tail, ".")

	protocolParts := strings.Split(pathParts[0], ":")

	r.Name = nameParts[0]
	r.Protocol = protocolParts[0]
}

func (r RepoTrack) String() string {
	return fmt.Sprintf("\"%s\" <%s> %t", r.Name, r.URL, r.WebhookSecretRequired)
}
