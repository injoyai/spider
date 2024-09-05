package tool

import (
	"encoding/json"
	"github.com/injoyai/goutil/oss"
	"net/http"
)

// LoadingCookies 从文件中加载cookie
func LoadingCookies(filename string) ([]*http.Cookie, error) {
	cookies := []*http.Cookie(nil)
	bs, err := oss.ReadBytes(filename)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(bs, &cookies); err != nil {
		return nil, err
	}
	return cookies, nil
}
