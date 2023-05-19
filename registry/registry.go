package registry

import (
	"fmt"
	"sync"
)

type registry struct {
	registrations []Registration
	mutex         *sync.Mutex
}

func (r *registry) add(reg Registration) error {
	r.mutex.Lock()
	r.registrations = append(r.registrations, reg)
	r.mutex.Unlock()
	return nil
}
func (r *registry) remove(name ServiceName) error {
	for i := range r.registrations {
		if r.registrations[i].ServiceName == name {
			r.mutex.Lock()
			r.registrations = append(r.registrations[:i], r.registrations[i+1:]...)
			r.mutex.Unlock()
			return nil
		}
	}
	return fmt.Errorf("%v not found", name)
}
