package utils

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/bililive-go/bililive-go/src/configs"
	"github.com/sirupsen/logrus"
)

type ByteCounter struct {
	ReadBytes  int64
	WriteBytes int64
}

type connCounter struct {
	net.Conn
	ByteCounter *ByteCounter
}

func (c *connCounter) Read(b []byte) (n int, err error) {
	n, err = c.Conn.Read(b)
	c.ByteCounter.ReadBytes += int64(n)
	return
}

func (c *connCounter) Write(b []byte) (n int, err error) {
	n, err = c.Conn.Write(b)
	c.ByteCounter.WriteBytes += int64(n)
	return
}

type ConnCounterManagerType struct {
	mapLock sync.Mutex
	bcMap   map[string]*ByteCounter
}

var ConnCounterManager ConnCounterManagerType

func (m *ConnCounterManagerType) SetConn(url string, bc *ByteCounter) {
	m.mapLock.Lock()
	defer m.mapLock.Unlock()
	m.bcMap[url] = bc
}

func (m *ConnCounterManagerType) GetConnCounter(url string) *ByteCounter {
	m.mapLock.Lock()
	defer m.mapLock.Unlock()
	bc, ok := m.bcMap[url]
	if !ok {
		return nil
	}
	return bc
}

func (m *ConnCounterManagerType) PrintMap() {
	m.mapLock.Lock()
	defer m.mapLock.Unlock()
	for url, counter := range m.bcMap {
		logrus.Infof("host[%s] TCP bytes received: %s, sent: %s", url,
			FormatBytes(counter.ReadBytes), FormatBytes(counter.WriteBytes))
	}
}

func CreateConnCounterClient() (*http.Client, error) {
	dialer := func(network, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(network, addr, 10*time.Second)
		if err != nil {
			return nil, err
		}

		byteCounter := ConnCounterManager.GetConnCounter(addr)
		if byteCounter == nil {
			byteCounter = &ByteCounter{}
			ConnCounterManager.SetConn(addr, byteCounter)
		}
		bc := &connCounter{Conn: conn, ByteCounter: byteCounter}
		return bc, nil
	}
	transport := &http.Transport{
		Dial: dialer,
	}
	return &http.Client{Transport: transport}, nil
}

// CreateProxyClient 创建支持代理的HTTP客户端
func CreateProxyClient(proxyConfig *configs.Proxy) (*http.Client, error) {
	if !proxyConfig.Enable {
		return &http.Client{Timeout: 30 * time.Second}, nil
	}

	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			var proxyURL string
			
			// 根据请求的协议选择代理
			if req.URL.Scheme == "https" && proxyConfig.HttpsUrl != "" {
				proxyURL = proxyConfig.HttpsUrl
			} else if proxyConfig.HttpUrl != "" {
				proxyURL = proxyConfig.HttpUrl
			} else {
				return nil, nil // 不使用代理
			}

			// 解析代理URL
			parsedURL, err := url.Parse(proxyURL)
			if err != nil {
				return nil, err
			}

			// 如果提供了用户名和密码，添加到代理URL
			if proxyConfig.Username != "" && proxyConfig.Password != "" {
				parsedURL.User = url.UserPassword(proxyConfig.Username, proxyConfig.Password)
			}

			return parsedURL, nil
		},
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: false},
		TLSHandshakeTimeout:   30 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	return &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}, nil
}

// CreateConnCounterProxyClient 创建支持代理和连接计数的HTTP客户端
func CreateConnCounterProxyClient(proxyConfig *configs.Proxy) (*http.Client, error) {
	dialer := func(network, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(network, addr, 10*time.Second)
		if err != nil {
			return nil, err
		}

		byteCounter := ConnCounterManager.GetConnCounter(addr)
		if byteCounter == nil {
			byteCounter = &ByteCounter{}
			ConnCounterManager.SetConn(addr, byteCounter)
		}
		bc := &connCounter{Conn: conn, ByteCounter: byteCounter}
		return bc, nil
	}

	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			var proxyURL string
			
			if req.URL.Scheme == "https" && proxyConfig.HttpsUrl != "" {
				proxyURL = proxyConfig.HttpsUrl
			} else if proxyConfig.HttpUrl != "" {
				proxyURL = proxyConfig.HttpUrl
			} else {
				return nil, nil
			}

			parsedURL, err := url.Parse(proxyURL)
			if err != nil {
				return nil, err
			}

			if proxyConfig.Username != "" && proxyConfig.Password != "" {
				parsedURL.User = url.UserPassword(proxyConfig.Username, proxyConfig.Password)
			}

			return parsedURL, nil
		},
		Dial: dialer,
	}

	return &http.Client{Transport: transport}, nil
}
