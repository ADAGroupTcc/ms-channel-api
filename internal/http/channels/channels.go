package channels

import (
	"net/http"

	"github.com/ADAGroupTcc/ms-channels-api/exceptions"
	"github.com/ADAGroupTcc/ms-channels-api/internal/domain"
	"github.com/ADAGroupTcc/ms-channels-api/internal/helpers"
	"github.com/ADAGroupTcc/ms-channels-api/internal/services/channels"
	"github.com/labstack/echo/v4"
)

type Handler interface {
	Create(c echo.Context) error
	Get(c echo.Context) error
	List(c echo.Context) error
	Update(c echo.Context) error
	Delete(c echo.Context) error
}

type channelsHandler struct {
	channelsService channels.Service
}

func New(channelsService channels.Service) Handler {
	return &channelsHandler{
		channelsService,
	}
}
func (h *channelsHandler) Create(c echo.Context) error {
	ctx := c.Request().Context()

	var channelRequest domain.ChannelRequest
	if err := c.Bind(&channelRequest); err != nil {
		return exceptions.New(exceptions.ErrInvalidPayload, err)
	}

	channel, err := h.channelsService.Create(ctx, channelRequest)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, channel)
}

func (h *channelsHandler) Get(c echo.Context) error {
	ctx := c.Request().Context()

	id := c.Param("id")
	channel, err := h.channelsService.Get(ctx, id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, channel)
}

func (h *channelsHandler) List(c echo.Context) error {
	ctx := c.Request().Context()

	var queryParams helpers.QueryParams
	err := helpers.BindQueryParams(c, &queryParams)
	if err != nil {
		return exceptions.New(exceptions.ErrInvalidPayload, err)
	}

	channels, err := h.channelsService.List(ctx, queryParams)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, channels)
}

func (h *channelsHandler) Update(c echo.Context) error {
	ctx := c.Request().Context()

	var channelRequest domain.ChannelPatchRequest
	if err := c.Bind(&channelRequest); err != nil {
		return exceptions.New(exceptions.ErrInvalidPayload, err)
	}
	id := c.Param("id")
	if id == "" {
		return exceptions.New(exceptions.ErrInvalidID, nil)
	}

	err := h.channelsService.Update(ctx, id, channelRequest)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *channelsHandler) Delete(c echo.Context) error {
	ctx := c.Request().Context()

	id := c.Param("id")
	err := h.channelsService.Delete(ctx, id)
	if err != nil {
		return c.NoContent(http.StatusNoContent)
	}

	return c.NoContent(http.StatusNoContent)
}
