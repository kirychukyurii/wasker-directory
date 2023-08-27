package model

import "time"

type (
	Model struct {
		Id int64 `json:"id" storage:"id"`

		CreatedAt time.Time    `json:"created_at" storage:"created_at"`
		CreatedBy LookupEntity `json:"created_by" storage:"created_by"`

		UpdatedAt time.Time    `json:"updated_at" storage:"updated_at"`
		UpdatedBy LookupEntity `json:"updated_by" storage:"updated_by"`

		DeletedAt *time.Time    `json:"-" storage:"-"`
		DeletedBy *LookupEntity `json:"-" storage:"-"`
	}

	LookupEntity struct {
		Id   int64  `json:"id,omitempty" storage:"id"`
		Name string `json:"name,omitempty" storage:"name"`
	}
)
