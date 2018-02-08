package shadow_gonum.org/v1/gonum/graph

import "gonum.org/v1/gonum/graph"

var Pkg = make(map[string]interface{})
func init() {
    Pkg["Builder"] = GijitShadow_InterfaceConvertTo2_Builder
    Pkg["Copy"] = graph.Copy
    Pkg["CopyWeighted"] = graph.CopyWeighted
    Pkg["Directed"] = GijitShadow_InterfaceConvertTo2_Directed
    Pkg["DirectedBuilder"] = GijitShadow_InterfaceConvertTo2_DirectedBuilder
    Pkg["DirectedMultigraph"] = GijitShadow_InterfaceConvertTo2_DirectedMultigraph
    Pkg["DirectedMultigraphBuilder"] = GijitShadow_InterfaceConvertTo2_DirectedMultigraphBuilder
    Pkg["DirectedWeightedBuilder"] = GijitShadow_InterfaceConvertTo2_DirectedWeightedBuilder
    Pkg["DirectedWeightedMultigraphBuilder"] = GijitShadow_InterfaceConvertTo2_DirectedWeightedMultigraphBuilder
    Pkg["Edge"] = GijitShadow_InterfaceConvertTo2_Edge
    Pkg["EdgeAdder"] = GijitShadow_InterfaceConvertTo2_EdgeAdder
    Pkg["EdgeRemover"] = GijitShadow_InterfaceConvertTo2_EdgeRemover
    Pkg["Graph"] = GijitShadow_InterfaceConvertTo2_Graph
    Pkg["Line"] = GijitShadow_InterfaceConvertTo2_Line
    Pkg["LineAdder"] = GijitShadow_InterfaceConvertTo2_LineAdder
    Pkg["LineRemover"] = GijitShadow_InterfaceConvertTo2_LineRemover
    Pkg["Multigraph"] = GijitShadow_InterfaceConvertTo2_Multigraph
    Pkg["MultigraphBuilder"] = GijitShadow_InterfaceConvertTo2_MultigraphBuilder
    Pkg["Node"] = GijitShadow_InterfaceConvertTo2_Node
    Pkg["NodeAdder"] = GijitShadow_InterfaceConvertTo2_NodeAdder
    Pkg["NodeRemover"] = GijitShadow_InterfaceConvertTo2_NodeRemover
    Pkg["Undirected"] = GijitShadow_InterfaceConvertTo2_Undirected
    Pkg["UndirectedBuilder"] = GijitShadow_InterfaceConvertTo2_UndirectedBuilder
    Pkg["UndirectedMultigraph"] = GijitShadow_InterfaceConvertTo2_UndirectedMultigraph
    Pkg["UndirectedMultigraphBuilder"] = GijitShadow_InterfaceConvertTo2_UndirectedMultigraphBuilder
    Pkg["UndirectedWeightedBuilder"] = GijitShadow_InterfaceConvertTo2_UndirectedWeightedBuilder
    Pkg["UndirectedWeightedMultigraphBuilder"] = GijitShadow_InterfaceConvertTo2_UndirectedWeightedMultigraphBuilder
    Pkg["Weighted"] = GijitShadow_InterfaceConvertTo2_Weighted
    Pkg["WeightedBuilder"] = GijitShadow_InterfaceConvertTo2_WeightedBuilder
    Pkg["WeightedDirected"] = GijitShadow_InterfaceConvertTo2_WeightedDirected
    Pkg["WeightedDirectedMultigraph"] = GijitShadow_InterfaceConvertTo2_WeightedDirectedMultigraph
    Pkg["WeightedEdge"] = GijitShadow_InterfaceConvertTo2_WeightedEdge
    Pkg["WeightedEdgeAdder"] = GijitShadow_InterfaceConvertTo2_WeightedEdgeAdder
    Pkg["WeightedLine"] = GijitShadow_InterfaceConvertTo2_WeightedLine
    Pkg["WeightedLineAdder"] = GijitShadow_InterfaceConvertTo2_WeightedLineAdder
    Pkg["WeightedMultigraph"] = GijitShadow_InterfaceConvertTo2_WeightedMultigraph
    Pkg["WeightedMultigraphBuilder"] = GijitShadow_InterfaceConvertTo2_WeightedMultigraphBuilder
    Pkg["WeightedUndirected"] = GijitShadow_InterfaceConvertTo2_WeightedUndirected
    Pkg["WeightedUndirectedMultigraph"] = GijitShadow_InterfaceConvertTo2_WeightedUndirectedMultigraph

}
func GijitShadow_InterfaceConvertTo2_Builder(x interface{}) (y graph.Builder, b bool) {
	y, b = x.(graph.Builder)
	return
}

