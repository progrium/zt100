package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	L "github.com/progrium/zt100/pkg/misc/logging"
	"github.com/progrium/zt100/pkg/misc/reflectutil"
)

type Service interface {
	Serve(context.Context)
}

type ShutdownService interface {
	Service
	Shutdown(context.Context) error
}

type Framework struct {
	Logger L.Logger

	services map[Service]func()

	ctx context.Context
	mu  sync.Mutex
}

func (f *Framework) InitializeDaemon() error {
	f.ctx = context.Background()
	f.services = make(map[Service]func())
	return nil
}

func (f *Framework) Start(s Service) {
	f.mu.Lock()
	defer f.mu.Unlock()

	_, exists := f.services[s]
	if exists {
		return
	}

	ctx, stopper := context.WithCancel(f.ctx)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				L.Error(f.Logger, reflectutil.TypeName(s), ": ", r)
			}
		}()
		L.Info(f.Logger, "starting: ", reflectutil.TypeName(s))
		s.Serve(ctx)
	}()

	f.services[s] = stopper
}

func (f *Framework) Stop(s Service) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	stopper, exists := f.services[s]
	if !exists {
		return fmt.Errorf("service not running: %s", reflectutil.TypeName(s))
	}

	L.Info(f.Logger, "stopping: ", reflectutil.TypeName(s))

	if ss, ok := s.(ShutdownService); ok {
		ctx, cancel := context.WithTimeout(f.ctx, 1*time.Second)
		if err := ss.Shutdown(ctx); err != nil {
			cancel()
			return err
		}
		cancel()
	}

	stopper()
	delete(f.services, s)

	return nil
}
