// +build testdata
package testdata

type Type1 struct {
}

func (r Type1) Func2a(arg1 *int) {
	if *arg1 != 0 {
		*arg1 = 1
	}
}

func (r *Type1) Func2b(arg1 *int) {
}

func (r *Type1) Func2c(arg1 *int) {
}
