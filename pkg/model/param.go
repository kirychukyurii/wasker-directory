package model

import "time"

type QueryParam struct {
	Id        []int64      `query:"id"`
	Name      string       `query:"name"`
	CreatedAt time.Time    `query:"created_at"`
	CreatedBy LookupEntity `query:"created_by"`

	Query string `query:"query"`
}
