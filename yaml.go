package viz

import (
	"bufio"
	"errors"
	"io"

	"github.com/cybozu-go/placemat"
	k8sYaml "github.com/kubernetes/apimachinery/pkg/util/yaml"
	"gopkg.in/yaml.v2"
)

type baseConfig struct {
	Kind string `yaml:"kind"`
}

// ClusterSpec is a set of resources in a virtual data center.
type ClusterSpec struct {
	Networks    []*placemat.NetworkSpec
	Images      []*placemat.ImageSpec
	DataFolders []*placemat.DataFolderSpec
	Nodes       []*placemat.NodeSpec
	Pods        []*placemat.PodSpec
}

// ReadYAML reads a Placemat YAML file and constructs Cluster
func ReadYAML(r *bufio.Reader) (*ClusterSpec, error) {
	var cluster ClusterSpec

	y := k8sYaml.NewYAMLReader(r)
	for {
		data, err := y.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		var c baseConfig
		err = yaml.Unmarshal(data, &c)
		if err != nil {
			return nil, err
		}

		switch c.Kind {
		case "Network":
			spec := new(placemat.NetworkSpec)
			err = yaml.Unmarshal(data, spec)
			if err != nil {
				return nil, err
			}

			cluster.Networks = append(cluster.Networks, spec)
		case "Image":
			spec := new(placemat.ImageSpec)
			err = yaml.Unmarshal(data, spec)
			if err != nil {
				return nil, err
			}
			cluster.Images = append(cluster.Images, spec)
		case "DataFolder":
			spec := new(placemat.DataFolderSpec)
			err = yaml.Unmarshal(data, spec)
			if err != nil {
				return nil, err
			}
			cluster.DataFolders = append(cluster.DataFolders, spec)
		case "Node":
			spec := new(placemat.NodeSpec)
			err = yaml.Unmarshal(data, spec)
			if err != nil {
				return nil, err
			}
			cluster.Nodes = append(cluster.Nodes, spec)
		case "Pod":
			spec := new(placemat.PodSpec)
			err = yaml.Unmarshal(data, spec)
			if err != nil {
				return nil, err
			}
			cluster.Pods = append(cluster.Pods, spec)
		default:
			return nil, errors.New("unknown resource: " + c.Kind)
		}
	}
	return &cluster, nil
}
