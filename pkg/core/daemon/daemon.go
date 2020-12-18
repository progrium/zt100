package daemon

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"reflect"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	L "github.com/progrium/zt100/pkg/misc/logging"
	"github.com/progrium/zt100/pkg/misc/registry"
)

// Initializer is initialized before services are started. Returning
// an error will cancel the start of daemon services.
type Initializer interface {
	InitializeDaemon() error
}

// Terminator is terminated when the daemon gets a stop signal.
type Terminator interface {
	TerminateDaemon(ctx context.Context) error
}

// Service is run after the daemon is initialized.
type Service interface {
	Serve(ctx context.Context)
}

// Framework is a top-level daemon lifecycle manager runs services given to it.
type Framework struct {
	Initializers []Initializer
	Services     []Service
	Terminators  []Terminator
	Logger       L.Logger
	Context      context.Context
	OnFinished   func()
	running      int32
	cancel       context.CancelFunc
	termErrs     chan []error
}

// New builds a daemon configured to run a set of services. The services
// are populated with each other if they have fields that match anything
// that was passed in.
func New(services ...Service) *Framework {
	d := &Framework{}
	d.AddServices(services...)
	return d
}

func Run(reg *registry.Registry) error {
	rv := reflect.ValueOf(&Framework{})
	reg.SelfPopulate()
	reg.ValueTo(rv)
	d := rv.Interface().(*Framework)
	return d.Run(context.Background())
}

// AddServices appends Service and Terminators to daemon
func (d *Framework) AddServices(services ...Service) {
	r, _ := registry.New(d)
	for _, s := range d.Services {
		r.Register(s)
	}
	for _, s := range services {
		r.Register(s)
		d.Services = append(d.Services, s)
		if t, ok := s.(Terminator); ok {
			d.Terminators = append(d.Terminators, t)
		}
	}
	r.SelfPopulate()
}

// Run executes the daemon lifecycle
func (d *Framework) Run(ctx context.Context) error {
	if !atomic.CompareAndSwapInt32(&d.running, 0, 1) {
		return errors.New("already running")
	}

	// call initializers
	for _, i := range d.Initializers {
		if err := i.InitializeDaemon(); err != nil {
			return err
		}
	}

	// finish if no services
	if len(d.Services) == 0 {
		return errors.New("no services to run")
	}

	if ctx == nil {
		ctx = context.Background()
	}
	ctx, cancelFunc := context.WithCancel(ctx)
	d.Context = ctx
	d.cancel = cancelFunc
	d.termErrs = make(chan []error)

	// setup terminators on stop signals
	go TerminateOnSignal(d)
	go TerminateOnContextDone(d)

	var wg sync.WaitGroup
	var running sync.Map
	for _, service := range d.Services {
		running.Store(service, nil)
		wg.Add(1)
		go func(s Service) {
			defer func() {
				if r := recover(); r != nil {
					// needs extra caller skip
					L.Error(d.Logger, "serve: ", r)
				}
			}()
			defer wg.Done()
			s.Serve(d.Context)
			running.Delete(s)
		}(service)
	}

	finished := make(chan bool)
	go func() {
		wg.Wait()
		close(finished)
	}()

	var errs []error
	select {
	case <-finished:
		select {
		case errs = <-d.termErrs:
		default:
		}
	case errs = <-d.termErrs:
		var waiting []reflect.Type
		running.Range(func(k, v interface{}) bool {
			waiting = append(waiting, reflect.ValueOf(k).Type())
			return true
		})
		log.Println("waiting on:", waiting)
		select {
		case <-finished:
		case <-time.After(2 * time.Second):
			// TODO: show/track what servies
			//L.Info(d.Logger, "warning: unfinished services")
			log.Println("warning: unfinished services")
		}
	}

	if d.OnFinished != nil {
		d.OnFinished()
	}

	if len(errs) > 0 {
		return errs[0]
	}

	return nil
}

// Terminate cancels the daemon context and calls Terminators in reverse order
func (d *Framework) Terminate() {
	if d == nil {
		// find these cases and prevent them!
		panic("daemon reference used to Terminate but daemon pointer is nil")
	}

	if !atomic.CompareAndSwapInt32(&d.running, 1, 0) {
		return
	}

	if d.cancel != nil {
		d.cancel()
	}

	// TODO: use for terminate timeout
	ctx := context.Background()

	var errs []error
	for i := len(d.Terminators) - 1; i >= 0; i-- {
		if err := d.Terminators[i].TerminateDaemon(ctx); err != nil {
			errs = append(errs, err)
		}
	}
	d.termErrs <- errs
}

// TerminateOnSignal waits for SIGINT, SIGHUP, SIGTERM, SIGKILL(?) to terminate the daemon.
func TerminateOnSignal(d *Framework) {
	termSigs := make(chan os.Signal, 1)
	signal.Notify(termSigs, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGTERM)
	<-termSigs
	d.Terminate()
}

// TerminateOnContextDone waits for the deamon's context to be canceled.
func TerminateOnContextDone(d *Framework) {
	<-d.Context.Done()
	d.Terminate()
}
