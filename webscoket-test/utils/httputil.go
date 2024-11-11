package utils

import (
	"bytes"
	"context"
	"fmt"
	"golang.org/x/net/proxy"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
)

type ProxyInfo struct {
	IsHttp        bool
	ProxyIp       string
	ProxyUser     string
	ProxyPassword string
}

func HttpRequset(reqUrl string, method string, params map[string]string, body []byte, headers map[string]string, reqProxy *ProxyInfo) ([]byte, error) {

	var client *http.Client
	var err error

	if reqProxy != nil {
		if reqProxy.IsHttp {

			urlProxy, _ := url.Parse(reqProxy.ProxyIp)
			client = &http.Client{
				Transport: &http.Transport{
					Proxy: http.ProxyURL(urlProxy),
				},
			}
		} else {
			var proxyUser *proxy.Auth
			//设定账号和用户名
			if reqProxy.ProxyUser != "" && reqProxy.ProxyPassword != "" {
				proxyUser = &proxy.Auth{
					User:     reqProxy.ProxyUser,
					Password: reqProxy.ProxyPassword,
				}
			} else {
				proxyUser = nil
			}

			dialer, err := proxy.SOCKS5("tcp", reqProxy.ProxyIp,
				proxyUser,
				&net.Dialer{
					Timeout:  10 * time.Second,
					Deadline: time.Now().Add(time.Second * 10),
				},
			)
			if err != nil {
				return nil, err
			}

			transport := &http.Transport{
				Proxy: nil,
				DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					return dialer.Dial(network, addr)
				},
				TLSHandshakeTimeout: 10 * time.Second,
				MaxIdleConnsPerHost: -1,   //连接池禁用缓存
				DisableKeepAlives:   true, //禁用客户端连接缓存到连接池
			}

			client = &http.Client{Transport: transport}
			if err != nil {
				return nil, err
			}
		}
	} else {
		client = &http.Client{
			Transport: &http.Transport{
				DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					conn, err := net.DialTimeout(network, addr, time.Second*10) //设置建立连接超时
					if err != nil {
						return nil, err
					}
					conn.SetDeadline(time.Now().Add(time.Second * 10)) //设置发送接受数据超时
					return conn, nil
				},
				ResponseHeaderTimeout: time.Second * 10,
				MaxIdleConnsPerHost:   -1,   //禁用连接池缓存
				DisableKeepAlives:     true, //禁用客户端连接缓存到连接池
			},
		}
	}
	if params != nil {
		urlValues := url.Values{}
		for key, value := range params {
			urlValues.Add(key, value)
		}
		reqUrl = fmt.Sprintf("%s?%s", reqUrl, urlValues.Encode())
	}

	req, err := http.NewRequest(method, reqUrl, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	if headers != nil {
		for key, value := range headers {
			req.Header.Set(key, value)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(resp.Body)
}
