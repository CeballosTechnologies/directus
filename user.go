package directus

import "time"

type User struct {
	AuthData           interface{} `json:"auth_data,omitempty"`
	Avatar             interface{} `json:"avatar,omitempty"`
	Description        interface{} `json:"description,omitempty"`
	Email              string      `json:"email,omitempty"`
	EmailNotifications bool        `json:"email_notifications,omitempty"`
	ExternalIdentifier interface{} `json:"external_identifier,omitempty"`
	FirstName          string      `json:"first_name,omitempty"`
	Id                 string      `json:"id,omitempty"`
	Language           interface{} `json:"language,omitempty"`
	LastAccess         time.Time   `json:"last_access,omitempty"`
	LastName           string      `json:"last_name,omitempty"`
	LastPage           string      `json:"last_page,omitempty"`
	Location           interface{} `json:"location,omitempty"`
	Password           string      `json:"password,omitempty"`
	Provider           string      `json:"provider,omitempty"`
	Role               string      `json:"role,omitempty"`
	Status             string      `json:"status,omitempty"`
	Tags               []string    `json:"tags,omitempty"`
	TfaSecret          interface{} `json:"tfa_secret,omitempty"`
	Theme              string      `json:"theme,omitempty"`
	Title              interface{} `json:"title,omitempty"`
	Token              interface{} `json:"token,omitempty"`
}
