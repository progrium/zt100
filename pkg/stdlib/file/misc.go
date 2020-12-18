package file

type Disposable interface {
	Dispose()
}

type disposer func()

func (d disposer) Dispose() {
	d()
}