func GijitShadow_InterfaceConvertTo1_Builder(x interface{}) graph.Builder {
	return x.(graph.Builder)
}


func GijitShadow_InterfaceConvertTo2_Directed(x interface{}) (y graph.Directed, b bool) {
	y, b = x.(graph.Directed)
	return
}

func GijitShadow_InterfaceConvertTo1_Directed(x interface{}) graph.Directed {
	return x.(graph.Directed)
}


func GijitShadow_InterfaceConvertTo2_DirectedBuilder(x interface{}) (y graph.DirectedBuilder, b bool) {
	y, b = x.(graph.DirectedBuilder)
	return
}

func GijitShadow_InterfaceConvertTo1_DirectedBuilder(x interface{}) graph.DirectedBuilder {
	return x.(graph.DirectedBuilder)
}


func GijitShadow_InterfaceConvertTo2_DirectedMultigraph(x interface{}) (y graph.DirectedMultigraph, b bool) {
	y, b = x.(graph.DirectedMultigraph)
	return
}

func GijitShadow_InterfaceConvertTo1_DirectedMultigraph(x interface{}) graph.DirectedMultigraph {
	return x.(graph.DirectedMultigraph)
}


func GijitShadow_InterfaceConvertTo2_DirectedMultigraphBuilder(x interface{}) (y graph.DirectedMultigraphBuilder, b bool) {
	y, b = x.(graph.DirectedMultigraphBuilder)
	return
}

func GijitShadow_InterfaceConvertTo1_DirectedMultigraphBuilder(x interface{}) graph.DirectedMultigraphBuilder {
	return x.(graph.DirectedMultigraphBuilder)
}


func GijitShadow_InterfaceConvertTo2_DirectedWeightedBuilder(x interface{}) (y graph.DirectedWeightedBuilder, b bool) {
	y, b = x.(graph.DirectedWeightedBuilder)
	return
}

func GijitShadow_InterfaceConvertTo1_DirectedWeightedBuilder(x interface{}) graph.DirectedWeightedBuilder {
	return x.(graph.DirectedWeightedBuilder)
}


func GijitShadow_InterfaceConvertTo2_DirectedWeightedMultigraphBuilder(x interface{}) (y graph.DirectedWeightedMultigraphBuilder, b bool) {
	y, b = x.(graph.DirectedWeightedMultigraphBuilder)
	return
}

func GijitShadow_InterfaceConvertTo1_DirectedWeightedMultigraphBuilder(x interface{}) graph.DirectedWeightedMultigraphBuilder {
	return x.(graph.DirectedWeightedMultigraphBuilder)
}


func GijitShadow_InterfaceConvertTo2_Edge(x interface{}) (y graph.Edge, b bool) {
	y, b = x.(graph.Edge)
	return
}

func GijitShadow_InterfaceConvertTo1_Edge(x interface{}) graph.Edge {
	return x.(graph.Edge)
}


func GijitShadow_InterfaceConvertTo2_EdgeAdder(x interface{}) (y graph.EdgeAdder, b bool) {
	y, b = x.(graph.EdgeAdder)
	return
}

func GijitShadow_InterfaceConvertTo1_EdgeAdder(x interface{}) graph.EdgeAdder {
	return x.(graph.EdgeAdder)
}


func GijitShadow_InterfaceConvertTo2_EdgeRemover(x interface{}) (y graph.EdgeRemover, b bool) {
	y, b = x.(graph.EdgeRemover)
	return
}

func GijitShadow_InterfaceConvertTo1_EdgeRemover(x interface{}) graph.EdgeRemover {
	return x.(graph.EdgeRemover)
}


func GijitShadow_InterfaceConvertTo2_Graph(x interface{}) (y graph.Graph, b bool) {
	y, b = x.(graph.Graph)
	return
}

func GijitShadow_InterfaceConvertTo1_Graph(x interface{}) graph.Graph {
	return x.(graph.Graph)
}


func GijitShadow_InterfaceConvertTo2_Line(x interface{}) (y graph.Line, b bool) {
	y, b = x.(graph.Line)
	return
}

func GijitShadow_InterfaceConvertTo1_Line(x interface{}) graph.Line {
	return x.(graph.Line)
}


func GijitShadow_InterfaceConvertTo2_LineAdder(x interface{}) (y graph.LineAdder, b bool) {
	y, b = x.(graph.LineAdder)
	return
}

func GijitShadow_InterfaceConvertTo1_LineAdder(x interface{}) graph.LineAdder {
	return x.(graph.LineAdder)
}


func GijitShadow_InterfaceConvertTo2_LineRemover(x interface{}) (y graph.LineRemover, b bool) {
	y, b = x.(graph.LineRemover)
	return
}

