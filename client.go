package ckRequest

import (
	"net/http"
	"net/url"
)

type HttpClient struct {
	Cl *http.Client
}

func NewClient() *HttpClient {
	c := &HttpClient{
		Cl: &http.Client{
			Jar: newSimpJar(),
		},
	}
	c.Cl.CheckRedirect = DisableRedirectCb
	return c
}

func (c *HttpClient) Do(req *HttpReq) *HttpResp {
	return c.autoRedirect(req)
}

func (c *HttpClient) DisableKeepAlive() {
	c.Cl.Transport.(*http.Transport).DisableKeepAlives = true
}

func (c *HttpClient) SetProxy(addr string) {
	c.Cl.Transport.(*http.Transport).Proxy = func(r *http.Request) (*url.URL, error) {
		return url.Parse(addr)
	}
}

func (c *HttpClient) StoreCookies(filepath string) error {
	return c.Cl.Jar.(*simpJar).StoreCookies(filepath)
}
func (c *HttpClient) LoadCookies(filepath string) error {
	return c.Cl.Jar.(*simpJar).LoadCookies(filepath)
}
