package request

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

type HttpResp struct {
	*http.Response
	copy   []byte
	isRead bool
	Err    error
}

func newResponse(resp *http.Response, err error) *HttpResp {
	return &HttpResp{resp, nil, false, err}
}

func (resp *HttpResp) Bytes() []byte {
	if resp.isRead {
		return resp.copy
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	PanicIf(err)
	resp.Body.Close()
	resp.copy = bytes
	resp.isRead = true
	return bytes
}

func (resp *HttpResp) String() string {
	if !resp.isRead {
		return string(resp.Bytes())
	}

	return string(resp.copy)
}

func (resp *HttpResp) Reader() *bytes.Reader {
	if !resp.isRead {
		resp.Bytes()
	}

	return bytes.NewReader(resp.copy)
}
