package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

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
					IsBase64Encoded:   false,
					Body:              `{"foo": "bar"}`,
				},
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
	type args struct {
		r events.ALBTargetGroupRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *http.Request
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := eventToHttpRequest(tt.args.r)
			if !tt.wantErr(t, err, fmt.Sprintf("eventToHttpRequest(%v)", tt.args.r)) {
				return
			}
			assert.Equalf(t, tt.want, got, "eventToHttpRequest(%v)", tt.args.r)
		})
	}
}
