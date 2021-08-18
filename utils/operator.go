package utils

func Ternary(c bool, t interface{}, f interface{}) interface{} {
	if c {
		return t
	}
	return f
}

func StrTernary(c bool, t string, f string) string {
	if c {
		return t
	}
	return f
}
