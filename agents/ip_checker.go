package agents

import (
	// "log"
	"net"

	"github.com/chimaera/prototype/core"
)

type IPChecker struct {
	state        *core.State
	orchestrator *core.Orchestrator
}

func NewIPChecker() *IPChecker {
	return &IPChecker{
		state: core.NewState(),
	}
}

func (c *IPChecker) ID() string {
	return "ip:checker"
}

func (c *IPChecker) Register(o *core.Orchestrator) error {
	o.Subscribe("new:hostname", c.onEndpoint)
	o.Subscribe("new:subdomain", c.onEndpoint)

	c.orchestrator = o

	// log.Printf("subscribed %s to `new:hostname` and `new:subdomain` events", c.ID())

	return nil
}

func (c *IPChecker) onEndpoint(hostname string) {
	if c.state.DidProcess(hostname) {
		return
	}

	c.state.Add(hostname)

	// log.Printf("got new endpoint to scan for ip addresses: %s", hostname)

	c.orchestrator.RunTask(func() {
		if addrs, err := net.LookupHost(hostname); err == nil {
			for _, addr := range addrs {
				c.orchestrator.Publish("new:ip", addr)
			}
		}
	})
}