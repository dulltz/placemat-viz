package viz

import (
	"strings"

	"github.com/awalterschulze/gographviz"
	"github.com/cybozu-go/placemat"
)

var (
	podsLookUp = make(map[string][]*placemat.PodSpec)
	addrLookUp = make(map[string]string)
)

// GenerateDots generates GraphViz string in dot language
func GenerateDots(cluster *ClusterSpec) (string, error) {
	g := gographviz.NewEscape()
	err := g.SetName("G")
	if err != nil {
		return "", err
	}
	err = g.SetDir(false)
	if err != nil {
		return "", err
	}

	err = connectPods(g, cluster)
	if err != nil {
		return "", err
	}

	err = connectNodesAndPods(g, cluster)

	return g.String(), nil
}

func connectPods(graph *gographviz.Escape, cluster *ClusterSpec) error {
	for _, pod := range cluster.Pods {
		attrs := make(map[string]string)
		attrs[string(gographviz.Shape)] = "diamond"
		err := graph.AddNode("G", pod.Name, attrs)
		if err != nil {
			return err
		}
		for _, nic := range pod.Interfaces {
			podsLookUp[nic.Network] = append(podsLookUp[nic.Network], pod)
			addrLookUp[nic.Network] = nic.Addresses[0]
		}
	}
	for networkName, pods := range podsLookUp {
		if len(pods) != 2 {
			continue
		}
		attrs := make(map[string]string)
		attrs[string(gographviz.Label)] = addrLookUp[networkName]
		err := graph.AddEdge(pods[0].Name, pods[1].Name, false, attrs)
		if err != nil {
			return err
		}
	}
	return nil
}

func connectNodesAndPods(graph *gographviz.Escape, cluster *ClusterSpec) error {
	err := prepareSubGraphs(graph, cluster)
	if err != nil {
		return err
	}
	for _, node := range cluster.Nodes {
		graphName := "cluster-" + strings.Join(node.Interfaces, "_")
		for _, networkName := range node.Interfaces {
			err := graph.AddNode(graphName, node.Name, nil)
			if err != nil {
				return err
			}
			if pods, ok := podsLookUp[networkName]; ok {
				err := graph.AddEdge(pods[0].Name, node.Name, false, nil)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func prepareSubGraphs(graph *gographviz.Escape, cluster *ClusterSpec) error {
	for _, node := range cluster.Nodes {
		interfacesName := strings.Join(node.Interfaces, "_")
		attrs := make(map[string]string)

		var addresses []string
		for _, networkName := range node.Interfaces {
			addresses = append(addresses, addrLookUp[networkName])
		}
		labelName := strings.Join(addresses, "\n")
		attrs[string(gographviz.Label)] = labelName
		graphName := "cluster-" + interfacesName
		err := graph.AddSubGraph("G", graphName, attrs)
		if err != nil {
			return err
		}
	}
	return nil
}
