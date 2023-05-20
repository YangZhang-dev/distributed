package log

import (
	"io"
	stlog "log"
	"net/http"
)

var log *stlog.Logger

const (
	handlePath  = "/log"
	prefix      = "[go] - "
	logFileName = "./cmd/logservice/distributed.log"
)

func Run() {
	log = stlog.New(fileLog(logFileName), prefix, stlog.LstdFlags)
}

func RegisterHandlers() {
	http.HandleFunc(handlePath, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			msg, err := io.ReadAll(r.Body)
			if err != nil || len(msg) == 0 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			write(string(msg))
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	})
}

func write(message string) {
	log.Printf("%v\n", message)
}
