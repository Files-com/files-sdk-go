package lib

import (
	"context"
	"net"
	"net/http"
	"sync"
)

type Transport struct {
	*http.Transport
	*net.Dialer
	Connections map[string]int
	mu          sync.Mutex
}

func (t *Transport) GetConnectionStats() map[string]int {
	t.mu.Lock()
	defer t.mu.Unlock()
	copiedMap := make(map[string]int)

	for key, value := range t.Connections {
		copiedMap[key] = value
	}
	return copiedMap
}

func (t *Transport) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	conn, err := t.Dialer.DialContext(ctx, network, address)

	t.mu.Lock()
	defer t.mu.Unlock()
	if err == nil {
		t.Connections[address]++
		return &Conn{Conn: conn, Transport: t, address: address}, err
	}

	return conn, err
}

type Conn struct {
	net.Conn
	*Transport
	address string
	sync.Once
}

func (c *Conn) Close() (err error) {
	err = c.Conn.Close()
	c.Do(func() {
		c.mu.Lock()
		defer c.mu.Unlock()
		c.Transport.Connections[c.address]--
	})

	return
}

func GetConnectionStatsFromClient(client *http.Client) (map[string]int, bool) {
	transport, ok := client.Transport.(*Transport)
	if !ok {
		return nil, false
	}
	return transport.GetConnectionStats(), true
}
