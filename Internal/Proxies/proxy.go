package Proxies

import (
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"strings"
	"sync"
	"time"
)

type SafeProxyHandler struct {
	cacheMutex       sync.Mutex
	proxyCache       map[string]Proxy
	usedProxiesMutex sync.Mutex
	usedProxies      map[string]bool
}

var SafeProxy *SafeProxyHandler
var once sync.Once

func init() {
	once.Do(func() {
		SafeProxy = &SafeProxyHandler{}
		SafeProxy.proxyCache = map[string]Proxy{}
		SafeProxy.usedProxies = map[string]bool{}
	})
}

type Proxy struct {
	GroupID   int
	GroupName string
	Ip           string
	Port         string
	Username     string
	Password     string
	RequiresAuth bool
	Formatted    *url.URL
	Raw          string
}

type ProxyGroup struct {
	GroupName           string
	GroupID             int
	IsMonitoringProxies bool
	Proxies             []string
}

var illegalProxies = []string{
	"127.0.0.",
	"localhost",
	"0.0.0.0",
}

func (sph *SafeProxyHandler) ParseProxy(proxy string) (Proxy, error) {
	sph.cacheMutex.Lock()
	defer func() {
		sph.cacheMutex.Unlock()
	}()

	if val, ok := sph.proxyCache[proxy]; ok {
		return val, nil
	}

	var proxyInstance Proxy
	var proxyUrl *url.URL
	var err error

	switch split := strings.Split(proxy, ":"); len(split) {
	case 2:

		proxyInstance = Proxy{
			Ip:   split[0],
			Port: split[1],
			Raw:  proxy,
		}

		proxyUrl, err = url.Parse(fmt.Sprintf("http://%s:%s", proxyInstance.Ip, proxyInstance.Port))
		if err != nil {
			return Proxy{}, errors.New("failed to parse proxy")
		}

		proxyInstance.Formatted = proxyUrl
	case 4:

		proxyInstance = Proxy{
			Raw:          proxy,
			Ip:           split[0],
			Port:         split[1],
			Username:     split[2],
			Password:     split[3],
			RequiresAuth: true,
		}

		proxyUrl, err = url.Parse(fmt.Sprintf("http://%s:%s@%s:%s", proxyInstance.Username, proxyInstance.Password, proxyInstance.Ip, proxyInstance.Port))
		if err != nil {
			return Proxy{}, errors.New("failed to parse proxy")
		}

		proxyInstance.Formatted = proxyUrl
	case 1:
		proxyUrl = nil
	default:
		return Proxy{}, errors.New("failed to parse proxy")
	}

	sph.proxyCache[proxy] = proxyInstance

	return proxyInstance, nil
}

func (sph *SafeProxyHandler) GetProxy(proxyList []string) string {

	sph.usedProxiesMutex.Lock()
	rand.Seed(time.Now().UnixNano())

	defer func() {
		sph.usedProxiesMutex.Unlock()
	}()

	rand.Shuffle(len(proxyList), func(i, j int) { proxyList[i], proxyList[j] = proxyList[j], proxyList[i] })

	var unusedProxyList = []string{}
	for _, proxy := range proxyList {
		if _, ok := sph.usedProxies[proxy]; !ok {
			unusedProxyList = append(unusedProxyList, proxy)
		}
	}

	if len(unusedProxyList) == 0 {
		return ""
	}

	proxyToUse := unusedProxyList[rand.Intn(len(unusedProxyList)-0)+0]
	sph.usedProxies[proxyToUse] = true
	return proxyToUse
}

func (sph *SafeProxyHandler) ReleaseProxy(proxy string) {
	sph.usedProxiesMutex.Lock()

	defer func() {
		sph.usedProxiesMutex.Unlock()
	}()

	delete(sph.usedProxies, proxy)
}
