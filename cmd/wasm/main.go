package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall/js"
	"time"
)

type Node struct {
	Name     string
	Neighbor []string
	LowLink  int
	Index    int
}

type Graph struct {
	Nodes   []Node
	SCCs    [][]string
	Bridges [][]string
	Visited []string
	Index   int
	Stack   []string
	Reach   int
}

func NewNode(name string) *Node {
	return &Node{
		Name:     name,
		Neighbor: nil,
		LowLink:  -1,
		Index:    -1,
	}
}

func NewGraph() *Graph {
	return &Graph{
		Nodes:   nil,
		SCCs:    nil,
		Bridges: nil,
		Visited: nil,
		Index:   -1,
		Stack:   nil,
		Reach:   0,
	}
}

func (g *Graph) getNodeByName(name string) int {
	for i, node := range g.Nodes {
		if node.Name == name {
			return i
		}
	}
	return -1
}

func (g *Graph) getNodeByIndex(index int) int {
	for i, node := range g.Nodes {
		if node.Index == index {
			return i
		}
	}
	return -1
}

func (g *Graph) AddEdge(from string, to string) {
	fromNodeIndex := g.getNodeByName(from)

	if fromNodeIndex == -1 {
		fromNode := NewNode(from)
		fromNode.Neighbor = append(fromNode.Neighbor, to)
		g.Nodes = append(g.Nodes, *fromNode)
	} else {
		g.Nodes[fromNodeIndex].Neighbor = append(g.Nodes[fromNodeIndex].Neighbor, to)
	}

	toNodeIndex := g.getNodeByName(to)
	if toNodeIndex == -1 {
		toNode := NewNode(to)
		g.Nodes = append(g.Nodes, *toNode)
	}
}

func (g *Graph) isInStack(name string) bool {
	for _, nodename := range g.Stack {
		if nodename == name {
			return true
		}
	}
	return false
}

