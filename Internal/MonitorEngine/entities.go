package MonitorEngine

type MonitorResp struct {
	CurrentStock bool
	Info         interface{} // Modules may need info returned via monitor (ie PID / ATC ids, so they can just pass int)
}

type MonitorInfo struct {
	Site        string
	Identifier  string
	Task        MonitorTask
}

type MonitorTask interface {
	Create(Site, Identifier string)
	Initialize()
	CheckStock()
	Start(checkStock func())
	Stop()
}
