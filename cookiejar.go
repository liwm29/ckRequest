/*
由于标准库的http/cookiejar的jar结构体字段均未导出,只导出了两个函数用于jar的接口,
这使得我们无法dump cookie,因此只能自己实现一个简单的cookiejar作为代替.
这个cookiejar并没有模拟浏览器实现cookie的管理,而只是简单的存储.
如果想实现一个完整的cookie管理模块,将标准库的实现拷贝过来即可

最好的实现应该是container/list,但是list没办法json.marshall
*/

package request

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// todo: 并发加锁

type simpJar struct {
	DB map[string][]*http.Cookie
}

func newSimpJar() *simpJar {
	return &simpJar{DB: make(map[string][]*http.Cookie)}
}

func (j *simpJar) SetCookies(url *url.URL, cookies []*http.Cookie) {
	// we don't want host:port,so strip it
	host := url.Hostname()

	if j.DB[host] == nil {
		j.DB[host] = cookies
		return
	}

	store := j.DB[host]
	for _, v := range cookies {
		if exist, pos := sliceExist(store, v); exist {
			store = sliceDelete(store, pos)
		}
		store = sliceAppend(store, v)
	}
	j.DB[host] = store
	// j.DB[url.Host] = append(j.DB[url.Host], cookies...)
}

func (j *simpJar) Cookies(url *url.URL) []*http.Cookie {
	matchedCookies := make([]*http.Cookie, 0, 10)
	hostNamse := getMatchedHost(url.Hostname())
	for _, v := range hostNamse {
		matchedCookies = append(matchedCookies, j.DB[v]...)
	}
	return matchedCookies
}

func (j *simpJar) StoreCookies(filepath string) error {
	cookieJson, err := json.MarshalIndent(j.DB, "", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath, cookieJson, 0666)
}

func (j *simpJar) LoadCookies(filepath string) error {
	cookieJson, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}
	return json.Unmarshal(cookieJson, &j.DB)
}

func (j *simpJar) Clear() {
	// j.DB = make(map[string][]*http.Cookie)
	for k := range j.DB {
		delete(j.DB, k)
	}
}

func (j *simpJar) AllCookieNames() map[string][]string {
	ret := make(map[string][]string)
	for k, v := range j.DB {
		ret[k] = cookie2names(v)
	}
	return ret
}

func sliceDelete(slice []*http.Cookie, i int) []*http.Cookie {
	return append(slice[:i], slice[i+1:]...)
}

func sliceExist(slice []*http.Cookie, cookie *http.Cookie) (bool, int) {
	for i := range slice {
		if slice[i].Name == cookie.Name {
			return true, i
		}
	}
	return false, -1
}

func sliceAppend(slice []*http.Cookie, cookie ...*http.Cookie) []*http.Cookie {
	return append(slice, cookie...)
}

func getMatchedHost(host string) []string {
	ret := []string{host}
	for {
		pos := strings.Index(host, ".")
		if pos == -1 {
			break
		}
		host = host[pos+1:]
		ret = append(ret, host)
	}
	return ret
}

func cookie2names(cookies []*http.Cookie) []string {
	ret := make([]string, 0, len(cookies))
	for _, v := range cookies {
		ret = append(ret, v.Name)
	}
	return ret
}
