package generics

func Throw(a any) {
	if a != nil {
		panic(a)
	}
}
func Must[R any](r R, e error) R {
	Throw(e)
	return r
}
