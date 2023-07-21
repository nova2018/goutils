package goutils

import (
	"math/rand"
	"time"
)

func Shuffle[T any](arr []T) []T {
	l := len(arr)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := l - 1; i > 0; i-- {
		tmp := r.Intn(i)
		arr[tmp], arr[i] = arr[i], arr[tmp]
	}
	return arr
}

func Filter[T any](arr []T, predicate func(T) bool) []T {
	result := make([]T, 0, len(arr))
	for _, elem := range arr {
		if predicate(elem) {
			result = append(result, elem)
		}
	}
	return result
}

func Diff[T comparable](in []T, slices ...[]T) []T {
	var (
		length  = len(in)
		seen    = make(map[T]bool, length)
		results = make([]T, 0, length)
	)
	for i := 0; i < length; i++ {
		seen[in[i]] = true
	}
	for _, slice := range slices {
		for _, s := range slice {
			if seen[s] {
				seen[s] = false
			}
		}
	}

	for i := 0; i < length; i++ {
		v := in[i]
		if seen[v] {
			results = append(results, v)
		}
	}
	return results
}

func Unique[T comparable](in []T) []T {
	var (
		length  = len(in)
		seen    = make(map[T]struct{}, length)
		results = make([]T, 0, length)
	)

	for i := 0; i < length; i++ {
		v := in[i]

		if _, ok := seen[v]; ok {
			continue
		}

		seen[v] = struct{}{}
		results = append(results, v)
	}

	return results
}

func Intersect[T comparable](in []T, slices ...[]T) []T {
	var (
		length   = len(in)
		sliceLen = len(slices)
		seen     = make(map[T]int, len(slices))
		results  = make([]T, 0, length)
	)
	for _, s := range in {
		seen[s] = 0
	}
	for _, slice := range slices {
		for _, r := range slice {
			if _, ok := seen[r]; ok {
				seen[r] = seen[r] + 1
			}
		}
	}
	for i := 0; i < length; i++ {
		v := in[i]
		if seen[v] == sliceLen {
			results = append(results, v)
		}
	}
	return results
}

func Map[T1, T2 any](list []T1, fn func(T1) T2) []T2 {
	l := len(list)
	r := make([]T2, l)
	for i := 0; i < l; i++ {
		r[i] = fn(list[i])
	}
	return r
}

func Merge[T any](list []T, arrays ...[]T) []T {
	var length = len(list)
	for _, arr := range arrays {
		length += len(arr)
	}
	result := make([]T, 0, length)
	result = append(result, list...)
	for _, arr := range arrays {
		result = append(result, arr...)
	}
	return result
}
