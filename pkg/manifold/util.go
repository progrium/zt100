package manifold

func ExpandPath(o Object, path string) string {
	obj := o.FindChild(path)
	if obj == nil {
		return path
	}
	return obj.Path()
}

func Walk(o Object, fn func(Object) error) error {
	if o.Parent() != nil {
		if err := fn(o); err != nil {
			return err
		}
	}
	for _, child := range o.Children() {
		if err := Walk(child, fn); err != nil {
			return err
		}
	}
	return nil
}
