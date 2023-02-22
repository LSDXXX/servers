package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/LSDXXX/libs/pkg/log"
	"github.com/LSDXXX/libs/pkg/servercontext"
)

// Request description
type Request struct {
	Method        string
	URL           string
	RequestParams map[string]string
	Headers       map[string]string
	PostForm      map[string]string
	Body          interface{}
}

// EncodeURL add query
func EncodeURL(rawURL string, querys map[string]string) (*url.URL, error) {
	surl, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	query := surl.Query()
	for k, v := range querys {
		query.Set(k, v)
	}
	surl.RawQuery = query.Encode()
	return surl, nil
}

// DoHttpRequest description
// @param ctx
// @param request
// @return *http.Response
// @return error
func DoHttpRequest(ctx context.Context, request Request) (*http.Response, error) {
	url := request.URL
	surl, err := EncodeURL(url, request.RequestParams)
	if err != nil {
		return nil, err
	}
	log.WithContext(ctx).Debugf("do http request, url: %s, request: %+v", surl.String(), request.Body)
	var buf *bytes.Buffer
	if request.Body != nil {
		switch body := request.Body.(type) {
		case string:
			buf = bytes.NewBufferString(body)
		case []byte:
			buf = bytes.NewBuffer(body)
		default:
			data, _ := json.Marshal(request.Body)
			buf = bytes.NewBuffer(data)
		}
	}
	req, err := http.NewRequest(strings.ToUpper(request.Method), surl.String(), buf)
	if err != nil {
		return nil, err
	}
	for k, v := range request.Headers {
		req.Header.Set(k, v)
	}
	for k, v := range request.PostForm {
		req.PostForm.Add(k, v)
	}
	return Do(ctx, req)
}

// Do description
// @param ctx
// @param request
// @return *http.Response
// @return error
func Do(ctx context.Context, request *http.Request) (*http.Response, error) {
	servercontext.ContextToHTTP(ctx, request)
	return http.DefaultClient.Do(request)
}

// DoWithClient description
// @param ctx
// @param client
// @param request
// @return *http.Response
// @return error
func DoWithClient(ctx context.Context, client *http.Client, request *http.Request) (*http.Response, error) {
	servercontext.ContextToHTTP(ctx, request)
	logged, _ := httputil.DumpRequest(request, true)
	log.WithContext(ctx).Debugf("dump http request: %s", string(logged))
	return client.Do(request)
}
