package layer

var Layers = []string{"database", "server", "browser", "ios", "android"}

func GetCoreLatest(layer string, core string) string {
	if core == "" || layer == "" {
		return ""
	}

	return "1.0.0"
}
