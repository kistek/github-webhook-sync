package webhook

import "fmt"

// Webhook represents the webhook payload sent from github
// We ignore fields we don't expect to use
type Webhook struct {
	Repository Repository
	Before     string
	After      string
}

// Repository represents a git repository as in github webhook schema
type Repository struct {
	Name     string
	FullName string
	ID       int
	Private  bool
}

// User represents a user as in github webhook schema
type User struct {
	Name  string
	Email string
}

func (w Webhook) String() string {
	return fmt.Sprintf(`name:%s,
before:%s
after: %s`,
		w.Repository.Name,
		w.Before,
		w.After,
	)
}
