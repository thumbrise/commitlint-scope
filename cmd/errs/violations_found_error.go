package errs

import (
	"fmt"
)

type ViolationsFoundError struct {
	num int
}

func NewViolationsFound(num int) *ViolationsFoundError {
	return &ViolationsFoundError{num: num}
}

func (v ViolationsFoundError) Error() string {
	return fmt.Sprintf("%d violations found", v.num)
}

func (v ViolationsFoundError) ExitCode() int {
	return 2
}
