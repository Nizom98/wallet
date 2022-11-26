package utils

func Ptr[E any](in E) *E {
	return &in
}
