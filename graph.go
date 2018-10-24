package viz

import (
	"bytes"
	"github.com/bradleyjkemp/memviz"
)

// GenerateDots generate GraphViz string
func GenerateDots(cluster interface{}) (string, error) {
	buf := &bytes.Buffer{}
	memviz.Map(buf, cluster)
	return buf.String(), nil
}
