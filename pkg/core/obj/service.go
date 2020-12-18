package obj

import (
	"context"

	"github.com/progrium/zt100/pkg/manifold"
	"github.com/progrium/zt100/pkg/manifold/image"
	"github.com/progrium/zt100/pkg/misc/logging"
	"github.com/progrium/zt100/pkg/ui/state"
)

type MountContributor interface {
	MountedComponent(com manifold.Component)
}

type Mountee interface {
	Mounted(obj manifold.Object) error
}

type Service struct {
	Protocol      string
	ListenAddr    string
	MountContribs []MountContributor

	Log   logging.Logger
	Root  manifold.Object
	Image *image.Image
	UI    *state.Framework

	view *ViewState
}

func (s *Service) Serve(ctx context.Context) {
	<-ctx.Done()
}

func (s *Service) Snapshot() error {
	return s.Image.Write(s.Root)
}

func (s *Service) MountedComponent(com manifold.Component, obj manifold.Object) error {
	if initializer, ok := com.Pointer().(Mountee); ok {
		if err := initializer.Mounted(obj); err != nil {
			return err
		}
	}
	for _, contrib := range s.MountContribs {
		contrib.MountedComponent(com)
	}
	return nil
}

func (s *Service) InitializeObject(obj manifold.Object) error {
	for _, com := range obj.Components() {
		if err := s.MountedComponent(com, obj); err != nil {
			return err
		}
	}
	return nil
}
