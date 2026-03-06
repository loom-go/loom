package loom

type Applier interface {
	Apply(parent any) (remove func() error, err error)
}
