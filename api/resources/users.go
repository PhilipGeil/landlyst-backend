package resources

type User struct {
	ID       int    `json:"id" db:"id"`
	FName    string `json:"fName" db:"fname"`
	LName    string `json:"lName" db:"lName"`
	Address  string `json:"address" db:"address"`
	City     string `json:"city" db:"city"`
	Zip_code int    `json:"zip_code" db:"zip_code"`
	Phone    string `json:"phone" db:"phone"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
