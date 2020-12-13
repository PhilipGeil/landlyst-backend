package resources

type User struct {
	ID       int    `json:"user_id" db:"id"`
	Fname    string `json:"fname" db:"fname"`
	Lname    string `json:"lname" db:"lName"`
	Address  string `json:"address" db:"address"`
	Zip_code string `json:"zip_code" db:"zip_code"`
	Phone    string `json:"phone" db:"phone"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
