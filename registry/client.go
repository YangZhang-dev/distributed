package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func RegisterService(r Registration) error {
	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	err := encoder.Encode(r)
	if err != nil {
		return err
	}
	resp, err := http.Post(ServicesURL, "application/json", buffer)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed register service and return code is %v", resp.StatusCode)
	}
	return nil
}
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
		return fmt.Errorf("failed deregister service and return code is %v", resp.StatusCode)
	}
	return nil
}
