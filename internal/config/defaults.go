package config

const (
	dbPathEnv   = "TODO_DBPATH"
	webDirEnv   = "TODO_WEBDIR"
	portEnv     = "TODO_PORT"
	passwordEnv = "TODO_PASSWORD"

	defaultDBPath   = "./scheduler.db"
	defaultWebDir   = "web"
	defaultPort     = "7540"
	defaultPassword = ""

	TaskReturnLimit = 50
	DBDateFormat    = "20060102"
	WebDateFormat   = "02.01.2006"
)
