package depman

var container map[string]interface{} = make(map[string]interface{})

func RegisterItem(name string, value interface{}) {
	container[name] = value
}
func GetItem(name string) (interface{}, bool) {
	v, ok := container[name]
	return v, ok
}
