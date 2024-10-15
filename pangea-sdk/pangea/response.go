package pangea

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"

	pu "github.com/pangeacyber/pangea-go/pangea-sdk/v4/internal/pangeautil"
)

type ResponseHeader struct {
	// The request ID
	RequestID *string `json:"request_id"`

	// The time the request was issued, ISO8601
	RequestTime *string `json:"request_time"`

	// The time the response was issued, ISO8601
	ResponseTime *string `json:"response_time"`

	// The HTTP status code msg
	Status *string `json:"status"`

	// The summary of the response
	Summary *string `json:"summary"`
}

type Response struct {
	ResponseHeader
	HTTPResponse  *http.Response
	AttachedFiles []AttachedFile

	// Query raw result
	RawResult json.RawMessage `json:"result"`

	rawResponse []byte
}

func (r *ResponseHeader) String() string {
	if r == nil {
		return ""
	}
	return fmt.Sprintf("request_id: %v, request_time: %v, response_time: %v, status: %v, summary: %v",
		StringValue(r.RequestID), StringValue(r.RequestTime), StringValue(r.ResponseTime),
		StringValue(r.Status), StringValue(r.Summary))
}

type PangeaResponse[T any] struct {
	Response
	Result         *T
	AcceptedResult *AcceptedResult
}

func (r *Response) UnmarshalResult(target interface{}) error {
	return json.Unmarshal(r.RawResult, target)
}

// newResponse takes a http.Response and tries to parse the body into a base pangea API response.
func newResponse(r *http.Response) (*Response, error) {
	response := &Response{HTTPResponse: r}
	defer r.Body.Close()
	var data []byte
	if isMultipart(r.Header) {
		boundary, err := pu.GetBoundary(r.Header.Get("Content-Type"))
		if err != nil {
			return nil, err
		}
		// Create a multipart reader
		multipartReader := multipart.NewReader(r.Body, boundary)
		var n = 0

		// Iterate through each part in the multipart response
		for {
			part, err := multipartReader.NextRawPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, err
			}
			defer part.Close()
			if n == 0 {
				// Read the part's content
				data, err = io.ReadAll(part)
				if err != nil {
					return nil, err
				}
			} else {
				// Read the part's content
				fb, err := io.ReadAll(part)
				if err != nil {
					return nil, err
				}
				af := AttachedFile{
					Filename:    part.FileName(),
					File:        fb,
					ContentType: part.Header.Get("Content-Type"),
				}
				response.AttachedFiles = append(response.AttachedFiles, af)
			}

			n += 1
		}
	} else {
		var err error
		data, err = io.ReadAll(r.Body)
		if err != nil {
			return nil, NewUnmarshalError(err, []byte{}, r)
		}
	}
	response.rawResponse = data
	if err := json.Unmarshal(data, response); err != nil {
		return nil, NewUnmarshalError(err, data, r)
	}
	return response, nil
}

// Check if the response is multipart
func isMultipart(header http.Header) bool {
	contentType := header.Get("Content-Type")
	return len(contentType) > 0 && contentType[:10] == "multipart/"
}

type AttachedFile struct {
	Filename    string
	File        []byte
	ContentType string
}

type AttachedFileSaveInfo struct {
	Filename string
	Folder   string
}

func (af AttachedFile) Save(info AttachedFileSaveInfo) error {
	folder := "./"
	if info.Folder != "" {
		folder = info.Folder
	}

	filename := "defaultFilename"
	if af.Filename != "" {
		filename = af.Filename
	}
	if info.Filename != "" {
		filename = info.Filename
	}

	if _, err := os.Stat(folder); os.IsNotExist(err) {
		// Directory does not exist, create it
		err := os.MkdirAll(folder, os.ModePerm)
		if err != nil {
			return err
		}
	}

	filePath := path.Join(folder, filename)

	err := os.WriteFile(filePath, af.File, 0644)
	if err != nil {
		return err
	}

	return nil
}

// MarshalJSON implements the json.Marshaler interface for CustomType.
func (r Response) MarshalJSON() ([]byte, error) {
	if r.rawResponse == nil {
		return nil, errors.New("unable to read response body")
	}
	b := make([]byte, len(r.rawResponse))
	nc := copy(b, r.rawResponse)
	if nc != len(r.rawResponse) {
		return nil, errors.New("unable to copy raw response")
	}
	return b, nil
}
