package service

import (
	"bufio"
	"bytes"
	"io"
	"net/http"
	"private/autoread/reactor/global"
	"strings"
)

type WithToken struct {
	Token  string
	Cookie []*http.Cookie
}

func (t *WithToken) P(uri string, body []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	c := "Verification=1;"
	for _, k := range t.Cookie {
		c += k.Raw
	}

	global.AddDefaultHeaders(req)
	global.AddOptHeaders(req, map[string]string{"authorization": "Bearer " + t.Token, "Cookie": c})
	resp, err := global.C.Do(req)
	if err != nil {
		return nil, err
	}

	if len(resp.Cookies()) > 0 {
		t.Cookie = resp.Cookies()
	}
	defer resp.Body.Close()
	buf := bufio.NewReader(resp.Body)
	bytearray, err := io.ReadAll(buf)
	if err != nil {
		return nil, err
	}

	return bytearray, err
}
func (t *WithToken) G(uri string, query map[string]string) ([]byte, error) {
	querystr := "?"
	if strings.Contains(uri, "?") {
		querystr = "&"
	}

	for k, v := range query {
		querystr += k + "=" + v + "&"
	}
	req, err := http.NewRequest("GET", uri+querystr, nil)
	if err != nil {
		return nil, err
	}
	c := "Verification=1;"
	for _, k := range t.Cookie {
		c += k.Raw
	}

	global.AddDefaultHeaders(req)
	req.Header.Del("content-type")
	global.AddOptHeaders(req, map[string]string{"authorization": "Bearer " + t.Token, "Cookie": c})
	resp, err := global.C.Do(req)
	if err != nil {
		return nil, err
	}

	if len(resp.Cookies()) > 0 {
		t.Cookie = resp.Cookies()
	}
	defer resp.Body.Close()
	buf := bufio.NewReader(resp.Body)
	bytearray, err := io.ReadAll(buf)
	if err != nil {
		return nil, err
	}

	return bytearray, err
}
