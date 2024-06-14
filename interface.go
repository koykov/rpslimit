package rpslimit

type RPSLimiter interface {
	Allow() bool
}
