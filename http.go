package goeurekaclient

import (
	"bytes"
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"time"
)

// httpRoundtrip http传输通信配置
var httpRoundtrip = &http.Transport{
	TLSClientConfig: &tls.Config{
		InsecureSkipVerify: true,
	},
	DisableKeepAlives: false,
	Proxy:             http.ProxyFromEnvironment,
	DialContext: (&net.Dialer{
		Timeout:   5 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext,
	MaxIdleConns:        100,
	MaxIdleConnsPerHost: 10,
	IdleConnTimeout:     30 * time.Second,
	TLSHandshakeTimeout: 2 * time.Second,
}

// HttpGet http Get 请求
func HttpGet(ul string, header http.Header, params url.Values, tmout int64) (*http.Response, error) {
	// 参数处理
	query := ""
	if params != nil {
		query += "?" + params.Encode()
	}

	// 实例化请求配置
	request, err := http.NewRequest("GET", ul+query, nil)
	if err != nil {
		return nil, err
	}

	//请求头处理
	if header != nil {
		request.Header = header
	}

	cli := http.Client{
		Timeout:   time.Second * time.Duration(tmout),
		Transport: httpRoundtrip,
	}

	//发起请求
	return cli.Do(request)
}

// HttpPost http Post 请求
func HttpPost(ul string, header http.Header, data []byte, tmout int64) (*http.Response, error) {
	// 实例化请求配置
	request, err := http.NewRequest("POST", ul, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	//请求头处理
	if header != nil {
		request.Header = header
	}

	cli := http.Client{
		Timeout:   time.Second * time.Duration(tmout),
		Transport: httpRoundtrip,
	}

	//发起请求
	return cli.Do(request)
}

// HttpPut http put 请求
func HttpPut(ul string, header http.Header, data []byte, tmout int64) (*http.Response, error) {
	// 实例化请求配置
	request, err := http.NewRequest("PUT", ul, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	//请求头处理
	if header != nil {
		request.Header = header
	}

	cli := http.Client{
		Timeout:   time.Second * time.Duration(tmout),
		Transport: httpRoundtrip,
	}

	//发起请求
	return cli.Do(request)
}

// HttpDelete http Delete 请求
func HttpDelete(ul string, header http.Header, params url.Values, tmout int64) (*http.Response, error) {
	// 参数处理
	query := ""
	if params != nil {
		query += "?" + params.Encode()
	}

	// 实例化请求配置
	request, err := http.NewRequest("DELETE", ul+query, nil)
	if err != nil {
		return nil, err
	}

	//请求头处理
	if header != nil {
		request.Header = header
	}

	cli := http.Client{
		Timeout:   time.Second * time.Duration(tmout),
		Transport: httpRoundtrip,
	}

	//发起请求
	return cli.Do(request)
}
