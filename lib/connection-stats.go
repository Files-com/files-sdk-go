package lib

import (
	"context"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
)

type Transport struct {
	*http.Transport
	*net.Dialer
	stats   *connectionStats
	statsMu sync.Mutex
}

type connectionStats struct {
	connections map[string]*int32
	mu          sync.RWMutex
}

func (t *Transport) GetConnectionStats() map[string]int {
	stats := t.connectionStats()
	stats.mu.RLock()
	defer stats.mu.RUnlock() // Keep the read lock for the entire function

	copiedMap := make(map[string]int)
	for key, value := range stats.connections {
		copiedMap[key] = int(atomic.LoadInt32(value))
	}
	return copiedMap
}

func (t *Transport) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	conn, err := t.Dialer.DialContext(ctx, network, address)

	if err == nil {
		stats := t.connectionStats()
		stats.mu.Lock()
		counter, ok := stats.connections[address]
		if !ok {
			intCounter := int32(0)
			counter = &intCounter
			stats.connections[address] = counter
		}
		stats.mu.Unlock()
		atomic.AddInt32(counter, 1)
		return &Conn{Conn: conn, counter: counter}, err
	}

	return conn, err
}

func (t *Transport) connectionStats() *connectionStats {
	t.statsMu.Lock()
	defer t.statsMu.Unlock()
	if t.stats == nil {
		t.stats = newConnectionStats()
	}
	return t.stats
}

func newConnectionStats() *connectionStats {
	return &connectionStats{connections: make(map[string]*int32)}
}

type Conn struct {
	net.Conn
	counter *int32
	sync.Once
}

func (c *Conn) Close() (err error) {
	err = c.Conn.Close()
	c.Do(func() { atomic.AddInt32(c.counter, -1) })

	return
}

func GetConnectionStatsFromClient(client *http.Client) (map[string]int, bool) {
	transport, ok := client.Transport.(*Transport)
	if !ok {
		return nil, false
	}
	return transport.GetConnectionStats(), true
}

func CloneHTTPClientWithMaxConnsPerHost(client *http.Client, maxConnsPerHost int) (*http.Client, bool) {
	if client == nil || maxConnsPerHost <= 0 {
		return client, false
	}

	cloned := *client
	switch transport := client.Transport.(type) {
	case *Transport:
		cloned.Transport = transport.cloneWithMaxConnsPerHost(maxConnsPerHost)
	case *http.Transport:
		cloned.Transport = cloneHTTPTransportWithMaxConnsPerHost(transport, maxConnsPerHost)
	case nil:
		defaultTransport, ok := http.DefaultTransport.(*http.Transport)
		if !ok {
			return client, false
		}
		cloned.Transport = cloneHTTPTransportWithMaxConnsPerHost(defaultTransport, maxConnsPerHost)
	default:
		return client, false
	}
	return &cloned, true
}

func CloneHTTPClientWithExactMaxConnsPerHost(client *http.Client, maxConnsPerHost int) (*http.Client, bool) {
	if client == nil || maxConnsPerHost <= 0 {
		return client, false
	}

	cloned := *client
	switch transport := client.Transport.(type) {
	case *Transport:
		cloned.Transport = transport.cloneWithExactMaxConnsPerHost(maxConnsPerHost)
	case *http.Transport:
		cloned.Transport = cloneHTTPTransportWithExactMaxConnsPerHost(transport, maxConnsPerHost)
	case nil:
		defaultTransport, ok := http.DefaultTransport.(*http.Transport)
		if !ok {
			return client, false
		}
		cloned.Transport = cloneHTTPTransportWithExactMaxConnsPerHost(defaultTransport, maxConnsPerHost)
	default:
		return client, false
	}
	return &cloned, true
}

func (t *Transport) cloneWithMaxConnsPerHost(maxConnsPerHost int) *Transport {
	var base *http.Transport
	if t.Transport != nil {
		base = t.Transport
	} else if defaultTransport, ok := http.DefaultTransport.(*http.Transport); ok {
		base = defaultTransport
	} else {
		base = &http.Transport{}
	}

	cloned := &Transport{
		Transport: cloneHTTPTransportWithMaxConnsPerHost(base, maxConnsPerHost),
		Dialer:    t.Dialer,
		stats:     t.connectionStats(),
	}
	if cloned.Dialer == nil {
		cloned.Dialer = &net.Dialer{}
	}
	cloned.Transport.DialContext = cloned.DialContext
	return cloned
}

func (t *Transport) cloneWithExactMaxConnsPerHost(maxConnsPerHost int) *Transport {
	var base *http.Transport
	if t.Transport != nil {
		base = t.Transport
	} else if defaultTransport, ok := http.DefaultTransport.(*http.Transport); ok {
		base = defaultTransport
	} else {
		base = &http.Transport{}
	}

	cloned := &Transport{
		Transport: cloneHTTPTransportWithExactMaxConnsPerHost(base, maxConnsPerHost),
		Dialer:    t.Dialer,
		stats:     t.connectionStats(),
	}
	if cloned.Dialer == nil {
		cloned.Dialer = &net.Dialer{}
	}
	cloned.Transport.DialContext = cloned.DialContext
	return cloned
}

func cloneHTTPTransportWithMaxConnsPerHost(transport *http.Transport, maxConnsPerHost int) *http.Transport {
	cloned := transport.Clone()
	applyHTTPTransportMaxConnsPerHost(cloned, maxConnsPerHost)
	return cloned
}

func cloneHTTPTransportWithExactMaxConnsPerHost(transport *http.Transport, maxConnsPerHost int) *http.Transport {
	cloned := transport.Clone()
	applyHTTPTransportExactMaxConnsPerHost(cloned, maxConnsPerHost)
	return cloned
}

func applyHTTPTransportMaxConnsPerHost(transport *http.Transport, maxConnsPerHost int) {
	if maxConnsPerHost <= 0 {
		return
	}
	if transport.MaxConnsPerHost == 0 || transport.MaxConnsPerHost < maxConnsPerHost {
		transport.MaxConnsPerHost = maxConnsPerHost
	}
	if transport.MaxIdleConns < maxConnsPerHost {
		transport.MaxIdleConns = maxConnsPerHost
	}
	if transport.MaxIdleConnsPerHost < maxConnsPerHost {
		transport.MaxIdleConnsPerHost = maxConnsPerHost
	}
}

func applyHTTPTransportExactMaxConnsPerHost(transport *http.Transport, maxConnsPerHost int) {
	if maxConnsPerHost <= 0 {
		return
	}
	transport.MaxConnsPerHost = maxConnsPerHost
	if transport.MaxIdleConns < maxConnsPerHost {
		transport.MaxIdleConns = maxConnsPerHost
	}
	if transport.MaxIdleConnsPerHost < maxConnsPerHost {
		transport.MaxIdleConnsPerHost = maxConnsPerHost
	}
}
