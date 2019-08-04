package proxy

import (
	"crypto/tls"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"

	"domain-proxy/src/base/logger"
)

var log *logger.Entry
var hostPortArr = strings.Split("192.168.1.10.80", ".")

func init() {
	log = logger.NewModuleLogger("proxy")
}

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
	log.Debug("NewReverseProxyIns: ", remote)
	if err != nil {
		return nil
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)
	if protocol == "https" {
		proxy.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
	return proxy
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
	arr = append(append([]string{}, hostPortArr[:5-len(arr)]...), arr...)
	host := strings.Join(arr[:4], ".")
	log.Tracef("key: %s, arr: %v", key, arr)
	return host, arr[4]
}

func NewReverseProxyInsFromKey(key string) *httputil.ReverseProxy {
	host, port := getHostAndPortFromKey(key)
	log.Infof("NewReverseProxyInsFromKey %s -> %s:%s", key, host, port)
	return NewReverseProxyIns(host, port)
}

func (rp *ReverseProxyPool) getProxy(key string) *httputil.ReverseProxy {
	rp.lock.RLock()
	defer rp.lock.RUnlock()
	if rp.cache[key] != nil {
		log.Tracef("getProxy: %s", key)
		return rp.cache[key]
	}
	return nil
}

func (rp *ReverseProxyPool) addProxy(key string) *httputil.ReverseProxy {
	rp.lock.Lock()
	defer rp.lock.Unlock()
	if rp.cache[key] == nil {
		log.Tracef("addProxy: %s", key)
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
