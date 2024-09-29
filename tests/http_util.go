//nolint:unused
package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func httpGet(url string, response any) error {
	return httpReq(http.MethodGet, url, nil, response)
}

func httpReq(method string, url string, data any, response any) error {
	var bodyReader io.Reader
	if data != nil {
		bodyReader = getJSONReader(data)
	}
	req, err := http.NewRequestWithContext(context.Background(), method, url, bodyReader)
	if err != nil {
		return err
	}
	return execute(req, response)
}

func execute(req *http.Request, response any) error {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		errorBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("request url: %q.\nserver response: %d, %q", req.URL, resp.StatusCode, string(errorBytes))
	}
	if response != nil {
		json.NewDecoder(resp.Body).Decode(response)
	}

	return nil
}

func getJSONReader(object any) io.Reader {
	jsonBody := &bytes.Buffer{}
	json.NewEncoder(jsonBody).Encode(object)
	bodyReader := bytes.NewReader(jsonBody.Bytes())

	return bodyReader
}
