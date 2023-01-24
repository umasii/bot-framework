package MonitorEngine

import (
	"fmt"
	"sync"
)

type MonitorEngine struct {
	CurrentlyMonitoring map[string][]chan MonitorResp
}

type StockAlert struct {
	StockStatus  bool   `json:"StockStatus"`
	Pid          string `json:"Pid"`
	Site         string `json:"Site"`
	Endpoint     string `json:"Endpoint"`
	FormData     string `json:"FormData"`
	Method       string `json:"Method"`
	GoodResponse string `json:"Good Response"`
}

var MEngine *MonitorEngine
var once sync.Once

func init() {
	once.Do(func() {
		MEngine = &MonitorEngine{}
		MEngine.Initialize()
	})
}

func (m *MonitorEngine) Initialize() {
	m.CurrentlyMonitoring = make(map[string][]chan MonitorResp)
}

func (m *MonitorEngine) Monitor(info MonitorInfo) chan MonitorResp {
	// unique identifier to lookup tasks in future
	key := fmt.Sprintf("%s:%s", info.Site, info.Identifier)
	fmt.Println(key)

	newChan :=  make(chan MonitorResp)

	if _, ok := m.CurrentlyMonitoring[key]; ok {
		fmt.Println(fmt.Sprintf("Already monitoring for [%s]!", key))

		m.CurrentlyMonitoring[key] = append(m.CurrentlyMonitoring[key], newChan)
		return newChan

	}

	m.CurrentlyMonitoring[key] = append(m.CurrentlyMonitoring[key], newChan)

	info.Task.Create(info.Identifier, info.Site)

	go info.Task.Start(info.Task.CheckStock)



	fmt.Println(fmt.Sprintf("Started new Monitor Task for [%s]!", key))

	return newChan
}

func (m *MonitorEngine) CheckForMonitoring(site, identifier string) bool {
	key := fmt.Sprintf("%s:%s", site, identifier)

	fmt.Println(m.CurrentlyMonitoring[key])

	if _, ok := m.CurrentlyMonitoring[key]; ok {
		return true
	}
	return false
}

func (m *MonitorEngine) AlertStock(site, identifier string, status bool, info StockAlert) {
	pushDown := MonitorResp{
		CurrentStock: status,
		Info: info,
	}
	key := fmt.Sprintf("%s:%s", site, identifier)

	fmt.Printf("LEN: %d \n ", len(m.CurrentlyMonitoring[key]))
	fmt.Println(m.CurrentlyMonitoring)
	for iter, task := range m.CurrentlyMonitoring[key] {
		fmt.Printf("iter %d \n", iter)
		taskX := task
		go func () {
			taskX <- pushDown
		} ()

	}
}