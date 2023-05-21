package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"
)

// RegisterService 在注册中心进行注册
func RegisterService(r Registration) error {
	// 注册每个服务自己的心跳url
	heartbeatURL, err := url.Parse(r.HeartbeatURL)
	if err != nil {
		return err
	}
	http.HandleFunc(heartbeatURL.Path, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// 注册每个服务自己的更新url
	serviceUpdateURL, err := url.Parse(r.ServiceUpdateURL)
	if err != nil {
		return err
	}
	http.Handle(serviceUpdateURL.Path, &serviceUpdateHandler{})

	// 向注册服务发送注册请求
	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	err = encoder.Encode(r)
	if err != nil {
		return err
	}
	resp, err := http.Post(ServicesURL, "application/json", buffer)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed register services and return code is %v", resp.StatusCode)
	}
	return nil
}

// ShutdownService 终止服务
func ShutdownService(r Registration) error {
	request, err := http.NewRequest(http.MethodDelete, ServicesURL, bytes.NewBuffer([]byte(r.ServiceName)))
	if err != nil {
		return err
	}
	request.Header.Add("Content-Type", "text/plain")
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed deregister services and return code is %v", resp.StatusCode)
	}
	return nil
}

type serviceUpdateHandler struct{}

// 当依赖项更新后注册服务调用的接口
func (suh serviceUpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	dec := json.NewDecoder(r.Body)
	var p patch
	err := dec.Decode(&p)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Printf("Updated received %v \n", p)
	prov.Update(p)
}

// registry client记录的服务信息
type providers struct {
	services map[ServiceName][]string
	mutex    *sync.RWMutex
}

// 所有服务共享一个providers
var prov = providers{
	services: make(map[ServiceName][]string, 0),
	mutex:    new(sync.RWMutex),
}

// Update 根据patch来更新client的服务信息
func (p *providers) Update(pat patch) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	for _, patchEntry := range pat.Added {
		serviceName := patchEntry.Name
		if _, ok := p.services[serviceName]; !ok {
			p.services[serviceName] = make([]string, 0)
		}
		p.services[serviceName] = append(p.services[serviceName], patchEntry.URL)
	}

	for _, patchEntry := range pat.Removed {
		serviceName := patchEntry.Name
		if providerURLs, ok := p.services[serviceName]; ok {
			for i := range providerURLs {
				if providerURLs[i] == patchEntry.URL {
					p.services[serviceName] = append(providerURLs[:i], providerURLs[i+1:]...)
				}
			}
		}
	}
}

// 在client端获取指定服务的URLs
func (p *providers) get(name ServiceName) []string {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	urls := p.services[name]
	return urls
}

// GetProvider 在client端获取指定服务的URLs
func GetProvider(name ServiceName) []string {
	return prov.get(name)
}
