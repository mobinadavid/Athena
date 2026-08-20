package utils

func CheckError(in, out error) bool {
	if in == nil || out == nil {
		return false
	}
	if in == out {
		return true
	} else if in.Error() == out.Error() {
		return true
	}
	return false
}
