package request

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/spf13/cast"
)

type HttpReq struct {
	*http.Request
}

func NewRequest(method, url string, body io.Reader) *HttpReq {
	req, _ := http.NewRequest(method, url, body)
	r := &HttpReq{req}
	initDefaultHdr(r)
	return r
}

func Get(url string) *HttpReq {
	return NewRequest("GET", url, nil)
}

func Post(url string, body io.Reader) *HttpReq {
	return NewRequest("POST", url, body)
}

func PostJson(url string, body string) *HttpReq {
	return NewRequest("POST", url, strings.NewReader(body)).Json()
}

// body can be built using url.Values{}.encode
func PostForm(url string, body string) *HttpReq {
	return NewRequest("POST", url, strings.NewReader(body)).Form()
}

// body should be a struct{anyType} , map[interface{}]interface{},map[interface{}][]interface{}
func PostFormAny(url_ string, body interface{}) *HttpReq {

	data := url.Values{}
	v := reflect.ValueOf(body)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			iFieldT := v.Type().Field(i)
			iFieldV := v.Field(i)

			name, ok := iFieldT.Tag.Lookup("form")
			if !ok {
				name = strings.ToLower(iFieldT.Name)
			}

			if iFieldV.CanInterface() {
				data[name] = []string{cast.ToString(iFieldV.Interface())}
			}
		}
	case reflect.Map:
		mapStrSliceStr, err := cast.ToStringMapStringSliceE(v.Interface())
		if err != nil {
			PanicIf(wrapError("PostFormStruct:", err))
		}
		data = mapStrSliceStr
	}

	return PostForm(url_, data.Encode())
}

// use "mime/multipart"
func PostMultiPart(url string, body interface{}) *HttpReq {
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	file, errFile1 := os.Open("filepath")
	defer file.Close()
	part1,
		errFile1 := writer.CreateFormFile("asdf", filepath.Base("/C:/Users/salvare000/Desktop/信息安全技术/18341018_李伟铭_lab5.pdf"))
	_, errFile1 = io.Copy(part1, file)
	if errFile1 != nil {
		return nil
	}
	_ = writer.WriteField("asdf", "fdsa")
	err := writer.Close()
	if err != nil {
		return nil
	}

	req := PostForm(url, payload.String())
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}

func PostMIME(url string, img []byte) *HttpReq {
	return nil
}

type Clienter interface {
	Do(req *HttpReq) *HttpResp
}

func (r *HttpReq) Do(c Clienter) *HttpResp {
	if c == nil {
		c = DefaultClient
	}
	return c.Do(r)
}

// it will just show less than 1024B , the overflow will be discarded
func (r *HttpReq) string() string {
	if r.GetBody == nil {
		return ""
	}
	reader, err := r.GetBody()
	PanicIf(err)
	buf := [1024]byte{}
	n, err := reader.Read(buf[:])
	PanicIf(err)
	// if n == 1024,we just discard the left
	reader.Close()
	return string(buf[:n])
}
