// Sometimes we'll want to sort a collection by something
// other than its natural order. For example, suppose we
// wanted to sort strings by their length instead of
// alphabetically. Here's an example of custom sorts
// in Go.

package main

import "strings"

// ByInstanceNameAsc Custom type to sort our data slice
// In order to sort by a custom function in Go, we need a
// corresponding type. Here we've created a `ByLength`
// type that is just an alias for the builtin `[]string`
// type.
type ByInstanceNameAsc [][]string

// We implement `sort.Interface` - `Len`, `Less`, and
// `Swap` - on our type so we can use the `sort` package's
// generic `Sort` function. `Len` and `Swap`
// will usually be similar across types and `Less` will
// hold the actual custom sorting logic. In our case we
// want to sort in order of increasing string length, so
// we use `len(s[i])` and `len(s[j])` here.
func (s ByInstanceNameAsc) Len() int {
	return len(s)
}
func (s ByInstanceNameAsc) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByInstanceNameAsc) Less(i, j int) bool {
	return strings.Compare(strings.ToLower(s[i][1]), strings.ToLower(s[j][1])) == -1
}

// With all of this in place, we can now implement our
// custom sort by casting the original `fruits` slice to
// `ByLength`, and then use `sort.Sort` on that typed
// slice.

type ByInstanceNameDesc [][]string

func (s ByInstanceNameDesc) Len() int {
	return len(s)
}
func (s ByInstanceNameDesc) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByInstanceNameDesc) Less(i, j int) bool {
	return strings.Compare(strings.ToLower(s[i][1]), strings.ToLower(s[j][1])) == 1
}

type ByInstanceIDAsc [][]string

func (s ByInstanceIDAsc) Len() int {
	return len(s)
}
func (s ByInstanceIDAsc) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByInstanceIDAsc) Less(i, j int) bool {
	return strings.Compare(strings.ToLower(s[i][0]), strings.ToLower(s[j][0])) == -1
}

type ByInstanceIDDesc [][]string

func (s ByInstanceIDDesc) Len() int {
	return len(s)
}
func (s ByInstanceIDDesc) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByInstanceIDDesc) Less(i, j int) bool {
	return strings.Compare(strings.ToLower(s[i][0]), strings.ToLower(s[j][0])) == 1
}
