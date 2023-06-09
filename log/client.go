package log

import (
	"bytes"
	"distributed/registry"
	"fmt"
	"net/http"

	stlog "log"
)

// SetClientLogger 将服务的log的Write方法重新为向log service发送请求
func SetClientLogger(serviceURL string, clientService registry.ServiceName) {
	stlog.SetPrefix(fmt.Sprintf("[%v] - ", clientService))
	stlog.SetFlags(0)
	stlog.SetOutput(&clientLogger{url: serviceURL})
}

type clientLogger struct {
	url string
}

func (cl clientLogger) Write(data []byte) (int, error) {
	b := bytes.NewBuffer(data)
	res, err := http.Post(cl.url+"/log", "text/plain", b)
	if err != nil {
		return 0, err
	}
	if res.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("Failed to send log message. Service responded with %d - %s", res.StatusCode, res.Status)
	}
	return len(data), nil
}
