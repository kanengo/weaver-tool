package callgraph

import (
	"fmt"
	"strings"

	"github.com/ServiceWeaver/weaver/runtime/bin"
	"github.com/ServiceWeaver/weaver/runtime/graph"
	"github.com/ServiceWeaver/weaver/runtime/logging"
)

// Mermaid returns a Mermaid diagram, https://mermaid.js.org/, of the component
// call graph embedded in the provided Service Weaver binary.
func Mermaid(binary string) (string, error) {
	components, g, err := bin.ReadComponentGraph(binary)
	if err != nil {
		return "", err
	}
	return mermaid(components, g), nil
}

// mermaid returns a Mermaid diagram of the given component graph.
func mermaid(components []string, g graph.Graph) string {
	// See https://mermaid.js.org/syntax/flowchart.html for details.
	var b strings.Builder
	fmt.Fprintln(&b, `%%{init: {"flowchart": {"defaultRenderer": "elk"}} }%%`)
	fmt.Fprintln(&b, "graph TD")

	// Nodes.
	fmt.Fprintln(&b, `    %% Nodes.`)
	g.PerNode(func(n graph.Node) {
		fmt.Fprintf(&b, "    %d(%s)\n", n, logging.ShortenComponent(components[n]))
	})
	fmt.Fprintln(&b, "")

	// Edges.
	fmt.Fprintln(&b, `    %% Edges.`)
	graph.PerEdge(g, func(e graph.Edge) {
		fmt.Fprintf(&b, "    %d --> %d\n", e.Src, e.Dst)
	})
	return b.String()
}
