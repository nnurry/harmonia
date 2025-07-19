package types

import (
	"github.com/nnurry/harmonia/pkg/errors"
)

type BuilderFlag interface {
	Name() string
}

type BuilderFlagMap struct {
	m map[BuilderFlag]bool
}

func NewEmptyFlapMap() (*BuilderFlagMap, error) {
	return &BuilderFlagMap{m: make(map[BuilderFlag]bool)}, nil
}

func NewFlagMapFromBuilderFlags(flags []BuilderFlag, defaultFlags []BuilderFlag, fallbackDefault bool) (*BuilderFlagMap, error) {
	if (len(flags) < 1) && fallbackDefault {
		return nil, &errors.EmptyBuilderFlagListError{}
	}

	builderFlagMap, _ := NewEmptyFlapMap()

	if len(defaultFlags) < 1 {
		return builderFlagMap, nil
	}

	for _, flag := range defaultFlags {
		builderFlagMap.m[flag] = false
	}

	return builderFlagMap, nil
}

func (bfm BuilderFlagMap) mark(flag BuilderFlag, value bool) {
	bfm.m[flag] = value
}

func (bfm BuilderFlagMap) MarkAsChecked(flag BuilderFlag) {
	bfm.mark(flag, true)
}

func (bfm BuilderFlagMap) MarkAsUnchecked(flag BuilderFlag) {
	bfm.mark(flag, false)
}

func (bfm BuilderFlagMap) Verify() error {
	unsatisfiedFlags := []string{}
	for flag, satisfied := range bfm.m {
		if !satisfied {
			unsatisfiedFlags = append(unsatisfiedFlags, flag.Name())
		}
	}

	if len(unsatisfiedFlags) > 0 {
		return &errors.NotAllFlagsSatisfiedBuilderError{Flags: unsatisfiedFlags}
	}

	return nil
}
