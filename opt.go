package goutils

import "golang.org/x/exp/constraints"

func HasOpt[T constraints.Integer](owner, target T) bool {
	return (owner & target) == target
}

func AttachOpt[T constraints.Integer](owner, target T) T {
	return owner | target
}

func DetachOpt[T constraints.Integer](owner, target T) T {
	return owner & ^target
}
