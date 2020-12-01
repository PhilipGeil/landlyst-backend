package resources

type Room struct {
	ID             int  `json:"id" db:"id"`
	Cleaned        bool `json:"cleaned" db:"cleaned"`
	CurrentStaying int  `json:"current_staying" db:"current_staying"`
	RoomNumber     int  `json:"room_number" db:"room_number"`
}

type RoomAdditions struct {
	ID          int    `json:"id" db:"id"`
	Item        string `json:"item" db:"item"`
	Price       int    `json:"price" db:"price"`
	Description string `json:"description" db:"description"`
}

type RoomRoomAdditions struct {
	Room          Room
	RoomAdditions RoomAdditions
}
