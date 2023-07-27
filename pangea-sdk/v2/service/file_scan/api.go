package file_scan

import (
	"context"
	"io"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/pangea"
)

func (e *FileScan) Scan(ctx context.Context, input *FileScanRequest, file io.Reader) (*pangea.PangeaResponse[FileScanResult], error) {
	req, err := e.Client.NewRequestMultPart("POST", "v1/scan", input, file)
	if err != nil {
		return nil, err
	}
	out := FileScanResult{}
	resp, err := e.Client.Do(ctx, req, &out)

	if resp == nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[FileScanResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, err
}

type FileScanRequest struct {
	Verbose  bool   `json:"verbose,omitempty"`
	Raw      bool   `json:"raw,omitempty"`
	Provider string `json:"provider,omitempty"`
}

type FileScanData struct {
	Category []string `json:"category"`
	Score    int      `json:"score"`
	Verdict  string   `json:"verdict"`
}

type FileScanResult struct {
	Data       FileScanData `json:"data"`
	Parameters interface{}  `json:"parameters,omitempty"`
	RawData    interface{}  `json:"raw_data,omitempty"`
}