func GijitShadow_InterfaceConvertTo1_LineRemover(x interface{}) graph.LineRemover {
	return x.(graph.LineRemover)
}


func GijitShadow_InterfaceConvertTo2_Multigraph(x interface{}) (y graph.Multigraph, b bool) {
	y, b = x.(graph.Multigraph)
	return
}

func GijitShadow_InterfaceConvertTo1_Multigraph(x interface{}) graph.Multigraph {
	return x.(graph.Multigraph)
}


func GijitShadow_InterfaceConvertTo2_MultigraphBuilder(x interface{}) (y graph.MultigraphBuilder, b bool) {
	y, b = x.(graph.MultigraphBuilder)
	return
}

func GijitShadow_InterfaceConvertTo1_MultigraphBuilder(x interface{}) graph.MultigraphBuilder {
	return x.(graph.MultigraphBuilder)
}


func GijitShadow_InterfaceConvertTo2_Node(x interface{}) (y graph.Node, b bool) {
	y, b = x.(graph.Node)
	return
}

func GijitShadow_InterfaceConvertTo1_Node(x interface{}) graph.Node {
	return x.(graph.Node)
}


func GijitShadow_InterfaceConvertTo2_NodeAdder(x interface{}) (y graph.NodeAdder, b bool) {
	y, b = x.(graph.NodeAdder)
	return
}

func GijitShadow_InterfaceConvertTo1_NodeAdder(x interface{}) graph.NodeAdder {
	return x.(graph.NodeAdder)
}


func GijitShadow_InterfaceConvertTo2_NodeRemover(x interface{}) (y graph.NodeRemover, b bool) {
	y, b = x.(graph.NodeRemover)
	return
}

func GijitShadow_InterfaceConvertTo1_NodeRemover(x interface{}) graph.NodeRemover {
	return x.(graph.NodeRemover)
}


func GijitShadow_NewStruct_Undirect() *graph.Undirect {
	return &graph.Undirect{}
}


func GijitShadow_NewStruct_UndirectWeighted() *graph.UndirectWeighted {
	return &graph.UndirectWeighted{}
}


func GijitShadow_InterfaceConvertTo2_Undirected(x interface{}) (y graph.Undirected, b bool) {
	y, b = x.(graph.Undirected)
	return
}

func GijitShadow_InterfaceConvertTo1_Undirected(x interface{}) graph.Undirected {
	return x.(graph.Undirected)
}


func GijitShadow_InterfaceConvertTo2_UndirectedBuilder(x interface{}) (y graph.UndirectedBuilder, b bool) {
	y, b = x.(graph.UndirectedBuilder)
	return
}

func GijitShadow_InterfaceConvertTo1_UndirectedBuilder(x interface{}) graph.UndirectedBuilder {
	return x.(graph.UndirectedBuilder)
}


func GijitShadow_InterfaceConvertTo2_UndirectedMultigraph(x interface{}) (y graph.UndirectedMultigraph, b bool) {
	y, b = x.(graph.UndirectedMultigraph)
	return
}

func GijitShadow_InterfaceConvertTo1_UndirectedMultigraph(x interface{}) graph.UndirectedMultigraph {
	return x.(graph.UndirectedMultigraph)
}


func GijitShadow_InterfaceConvertTo2_UndirectedMultigraphBuilder(x interface{}) (y graph.UndirectedMultigraphBuilder, b bool) {
	y, b = x.(graph.UndirectedMultigraphBuilder)
	return
}

func GijitShadow_InterfaceConvertTo1_UndirectedMultigraphBuilder(x interface{}) graph.UndirectedMultigraphBuilder {
	return x.(graph.UndirectedMultigraphBuilder)
}


func GijitShadow_InterfaceConvertTo2_UndirectedWeightedBuilder(x interface{}) (y graph.UndirectedWeightedBuilder, b bool) {
	y, b = x.(graph.UndirectedWeightedBuilder)
	return
}

func GijitShadow_InterfaceConvertTo1_UndirectedWeightedBuilder(x interface{}) graph.UndirectedWeightedBuilder {
	return x.(graph.UndirectedWeightedBuilder)
}


func GijitShadow_InterfaceConvertTo2_UndirectedWeightedMultigraphBuilder(x interface{}) (y graph.UndirectedWeightedMultigraphBuilder, b bool) {
	y, b = x.(graph.UndirectedWeightedMultigraphBuilder)
	return
}

func GijitShadow_InterfaceConvertTo1_UndirectedWeightedMultigraphBuilder(x interface{}) graph.UndirectedWeightedMultigraphBuilder {
	return x.(graph.UndirectedWeightedMultigraphBuilder)
}


func GijitShadow_InterfaceConvertTo2_Weighted(x interface{}) (y graph.Weighted, b bool) {
	y, b = x.(graph.Weighted)
	return
}

