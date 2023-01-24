package ProxyStore

import (
	"errors"
	"fmt"

	"github.com/cicadaaio/LVBot/Internal/Errors"
	store "github.com/cicadaaio/LVBot/Internal/Helpers/DataStore"
	"github.com/cicadaaio/LVBot/Internal/Proxies"
)

var proxyGroupList []*Proxies.ProxyGroup

func AddProxy(groupID int, proxy string) error {
	proxyGroup, err := GetProxyGroupByID(groupID)
	if err != nil {
		return Errors.Handler(errors.New("Failed to Add Proxy"))
	}
	proxyGroup.Proxies = append(proxyGroup.Proxies, proxy)
	store.Write(&proxyGroupList, "proxies")
	return nil

}

func AddProxyGroup(proxyGroup *Proxies.ProxyGroup) {
	proxyGroups := GetProxyGroups()
	proxyGroup.GroupID = getNextProxyGroupID(&proxyGroups)
	if proxyGroup.Proxies == nil {
		proxyGroup.Proxies = []string{}
	}
	proxyGroups = append(proxyGroups, *proxyGroup)
	store.Write(&proxyGroups, "proxies")
}

func GetProxyGroups() []Proxies.ProxyGroup {
	var existingProxyGroups []Proxies.ProxyGroup
	store.Read(&existingProxyGroups, "proxies")
	return existingProxyGroups
}

func GetProxyGroupByID(proxyGroupID int) (*Proxies.ProxyGroup, error) {
	store.Read(&proxyGroupList, "proxies")
	for i := range proxyGroupList {
		if proxyGroupList[i].GroupID == proxyGroupID {
			return proxyGroupList[i], nil
		}
	}
	return nil, Errors.Handler(errors.New(fmt.Sprintf("Proxy Group with ID %v was not found", proxyGroupID)))
}

func GetMonitoringProxies() []string {
	var existingProxyGroups []Proxies.ProxyGroup
	var monitoringProxies []string
	store.Read(&existingProxyGroups, "proxies")

	for i := range existingProxyGroups {
		if existingProxyGroups[i].IsMonitoringProxies {
			monitoringProxies = append(monitoringProxies, existingProxyGroups[i].Proxies...)
		}
	}

	return monitoringProxies
}

func getNextProxyGroupID(proxyGroups *[]Proxies.ProxyGroup) int {
	if len(*proxyGroups) > 0 {
		lastProxyGroup := (*proxyGroups)[len(*proxyGroups)-1]
		return lastProxyGroup.GroupID + 1
	}
	return 1
}
