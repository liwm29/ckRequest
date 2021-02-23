package ckRequest

// "net/http"

func initDefaultHdr(r *HttpReq) {
	r.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.130 Safari/537.36")
}

func (r *HttpReq) Referer(referer string) *HttpReq {
	r.Header.Set("Referer", referer)
	return r
}

func (r *HttpReq) Origin(origin string) *HttpReq {
	r.Header.Set("Origin", origin)
	return r
}

func (r *HttpReq) Json() *HttpReq {
	r.Header.Set("Content-Type", "application/json;charset=UTF-8")
	return r
}

func (r *HttpReq) Form() *HttpReq {
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return r
}
func (r *HttpReq) Jpeg() *HttpReq {
	r.Header.Set("Content-Type", "image/jpeg")
	return r
}
