package helpers

import (
	"strings"

	"github.com/labstack/echo/v4"
)

func BindQueryParams(c echo.Context, queryParams *QueryParams) error {
	if err := c.Bind(queryParams); err != nil {
		return err
	}
	queryParams.normalize()
	return nil
}

type QueryParams struct {
	RawChannelIds string `json:"channel_ids" query:"channel_ids"`
	RawUserIds    string `json:"user_ids" query:"user_ids"`
	ChannelIDs    []string
	UserIds       []string
	Limit         int64 `json:"limit" query:"limit"`
	Offset        int64 `json:"next_page" query:"next_page"`
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
