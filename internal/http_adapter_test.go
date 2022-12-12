package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

type failReader struct{}

func (failReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("simulated error")
}

func validLambdaRequest() *events.ALBTargetGroupRequest {
	return &events.ALBTargetGroupRequest{
		HTTPMethod:                      "POST",
		Path:                            "/shining",
		QueryStringParameters:           nil,
		MultiValueQueryStringParameters: nil,
		Headers: map[string]string{
			"Accept":       "application/json;v=1",
			"Content-Type": "application/json",
		},
		MultiValueHeaders: nil,
		RequestContext:    events.ALBTargetGroupRequestContext{},
		IsBase64Encoded:   true,
		Body:              "eyJmb28iOiAiYmFyIn0=",
	}
}

func validHttpRequest(t *testing.T) *http.Request {
	body := `{"foo": "bar"}`
	req, err := http.NewRequest("POST", "https://localhost/shining", bytes.NewBufferString(body))
	assert.NoError(t, err)
	req.Header.Set("Accept", "application/json;v=1")
	req.Header.Set("Content-Type", "application/json")
	return req
}

func validHTTPResponse() *http.Response {
	body := `{"foo": "bar"}`
	headers := http.Header{
		"Content-Type": []string{"application/json"},
	}
	return &http.Response{
		Status:           "200 OK",
		StatusCode:       200,
		Proto:            "HTTP/1.1",
		ProtoMajor:       1,
		ProtoMinor:       1,
		Header:           headers,
		Body:             ioutil.NopCloser(bytes.NewBufferString(body)),
		ContentLength:    int64(len(body)),
		TransferEncoding: nil,
		Close:            false,
		Uncompressed:     false,
		Trailer:          nil,
		Request:          nil,
		TLS:              nil,
	}
}

func Test_decodeLambdaBody(t *testing.T) {
	type args struct {
		body string
		out  *bytes.Buffer
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "correct encoding",
			args: args{
				body: "ewogICAgInJlcXVlc3RDb250ZXh0IjogewogICAgICAgICJlbGIiOiB7CiAgICAgICAgICAgICJ0YXJnZXRHcm91cEFybiI6ICJhcm46YXdzOmVsYXN0aWNsb2FkYmFsYW5jaW5nOnJlZ2lvbjoxMjM0NTY3ODkwMTI6dGFyZ2V0Z3JvdXAvbXktdGFyZ2V0LWdyb3VwLzZkMGVjZjgzMWVlYzlmMDkiCiAgICAgICAgfQogICAgfSwKICAgICJodHRwTWV0aG9kIjogIkdFVCIsICAKICAgICJwYXRoIjogIi8iLCAgCiAgICAicXVlcnlTdHJpbmdQYXJhbWV0ZXJzIjoge30sICAKICAgICJoZWFkZXJzIjogewogICAgICAgICJ1c2VyLWFnZW50IjogIkVMQi1IZWFsdGhDaGVja2VyLzIuMCIKICAgIH0sICAKICAgICJib2R5IjogIiIsICAKICAgICJpc0Jhc2U2NEVuY29kZWQiOiB0cnVlCn0=",
				out:  bytes.NewBuffer([]byte{}),
			},
		},
		{
			name: "invalid b64 encoding",
			args: args{
				body: "foobarbaz",
				out:  bytes.NewBuffer([]byte{}),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := decodeLambdaBody(tt.args.body, tt.args.out)
			if (err != nil) != tt.wantErr {
				t.Errorf("decodeLambdaBody() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				var event events.ALBTargetGroupRequest
				assert.NoError(t, json.Unmarshal(tt.args.out.Bytes(), &event))
			}
		})
	}
}

func Test_getRequestBody(t *testing.T) {
	notB64 := validLambdaRequest()
	notB64.Body = `{"foo": "bar"}`
	notB64.IsBase64Encoded = false
	type args struct {
		r *events.ALBTargetGroupRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *bytes.Buffer
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "not base64 encoded",
			args: args{
				r: notB64,
			},
			want:    bytes.NewBufferString(`{"foo": "bar"}`),
			wantErr: assert.NoError,
		},
		{
			name: "base64 encoded",
			args: args{
				r: &events.ALBTargetGroupRequest{
					HTTPMethod:                      "POST",
					Path:                            "/shining",
					QueryStringParameters:           nil,
					MultiValueQueryStringParameters: nil,
					Headers: map[string]string{
						"Accept":       "application/json;v=1",
						"Content-Type": "application/json",
					},
					MultiValueHeaders: nil,
					RequestContext:    events.ALBTargetGroupRequestContext{},
					IsBase64Encoded:   true,
					Body:              "eyJmb28iOiAiYmFyIn0=",
				},
			},
			want:    bytes.NewBufferString(`{"foo": "bar"}`),
			wantErr: assert.NoError,
		},
		{
			name: "invalid base64 encoded",
			args: args{
				r: &events.ALBTargetGroupRequest{
					HTTPMethod:                      "POST",
					Path:                            "/shining",
					QueryStringParameters:           nil,
					MultiValueQueryStringParameters: nil,
					Headers: map[string]string{
						"Accept":       "application/json;v=1",
						"Content-Type": "application/json",
					},
					MultiValueHeaders: nil,
					RequestContext:    events.ALBTargetGroupRequestContext{},
					IsBase64Encoded:   true,
					Body:              `{"foo": "bar"}`,
				},
			},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getRequestBody(tt.args.r)
			if !tt.wantErr(t, err, fmt.Sprintf("getRequestBody(%v)", tt.args.r)) {
				return
			}
			assert.Equalf(t, tt.want, got, "getRequestBody(%v)", tt.args.r)
		})
	}
}

