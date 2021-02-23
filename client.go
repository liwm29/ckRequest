package request

import (
	"io"
	"net/http"
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

func (c *HttpClient) Get(url string) *HttpResp {
	return Get(url).Do(c)
}

func (c *HttpClient) Post(url string, body io.Reader) *HttpResp {
	return Post(url, body).Do(c)
}

func (c *HttpClient) PostJson(url, body string) *HttpResp {
	return PostJson(url, body).Do(c)
}

func (c *HttpClient) PostForm(url, body string) *HttpResp {
	return PostForm(url, body).Do(c)
}

func (c *HttpClient) StoreCookies(filepath string) error {
	return c.Cl.Jar.(*simpJar).StoreCookies(filepath)
}
func (c *HttpClient) LoadCookies(filepath string) error {
	return c.Cl.Jar.(*simpJar).LoadCookies(filepath)
}
