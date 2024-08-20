package balancer

import (
	"net/url"
	"sync/atomic"
)

type ServerPool struct {
	backends []*Backend
	current  uint64
}

func NewServerPool() *ServerPool {
	return &ServerPool{
		backends: make([]*Backend, 0),
		current:  0,
	}
}

func (p *ServerPool) RegisterBackend(url *url.URL) {
	backend := NewBackend(url)
	p.backends = append(p.backends, backend)
}

func (p *ServerPool) NextIndex() int {
	if len(p.backends) == 0 {
		return -1
	}
	next := atomic.AddUint64(&p.current, 1) % uint64(len(p.backends))
	return int(next)
}

func (p *ServerPool) NextPeer() *Backend {
	// next might not be actually alive
	next := p.NextIndex()
	if next == -1 {
		return nil
	}
	stop := next + len(p.backends)
	for i := next; i <= stop; i++ {
		if p.backends != nil && p.backends[i%len(p.backends)].Alive {
			atomic.StoreUint64(&p.current, uint64(i))
			return p.backends[i]
		}
	}
	return nil
}
