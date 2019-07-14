// +build testdata
package testdata

func Func4(arg1 *int) {
	if *arg1 != 0 {
		*arg1 = 1
	}
}
