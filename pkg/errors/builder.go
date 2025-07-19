package errors

import (
	"fmt"
)

type EmptyBuilderFlagListError struct {
}

func (e *EmptyBuilderFlagListError) Error() string { return "empty builder flag list" }

type NotAllFlagsSatisfiedBuilderError struct {
	Flags []string
}

func (e *NotAllFlagsSatisfiedBuilderError) Error() string {
	return fmt.Sprintf("flag [%v] not satisfied", e.Flags)
}
