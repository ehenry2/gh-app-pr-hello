package internal

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/rs/zerolog/log"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
)

func decodeLambdaBody(body string, out *bytes.Buffer) error {
	dec := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(body))
	_, err := io.Copy(out, dec)
	return err
}

func getRequestBody(r *events.ALBTargetGroupRequest) (*bytes.Buffer, error) {
	reqBody := bytes.NewBuffer([]byte{})
	if r.IsBase64Encoded {
		if err := decodeLambdaBody(r.Body, reqBody); err != nil {
			return nil, err
		}
	} else {
		reqBody = bytes.NewBufferString(r.Body)
	}
	return reqBody, nil
}

func eventToHttpRequest(r events.ALBTargetGroupRequest) (*http.Request, error) {
	rawURL := fmt.Sprintf("https://localhost%s", r.Path)
	body, err := getRequestBody(&r)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(r.HTTPMethod, rawURL, body)
	if err != nil {
		return nil, err
	}
	for k, v := range r.Headers {
		req.Header.Set(k, v)
	}
	return req, nil
}

func responseToEvent(resp *http.Response) (events.ALBTargetGroupResponse, error) {
	headers := make(map[string]string)
	for k, v := range resp.Header {
		headers[k] = strings.Join(v, ",")
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return events.ALBTargetGroupResponse{}, err
	}
	event := events.ALBTargetGroupResponse{
		StatusCode:        resp.StatusCode,
		StatusDescription: resp.Status,
		Headers:           headers,
		Body:              string(b),
		IsBase64Encoded:   false,
	}
	return event, nil
}

type AlbHandler struct{}

func (h *AlbHandler) ProxyWithContext(ctx context.Context, r events.ALBTargetGroupRequest) (events.ALBTargetGroupResponse, error) {
	ctx = log.Logger.WithContext(ctx)
	req, err := eventToHttpRequest(r)
	if err != nil {
		return events.ALBTargetGroupResponse{}, err
	}
	req = req.WithContext(ctx)
	recorder := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(recorder, req)
	return responseToEvent(recorder.Result())
}
