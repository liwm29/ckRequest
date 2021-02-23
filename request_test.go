package request

import (
	"fmt"
	"testing"
)

func TestAll(t *testing.T) {
	c := NewClient()
	resp := Get("https://portal.sysu.edu.cn/").Do(c)
	fmt.Println(resp.String())
}
