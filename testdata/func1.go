// +build testdata
package testdata

func Func1(arg1 *int) {
	if *arg1 != 0 {
		*arg1 = 1
	}
}

type Type1 struct {
}

func (r Type1) Func2(arg1 *int) {
}

func (r *Type1) Func3(arg1 *int) {
}

func (r  * Type1 ) Func4(arg1 *int) {
}
