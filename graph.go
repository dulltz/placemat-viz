package viz

import (
	"strings"

	"github.com/awalterschulze/gographviz"
	"github.com/cybozu-go/placemat"
)

type Visualizer struct {
	podsLookUp map[string][]*placemat.PodSpec
	addrLookUp map[string]string
}

// NewVisualizer initialize Visualizer
func NewVisualizer() *Visualizer {
	return &Visualizer{
		make(map[string][]*placemat.PodSpec),
		make(map[string]string),
	}
}

// Generate generates a Graphviz graph
func (v *Visualizer) Generate(cluster *ClusterSpec) (*gographviz.Escape, error) {
	res := gographviz.NewEscape()
	err := res.SetName("G")
	if err != nil {
		return nil, err
	}
	err = res.SetDir(false)
	if err != nil {
		return nil, err
	}

	err = v.connectPods(res, cluster)
	if err != nil {
		return nil, err
	}

	err = v.connectNodesAndPods(res, cluster)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (v *Visualizer) connectPods(graph *gographviz.Escape, cluster *ClusterSpec) error {
	for _, pod := range cluster.Pods {
		attrs := make(map[string]string)
		attrs[string(gographviz.Shape)] = "diamond"
		err := graph.AddNode("G", pod.Name, attrs)
		if err != nil {
			return err
		}
		for _, nic := range pod.Interfaces {
			v.podsLookUp[nic.Network] = append(v.podsLookUp[nic.Network], pod)
			v.addrLookUp[nic.Network] = nic.Addresses[0]
		}
	}
	for networkName, pods := range v.podsLookUp {
		if len(pods) != 2 {
			continue
		}
		attrs := make(map[string]string)
		attrs[string(gographviz.Label)] = v.addrLookUp[networkName]
		err := graph.AddEdge(pods[0].Name, pods[1].Name, false, attrs)
		if err != nil {
			return err
		}
	}
	return nil
}

func (v *Visualizer) connectNodesAndPods(graph *gographviz.Escape, cluster *ClusterSpec) error {
	err := v.prepareSubGraphs(graph, cluster)
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
			if pods, ok := v.podsLookUp[networkName]; ok {
				err := graph.AddEdge(pods[0].Name, node.Name, false, nil)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (v *Visualizer) prepareSubGraphs(graph *gographviz.Escape, cluster *ClusterSpec) error {
	for _, node := range cluster.Nodes {
		interfacesName := strings.Join(node.Interfaces, "_")
		attrs := make(map[string]string)

		var addresses []string
		for _, networkName := range node.Interfaces {
			addresses = append(addresses, v.addrLookUp[networkName])
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
