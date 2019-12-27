package main

import (
	"fmt"
	"log"
	"math/rand"

	. "gorgonia.org/gorgonia"
	"gorgonia.org/tensor"
)

var err error
var dt tensor.Dtype = tensor.Float64

type nn struct {
	g      *ExprGraph
	w0, w1 *Node

	pred    *Node
	predVal Value
}

func newNN(g *ExprGraph) *nn {
	// Create node for w/weight
	w0 := NewMatrix(g, dt, WithShape(2, 2), WithName("w0"), WithInit(GlorotN(1.0)))
	w1 := NewMatrix(g, dt, WithShape(2, 1), WithName("w1"), WithInit(GlorotN(1.0)))
	return &nn{
		g:  g,
		w0: w0,
		w1: w1}
}

func (m *nn) learnables() Nodes {
	return Nodes{m.w0, m.w1}
}

func (m *nn) fwd(x *Node) (err error) {
	var l0, l1, l2 *Node
	var l0dot, l1dot *Node


	// Camada de input
	l0 = x

	// Multiplicação pelos pesos e sigmoid
	l0dot = Must(Mul(l0, m.w0))

	// Input para a hidden layer
	l1 = Must(Sigmoid(l0dot))

	// Multiplicação pelos pesos:
	l1dot = Must(Mul(l1, m.w1))

	// Camada de saída:
	l2 = Must(Sigmoid(l1dot))

	m.pred = l2
	Read(m.pred, &m.predVal)
	return nil

}


func main() {

	rand.Seed(31337)

	// Create graph and network
	g := NewGraph()
	m := newNN(g)

	// Set input x to network
	xB := []float64{1,0,0,1,1,1,0,0}
	xT := tensor.New(tensor.WithBacking(xB), tensor.WithShape(4, 2))
	x := NewMatrix(g,
		tensor.Float64,
		WithName("X"),
		WithShape(4, 2),
		WithValue(xT),
	)

	// Define validation data set
	yB := []float64{1, 1, 0, 0}
	yT := tensor.New(tensor.WithBacking(yB), tensor.WithShape(4, 1))
	y := NewMatrix(g,
		tensor.Float64,
		WithName("y"),
		WithShape(4, 1),
		WithValue(yT),
	)

	// Run forward pass
	if err := m.fwd(x); err != nil {
		log.Fatalf("%+v", err)
	}

	// Calculate Cost w/MSE
	losses := Must(Sub(y, m.pred))
	square := Must(Square(losses))
	cost := Must(Mean(square))

	// Do Gradient updates
	if _, err = Grad(cost, m.learnables()...); err != nil {
		log.Fatal(err)
	}

	// Instantiate VM and Solver
	vm := NewTapeMachine(g, BindDualValues(m.learnables()...))
	solver := NewVanillaSolver(WithLearnRate(0.1))

	for i := 0; i < 10000; i++ {
		vm.Reset()
		if err = vm.RunAll(); err != nil {
			log.Fatalf("Failed at inter  %d: %v", i, err)
		}
		solver.Step(NodesToValueGrads(m.learnables()))
		vm.Reset()
	}
	fmt.Println("\n\nOutput after Training: \n", m.predVal)
}