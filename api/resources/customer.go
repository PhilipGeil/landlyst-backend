package resources

type Customer struct {
	ID       int    `json:"id" db:"id"`
	FName    string `json:"fName" db:"fname"`
	LName    string `json:"lName" db:"lName"`
	Address  string `json:"address" db:"address"`
	Zip_code string `json:"zip_code" db:"zip_code"`
	Phone    string `json:"phone" db:"phone"`
	Email    string `json:"email" db:"email"`
	UserID   int    `json:"user_id"`
}
