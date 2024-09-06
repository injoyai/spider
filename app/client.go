package app

import (
	"crypto/tls"
	"golang.org/x/net/proxy"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	*http.Client
	Proxy string
}

func newClient(timeout time.Duration, proxyUrl string) *Client {
	transport := &http.Transport{
		//连接结束后会直接关闭,不复用
		DisableKeepAlives: true,
		TLSClientConfig: &tls.Config{
			//设置可以访问HTTPS
			InsecureSkipVerify: true,
		},
	}

	if len(proxyUrl) > 0 {
		p, err := url.Parse(proxyUrl)
		if err == nil {
			switch p.Scheme {
			case "http", "https":
				transport.Proxy = http.ProxyURL(p)
			case "socks5", "socks5h":
				dialer, err := proxy.FromURL(p, proxy.Direct)
				if err == nil {
					transport.Dial = dialer.Dial
				}
			default:
				transport.Proxy = http.ProxyURL(p)
			}
		}
	}

	return &Client{
		Client: &http.Client{
			Transport: transport,
			Timeout:   timeout,
		},
		Proxy: proxyUrl,
	}
}
