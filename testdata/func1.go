// +build testdata
package testdata

func Func1(arg1 *int) {
	if *arg1 != 0 {
		*arg1 = 1
	}
}
