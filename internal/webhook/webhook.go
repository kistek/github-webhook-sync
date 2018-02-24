package webhook

import "fmt"

// Webhook represents the webhook payload sent from github
// We ignore fields we don't expect to use
type Webhook struct {
	Repository Repository
	Before     string `json:"before"`
	After      string `json:"after"`
}

// Repository represents a git repository as in github webhook schema
type Repository struct {
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	SSHURL   string `json:"ssh_url"`
	HTTPURL  string `json:"clone_url"`
	ID       int    `json:"id"`
	Private  bool   `json:"private"`
}

// User represents a user as in github webhook schema
type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (w Webhook) String() string {
	return fmt.Sprintf("%s, ssh_url: <%s>, http_url: <%s>, after: %s",
		w.Repository.FullName,
		w.Repository.SSHURL,
		w.Repository.HTTPURL,
		w.After,
	)
}
