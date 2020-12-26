package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"onedrive/fileutil"
	"io"
	"net/http"
	"net/url"
	"os"
)

const (
	baseURL                       = "https://graph.microsoft.com/v1.0"
	statusInsufficientStorage int = 507
)

// OneDrive is the entry point for the client. It manages the communication with
// Microsoft OneDrive Graph API
type OneDrive struct {
	Client  *http.Client
	BaseURL string
}

// NewOneDrive returns a new OneDrive client to enable you to communicate with
// the API
func NewOneDriveClient(c *http.Client, debug bool) *OneDrive {
	drive := OneDrive{
		Client:  c,
		BaseURL: baseURL,
	}
	return &drive
}

func createRequestBody(body interface{}) (io.ReadWriter, error) {
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}
	return buf, nil
}

func getRequestBody(body interface{}) (io.ReadWriter, error) {
	var buf io.ReadWriter

	switch body.(type) {
	case string:
		if body != nil {
			buf = new(bytes.Buffer)
			err := json.NewEncoder(buf).Encode(body)
			if err != nil {
				return nil, err
			}
		}
		break
	case []byte:
		buf = bytes.NewBuffer(body.([]byte))
		break
	case *os.File:
		//file := &bytes.Buffer{}
		//_, err := io.Copy(file, body.(*os.File))
		//if err != nil {
		//	return nil, err
		//}
		bytesData, err := fileutil.ReadFile(body.(*os.File))
		if err != nil {
			return nil, err
		}
		buf = bytes.NewBuffer(bytesData)
		break
	}

	return buf, nil
}

func isValidUrl(uri string) bool {
	_, err := url.ParseRequestURI(uri)
	if err != nil {
		return false
	}

	u, err := url.Parse(uri)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}

// Generate request
func (od *OneDrive) NewRequest(method, uri string, requestHeaders map[string]string, body interface{}) (*http.Request, error) {
	reqBody, err := getRequestBody(body)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse the file into Bytes  reason: %v", err)
	}
	var req *http.Request
	if isValidUrl(uri) {
		req, err = http.NewRequest(method, uri, reqBody)
	} else {
		req, err = http.NewRequest(method, od.BaseURL+uri, reqBody)
	}
	if err != nil {
		return nil, err
	}
	//Adding default header
	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", getUserAgent())

	//Adding application specific Headers
	if requestHeaders != nil {
		for header, value := range requestHeaders {
			req.Header.Set(header, value)
		}
	}

	return req, nil
}

//Execute request
func (od *OneDrive) Do(req *http.Request) (*http.Response, error) {
	resp, err := od.Client.Do(req)
	if err != nil {
		return nil, err
	}
	//defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest && resp.StatusCode <= statusInsufficientStorage {
		newErr := new(Error)
		if err := json.NewDecoder(resp.Body).Decode(newErr); err != nil {
			return resp, err
		}
		return resp, newErr
	}
	return resp, err
}