func GijitShadow_InterfaceConvertTo1_Weighted(x interface{}) graph.Weighted {
	return x.(graph.Weighted)
}


func GijitShadow_InterfaceConvertTo2_WeightedBuilder(x interface{}) (y graph.WeightedBuilder, b bool) {
	y, b = x.(graph.WeightedBuilder)
	return
}

func GijitShadow_InterfaceConvertTo1_WeightedBuilder(x interface{}) graph.WeightedBuilder {
	return x.(graph.WeightedBuilder)
}


func GijitShadow_InterfaceConvertTo2_WeightedDirected(x interface{}) (y graph.WeightedDirected, b bool) {
	y, b = x.(graph.WeightedDirected)
	return
}

func GijitShadow_InterfaceConvertTo1_WeightedDirected(x interface{}) graph.WeightedDirected {
	return x.(graph.WeightedDirected)
}


func GijitShadow_InterfaceConvertTo2_WeightedDirectedMultigraph(x interface{}) (y graph.WeightedDirectedMultigraph, b bool) {
	y, b = x.(graph.WeightedDirectedMultigraph)
	return
}

func GijitShadow_InterfaceConvertTo1_WeightedDirectedMultigraph(x interface{}) graph.WeightedDirectedMultigraph {
	return x.(graph.WeightedDirectedMultigraph)
}


func GijitShadow_InterfaceConvertTo2_WeightedEdge(x interface{}) (y graph.WeightedEdge, b bool) {
	y, b = x.(graph.WeightedEdge)
	return
}

func GijitShadow_InterfaceConvertTo1_WeightedEdge(x interface{}) graph.WeightedEdge {
	return x.(graph.WeightedEdge)
}


func GijitShadow_InterfaceConvertTo2_WeightedEdgeAdder(x interface{}) (y graph.WeightedEdgeAdder, b bool) {
	y, b = x.(graph.WeightedEdgeAdder)
	return
}

func GijitShadow_InterfaceConvertTo1_WeightedEdgeAdder(x interface{}) graph.WeightedEdgeAdder {
	return x.(graph.WeightedEdgeAdder)
}


func GijitShadow_NewStruct_WeightedEdgePair() *graph.WeightedEdgePair {
	return &graph.WeightedEdgePair{}
}


func GijitShadow_InterfaceConvertTo2_WeightedLine(x interface{}) (y graph.WeightedLine, b bool) {
	y, b = x.(graph.WeightedLine)
	return
}

func GijitShadow_InterfaceConvertTo1_WeightedLine(x interface{}) graph.WeightedLine {
	return x.(graph.WeightedLine)
}


func GijitShadow_InterfaceConvertTo2_WeightedLineAdder(x interface{}) (y graph.WeightedLineAdder, b bool) {
	y, b = x.(graph.WeightedLineAdder)
	return
}

func GijitShadow_InterfaceConvertTo1_WeightedLineAdder(x interface{}) graph.WeightedLineAdder {
	return x.(graph.WeightedLineAdder)
}


func GijitShadow_InterfaceConvertTo2_WeightedMultigraph(x interface{}) (y graph.WeightedMultigraph, b bool) {
	y, b = x.(graph.WeightedMultigraph)
	return
}

func GijitShadow_InterfaceConvertTo1_WeightedMultigraph(x interface{}) graph.WeightedMultigraph {
	return x.(graph.WeightedMultigraph)
}


func GijitShadow_InterfaceConvertTo2_WeightedMultigraphBuilder(x interface{}) (y graph.WeightedMultigraphBuilder, b bool) {
	y, b = x.(graph.WeightedMultigraphBuilder)
	return
}

func GijitShadow_InterfaceConvertTo1_WeightedMultigraphBuilder(x interface{}) graph.WeightedMultigraphBuilder {
	return x.(graph.WeightedMultigraphBuilder)
}


func GijitShadow_InterfaceConvertTo2_WeightedUndirected(x interface{}) (y graph.WeightedUndirected, b bool) {
	y, b = x.(graph.WeightedUndirected)
	return
}

func GijitShadow_InterfaceConvertTo1_WeightedUndirected(x interface{}) graph.WeightedUndirected {
	return x.(graph.WeightedUndirected)
}


func GijitShadow_InterfaceConvertTo2_WeightedUndirectedMultigraph(x interface{}) (y graph.WeightedUndirectedMultigraph, b bool) {
	y, b = x.(graph.WeightedUndirectedMultigraph)
	return
}

func GijitShadow_InterfaceConvertTo1_WeightedUndirectedMultigraph(x interface{}) graph.WeightedUndirectedMultigraph {
	return x.(graph.WeightedUndirectedMultigraph)
}

