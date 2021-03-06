package cluster

import (
	"fmt"
	"log"

	"github.com/adyatlov/bunxp/xp"

	"github.com/mesosphere/bun/v2/bundle"
)

func init() {
	t := xp.ObjectType{
		Name:              "cluster",
		DisplayName:       "Cluster",
		PluralDisplayName: "Clusters",
		Description:       "DC/OS Cluster",
		Metrics:           []xp.MetricTypeName{"dcos-version"},
		DefaultMetrics:    []xp.MetricTypeName{"dcos-version"},
		FindObject:        findObject,
		GetChildren:       getChildren,
	}
	xp.RegisterObjectType(t)
}

func findObject(b *bundle.Bundle, id xp.ObjectId) (xp.ObjectName, error) {
	if len(b.Hosts) == 0 {
		return "", fmt.Errorf("cannot find any hosts")
	}
	var host bundle.Host
	for _, host = range b.Hosts {
		break
	}
	config := &struct {
		ClusterName string `json:"cluster_name"`
	}{}
	if err := host.ReadJSON("expanded-config", config); err != nil {
		return "", fmt.Errorf("cannot parse cluster name: %s",
			err.Error())
	}
	return xp.ObjectName(config.ClusterName), nil
}

func getChildren(b *bundle.Bundle, id xp.ObjectId) ([]xp.ObjectGroup, error) {
	groups := make([]xp.ObjectGroup, 0, 1)
	group, err := findAgents(b)
	if err != nil {
		return groups, err
	}
	groups = append(groups, group)
	return groups, nil
}

func findAgents(b *bundle.Bundle) (xp.ObjectGroup, error) {
	t := xp.MustGetObjectType("agent")
	group := xp.ObjectGroup{
		Type:    t.Name,
		Objects: make([]xp.Object, 0, len(b.Hosts)),
	}
	for ip, host := range b.Hosts {
		if host.Type != bundle.DTAgent && host.Type != bundle.DTPublicAgent {
			continue
		}
		obj, err := t.New(b, xp.ObjectId(ip))
		if err != nil {
			log.Printf("cannot create agent \"%s\": %s\n", ip, err.Error())
		}
		group.Objects = append(group.Objects, obj)
	}
	return group, nil
}
