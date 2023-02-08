package directus

type Role struct {
	AdminAccess bool        `json:"admin_access,omitempty"`
	AppAccess   bool        `json:"app_access,omitempty"`
	Description interface{} `json:"description,omitempty"`
	EnforceTfa  bool        `json:"enforce_tfa,omitempty"`
	Icon        string      `json:"icon,omitempty"`
	Id          string      `json:"id,omitempty"`
	IpAccess    interface{} `json:"ip_access,omitempty"`
	Name        string      `json:"name,omitempty"`
	Users       []string    `json:"users,omitempty"`
}
