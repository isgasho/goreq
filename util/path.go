package util

func NormalizePath(path string) string {
	if path == "/" {
		return ""
	}
	return path
}
