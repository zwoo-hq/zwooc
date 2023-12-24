package helper

func BuildName(parts ...string) string {
	name := ""
	if len(parts) > 0 {
		name = parts[0]
	}
	for i := 1; i < len(parts); i++ {
		name += "/" + parts[1]
	}
	return name
}
