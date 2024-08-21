package config

var (
	commit    = "unknown"
	buildDate = "unknown"
	version   = "unknown"
)

func PrintVersion() {
	println("Version:", version)
	println("Commit:", commit)
	println("Build Date:", buildDate)
}

func UserAgent() string {
	return "ewol/" + version
}
