package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"strings"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
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
	Edges   [][]string
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
		g.Edges = append(g.Edges, []string{from, to})
	} else {
		g.Nodes[fromNodeIndex].Neighbor = append(g.Nodes[fromNodeIndex].Neighbor, to)
		g.Edges = append(g.Edges, []string{from, to})
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

func drawGraphWithColors(graph *Graph) {
	canvasWidth, canvasHeight := 1920, 1080
	scaleX, scaleY := 50, 50

	canvas := image.NewRGBA(image.Rect(0, 0, canvasWidth, canvasHeight))
	draw.Draw(canvas, canvas.Bounds(), image.NewUniform(color.White), image.Point{}, draw.Src)

	colors := []color.Color{color.RGBA{255, 0, 0, 255}, color.RGBA{0, 255, 0, 255}, color.RGBA{0, 0, 255, 255}, color.RGBA{255, 255, 0, 255}, color.RGBA{255, 0, 255, 255}, color.RGBA{0, 255, 255, 255}}
	sccColors := make(map[string]color.Color)

	for i, scc := range graph.SCCs {
		sccColor := colors[i%len(colors)]
		for _, node := range scc {
			sccColors[node] = sccColor
		}
	}
	for _, edge := range graph.Edges {
		fromNodeIndex := graph.getNodeByName(edge[0])
		toNodeIndex := graph.getNodeByName(edge[1])
		if toNodeIndex < fromNodeIndex {
			temp := fromNodeIndex
			fromNodeIndex = toNodeIndex
			toNodeIndex = temp

		}
		if fromNodeIndex != -1 && toNodeIndex != -1 {
			fromX, fromY := graph.Nodes[fromNodeIndex].Index*scaleX+50, graph.Nodes[fromNodeIndex].LowLink*scaleY+50
			toX, toY := graph.Nodes[toNodeIndex].Index*scaleX+50, graph.Nodes[toNodeIndex].LowLink*scaleY+50
			drawEdge(canvas, fromX, fromY, toX, toY, color.Black)
		}
	}

	for _, node := range graph.Nodes {
		x, y := node.Index*scaleX+50, node.LowLink*scaleY+50
		drawNode(canvas, x, y, node.Name, sccColors[node.Name])
	}

	for _, bridge := range graph.Bridges {
		fromNodeIndex := graph.getNodeByName(bridge[0])
		toNodeIndex := graph.getNodeByName(bridge[1])
		if fromNodeIndex != -1 && toNodeIndex != -1 {
			fromX, fromY := graph.Nodes[fromNodeIndex].Index*scaleX+50, graph.Nodes[fromNodeIndex].LowLink*scaleY+50
			toX, toY := graph.Nodes[toNodeIndex].Index*scaleX+50, graph.Nodes[toNodeIndex].LowLink*scaleY+50
			drawEdge(canvas, fromX, fromY, toX, toY, color.RGBA{255, 0, 0, 255})
		}
	}

	file, err := os.Create("graph_with_colors.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	png.Encode(file, canvas)
}

func drawNode(canvas draw.Image, x, y int, label string, col color.Color) {
	nodeRadius := 20
	// Draw the node circle with a label
	for i := x - nodeRadius; i <= x+nodeRadius; i++ {
		for j := y - nodeRadius; j <= y+nodeRadius; j++ {
			if (i-x)*(i-x)+(j-y)*(j-y) <= nodeRadius*nodeRadius {
				canvas.Set(i, j, col)
			}
		}
	}
	drawLabel(canvas, x, y, label, color.Black)
}

func drawLabel(canvas draw.Image, x, y int, label string, col color.Color) {
	fontSize := 10

	d := &font.Drawer{
		Dst:  canvas,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  fixed.P(x-fontSize/2*len(label)+3, y-fontSize/2+10),
	}
	d.DrawString(label)
}

func drawEdge(canvas draw.Image, x1, y1, x2, y2 int, col color.Color) {
	drawLine(canvas, x1, y1, x2, y2, col)
}

func drawLine(canvas draw.Image, x1, y1, x2, y2 int, col color.Color) {
	dx, dy := x2-x1, y2-y1
	absDx, absDy := abs(dx), abs(dy)
	incrX, incrY := sign(dx), sign(dy)
	if absDy < absDx {
		if x2 < x1 {
			x1, y1, x2, y2 = x2, y2, x1, y1
		}
		acc := absDx / 2
		for x1 < x2 {
			canvas.Set(x1, y1, col)
			x1 += incrX
			acc += absDy
			if acc >= absDx {
				acc -= absDx
				y1 += incrY
			}
		}
	} else {
		if y2 < y1 {
			x1, y1, x2, y2 = x2, y2, x1, y1
		}
		acc := absDy / 2
		for y1 < y2 {
			canvas.Set(x1, y1, col)
			y1 += incrY
			acc += absDx
			if acc >= absDy {
				acc -= absDy
				x1 += incrX
			}
		}
	}
}

func fixedPoint(x, y, size int) image.Point {
	// Size is the font size.
	return image.Point{
		X: x - size/2*len("X"),
		Y: y - size/2,
	}
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func sign(a int) int {
	if a < 0 {
		return -1
	} else if a > 0 {
		return 1
	}
	return 0
}

func pow(a int, b int) int {
	if b == 0 {
		return 1
	}
	return a * pow(a, b-1)
}

func main() {
	graph := NewGraph()

	// Uncomment the following line to read input from a file
	inputFile, _ := os.Open("test/input.txt")
	scanner := bufio.NewScanner(inputFile)

	// Comment the following lines if reading from a file
	// fmt.Println("Enter the graph edges (type 'done' to finish input):")
	// scanner := bufio.NewScanner(os.Stdin)

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

	fmt.Println("Total time:", totalTime)
	// for _, node := range graph.Nodes {
	// 	fmt.Println(node)
	// }
	fmt.Println("Strongly Connected Components:")
	for _, scc := range graph.SCCs {
		fmt.Println(scc)
	}
	fmt.Println("\nBridges:")
	for _, bridge := range graph.Bridges {
		fmt.Println(bridge)
	}

	drawGraphWithColors(graph)
}
