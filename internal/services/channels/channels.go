package channels

import (
	"context"

	"github.com/ADAGroupTcc/ms-channels-api/exceptions"
	"github.com/ADAGroupTcc/ms-channels-api/internal/domain"
	"github.com/ADAGroupTcc/ms-channels-api/internal/helpers"
	"github.com/ADAGroupTcc/ms-channels-api/internal/repositories/channels"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service interface {
	Create(ctx context.Context, request domain.ChannelRequest) (*domain.Channel, error)
	Get(ctx context.Context, id string) (*domain.Channel, error)
	List(ctx context.Context, queryParams helpers.QueryParams) (*domain.ChannelResponse, error)
	Update(ctx context.Context, id string, request domain.ChannelPatchRequest) error
	Delete(ctx context.Context, id string) error
}

type ChannelService struct {
	channelRepository channels.Repository
}

func New(channelRepository channels.Repository) Service {
	return &ChannelService{
		channelRepository,
	}
}

func (h *ChannelService) Create(ctx context.Context, request domain.ChannelRequest) (*domain.Channel, error) {
	err := request.Validate()
	if err != nil {
		return nil, err
	}

	Channel := request.ToChannel()

	return h.channelRepository.Create(ctx, Channel)
}

func (h *ChannelService) Get(ctx context.Context, id string) (*domain.Channel, error) {
	parsedId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, exceptions.New(exceptions.ErrInvalidID, err)
	}
	return h.channelRepository.Get(ctx, parsedId)
}

func (h *ChannelService) List(ctx context.Context, queryParams helpers.QueryParams) (*domain.ChannelResponse, error) {
	parsedUserIds, err := h.parseObjectIdFromString(queryParams.UserIds)
	if err != nil {
		return nil, err
	}

	parsedHeaderUserId, err := primitive.ObjectIDFromHex(queryParams.HeaderUserId)
	if err != nil {
		return nil, err
	}

	var channels domain.ChannelResponseGeneral

	if queryParams.ShowMembers {
		channels, err = h.channelRepository.Aggregate(ctx, parsedUserIds, parsedHeaderUserId)
		if err != nil {
			return nil, err
		}
	} else {
		parsedChannelIds, err := h.parseObjectIdFromString(queryParams.ChannelIDs)
		if err != nil {
			return nil, err
		}

		channels, err = h.channelRepository.List(ctx, parsedChannelIds, parsedUserIds, parsedHeaderUserId, queryParams.Limit, queryParams.Offset)
		if err != nil {
			return nil, err
		}
	}

	response := &domain.ChannelResponse{
		Channels: channels,
	}

	if channelsParsed, ok := channels.([]interface{}); ok && len(channelsParsed) == int(queryParams.Limit) {
		response.NextPage = queryParams.Offset + 1
	}

	return response, nil
}

func (h *ChannelService) Update(ctx context.Context, id string, request domain.ChannelPatchRequest) error {
	parsedId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return exceptions.New(exceptions.ErrInvalidID, err)
	}

	err = request.Validate()
	if err != nil {
		return err
	}

	fieldsToUpdate := request.ToBsonM()

	return h.channelRepository.Update(ctx, parsedId, fieldsToUpdate)
}

func (h *ChannelService) Delete(ctx context.Context, id string) error {
	parsedId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return exceptions.New(exceptions.ErrInvalidID, err)
	}

	return h.channelRepository.Delete(ctx, parsedId)
}

func (*ChannelService) parseObjectIdFromString(ids []string) ([]primitive.ObjectID, error) {
	var parsedIds []primitive.ObjectID = make([]primitive.ObjectID, 0)
	for _, id := range ids {
		parsedId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			continue
		}
		parsedIds = append(parsedIds, parsedId)
	}
	return parsedIds, nil
}
