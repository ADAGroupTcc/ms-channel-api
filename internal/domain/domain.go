package domain

import (
	"fmt"
	"regexp"

	"github.com/ADAGroupTcc/ms-users-api/exceptions"
	"github.com/ADAGroupTcc/ms-users-api/pkg/mongorm"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	mongorm.Model `bson:",inline"`
	FirstName     string   `json:"first_name" bson:"first_name"`
	LastName      string   `json:"last_name" bson:"last_name"`
	Description   string   `json:"description,omitempty" bson:"description"`
	Nickname      string   `json:"nickname" bson:"nickname"`
	Email         string   `json:"email" bson:"email"`
	CPF           string   `json:"cpf" bson:"cpf"`
	Categories    []string `json:"categories" bson:"categories"`
	IsDenunciated bool     `json:"is_denunciated" bson:"is_denunciated"`
}

type UserResponse struct {
	Users    []*User `json:"users"`
	NextPage int64   `json:"next_page,omitempty"`
}

type UserV2 struct {
	mongorm.Model `bson:",inline"`
	FirstName     string               `json:"first_name" bson:"first_name"`
	LastName      string               `json:"last_name" bson:"last_name"`
	Description   string               `json:"description,omitempty" bson:"description"`
	Nickname      string               `json:"nickname" bson:"nickname"`
	Email         string               `json:"email" bson:"email"`
	CPF           string               `json:"cpf" bson:"cpf"`
	Categories    []primitive.ObjectID `json:"categories" bson:"categories"`
	IsDenunciated bool                 `json:"is_denunciated" bson:"is_denunciated"`
}

type UserRequest struct {
	FirstName   string   `json:"first_name"`
	LastName    string   `json:"last_name"`
	Description string   `json:"description,omitempty"`
	Nickname    *string  `json:"nickname"`
	Email       string   `json:"email"`
	CPF         string   `json:"cpf"`
	Categories  []string `json:"categories" `
}

func (u *UserRequest) Validate() error {
	if u.FirstName == "" || len(u.FirstName) < 3 {
		return exceptions.New(exceptions.ErrInvalidFirstName, nil)
	}
	if u.LastName == "" || len(u.LastName) < 3 {
		return exceptions.New(exceptions.ErrInvalidLastName, nil)
	}
	match, err := regexp.MatchString("^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$", u.Email)
	if err != nil || !match {
		return exceptions.New(exceptions.ErrInvalidEmail, nil)
	}
	match, err = regexp.MatchString("^[0-9]{3}.[0-9]{3}.[0-9]{3}-[0-9]{2}$", u.CPF)
	if err != nil || !match {
		return exceptions.New(exceptions.ErrInvalidCPF, nil)
	}

	return nil
}

func (u *UserRequest) ValidateV2() error {
	if u.FirstName == "" || len(u.FirstName) < 3 {
		return exceptions.New(exceptions.ErrInvalidFirstName, nil)
	}
	if u.LastName == "" || len(u.LastName) < 3 {
		return exceptions.New(exceptions.ErrInvalidLastName, nil)
	}
	match, err := regexp.MatchString("^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$", u.Email)
	if err != nil || !match {
		return exceptions.New(exceptions.ErrInvalidEmail, nil)
	}
	match, err = regexp.MatchString("^[0-9]{3}.[0-9]{3}.[0-9]{3}-[0-9]{2}$", u.CPF)
	if err != nil || !match {
		return exceptions.New(exceptions.ErrInvalidCPF, nil)
	}

	for _, category := range u.Categories {
		_, err := primitive.ObjectIDFromHex(category)
		if err != nil {
			return exceptions.New(exceptions.ErrInvalidCategories, nil)
		}
	}
	return nil
}

func (u *UserRequest) ToUser() *User {
	if u.Nickname == nil {
		newNickname := fmt.Sprintf("%s %s", u.FirstName, u.LastName)
		u.Nickname = &newNickname
	}

	return &User{
		FirstName:     u.FirstName,
		LastName:      u.LastName,
		Description:   u.Description,
		Nickname:      *u.Nickname,
		Email:         u.Email,
		CPF:           u.CPF,
		Categories:    u.Categories,
		IsDenunciated: false,
	}
}

func (u *UserRequest) ToUserV2() *UserV2 {
	if u.Nickname == nil {
		newNickname := fmt.Sprintf("%s %s", u.FirstName, u.LastName)
		u.Nickname = &newNickname
	}

	user := &UserV2{
		FirstName:     u.FirstName,
		LastName:      u.LastName,
		Description:   u.Description,
		Nickname:      *u.Nickname,
		Email:         u.Email,
		CPF:           u.CPF,
		Categories:    make([]primitive.ObjectID, len(u.Categories)),
		IsDenunciated: false,
	}

	for i, category := range u.Categories {
		parsedCategory, _ := primitive.ObjectIDFromHex(category)
		user.Categories[i] = parsedCategory
	}

	return user
}

type UserPatchRequest struct {
	FirstName     *string   `json:"first_name"`
	LastName      *string   `json:"last_name"`
	Description   *string   `json:"description"`
	Nickname      *string   `json:"nickname"`
	Email         *string   `json:"email"`
	Categories    *[]string `json:"categories"`
	IsDenunciated *bool     `json:"is_denunciated"`
}

func (u *UserPatchRequest) Validate() error {
	if u.FirstName != nil && (len(*u.FirstName) < 3) {
		return exceptions.New(exceptions.ErrInvalidFirstName, nil)
	}
	if u.LastName != nil && (len(*u.LastName) < 3) {
		return exceptions.New(exceptions.ErrInvalidLastName, nil)
	}
	if u.Email != nil {
		match, err := regexp.MatchString("^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$", *u.Email)
		if err != nil || !match {
			return exceptions.New(exceptions.ErrInvalidEmail, nil)
		}
	}

	return nil
}

func (u *UserPatchRequest) ToBsonM() bson.M {
	updateFields := bson.M{}
	if u.FirstName != nil {
		updateFields["first_name"] = *u.FirstName
	}
	if u.LastName != nil {
		updateFields["last_name"] = *u.LastName
	}
	if u.Description != nil {
		updateFields["description"] = *u.Description
	}
	if u.Nickname != nil {
		updateFields["nickname"] = *u.Nickname
	}
	if u.Email != nil {
		updateFields["email"] = *u.Email
	}
	if u.Categories != nil {
		updateFields["categories"] = *u.Categories
	}
	if u.IsDenunciated != nil {
		updateFields["is_denunciated"] = *u.IsDenunciated
	}
	return bson.M{"$set": updateFields}
}