package rpslimit

type Interface interface {
	Allow() bool
}
