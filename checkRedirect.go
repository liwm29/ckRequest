package request

import (
	"net/http"
	"net/url"
	"strings"
)

func DisableRedirectCb(req *http.Request, via []*http.Request) error {
	return http.ErrUseLastResponse
}

func DefaultRedirectCb(req *http.Request, via []*http.Request) error {
	if len(via) > 20 {
		return http.ErrUseLastResponse
	} else {
		if len(via) == 1 {
			logger.Println("[  Request  ] via ", via[0].URL.String(), "\t\t", "Cookie: ", cookie2names(via[0].Cookies()))
		}
		logger.Println("[Redirecting] via ", req.URL.String(), "\t\t", "Cookie: ", cookie2names(req.Cookies()))
	}
	return nil
}


func (c *HttpClient) autoRedirect(ireq *HttpReq) *HttpResp {
	req := ireq
	var respW *HttpResp
	for {
		resp, err := c.Cl.Do(req.Request)
		respW = newResponse(resp, err)

		PanicIf(err)

		redirectMethod, shouldRedirect, includeBody := redirectBehavior(req.Method, resp, ireq.Request)

		if !shouldRedirect {
			if ireq == req {
				logRequest(req, respW, "req&resp")
			} else {
				logRequest(nil, respW, "final")
			}
			break
		}
		if ireq == req{
			logRequest(req, nil, "ireq")
		}
		logRequest(nil, respW, "redirect")

		location, err := resp.Location()
		PanicIf(err)
		ref := refererForURL(resp.Request.URL, location)
		req = NewRequest(redirectMethod, location.String(), nil).Referer(ref)
		if includeBody && ireq.GetBody != nil {
			// GetBody是一个closure,所以即使读完再调用也可以获取最初的snapshot
			req.Body, err = ireq.GetBody()
			PanicIf(err)
		}
	}
	return respW
}

// refererForURL returns a referer without any authentication info or
// an empty string if lastReq scheme is https and newReq scheme is http.
func refererForURL(lastReq, newReq *url.URL) string {
	// https://tools.ietf.org/html/rfc7231#section-5.5.2
	//   "Clients SHOULD NOT include a Referer header field in a
	//    (non-secure) HTTP request if the referring page was
	//    transferred with a secure protocol."
	if lastReq.Scheme == "https" && newReq.Scheme == "http" {
		return ""
	}
	referer := lastReq.String()
	if lastReq.User != nil {
		// This is not very efficient, but is the best we can
		// do without:
		// - introducing a new method on URL
		// - creating a race condition
		// - copying the URL struct manually, which would cause
		//   maintenance problems down the line
		auth := lastReq.User.String() + "@"
		referer = strings.Replace(referer, auth, "", 1)
	}
	return referer
}

// redirectBehavior describes what should happen when the
// client encounters a 3xx status code from the server
func redirectBehavior(reqMethod string, resp *http.Response, ireq *http.Request) (redirectMethod string, shouldRedirect, includeBody bool) {
	switch resp.StatusCode {
	case 301, 302, 303:
		redirectMethod = reqMethod
		shouldRedirect = true
		includeBody = false

		// RFC 2616 allowed automatic redirection only with GET and
		// HEAD requests. RFC 7231 lifts this restriction, but we still
		// restrict other methods to GET to maintain compatibility.
		// See Issue 18570.
		if reqMethod != "GET" && reqMethod != "HEAD" {
			redirectMethod = "GET"
		}
	case 307, 308:
		redirectMethod = reqMethod
		shouldRedirect = true
		includeBody = true

		// Treat 307 and 308 specially, since they're new in
		// Go 1.8, and they also require re-sending the request body.
		if resp.Header.Get("Location") == "" {
			// 308s have been observed in the wild being served
			// without Location headers. Since Go 1.7 and earlier
			// didn't follow these codes, just stop here instead
			// of returning an error.
			// See Issue 17773.
			shouldRedirect = false
			break
		}
		if ireq.GetBody == nil {
			// We had a request body, and 307/308 require
			// re-sending it, but GetBody is not defined. So just
			// return this response to the user instead of an
			// error, like we did in Go 1.7 and earlier.
			shouldRedirect = false
		}
	}
	return redirectMethod, shouldRedirect, includeBody
}