func (g *Graph) isVisited(name string) bool {
	for _, visname := range g.Visited {
		if visname == name {
			return true
		}
	}
	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (g *Graph) FindSCC() {
	g.Visited = make([]string, len(g.Nodes))
	g.Index = -1
	g.Stack = nil

	for _, node := range g.Nodes {
		if !g.isVisited(node.Name) {
			g.DFSForSCCs(node.Name)
		}
	}
}

func contains(arr []int, i int) bool {
	for _, n := range arr {
		if n == i {
			return true
		}
	}
	return false
}

func (g *Graph) DFSForSCCs(name string) {
	node := g.getNodeByName(name)
	if node == -1 {
		return
	}

	g.Index++
	g.Nodes[node].LowLink = g.Index
	g.Nodes[node].Index = g.Index
	g.Stack = append(g.Stack, g.Nodes[node].Name)
	g.Visited = append(g.Visited, g.Nodes[node].Name)
	// fmt.Println(g.Stack)
	for _, neighbor := range g.Nodes[node].Neighbor {
		if !g.isVisited(neighbor) {
			g.DFSForSCCs(neighbor)
			g.Nodes[node].LowLink = min(g.Nodes[node].LowLink, g.Nodes[g.getNodeByName(neighbor)].LowLink)
		} else if g.isInStack(neighbor) {
			g.Nodes[node].LowLink = min(g.Nodes[node].LowLink, g.Nodes[g.getNodeByName(neighbor)].LowLink)
		}
	}

	if g.Nodes[node].LowLink == g.Nodes[node].Index {
		var scc []string
		for {
			nodeFromStack := g.Stack[len(g.Stack)-1]
			g.Stack = g.Stack[:len(g.Stack)-1]

			scc = append(scc, nodeFromStack)
			if nodeFromStack == name {
				break
			}
		}
		g.SCCs = append(g.SCCs, scc)
	}
}

func removeElement(arr []string, e string) []string {
	indexToRemove := -1
	for i, element := range arr {
		if element == e {
			indexToRemove = i
			break
		}
	}
	if indexToRemove == -1 {
		return arr // Element not found
	}
	return append(arr[:indexToRemove], arr[indexToRemove+1:]...)
}

func (g *Graph) FindBridges() {
	g.Bridges = nil
	g.Visited = nil
	g.Index = -1

	for _, node := range g.Nodes {
		if !g.isVisited(node.Name) {
			g.DFSForBridges(node.Name, node.Name)
		}
	}
}

func (g *Graph) DFSForBridges(currNode, parentNode string) {
	nodeIndex := g.getNodeByName(currNode)
	if nodeIndex == -1 {
		return
	}

	g.Index++
	g.Nodes[nodeIndex].LowLink = g.Index
	g.Nodes[nodeIndex].Index = g.Index
	g.Visited = append(g.Visited, g.Nodes[nodeIndex].Name)
	g.Stack = append(g.Stack, g.Nodes[nodeIndex].Name)

	for _, neighbor := range g.Nodes[nodeIndex].Neighbor {
		if neighbor == parentNode {
			continue
		}

		if !g.isVisited(neighbor) {
			g.DFSForBridges(neighbor, currNode)
			g.Nodes[nodeIndex].LowLink = min(g.Nodes[nodeIndex].LowLink, g.Nodes[g.getNodeByName(neighbor)].LowLink)
			if g.Nodes[g.getNodeByName(neighbor)].LowLink > g.Nodes[nodeIndex].LowLink {
				g.Bridges = append(g.Bridges, []string{currNode, neighbor})
			}
		} else {
			g.Nodes[nodeIndex].LowLink = min(g.Nodes[nodeIndex].LowLink, g.Nodes[g.getNodeByName(neighbor)].LowLink)
		}
	}
}

func mainFunc(this js.Value, p []js.Value) any {
	if len(p) < 1 {
		return nil
	}
	inputText := p[0].String()
	graph := NewGraph()
	var outputString string

	inputLines := strings.Split(inputText, "\n")
	for _, line := range inputLines {
		if line == "done" {
			break
		}
		edge := strings.Split(line, " ")
		graph.AddEdge(edge[0], edge[1])
	}
	start := time.Now()
	graph.FindSCC()
	graph.FindBridges()

	end := time.Now()

	totalTime := end.Sub(start)

	outputString = outputString + "Total time: " + totalTime.String() + "\n\n"

	var sccString string
	var bridgeString string
	outputString = outputString + "Strongly Connected Components:\n"
	for _, scc := range graph.SCCs {
		sccString = strings.Join(scc, " ")
		outputString = outputString + "[" + sccString + "]" + "\n"
		// fmt.Println(scc)
	}

	outputString = outputString + "\n\nBridges:\n"
	for _, bridge := range graph.Bridges {
		bridgeString = strings.Join(bridge, " ")
		outputString = outputString + "[" + bridgeString + "]" + "\n"
	}

	return js.ValueOf(outputString)
}

func mainFunc2(this js.Value, p []js.Value) any {
	if len(p) < 1 {
		return nil
	}
	inputText := p[0].String()
	graph := NewGraph()
	var outputString string

	inputFile, _ := os.Open(inputText)
	scanner := bufio.NewScanner(inputFile)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "done" {
			break
		}
		edge := strings.Split(line, " ")
		graph.AddEdge(edge[0], edge[1])
	}
	start := time.Now()
	graph.FindSCC()
	graph.FindBridges()

	end := time.Now()

	totalTime := end.Sub(start)

	outputString = outputString + "Total time: " + totalTime.String() + "\n\n"

	var sccString string
	var bridgeString string
	outputString = outputString + "Strongly Connected Components:\n"
	for _, scc := range graph.SCCs {
		sccString = strings.Join(scc, " ")
		outputString = outputString + sccString + "\n"
		// fmt.Println(scc)
	}

	outputString = outputString + "\n\nBridges:\n"
	for _, bridge := range graph.Bridges {
		bridgeString = strings.Join(bridge, " ")
		outputString = outputString + bridgeString + "\n"
	}

	return js.ValueOf(outputString)
}

func main() {
	c := make(chan struct{}, 0)
	fmt.Println("Hello WebAssembly from Go!")

	js.Global().Set("mainFunc", js.FuncOf(mainFunc))
	// js.Global().Set("mainFunc2", js.FuncOf(mainFunc2))
	<-c
}
