package router

var (
	openApi = make(map[string]bool)
)

func isOpen(path string) bool {
	if _, ok := openApi[path]; ok {
		return true
	}
	return false
}
