package image

import (
	"encoding/json"
	"fmt"
	"log"
	paths "path"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"sync"

	"github.com/progrium/zt100/pkg/manifold"
	"github.com/progrium/zt100/pkg/manifold/library"
	"github.com/progrium/zt100/pkg/manifold/object"
	"github.com/spf13/afero"
)

const (
	ObjectDir  = "obj"
	ObjectFile = "object.json"

	PackageDir = "pkg"
	UserPkgDir = "usr"
)

type componentInitializer interface {
	InitializeComponent(o manifold.Object)
}

type componentEnabler interface {
	ComponentEnable()
}
type componentDisabler interface {
	ComponentDisable()
}

type Image struct {
	fs       afero.Fs
	objFs    afero.Fs
	pkgFs    afero.Fs
	filepath string

	lastObjPath map[string]string
	writeMu     sync.Mutex
}

func New(filepath string) *Image {
	return &Image{
		filepath:    filepath,
		fs:          afero.NewBasePathFs(afero.NewOsFs(), filepath),
		lastObjPath: make(map[string]string),
	}
}

func ApplyRefs(root manifold.Object, refs []manifold.SnapshotRef) {
	for _, ref := range refs {
		src := root.FindID(ref.ObjectID)
		if src == nil {
			log.Printf("no object found for snapshot ref at %s", ref.ObjectID)
			continue
		}
		parts := strings.Split(ref.TargetID, "/")
		comKey := ""
		objID := parts[0]
		if len(parts) > 1 {
			comKey = parts[1]
		}
		dst := root.FindID(objID)
		if dst == nil {
			log.Printf("no object found for snapshot ref target at %s", ref.TargetID)
			continue
		}
		_, targetType, _ := src.GetField(ref.Path)
		var isPointer bool
		if targetType.Kind() == reflect.Ptr {
			targetType = targetType.Elem()
			isPointer = true
		}
		if targetType.Kind() == reflect.Interface {
			isPointer = true
		}
		ptr := reflect.New(targetType)
		if comKey == "" {
			dst.ValueTo(ptr)
		} else {
			com := dst.Component(comKey)
			ptr = reflect.ValueOf(com.Pointer())
		}
		if !isPointer {
			ptr = reflect.Indirect(ptr)
		}
		src.SetField(ref.Path, ptr.Interface())
	}
}

func (i *Image) Load() (manifold.Object, error) {
	i.objFs = afero.NewBasePathFs(i.fs, ObjectDir)

	if ok, err := afero.Exists(i.objFs, ObjectFile); !ok || err != nil {
		r := object.New("::root")
		r.AppendChild(object.New("System"))
		return r, nil
	}

	obj, refs, err := i.loadObject(i.objFs, "/")
	if err != nil {
		return nil, err
	}
	ApplyRefs(obj, refs)

	err = manifold.Walk(obj, func(o manifold.Object) error {
		if err := o.UpdateRegistry(); err != nil {
			return err
		}
		return nil
	})

	return obj, err
}

func FromSnapshot(snapshot manifold.ObjectSnapshot) (manifold.Object, []manifold.SnapshotRef, error) {
	obj := object.FromSnapshot(snapshot)
	refs := library.LoadComponents(obj, snapshot)
	err := manifold.Walk(obj, func(o manifold.Object) error {
		if err := o.UpdateRegistry(); err != nil {
			return err
		}
		return nil
	})
	return obj, refs, err
}

func (i *Image) loadObject(fs afero.Fs, path string) (manifold.Object, []manifold.SnapshotRef, error) {
	// TODO: Handle missing components?

	buf, err := afero.ReadFile(fs, ObjectFile)
	if err != nil {
		return nil, nil, err
	}

	var snapshot manifold.ObjectSnapshot
	err = json.Unmarshal(buf, &snapshot)
	if err != nil {
		return nil, nil, err
	}

	obj := object.FromSnapshot(snapshot)
	refs := library.LoadComponents(obj, snapshot)
	i.lastObjPath[obj.ID()] = path

	obj.Attrs().Set("_src", filepath.Join(i.filepath, ObjectDir, path, ObjectFile))

	for _, childInfo := range snapshot.Children {
		name := pathNameFromImage(childInfo)
		if ok, err := afero.Exists(fs, name); !ok || err != nil {
			continue
		}
		childFs := afero.NewBasePathFs(fs, name)
		childPath := paths.Join(path, name)
		child, childRefs, err := i.loadObject(childFs, childPath)
		if err != nil {
			return nil, nil, err
		}
		obj.AppendChild(child)
		refs = append(refs, childRefs...)
	}

	return obj, refs, obj.UpdateRegistry()
}

func (i *Image) Write(root manifold.Object) error {
	i.writeMu.Lock()
	defer i.writeMu.Unlock()

	i.objFs = afero.NewBasePathFs(i.fs, ObjectDir)

	if err := i.fs.MkdirAll(ObjectDir, 0755); err != nil {
		return err
	}

	return i.writeObject(i.objFs, "/", root)
}

func (i *Image) writeObject(fs afero.Fs, path string, obj manifold.Object) error {
	i.lastObjPath[obj.ID()] = path
	iobj := obj.Snapshot()

	buf, err := json.MarshalIndent(iobj, "", "  ")
	if err != nil {
		return err
	}
	if err := afero.WriteFile(fs, ObjectFile, buf, 0644); err != nil {
		return err
	}

	for _, child := range obj.Children() {
		childPath := paths.Join(path, pathName(child))
		oldPath := i.lastObjPath[child.ID()]
		if oldPath != "" && oldPath != childPath {
			if err := i.objFs.Rename(oldPath, childPath); err != nil {
				return err
			}
		}
		if err := fs.MkdirAll(pathName(child), 0755); err != nil {
			return err
		}
		childFs := afero.NewBasePathFs(fs, pathName(child))
		if err := i.writeObject(childFs, childPath, child); err != nil {
			return err
		}
	}

	return nil
}

func pathNameFromImage(parts []string) string {
	shortid := parts[0][len(parts[0])-8:]
	exp := regexp.MustCompile("[^a-zA-Z0-9]+")
	name := strings.ToLower(exp.ReplaceAllString(parts[1], ""))
	return fmt.Sprintf("%s-%s", name[:min(8, len(name))], shortid)
}

func pathName(obj manifold.Object) string {
	return pathNameFromImage([]string{obj.ID(), obj.Name()})
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
