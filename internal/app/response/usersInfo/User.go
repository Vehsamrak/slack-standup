package usersInfo

type User struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	RealName string `json:"real_name"`
	IsBot    bool   `json:"is_bot"`
}

func (user *User) NormalizedName() string {
	if user.RealName != "" {
		return user.RealName
	}

	return user.Name
}
