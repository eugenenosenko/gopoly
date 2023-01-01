package users

type User interface {
	IsUser()
}

type Contact interface {
	IsContact()
}

type RegularUser struct {
	ID       string    `json:"id"`
	Type     string    `json:"kind"`
	Name     string    `json:"name"`
	Address  string    `json:"address"`
	Contacts []Contact `json:"contacts"`
}

func (a RegularUser) IsUser() {}

type PrivilegedUser struct {
	ID         string    `json:"id"`
	Type       string    `json:"kind"`
	Name       string    `json:"name"`
	Address    string    `json:"address"`
	Contacts   []Contact `json:"contacts"`
	Privileges []string  `json:"privileges"`
}

func (a PrivilegedUser) IsUser() {}

type BannedUser struct {
	ID        string    `json:"id"`
	Type      string    `json:"kind"`
	Contacts  []Contact `json:"contacts"`
	BanReason string    `json:"ban_reason"`
}

func (o BannedUser) IsUser() {}

type BusinessContact struct {
	ID           string `json:"id"`
	BusinessName string `json:"business_name"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
}

func (c BusinessContact) IsContact() {}

type FullName struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

type PrivateContact struct {
	ID       string   `json:"id"`
	FullName FullName `json:"fullname"`
	Phone    string   `json:"phone"`
	Email    string   `json:"email"`
}

func (c PrivateContact) IsContact() {}
