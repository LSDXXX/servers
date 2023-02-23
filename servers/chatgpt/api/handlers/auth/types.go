package auth

type User struct {
	Id        int
	UserName  string
	FirstName string
	LastName  string
}

type Login struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type IdentityInfo struct {
	Id    int
	WSKey string
}
