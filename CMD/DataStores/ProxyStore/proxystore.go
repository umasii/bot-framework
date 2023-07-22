package proxystore

import (
	goErrors "errors"
	"fmt"

	errors "github.com/umasii/bot-framework/internal/errors"
	store "github.com/umasii/bot-framework/internal/helpers/datastore"
	proxies "github.com/umasii/bot-framework/internal/proxies"
)

var proxyGroupList []*proxies.ProxyGroup

func AddProxy(groupID int, proxy string) error {
	proxyGroup, err := GetProxyGroupByID(groupID)
	if err != nil {
		return errors.Handler(goErrors.New("Failed to Add Proxy"))
	}
	proxyGroup.Proxies = append(proxyGroup.Proxies, proxy)
	store.Write(&proxyGroupList, "proxies")
	return nil

}

func AddProxyGroup(proxyGroup *proxies.ProxyGroup) {
	proxyGroups := GetProxyGroups()
	proxyGroup.GroupID = getNextProxyGroupID(&proxyGroups)
	if proxyGroup.Proxies == nil {
		proxyGroup.Proxies = []string{}
	}
	proxyGroups = append(proxyGroups, *proxyGroup)
	store.Write(&proxyGroups, "proxies")
}

func GetProxyGroups() []proxies.ProxyGroup {
	var existingProxyGroups []proxies.ProxyGroup
	store.Read(&existingProxyGroups, "proxies")
	return existingProxyGroups
}

func GetProxyGroupByID(proxyGroupID int) (*proxies.ProxyGroup, error) {
	store.Read(&proxyGroupList, "proxies")
	for i := range proxyGroupList {
		if proxyGroupList[i].GroupID == proxyGroupID {
			return proxyGroupList[i], nil
		}
	}
	return nil, errors.Handler(goErrors.New(fmt.Sprintf("Proxy Group with ID %v was not found", proxyGroupID)))
}

func GetMonitoringProxies() []string {
	var existingProxyGroups []proxies.ProxyGroup
	var monitoringProxies []string
	store.Read(&existingProxyGroups, "proxies")

	for i := range existingProxyGroups {
		if existingProxyGroups[i].IsMonitoringProxies {
			monitoringProxies = append(monitoringProxies, existingProxyGroups[i].Proxies...)
		}
	}

	return monitoringProxies
}

func getNextProxyGroupID(proxyGroups *[]proxies.ProxyGroup) int {
	if len(*proxyGroups) > 0 {
		lastProxyGroup := (*proxyGroups)[len(*proxyGroups)-1]
		return lastProxyGroup.GroupID + 1
	}
	return 1
}
