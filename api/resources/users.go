package resources

type User struct {
	FName       string   `json:"fName"`
	LName       string   `json:"lName"`
	Address     string   `json:"address"`
	City        string   `json:"city"`
	Zip_code    int      `json:"zip_code"`
	Phone       string   `json:"phone"`
	Email       string   `json:"email"`
	Password    string   `json:"password"`
	Permissions []string `json:"permissions"`
}
