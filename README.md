# ckRequest

## Feature
- 支持重定向过程中逐步请求的cookie/set-cookie查看
- 支持cookie的导入导出
- 简单的api用于设置POST请求的body和content-type
  - 支持form,json,multipart,binary

## doc
```go
package ckRequest // import "ckRequest"


FUNCTIONS

func DefaultRedirectCb(req *http.Request, via []*http.Request) error
func DisableRedirectCb(req *http.Request, via []*http.Request) error

TYPES

type Clienter interface {
        Do(req *HttpReq) *HttpResp
}

type File struct {
        FilePath string
}

type HttpClient struct {
        Cl *http.Client
}

var DefaultClient *HttpClient
func NewClient() *HttpClient

func (c *HttpClient) Do(req *HttpReq) *HttpResp

func (c *HttpClient) Get(url string) *HttpResp

func (c *HttpClient) LoadCookies(filepath string) error

func (c *HttpClient) Post(url string, body io.Reader) *HttpResp

func (c *HttpClient) PostForm(url, body string) *HttpResp

func (c *HttpClient) PostJson(url, body string) *HttpResp

func (c *HttpClient) StoreCookies(filepath string) error

type HttpReq struct {
        *http.Request
}

func Get(url string) *HttpReq

func NewRequest(method, url string, body io.Reader) *HttpReq

func Post(url string, body io.Reader) *HttpReq

func PostForm(url string, body string) *HttpReq
    body can be built using url.Values{}.encode

func PostFormAny(url_ string, body interface{}) *HttpReq
    body should be a struct{anyType} ,
    map[interface{}]interface{},map[interface{}][]interface{}

func PostJson(url string, body string) *HttpReq

func PostMIMEJpg(url string, img []byte) *HttpReq

func PostMultiPart(url string, body map[string]interface{}) *HttpReq
    use "mime/multipart"

func (r *HttpReq) Do(c Clienter) *HttpResp

func (r *HttpReq) Form() *HttpReq

func (r *HttpReq) Jpeg() *HttpReq

func (r *HttpReq) Json() *HttpReq

func (r *HttpReq) Origin(origin string) *HttpReq

func (r *HttpReq) Referer(referer string) *HttpReq

type HttpResp struct {
        *http.Response
        Err error
        // Has unexported fields.
}

func (resp *HttpResp) Bytes() []byte

func (resp *HttpResp) Reader() *bytes.Reader

func (resp *HttpResp) String() string
```