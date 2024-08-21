package config

var (
	commit  = "unknown"
	date    = "unknown"
	version = "unknown"
)

func PrintVersion() {
	println("Version:", version)
	println("Commit:", commit)
	println("Date:", date)
}

func UserAgent() string {
	return "ewol/" + version
}
