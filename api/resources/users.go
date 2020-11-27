package resources

type User struct {
	FName    string `json:"fName" db:"fName"`
	LName    string `json:"lName" db:"lName"`
	Address  string `json:"address" db:"address"`
	City     string `json:"city" db:"city"`
	Zip_code int    `json:"zip_code" db:"zip_code"`
	Phone    string `json:"phone" db:"phone"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}
