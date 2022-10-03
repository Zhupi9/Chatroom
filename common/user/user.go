package user

type User struct {
	UserPwd     string `json:"userpwd"`
	UserName    string `json:"username"`
	PhoneNumber string `json:"phonenumber"`
}
