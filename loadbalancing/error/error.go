package error

type NoServersError struct {
}

func (e *NoServersError) Error() string {
	return "No servers found as their count is zero"
}
