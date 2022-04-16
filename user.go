package main

import (
	"log"
	"time"
)

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

func getEmail(email string) []Email {
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

func getUser() User {
	user, _, err := getClient().Users.CurrentUser()
	if err != nil {
		log.Fatalf("Failed to get user: %v", err)
	}

	return User{
		Type:      "user",
		URL:       user.WebURL,
		Login:     user.Username,
		Name:      user.Name,
		Company:   nil,
		Website:   user.WebsiteURL,
		Location:  nil,
		Emails:    getEmail(user.Email),
		CreatedAt: user.CreatedAt,
	}
}
