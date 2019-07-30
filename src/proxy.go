package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
)

type ReverseProxyPool struct {
	lock  sync.RWMutex
	cache map[string]*httputil.ReverseProxy
}

func NewReverseProxyIns(host string, port string) *httputil.ReverseProxy {
	if host == "" || port == "" {
		return nil
	}
	protocol := "http"
	if port == "443" {
		protocol = "https"
	}
	remote, err := url.Parse(protocol + "://" + host + ":" + port)
	if err != nil {
		return nil
	}
	return httputil.NewSingleHostReverseProxy(remote)
}

func NewReverseProxyPool() *ReverseProxyPool {
	return &ReverseProxyPool{
		cache: make(map[string]*httputil.ReverseProxy),
	}
}

func getHostAndPortFromKey(key string) (string, string) {
	arr := strings.Split(key, "-")
	if len(arr) < 2 {
		return "", ""
	}
	hostPortArr := strings.Split("192.168.1.10.80", ".")
	arr = append(hostPortArr[:len(arr)-1], arr...)
	host := strings.Join(arr[:4], ".")
	return host, arr[4]
}

func NewReverseProxyInsFromKey(key string) *httputil.ReverseProxy {
	host, port := getHostAndPortFromKey(key)
	return NewReverseProxyIns(host, port)
}

func (rp *ReverseProxyPool) getProxy(key string) *httputil.ReverseProxy {
	rp.lock.RLock()
	defer rp.lock.RUnlock()
	if rp.cache[key] != nil {
		return rp.cache[key]
	}
	return nil
}

func (rp *ReverseProxyPool) addProxy(key string) *httputil.ReverseProxy {
	rp.lock.Lock()
	defer rp.lock.Unlock()
	if rp.cache[key] == nil {
		rp.cache[key] = NewReverseProxyInsFromKey(key)
	}
	return rp.cache[key]
}

func (rp *ReverseProxyPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := strings.Split(r.Host, ".")
	proxy := rp.getProxy(key[0])
	if proxy == nil {
		proxy = rp.addProxy(key[0])
		if proxy == nil {
			_, _ = w.Write([]byte("ok"))
			return
		}
	}
	proxy.ServeHTTP(w, r)
}
