// +build testdata
package testdata

// This file is not in testdata_set.txt

func Func5(arg1 *int) {
	if *arg1 != 0 {
		*arg1 = 1
	}
}
