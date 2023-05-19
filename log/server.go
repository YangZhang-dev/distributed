package log

import (
	"io/ioutil"
	stlog "log"
	"net/http"
)

var log *stlog.Logger

const (
	handlePath = "/log"
	prefix     = "go: "
)

func Run(destination string) {
	log = stlog.New(fileLog(destination), prefix, stlog.LstdFlags)
}

func RegisterHandlers() {
	http.HandleFunc(handlePath, func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodPost:
			msg, err := ioutil.ReadAll(request.Body)
			if err != nil {
				writer.WriteHeader(http.StatusBadRequest)
				return
			}
			write(string(msg))
		default:
			writer.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	})
}

func write(message string) {
	log.Printf("%v\n", message)
}
