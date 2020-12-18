package state

import (
	"errors"
	"io"
	"log"
	"strings"
	"sync"

	"github.com/manifold/qtalk/golang/rpc"
)

type Contributor interface {
	InitializeState() (name string, state interface{})
	UpdateState() error
}

type Framework struct {
	Contributors []Contributor

	state   map[string]interface{}
	clients map[rpc.Caller]struct{}

	mu sync.Mutex
}

func (f *Framework) UpdateState() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	for _, contributor := range f.Contributors {
		if err := contributor.UpdateState(); err != nil {
			return err
		}
	}

	for client := range f.clients {
		_, err := client.Call("state.update", f.state, nil)
		if err != nil {
			if !errors.Is(err, io.EOF) && !IsErrClosed(err) {
				log.Println("UpdateState:", err)
			}
			delete(f.clients, client)
		}
	}

	return nil
}

func IsErrClosed(err error) bool {
	// TODO go1.16: use net.ErrClosed
	return strings.Contains(err.Error(), "use of closed network connection")
}