func Test_eventToHttpRequest(t *testing.T) {
	badRequest := validLambdaRequest()
	badRequest.Body = "asdkljflksdjf"
	badRequest.IsBase64Encoded = true
	badRequestUrl := validLambdaRequest()
	badRequestUrl.Path = "!@#$%^&**(("
	type args struct {
		r events.ALBTargetGroupRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *http.Request
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "valid request",
			args:    args{r: *validLambdaRequest()},
			want:    validHttpRequest(t),
			wantErr: assert.NoError,
		},
		{
			name:    "invalid b64 encoding",
			args:    args{r: *badRequest},
			wantErr: assert.Error,
		},
		{
			name:    "invalid request url",
			args:    args{r: *badRequestUrl},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := eventToHttpRequest(tt.args.r)
			if !tt.wantErr(t, err, fmt.Sprintf("eventToHttpRequest(%v)", tt.args.r)) {
				return
			}
			if tt.want != nil {
				defer got.Body.Close()
				assert.Equal(t, tt.want.Method, got.Method)
				assert.Equal(t, tt.want.URL, got.URL)
				gotBody, err := ioutil.ReadAll(got.Body)
				assert.NoError(t, err)
				wantBody, err := ioutil.ReadAll(tt.want.Body)
				assert.NoError(t, err)
				assert.Equal(t, wantBody, gotBody)
			}

			//assert.Equalf(t, tt.want, got, "eventToHttpRequest(%v)", tt.args.r)
		})
	}
}

func Test_responseToEvent(t *testing.T) {
	errHttpResp := validHTTPResponse()
	errHttpResp.Body = ioutil.NopCloser(failReader{})
	respMultiHeader := validHTTPResponse()
	respMultiHeader.Header = map[string][]string{
		"X-Foo-Bar": []string{"foo", "bar"},
	}
	type args struct {
		resp *http.Response
	}
	tests := []struct {
		name    string
		args    args
		want    events.ALBTargetGroupResponse
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "valid response",
			args: args{resp: validHTTPResponse()},
			want: events.ALBTargetGroupResponse{
				StatusCode:        200,
				StatusDescription: "200 OK",
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
				MultiValueHeaders: nil,
				Body:              `{"foo": "bar"}`,
				IsBase64Encoded:   false,
			},
			wantErr: assert.NoError,
		},
		{
			name: "multiheader",
			args: args{resp: respMultiHeader},
			want: events.ALBTargetGroupResponse{
				StatusCode:        200,
				StatusDescription: "200 OK",
				Headers: map[string]string{
					"X-Foo-Bar": "foo,bar",
				},
				MultiValueHeaders: nil,
				Body:              `{"foo": "bar"}`,
				IsBase64Encoded:   false,
			},
			wantErr: assert.NoError,
		},
		{
			name:    "response reading error",
			args:    args{resp: errHttpResp},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := responseToEvent(tt.args.resp)
			if !tt.wantErr(t, err, fmt.Sprintf("responseToEvent(%v)", tt.args.resp)) {
				return
			}
			assert.Equalf(t, tt.want, got, "responseToEvent(%v)", tt.args.resp)
		})
	}
}

func TestAlbHandler_ProxyWithContext(t *testing.T) {
	valid := validLambdaRequest()
	valid.Path = "/valid"
	fail500 := validLambdaRequest()
	fail500.Path = "/fail"
	invalidB64 := validLambdaRequest()
	invalidB64.IsBase64Encoded = true
	invalidB64.Body = "aslkdjflsjdfkdjsfkljsdf"
	http.Handle("/valid", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"foo": "bar"}`))
	}))
	http.Handle("/fail", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
	}))
	type args struct {
		ctx context.Context
		r   events.ALBTargetGroupRequest
	}
	tests := []struct {
		name    string
		args    args
		want    events.ALBTargetGroupResponse
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "valid request",
			args: args{
				ctx: context.Background(),
				r:   *valid,
			},
			want: events.ALBTargetGroupResponse{
				StatusCode:        200,
				StatusDescription: "200 OK",
				Headers:           map[string]string{"Content-Type": "application/json"},
				Body:              `{"foo": "bar"}`,
			},
			wantErr: assert.NoError,
		},
		{
			name: "500 error",
			args: args{
				ctx: context.Background(),
				r:   *fail500,
			},
			want: events.ALBTargetGroupResponse{
				StatusCode:        500,
				StatusDescription: "500 Internal Server Error",
				Headers:           map[string]string{"Content-Type": "application/json"},
				Body:              "",
			},
			wantErr: assert.NoError,
		},
		{
			name: "decoding error",
			args: args{
				ctx: context.Background(),
				r:   *invalidB64,
			},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &AlbHandler{}
			got, err := h.ProxyWithContext(tt.args.ctx, tt.args.r)
			if !tt.wantErr(t, err, fmt.Sprintf("ProxyWithContext(%v, %v)", tt.args.ctx, tt.args.r)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ProxyWithContext(%v, %v)", tt.args.ctx, tt.args.r)
		})
	}
}
