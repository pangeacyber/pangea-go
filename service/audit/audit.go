package audit

import (
	"context"

	"go-pangea/pangea"
)

type Audit struct {
	client *pangea.Client
}

func New() *Audit {
	return &Audit{}
}

type DataInput struct {
	Action  *string `json:"action"`
	Actor   *string `json:"actor"`
	Message *string `json:"message"`
	New     *string `json:"new"`
	Old     *string `json:"old"`
	Source  *string `json:"source"`
	Status  *string `json:"status"`
	Target  *string `json:"target"`
}

type LogInput struct {
	Data       *DataInput `json:"data"`
	ReturnHash *bool      `json:"return_hash"`
}

type LogOutput struct {
	Hash *string `json:"hash"`
}

func (a *Audit) Log(ctx context.Context, input *LogInput) (*LogOutput, *pangea.Response, error) {
	if input == nil {
		input = &LogInput{}
	}
	req, err := a.client.NewRequest("POST", "audit", "v1/audit/log", input)
	if err != nil {
		return nil, nil, err
	}

	var out LogOutput
	resp, err := a.client.Do(ctx, req, &out)
	if err != nil {
		return nil, resp, err
	}

	return &out, resp, nil
}
