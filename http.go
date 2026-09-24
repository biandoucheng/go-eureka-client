package goeurekaclient

import (
	"bytes"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// httpRoundtrip http传输通信配置
var (
	httpRoundtripMu sync.RWMutex
	httpRoundtrip   http.RoundTripper = newDefaultHTTPTransport()
)

// SetHTTPTransport 设置Eureka客户端使用的HTTP传输器。
// 传入nil时恢复为库的默认传输器。调用方可以传入自定义的
// *http.Transport，以配置代理、超时、TLS证书和连接池等参数。
func SetHTTPTransport(transport http.RoundTripper) {
	if transport == nil {
		transport = newDefaultHTTPTransport()
	}

	httpRoundtripMu.Lock()
	httpRoundtrip = transport
	httpRoundtripMu.Unlock()
}

// newDefaultHTTPTransport 创建库默认的HTTP传输器。
// 保留原有连接池、拨号和超时配置，但不跳过TLS证书校验。
func newDefaultHTTPTransport() http.RoundTripper {
	return &http.Transport{
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
}

// getHTTPTransport 获取当前Eureka客户端使用的HTTP传输器。
func getHTTPTransport() http.RoundTripper {
	httpRoundtripMu.RLock()
	transport := httpRoundtrip
	httpRoundtripMu.RUnlock()
	return transport
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
		Transport: getHTTPTransport(),
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
		Transport: getHTTPTransport(),
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
		Transport: getHTTPTransport(),
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
		Transport: getHTTPTransport(),
	}

	//发起请求
	return cli.Do(request)
}
