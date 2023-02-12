// Package util has utility functions and variables
package util

// EqualError compares two error args.
func EqualError(e1, e2 error) bool {
	if e1 == nil && e2 == nil {
		return true
	}
	if (e1 != nil && e2 != nil) && e1.Error() == e2.Error() {
		return true
	}
	return false
}
