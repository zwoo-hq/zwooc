package helper

func BuildName(parts ...string) string {
	name := ""
	for _, part := range parts {
		name += "/" + part
	}
	return name
}
