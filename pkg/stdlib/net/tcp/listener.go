package tcp

import (
	"log"
	"net"
	"strings"
)

type Listener struct {
	Address string

	l net.Listener
}

func (c *Listener) ComponentEnable() {
	if c.l != nil {
		c.l.Close()
	}
	log.Printf("new tcp listener: %s\n", c.Address)
	var err error
	c.l, err = net.Listen("tcp", c.Address)
	if err != nil {
		panic(err)
	}
}

func (c *Listener) ComponentDisable() {
	if c.l != nil {
		err := c.l.Close()
		if err != nil {
			if !IsErrClosed(err) {
				log.Println(err)
			}
		}
	}
}

func (c *Listener) Accept() (net.Conn, error) {
	return c.l.Accept()
}

func (c *Listener) Close() error {
	return c.l.Close()
}

func (c *Listener) Addr() net.Addr {
	return c.l.Addr()
}

func IsErrClosed(err error) bool {
	// TODO go1.16: use net.ErrClosed
	return strings.Contains(err.Error(), "use of closed network connection")
}
