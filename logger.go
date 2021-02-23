package request

import (
	"encoding/json"
	"strings"
)

type reqLogLine struct {
	Info                             string
	Target, Referer, Cookie, ReqBody string
	SetCookie, RespBody              string
}

// todo : json.Marshall is ugly, update it
func logRequest(req *HttpReq, resp *HttpResp, info string) {
	if resp != nil && resp.Err != nil {
		info += resp.Err.Error()
	}
	line := reqLogLine{Info: info}

	if req != nil {
		line.Target = req.Request.URL.String()
		line.Referer = req.Request.Referer()
		line.Cookie = strings.Join(cookie2names(req.Cookies()), ",")
		line.ReqBody = removeNewLine(req.string())
	}
	if resp != nil {
		line.Target = resp.Request.URL.String()
		line.SetCookie = strings.Join(cookie2names(resp.Cookies()), ",")
		line.RespBody = removeNewLine(resp.String())
	}
	b, _ := json.Marshal(&line)
	logger.Println(string(b))
}

func removeNewLine(multiLine string) string {
	tmp := strings.ReplaceAll(multiLine, "\n", "")
	return strings.ReplaceAll(tmp, "\r", "")
}
