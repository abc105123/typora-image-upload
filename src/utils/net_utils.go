package utils

import (
	"io"
	"net/http"
)

func ReadResponseBody(resp *http.Response) []byte {
	if resp == nil {
		return nil
	}

	body := resp.Body
	defer body.Close()

	data, _ := io.ReadAll(body)

	return data
}
