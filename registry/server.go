package registry

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

// ServerPort registry service 启动端口
const ServerPort = ":3000"

// ServicesURL 服务url
const ServicesURL = "http://127.0.0.1" + ServerPort + "/services"

var reg = registry{
	registrations: make([]Registration, 0),
	mutex:         new(sync.Mutex),
}

type RegistryService struct{}

func (rs *RegistryService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Request received!")
	switch r.Method {
	case http.MethodPost:
		decoder := json.NewDecoder(r.Body)
		var r Registration
		err := decoder.Decode(&r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Printf("Adding service: %v with URL: %v\n", r.ServiceName, r.ServiceURL)
		err = reg.add(r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	case http.MethodDelete:
		payload, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		serviceName := ServiceName(payload)
		log.Printf("Removing %v", serviceName)
		err = reg.remove(serviceName)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
