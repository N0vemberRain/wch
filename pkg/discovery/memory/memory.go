package memory

import (
	"context"
	"errors"
	"sync"
	"time"

	"wch/pkg/discovery"
)

type serviceName string
type instanceID string

// Registry defines an in-memory service registry
type Registry struct {
	sync.RWMutex
	serviceAddrs map[serviceName]map[instanceID]*serviceInstance
}

type serviceInstance struct {
	hostPort   string
	lastActive time.Time
}

func NewRegistry() *Registry {
	return &Registry{
		serviceAddrs: map[serviceName]map[instanceID]*serviceInstance{},
	}
}

// Register creaters a service record in the registry
func (r *Registry) Register(
	ctx context.Context,
	instID string,
	svcName string,
	hostPort string,
) error {
	r.Lock()
	defer r.Unlock()

	id := instanceID(instID)
	name := serviceName(svcName)

	if _, ok := r.serviceAddrs[name]; !ok {
		r.serviceAddrs[name] = map[instanceID]*serviceInstance{}
	}

	r.serviceAddrs[name][id] = &serviceInstance{
		hostPort:   hostPort,
		lastActive: time.Now(),
	}

	return nil
}

func (r *Registry) Deregister(
	ctx context.Context,
	instID string,
	svcName string,
) error {
	r.Lock()
	defer r.Unlock()

	id := instanceID(instID)
	name := serviceName(svcName)

	if _, ok := r.serviceAddrs[name]; !ok {
		return nil
	}

	delete(r.serviceAddrs[name], id)
	return nil
}

// ReportHealthState is a push mechanism for reporting healthy state to the registry
func (r *Registry) ReportHealthState(instID string, svcName string) error {
	r.Lock()
	defer r.Unlock()

	id := instanceID(instID)
	name := serviceName(svcName)

	if _, ok := r.serviceAddrs[name]; !ok {
		return errors.New("service is not registered yet")
	}

	if _, ok := r.serviceAddrs[name][id]; !ok {
		return errors.New("service instance is not registered yet")
	}

	r.serviceAddrs[name][id].lastActive = time.Now()
	return nil
}

// ServiceAddresses returns the list of addresses of
// active instances of the given service
func (r *Registry) ServiceAddresses(ctx context.Context, svcName string) ([]string, error) {
	r.RLock()
	defer r.RUnlock()

	name := serviceName(svcName)

	if len(r.serviceAddrs[name]) == 0 {
		return nil, discovery.ErrNotFound
	}

	var res []string

	for _, i := range r.serviceAddrs[name] {
		if i.lastActive.Before(time.Now().Add(-5 * time.Second)) {
			continue
		}

		res = append(res, i.hostPort)
	}

	return res, nil
}
