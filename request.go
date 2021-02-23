package ckRequest

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
			panicIf(wrapError("PostFormStruct:", err))
		}
		data = mapStrSliceStr
	}

	return PostForm(url_, data.Encode())
}

type File struct {
	FilePath string
}

// use "mime/multipart"
func PostMultiPart(url string, body map[string]interface{}) *HttpReq {
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	for k, v := range body {
		switch z := v.(type) {
		case File:
			file, errFile1 := os.Open(z.FilePath)
			defer file.Close()
			part1,
				errFile1 := writer.CreateFormFile(k, filepath.Base(z.FilePath))
			_, errFile1 = io.Copy(part1, file)
			panicIf(wrapError("postMultiPart:", errFile1))
		case string:
			_ = writer.WriteField(k, z)
		}
	}
	err := writer.Close()
	panicIf(wrapError("close multipartWriter:", err))

	req := PostForm(url, payload.String())
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}

func PostMIMEJpg(url string, img []byte) *HttpReq {
	return Post(url, bytes.NewReader(img)).Jpeg()
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
	panicIf(err)
	buf := [1024]byte{}
	n, err := reader.Read(buf[:])
	panicIf(err)
	// if n == 1024,we just discard the left
	reader.Close()
	return string(buf[:n])
}
