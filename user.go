package main

import (
	"log"
	"time"
)

type UserService struct {
	*BaseService
}

type User struct {
	Type      string     `json:"type"`
	URL       string     `json:"url"`
	Login     string     `json:"login"`
	Name      string     `json:"name"`
	Company   *string    `json:"company"`
	Website   string     `json:"website"`
	Location  *string    `json:"location"`
	Emails    []Email    `json:"emails"`
	CreatedAt *time.Time `json:"created_at"`
}

type Email struct {
	Address string `json:"address"`
	Primary bool   `json:"primary"`
}

func NewUserService(e *Exporter) *UserService {
	return &UserService{
		BaseService: &BaseService{
			exporter: e,
			filename: "users.json",
		},
	}
}

func (u *UserService) GetEmail(email string) []Email {
	if email != "" {
		// TODO: This is kludgy!
		var e []Email
		e = append(e, Email{
			Address: email,
			Primary: true,
		})
		return e
	} else {
		return []Email{}
	}
}

func (u *UserService) GetUser() (*User, error) {
	user, _, err := u.exporter.Client.Users.CurrentUser()
	if err != nil {
		return nil, err
	}

	return &User{
		Type:      "user",
		URL:       user.WebURL,
		Login:     user.Username,
		Name:      user.Name,
		Company:   nil,
		Website:   user.WebsiteURL,
		Location:  nil,
		Emails:    u.GetEmail(user.Email),
		CreatedAt: user.CreatedAt,
	}, nil
}

func (u *UserService) Export() {
	user, err := u.GetUser()
	if err != nil {
		log.Fatalf("Failed to get user: %v", err)
	}
	if user != nil {
		var contains bool
		for _, us := range u.exporter.State.Users {
			if us.Login == user.Login {
				contains = true
				break
			}
		}

		if !contains {
			u.exporter.State.Users = append(u.exporter.State.Users, *user)
		}
	}
}

func (u *UserService) WriteFile() error {
	return u.exporter.WriteJsonFile(u.filename, u.exporter.State.Users)
}
