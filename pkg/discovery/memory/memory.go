package discovery

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/karaMuha/go-movie/pkg/discovery"
)

type MemoryRegistry struct {
	sync.RWMutex
	serviceAddresses map[string]map[string]*serviceInstance
}

var _ discovery.Registry = (*MemoryRegistry)(nil)

type serviceInstance struct {
	hostPort   string
	lastActive time.Time
}

func NewMemoryRegistry() MemoryRegistry {
	return MemoryRegistry{
		serviceAddresses: make(map[string]map[string]*serviceInstance),
	}
}

// Register creates a service record in the registry.
func (r *MemoryRegistry) Register(ctx context.Context, instanceID string, serviceName string, hostPort string) error {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.serviceAddresses[serviceName]; !ok {
		r.serviceAddresses[serviceName] = map[string]*serviceInstance{}
	}
	r.serviceAddresses[serviceName][instanceID] = &serviceInstance{hostPort: hostPort, lastActive: time.Now()}
	return nil
}

// Deregister removes a service record from the registry.
func (r *MemoryRegistry) Deregister(ctx context.Context, instanceID string, serviceName string) error {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.serviceAddresses[serviceName]; !ok {
		return nil
	}
	delete(r.serviceAddresses[serviceName], instanceID)
	return nil
}

// ReportHealthyState is a push mechanism for reporting healthy state to the registry.
func (r *MemoryRegistry) ReportHealthyState(instanceID string, serviceName string) error {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.serviceAddresses[serviceName]; !ok {
		return errors.New("service is not registered yet")
	}
	if _, ok := r.serviceAddresses[serviceName][instanceID]; !ok {
		return errors.New("service instance is not registered yet")
	}
	r.serviceAddresses[serviceName][instanceID].lastActive = time.Now()
	return nil
}

// ServiceAddresses returns the list of addresses of active instances of the given service.
func (r *MemoryRegistry) ServiceAddresses(ctx context.Context, serviceName string) ([]string, error) {
	r.RLock()
	defer r.RUnlock()
	if len(r.serviceAddresses[serviceName]) == 0 {
		return nil, discovery.ErrNotFound
	}
	var res []string
	for _, i := range r.serviceAddresses[serviceName] {
		if i.lastActive.Before(time.Now().Add(-5 * time.Second)) {
			continue
		}
		res = append(res, i.hostPort)
	}
	return res, nil
}
