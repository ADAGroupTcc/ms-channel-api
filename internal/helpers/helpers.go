package helpers

import (
	"strings"

	"github.com/ADAGroupTcc/ms-channels-api/exceptions"
	"github.com/labstack/echo/v4"
)

func BindQueryParams(c echo.Context, queryParams *QueryParams) error {
	if err := c.Bind(queryParams); err != nil {
		return err
	}
	queryParams.HeaderUserId = c.Request().Header.Get("user_id")
	if queryParams.HeaderUserId == "" {
		return exceptions.ErrHeaderUserIdIsReq
	}

	queryParams.normalize()
	return nil
}

type QueryParams struct {
	RawChannelIds string `query:"channel_ids"`
	RawUserIds    string `query:"user_ids"`
	ShowMembers   bool   `query:"show_members"`
	HeaderUserId  string
	ChannelIDs    []string
	UserIds       []string
	Limit         int64 `query:"limit"`
	Offset        int64 `query:"next_page"`
}

func (q *QueryParams) normalize() {
	q.ChannelIDs = strings.Split(q.RawChannelIds, ",")
	q.UserIds = strings.Split(q.RawUserIds, ",")
	q.RawUserIds = ""
	q.RawChannelIds = ""
	if q.Limit < 1 {
		q.Limit = 10
	}
	if q.Offset < 0 {
		q.Offset = 0
	}
}
