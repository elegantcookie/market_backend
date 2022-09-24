package user

type User struct {
	ID            string   `json:"id" bson:"_id,omitempty"`
	Username      string   `json:"username" bson:"username"`
	Password      string   `json:"-" bson:"password"`
	HasFreeTicket bool     `json:"has_free_ticket" bson:"has_free_ticket"`
	Tickets       []string `json:"tickets" bson:"tickets"`
}

type UserDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ResponseUserDTO struct {
	ID            string   `json:"id"`
	Username      string   `json:"username"`
	HasFreeTicket bool     `json:"has_free_ticket"`
	Tickets       []string `json:"tickets"`
}
