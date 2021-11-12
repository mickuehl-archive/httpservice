package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"strconv"
	"strings"

	"github.com/txsvc/stdlib/observer"
)

// UnmashalJSONResponse unmarshalls a generic HTTP response body into a JSON structure
// Pass optionally a pointer to a byte array to get the raw body of the response object written back
func UnmashalJSONResponse(resp *http.Response, v interface{}, b *[]byte) error {
	// read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("unable to read response body: %v", err)
	}

	if b != nil {
		*b = body
	}

	// check http status ok
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s: %s", resp.Status, body)
	}

	// decode as json and return if ok
	err = json.Unmarshal(body, v)
	if err == nil {
		return nil
	}

	// check json response content type
	ct := resp.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(ct)
	if err == nil && mediaType == "application/json" {
		return fmt.Errorf("got Content-Type = application/json, but could not unmarshal as JSON: %v", err)
	}
	return fmt.Errorf("expected Content-Type = application/json, got %q: %v", ct, err)
}

func LogHttpRequest(ctx context.Context, req *http.Request) {
	observer.LogWithLevel(observer.LevelInfo, req.RequestURI, "user-agent", req.UserAgent())
}

// ParseRange extracts a byte range if specified. For specs see
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Range_requests
func ParseRange(r string) (int64, int64) {
	if r == "" {
		return 0, -1 // no range requested
	}
	parts := strings.Split(r, "=")
	if len(parts) != 2 {
		return 0, -1 // no range requested
	}
	// we simply assume that parts[0] == "bytes"
	ra := strings.Split(parts[1], "-")
	if len(ra) != 2 { // again a simplification, multiple ranges or overlapping ranges are not supported
		return 0, -1
	}

	start, err := strconv.ParseInt(ra[0], 10, 64)
	if err != nil {
		return 0, -1
	}
	end, err := strconv.ParseInt(ra[1], 10, 64)
	if err != nil {
		return 0, -1
	}

	return start, end - start
}
