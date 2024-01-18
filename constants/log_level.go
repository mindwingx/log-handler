package constants

const (
	Trace = "trace" //Only when I would be "tracing" the code and trying to find one part of a function specifically.
	Debug = "debug" //Information that is diagnostically helpful to people more than just developers ( IT, sysadmins, etc.).
	Info  = "info"  //Generally useful information to log (service start/stop, configuration assumptions, etc.).Info I want to always have available but usually don't care about under normal circumstances. This is my out-of-the-box config level.
	Warn  = "warn"  //Anything that can potentially cause application oddities, but for which I am automatically recovering.(Such as switching from a primary to backup server, retrying an operation, missing secondary data, etc.)
	Error = "error" //Any error which is fatal to the operation, but not the service or application (can't open a required file, missing data, etc.). These errors will force user (administrator, or direct user) intervention. These are usually reserved (in my apps) for incorrect connection strings, missing services, etc.
	Fatal = "fatal" //Any error that is forcing a shutdown of the service or application to prevent data loss (or further data loss).I reserve these only for the most heinous errors and situations where there is guaranteed to have been data corruption or loss.
)

var LogLevels = []string{Trace, Debug, Info, Warn, Error, Fatal}
