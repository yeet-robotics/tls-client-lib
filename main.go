//go:build linux || darwin || windows

package main

/*
#include <stdlib.h>
*/
import "C"
import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	tls "tls-client-go"
	"unsafe"
)

type JsonRequest struct {
	Url             string            `json:"url"`
	Method          string            `json:"method"`
	Proxy           string            `json:"proxy"`
	Body            []byte            `json:"body"`
	Headers         map[string]string `json:"headers"`
	Cookies         map[string]string `json:"cookies"`
	Timeout         int               `json:"timeout"`
	FollowRedirects bool              `json:"follow_redirects"`
}

type JsonResponse struct {
	Error      string              `json:"error,omitempty"`
	StatusCode int                 `json:"status_code,omitempty"`
	Body       []byte              `json:"body,omitempty"`
	Headers    map[string][]string `json:"headers,omitempty"`
	Url        string              `json:"url,omitempty"`
}

//export libRequest
func libRequest(reqJsonC *C.char) *C.char {
	reqJson := C.GoString(reqJsonC)

	var req *JsonRequest
	if err := json.Unmarshal([]byte(reqJson), &req); err != nil {
		return C.CString("invalid JSON")
	}

	resp := request(req)

	var buff bytes.Buffer
	enc := json.NewEncoder(&buff)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(resp); err != nil {
		return C.CString("failed encoding")
	}

	return C.CString(string(buff.Bytes()))
}

//export libFree
func libFree(str *C.char) {
	C.free(unsafe.Pointer(str))
}

var pool *tls.ClientsPool

func init() {
	pool = tls.NewPool()
	pool.Config(&tls.PoolConfig{
		MaxTries: 1,
		Client:   tls.Http2Client(),
	})
}

func request(req *JsonRequest) *JsonResponse {
	urlParsed, err := url.Parse(req.Url)
	if err != nil {
		return &JsonResponse{Error: fmt.Sprintf("invalid url: %s", err.Error())}
	}

	proxyParsed, err := url.Parse(req.Proxy)
	if req.Proxy == "" {
		proxyParsed = nil
	} else if err != nil {
		return &JsonResponse{Error: fmt.Sprintf("invalid proxy: %s", err.Error())}
	}

	jar, _ := cookiejar.New(nil)

	{
		var cookies []*http.Cookie
		for name, val := range req.Cookies {
			cookies = append(cookies, &http.Cookie{Name: name, Value: val, Domain: urlParsed.Host, Path: "/"})
		}

		jar.SetCookies(urlParsed, cookies)
	}

	_, resp, err := pool.Do(
		context.Background(),
		&tls.Request{
			RequestUrl:     urlParsed,
			ProxyUrl:       proxyParsed,
			Method:         req.Method,
			Body:           req.Body,
			Headers:        tls.CopyHeaders(req.Headers),
			CookieJar:      jar,
			Timeout:        req.Timeout,
			FollowRedirect: req.FollowRedirects,
			Intent:         tls.RequestIntentNone,
			Decompress:     true,
		},
		// tls.NewRandomizeHeaderOrderOption(rand.Int63()),
		tls.NewCacheClientOption(false),
	)
	if err != nil {
		return &JsonResponse{Error: fmt.Sprintf("request failed: %s", err.Error())}
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return &JsonResponse{Error: fmt.Sprintf("failed reading body: %s", err.Error())}
	}

	return &JsonResponse{StatusCode: resp.StatusCode, Headers: resp.Header, Body: respBody, Url: resp.Request.URL.String()}
}

func main() {
}
