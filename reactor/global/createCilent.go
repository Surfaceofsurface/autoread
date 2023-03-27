package global

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

var C *http.Client

type FPOST = func(uri string, body []byte) ([]byte, error)
type FGET = func(uri string, param map[string]string) ([]byte, error)
type ReqMethod struct {
	G FGET
	P FPOST
}
type EnvOpt = map[string]string

func init() {

	proxyExist, _ := http.Get(LOCAL_PROXY)
	if proxyExist == nil {
		C = http.DefaultClient
		return
	} else {
		C = &http.Client{
			Transport: &http.Transport{
				Proxy: func(r *http.Request) (*url.URL, error) {
					u, e := url.Parse(LOCAL_PROXY)
					return u, e
				},
			},
		}
	}
}
func POST(uri string, body []byte) (<-chan []byte, <-chan error) {
	_chan := make(chan []byte)
	_errchan := make(chan error)
	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(body))
	if err != nil {
		_chan <- nil
		_errchan <- err
		return _chan, _errchan
	}
	AddDefaultHeaders(req)
	go func() {
		resp, err := C.Do(req)
		if err != nil {
			_chan <- nil
			_errchan <- err
		}
		fmt.Printf("resp.Cookies(): %v\n", resp.Cookies())
		buf := bufio.NewReader(resp.Body)
		bytearray, err := io.ReadAll(buf)

		if err != nil {
			_chan <- nil
			_errchan <- err
			return
		}

		_chan <- bytearray
		_errchan <- nil
	}()

	return _chan, _errchan
}
func LPOST(uri string, body []byte) (*http.Response, error) {

	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	AddDefaultHeaders(req)

	resp, err := C.Do(req)
	return resp, err

}
func CGET(uri string, opt EnvOpt) (*http.Response, error) {

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	AddDefaultHeaders(req)
	AddOptHeaders(req, opt)
	resp, err := C.Do(req)
	return resp, err

}
func AddDefaultHeaders(req *http.Request) {
	req.Header.Add("accept", "application/json, text/plain, */*")
	req.Header.Add("accept-language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	req.Header.Add("content-type", "application/json;charset=UTF-8")
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-site", "same-origin")
	req.Header.Add("Referer", "https://m.cxstar.com/book/24026b13000001XXXX/read")
	req.Header.Add("Referrer-Policy", "strict-origin-when-cross-origin")
	req.Header.Add("Origin", "https://m.cxstar.com")
	req.Header.Add("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1 Edg/110.0.0.0")
}
func AddOptHeaders(req *http.Request, opt EnvOpt) {
	for k, v := range opt {
		req.Header.Add(k, v)
	}
}
func NewReqEnv(opt EnvOpt) ReqMethod {
	_GET :=
		func(uri string, query map[string]string) ([]byte, error) {
			var querystr string
			if query != nil {
				querystr = "?"
				for k, v := range query {
					querystr += k + "=" + v + "&"
				}
			} else {
				querystr = ""
			}
			req, err := http.NewRequest("GET", uri+querystr, nil)
			if err != nil {
				// _chan <- nil
				// _errchan <- err
				return nil, err
			}
			AddDefaultHeaders(req)
			req.Header.Del("content-type")
			AddOptHeaders(req, opt)
			resp, err := C.Do(req)
			if err != nil {
				return nil, err
			}
			buf := bufio.NewReader(resp.Body)
			bytearray, err := io.ReadAll(buf)
			if err != nil {
				return nil, err
			}
			defer resp.Body.Close()
			return bytearray, err
		}
	_POST :=
		func(uri string, body []byte) ([]byte, error) {
			// _chan := make(chan []byte)
			// _errchan := make(chan error)
			req, err := http.NewRequest("POST", uri, bytes.NewBuffer(body))
			if err != nil {
				// _chan <- nil
				// _errchan <- err
				return nil, err
			}
			AddDefaultHeaders(req)
			AddOptHeaders(req, opt)
			resp, err := C.Do(req)
			if err != nil {
				return nil, err
			}
			defer resp.Body.Close()
			buf := bufio.NewReader(resp.Body)
			bytearray, err := io.ReadAll(buf)
			if err != nil {
				return nil, err
			}

			return bytearray, err
		}
	return ReqMethod{
		_GET,
		_POST,
	}
}
