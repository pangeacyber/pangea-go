package audit

import (
	"context"

	"github.com/pangeacyber/go-pangea/pangea"
)

type LogInput struct {
	// A structured record describing that <actor> did <action> on <target>
	// changing it from <old> to <new> and the operation was <status>, and/or a free-form <message>.
	Data *Data `json:"data"`

	// If the response should include the hash.
	ReturnHash *bool `json:"return_hash"`
}

type Data struct {
	// An identifier for _who_ the audit record is about.
	Actor *string `json:"actor,omitempty"`

	// What action was performed on a record.
	// eg: "created", "deleted", "updated"
	Action *string `json:"action,omitempty"`

	// A free form text field describing the event.
	// Message is a required field.
	Message *string `json:"message"`

	// The value of a record _after_ it was changed.
	New *string `json:"new,omitempty"`

	// The value of a record _before_ it was changed.
	Old *string `json:"old,omitempty"`

	// The source of a record. Can be used to hard-split logged and searched data.
	Source *string `json:"source,omitempty"`

	// The status or result of the event
	// eg: "failure", "success"
	Status *string `json:"status,omitempty"`

	// An identifier for what the audit record is about.
	Target *string `json:"target,omitempty"`
}

type LogOutput struct {
	// The hash of the log.
	Hash *string `json:"hash"`
}

func (a *Audit) Log(ctx context.Context, input *LogInput) (*LogOutput, *pangea.Response, error) {
	if input == nil {
		input = &LogInput{}
	}
	req, err := a.Client.NewRequest("POST", "audit", "v1/audit/log", input)
	if err != nil {
		return nil, nil, err
	}

	var out LogOutput
	resp, err := a.Client.Do(ctx, req, &out)
	if err != nil {
		return nil, resp, err
	}

	return &out, resp, nil
}
