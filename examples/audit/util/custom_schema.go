package util

import (
	"time"
)

// Custom schema implementation

type CustomSchemaEvent struct {
	Message       string     `json:"message"`
	FieldInt      int        `json:"field_int,omitempty"`
	FieldBool     bool       `json:"field_bool,omitempty"`
	FieldStrShort string     `json:"field_str_short,omitempty"`
	FieldStrLong  string     `json:"field_str_long,omitempty"`
	FieldTime     *time.Time `json:"field_time,omitempty"`

	// TenantID field
	TenantID string `json:"tenant_id,omitempty"`
}

func (e *CustomSchemaEvent) GetTenantID() string {
	return e.TenantID
}

func (e *CustomSchemaEvent) SetTenantID(tid string) {
	e.TenantID = tid
}
