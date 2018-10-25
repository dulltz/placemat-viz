package viz

import (
	"testing"

	"github.com/cybozu-go/placemat"
)

func TestGenerate(t *testing.T) {
	graph := NewVisualizer()
	g, err := graph.Generate(&ClusterSpec{})
	if err != nil {
		t.Fatal(err)
	}
	if len(g.Edges.Edges) != 0 {
		t.Error(err)
	}

	net1 := "net1"
	net2 := "swith-to-nodes"
	pods := []*placemat.PodSpec{
		{
			Name: "switch1",
			Interfaces: []placemat.PodInterfaceSpec{
				{
					Network:   net1,
					Addresses: []string{"1.1.1.1/32"},
				},
			},
		},
		{
			Name: "switch2",
			Interfaces: []placemat.PodInterfaceSpec{
				{
					Network:   net1,
					Addresses: []string{"1.1.1.1/32"},
				},
				{
					Network:   net2,
					Addresses: []string{"10.0.10.0/26"},
				},
			},
		},
	}

	g, err = NewVisualizer().Generate(&ClusterSpec{Pods: pods})
	if err != nil {
		t.Fatal(err)
	}
	if len(g.Edges.Edges) != 1 {
		t.Error("should have one edge, actual:", len(g.Edges.Edges))
	}

	g, err = NewVisualizer().Generate(&ClusterSpec{Pods: pods[:1]})
	if err != nil {
		t.Fatal(err)
	}
	if len(g.Edges.Edges) != 0 {
		t.Error("should have no edge, actual:", len(g.Edges.Edges))
	}

	nodes := []*placemat.NodeSpec{
		{Name: "node1", Interfaces: []string{net2}},
		{Name: "node2", Interfaces: []string{net2}},
	}
	g, err = NewVisualizer().Generate(&ClusterSpec{Pods: pods, Nodes: nodes})
	if err != nil {
		t.Fatal(err)
	}
	if len(g.Edges.Edges) != 3 {
		t.Error("should have 3 edge, actual:", len(g.Edges.Edges))
	}
	if len(g.Nodes.Nodes) != 4 {
		t.Error("should have 4 nodes, actual:", len(g.Nodes.Nodes))
	}
	if len(g.SubGraphs.SubGraphs) != 1 {
		t.Error("should have 1 sub graphs, actual:", len(g.SubGraphs.SubGraphs))
	}

	g, err = NewVisualizer().Generate(&ClusterSpec{Nodes: nodes})
	if err != nil {
		t.Fatal(err)
	}
	if len(g.Edges.Edges) != 0 {
		t.Error("should have 0 edge, actual:", len(g.Edges.Edges))
	}
	if len(g.Nodes.Nodes) != 2 {
		t.Error("should have 2 nodes, actual:", len(g.Nodes.Nodes))
	}
}
