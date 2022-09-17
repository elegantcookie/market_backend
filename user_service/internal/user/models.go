package user

import (
	"math/rand"
	"strconv"
	"time"
)

type User struct {
	ID            string `json:"id" bson:"_id,omitempty"`
	Name          string `json:"name" bson:"name"`
	Surname       string `json:"surname" bson:"surname"`
	Email         string `json:"email" bson:"email"`
	PhoneNumber   string `json:"phone_number" bson:"phone_number"`
	VkID          string `json:"vk_id" bson:"vk_id"`
	City          string `json:"city" bson:"city"`
	PhoneApproved bool   `json:"phone_approved" bson:"phone_approved"`
	LastOnline    int64  `json:"last_online" bson:"last_online"`
	//AverageRating int    `json:"average_rating" bson:"average_rating"`
	//RateCount     int    `json:"rate_count" bson:"rate_count"`
	//PhoneToken string `json:"-" bson:"phone_token,omitempty"`
}

type UpdateUserDTO struct {
	ID            string `json:"id" bson:"_id,omitempty"`
	Name          string `json:"name" bson:"name"`
	Surname       string `json:"surname" bson:"surname"`
	Email         string `json:"email" bson:"email"`
	PhoneNumber   string `json:"phone_number" bson:"phone_number"`
	VkID          string `json:"vk_id" bson:"vk_id"`
	City          string `json:"city" bson:"city"`
	PhoneApproved bool   `json:"phone_approved" bson:"phone_approved"`
	LastOnline    int64  `json:"last_online" bson:"last_online"`
	//AverageRating int    `json:"average_rating" bson:"average_rating"`
	//RateCount     int    `json:"rate_count" bson:"rate_count"`
}

type CreateByPhoneDTO struct {
	PhoneNumber      string `json:"phone_number" example:"+79999999999"`
	VerificationCode string `json:"verification_code" example:"123321"`
}

type CreateByVkDTO struct {
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	VkID        string `json:"vk_id"`
	PhoneNumber string `json:"phone_number"`
}

func NewUserFromUpdateDTO(dto UpdateUserDTO) User {
	return User{
		Name:          dto.Name,
		Surname:       dto.Surname,
		Email:         dto.Email,
		PhoneNumber:   dto.PhoneNumber,
		VkID:          dto.VkID,
		City:          dto.City,
		PhoneApproved: dto.PhoneApproved,
		LastOnline:    dto.LastOnline,
	}
}

func NewUserByPhone(dto CreateByPhoneDTO) User {
	return User{
		Name:          "",
		Surname:       "",
		Email:         "",
		PhoneNumber:   dto.PhoneNumber,
		VkID:          "",
		City:          "",
		PhoneApproved: true,
		//PhoneToken:    dto.PhoneToken,
	}
}

func NewUserByVkID(dto CreateByVkDTO) User {
	return User{
		Name:          dto.Name,
		Surname:       dto.Surname,
		Email:         "",
		PhoneNumber:   dto.PhoneNumber,
		VkID:          dto.VkID,
		City:          "",
		PhoneApproved: true,
	}
}

// Returns string code with n digits
func generateVerificationCode(n int) string {
	rand.Seed(time.Now().Unix())
	str := ""
	for i := 0; i < n; i++ {
		str += strconv.Itoa(rand.Int() % 10)
	}
	return str
}

//func (u *User) CheckPhoneToken(password string) error {
//	err := bcrypt.CompareHashAndPassword([]byte(u.PhoneNumber), []byte(password))
//	if err != nil {
//		return fmt.Errorf("phone token does not match")
//	}
//	return nil
//}
//
//func (u *User) GeneratePhoneTokenHash() error {
//	phoneToken, err := generatePhoneTokenHash(u.PhoneNumber)
//	if err != nil {
//		return err
//	}
//	u.PhoneNumber = phoneToken
//	return nil
//}
//
//func generatePhoneTokenHash(token string) (string, error) {
//	hash, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.MinCost)
//	if err != nil {
//		return "", fmt.Errorf("failed to hash token due to error %w", err)
//	}
//	return string(hash), nil
//}
