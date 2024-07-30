package channels

import (
	"context"
	"errors"

	"github.com/ADAGroupTcc/ms-channels-api/exceptions"
	"github.com/ADAGroupTcc/ms-channels-api/internal/domain"
	"github.com/ADAGroupTcc/ms-channels-api/pkg/mongorm"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const CHANNEL_COLLECTION = "channels"

type Repository interface {
	Create(ctx context.Context, Channel *domain.Channel) (*domain.Channel, error)
	Get(ctx context.Context, id primitive.ObjectID) (*domain.Channel, error)
	List(ctx context.Context, channelIds []primitive.ObjectID, userIds []primitive.ObjectID, limit int64, offset int64) ([]*domain.Channel, error)
	Update(ctx context.Context, id primitive.ObjectID, fields bson.M) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}

type ChannelRepository struct {
	db *mongo.Database
}

func New(db *mongo.Database) Repository {
	return &ChannelRepository{db}
}

func (h *ChannelRepository) Create(ctx context.Context, channel *domain.Channel) (*domain.Channel, error) {
	filter := bson.M{"members": bson.M{"$all": channel.Members}}
	err := channel.Read(ctx, h.db, CHANNEL_COLLECTION, filter, channel)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			err = channel.Create(ctx, h.db, CHANNEL_COLLECTION, channel)
			if err != nil {
				return nil, exceptions.New(exceptions.ErrDatabaseFailure, err)
			}
			return channel, nil
		}
	}

	return nil, exceptions.New(exceptions.ErrChannelAlreadyExists, err)
}

func (h *ChannelRepository) Get(ctx context.Context, id primitive.ObjectID) (*domain.Channel, error) {
	Channel := &domain.Channel{}
	err := Channel.Read(ctx, h.db, CHANNEL_COLLECTION, bson.M{"_id": id}, Channel)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, exceptions.New(exceptions.ErrChannelNotFound, err)
		}

		return nil, exceptions.New(exceptions.ErrDatabaseFailure, err)
	}
	return Channel, nil
}

func (h *ChannelRepository) List(ctx context.Context, channelIds []primitive.ObjectID, userIds []primitive.ObjectID, limit int64, offset int64) ([]*domain.Channel, error) {
	var channels []*domain.Channel = make([]*domain.Channel, 0)
	var filter bson.M
	if len(channelIds) > 0 {
		filter = bson.M{"_id": bson.M{"$in": channelIds}}
	}
	if len(userIds) > 0 {
		if filter == nil {
			filter = bson.M{"members": bson.M{"$all": userIds}}
		} else {
			filter["members"] = bson.M{"$all": userIds}
		}
	}
	err := mongorm.List(ctx, h.db, CHANNEL_COLLECTION, filter, &channels, options.Find().SetLimit(limit).SetSkip(offset*limit))
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return channels, nil
		}
		return nil, exceptions.New(exceptions.ErrDatabaseFailure, err)
	}
	return channels, nil
}

func (h *ChannelRepository) Update(ctx context.Context, id primitive.ObjectID, fields bson.M) error {
	Channel := &domain.Channel{}
	err := Channel.Update(ctx, h.db, CHANNEL_COLLECTION, bson.M{"_id": id}, fields, options.FindOneAndUpdate().SetReturnDocument(options.After))
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return exceptions.New(exceptions.ErrChannelNotFound, err)
		}
		return exceptions.New(exceptions.ErrDatabaseFailure, err)
	}
	return nil
}

func (h *ChannelRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	Channel := &domain.Channel{}
	err := Channel.Delete(ctx, h.db, CHANNEL_COLLECTION, bson.M{"_id": id})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return exceptions.New(exceptions.ErrChannelNotFound, err)
		}
		return exceptions.New(exceptions.ErrDatabaseFailure, err)
	}
	return nil
}
