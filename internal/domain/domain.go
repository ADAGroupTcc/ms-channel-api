package domain

import (
	"github.com/ADAGroupTcc/ms-channels-api/exceptions"
	"github.com/ADAGroupTcc/ms-channels-api/pkg/mongorm"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChannelResponse struct {
	Channels []*Channel `json:"channels"`
	NextPage int64      `json:"next_page,omitempty"`
}

type Channel struct {
	mongorm.Model `bson:",inline"`
	Name          string               `json:"name" bson:"name"`
	Description   string               `json:"description" bson:"description"`
	Members       []primitive.ObjectID `json:"members" bson:"members"`
	Admins        []primitive.ObjectID `json:"admins" bson:"admins"`
}

type ChannelRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Members     []string `json:"members"`
	Admins      []string `json:"admins"`
}

func (r *ChannelRequest) Validate() error {
	if r.Name == "" || len(r.Name) < 3 {
		return exceptions.New(exceptions.ErrInvalidNameField, nil)
	}
	if r.Members == nil || len(r.Members) < MEMBERS_MINIMUM {
		return exceptions.New(exceptions.ErrInvalidMembersField, nil)
	}
	if r.Admins == nil || len(r.Admins) < ADMINS_MINIMUM {
		return exceptions.New(exceptions.ErrInvalidAdminsField, nil)
	}

	err := r.ValidateMembersAndAdmins()
	if err != nil {
		return err
	}

	return err
}

func (r *ChannelRequest) ValidateMembersAndAdmins() error {
	var err error
	var countAdmins int
	for _, member := range r.Members {
		err = ValidateUserId(member)
		for _, admin := range r.Admins {
			if member == admin {
				countAdmins++
			}
		}
	}

	if countAdmins < len(r.Admins) {
		return exceptions.New(exceptions.ErrInvalidAdminsField, nil)
	}

	for _, admin := range r.Admins {
		err = ValidateUserId(admin)
	}
	return err
}

func ValidateUserId(userId string) error {
	_, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return exceptions.New(exceptions.ErrInvalidUserIdSent, err)
	}
	return nil
}

func ParseUserIds(userIds []string) ([]primitive.ObjectID, error) {
	var parsedUserIds []primitive.ObjectID
	for _, userId := range userIds {
		parsedId, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			return nil, exceptions.New(exceptions.ErrInvalidMembersField, err)
		}
		parsedUserIds = append(parsedUserIds, parsedId)
	}

	return parsedUserIds, nil
}

func (r *ChannelRequest) ToChannel() *Channel {
	members, _ := ParseUserIds(r.Members)
	admins, _ := ParseUserIds(r.Admins)
	return &Channel{
		Name:        r.Name,
		Description: r.Description,
		Members:     members,
		Admins:      admins,
	}
}

type ChannelPatchRequest struct {
	Name        *string   `json:"name"`
	Description *string   `json:"description"`
	Members     *[]string `json:"members"`
	Admins      *[]string `json:"admins"`
}

func (r *ChannelPatchRequest) Validate() error {
	if r.Name != nil && len(*r.Name) < 3 {
		return exceptions.New(exceptions.ErrInvalidNameField, nil)
	}
	if r.Members != nil && len(*r.Members) < MEMBERS_MINIMUM {
		return exceptions.New(exceptions.ErrInvalidMembersField, nil)
	}
	if r.Admins != nil && len(*r.Admins) < ADMINS_MINIMUM {
		return exceptions.New(exceptions.ErrInvalidAdminsField, nil)
	}

	var err error
	members := r.Members
	admins := r.Admins
	if members != nil {
		for _, member := range *r.Members {
			err = ValidateUserId(member)
		}
	}

	if admins != nil {
		for _, admin := range *r.Admins {
			err = ValidateUserId(admin)
		}
	}

	return err
}

func (r *ChannelPatchRequest) ToBsonM() bson.M {
	fields := bson.M{}
	if r.Name != nil {
		fields["name"] = *r.Name
	}
	if r.Description != nil {
		fields["description"] = *r.Description
	}
	if r.Members != nil {
		parsedMembers, _ := ParseUserIds(*r.Members)
		fields["members"] = parsedMembers
	}
	if r.Admins != nil {
		parsedAdmins, _ := ParseUserIds(*r.Admins)
		fields["admins"] = parsedAdmins
	}

	response := bson.M{"$set": fields}
	return response
}

const (
	MEMBERS_MINIMUM = 2
	ADMINS_MINIMUM  = 1
)
