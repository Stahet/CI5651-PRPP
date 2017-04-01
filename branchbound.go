package main

import (
	"fmt"
	"math"
	"sort"
)

func estaEnSolucionParcial(e *Edge, solParcial []*Edge) bool {
	ocurr := 0
	for _, edge := range solParcial {
		if edge.equals(e) {
			ocurr = ocurr + 1
		}
		if ocurr >= 2 {
			return true
		}
	}
	return false
}

func (g *Graph) obtenerListaSucesores(v int) []*Edge {
	lista := Edges{}
	newlist := make([]*Edge, 0, len(g.edges[v])*2)
	for _, edge := range g.edges[v] {
		lista = append(lista, Edge{edge.start, edge.end, edge.cost, edge.benefit, 0})
		lista = append(lista, Edge{edge.start, edge.end, edge.cost, 0, 0})
	}
	sort.Sort(sort.Reverse(lista))
	for _, edge := range lista {
		newlist = append(newlist, g.edges[edge.start][edge.end])
	}
	return newlist
}

func (g *Graph) cumpleAcotamiento(e *Edge, solParcial []*Edge, mejorSol []*Edge, beneficioDisponible int) bool {
	beneficioE := g.NetBenefit(e.start, e.end)
	//fmt.Println(e, "Beneficio E: ", beneficioE)
	beneficioSolParcial := getPathBenefit(solParcial) + beneficioE
	//fmt.Println(mejorSol)
	//fmt.Println("SolParcialBenefit", beneficioSolParcial, "beneficioMejor:", getPathBenefit(mejorSol))
	maxBeneficio := beneficioDisponible - int(math.Max(0, float64(beneficioE))) + beneficioSolParcial
	if maxBeneficio <= getPathBenefit(mejorSol) {
		return false
	}
	return true
}

// Global variables defined in Main
// var mejorSol []*Edge
// var beneficioDisponible int
// var solParcial []*Edge
//
//
func (g *Graph) branchAndBound() {
	v := 1
	if len(solParcial) > 0 {
		v = solParcial[len(solParcial)-1].end
	}
	if v == 1 {
		if getPathBenefit(solParcial) > getPathBenefit(mejorSol) {
			mejorSol = make([]*Edge, 0)
			for _, elem := range solParcial {
				mejorSol = append(mejorSol, elem)
			}
			fmt.Println("nueva mejor sol: ", mejorSol)
		}
	}
	//fmt.Println(mejorSol)
	sucesores := g.obtenerListaSucesores(v)
	estaSolParcial, cumpleAco, cicloNeg := true, true, true
	fmt.Print("solParcial: ", solParcial, " benef: ", getPathBenefit(solParcial))
	//fmt.Println(" sucesores: ", sucesores)
	fmt.Println(" Benef mejorSol: ", getPathBenefit(mejorSol))
	for _, edge := range sucesores {
		estaSolParcial = estaEnSolucionParcial(edge, solParcial)
		cumpleAco = g.cumpleAcotamiento(edge, solParcial, mejorSol, beneficioDisponible)
		cicloNeg = checkNegativeCycle(edge, solParcial)
		fmt.Println("nodo:", v, "lado: ", edge, "| BenefLado:", g.NetBenefit(edge.start, edge.end), "| estaSolPar:", estaSolParcial, "| cumpleAc:", cumpleAco, "| NegCycle:", cicloNeg, "| benefDisponible:", beneficioDisponible, "| cond:", !estaSolParcial && cumpleAco && !cicloNeg)

		if !estaSolParcial && cumpleAco && !cicloNeg {
			fmt.Println()
			solParcial = append(solParcial, edge)
			beneficioDisponible = beneficioDisponible - int(math.Max(0, float64(g.NetBenefit(edge.start, edge.end))))
			g.AddOcurr(edge.start, edge.end)
			g.branchAndBound()
			g.RemoveOcurr(edge.start, edge.end)
		}
	}
	if len(solParcial) > 0 {
		fmt.Println("solParcial: ", solParcial)
		ultimo := solParcial[len(solParcial)-1]
		solParcial = solParcial[:len(solParcial)-1] // edge = eliminarUltimoLado(solParcial)
		g.RemoveOcurr(ultimo.start, ultimo.end)
		fmt.Println("Return Quitar ultimo:", ultimo, "ocurr:", ultimo.ocurr, "netB:", g.NetBenefit(ultimo.start, ultimo.end))
		beneficioDisponible = beneficioDisponible + int(math.Max(0, float64(g.NetBenefit(ultimo.start, ultimo.end))))
		g.AddOcurr(ultimo.start, ultimo.end)
	}
}

// Get path total benefit
func getPathBenefit(path []*Edge) int {
	seen := make(map[int]map[int]bool)
	total := 0
	for _, edge := range path {
		if len(seen[edge.start]) != 0 {
			if seen[edge.start][edge.end] {
				total = total - edge.cost
			} else {
				total = total + edge.benefit - edge.cost
			}

		} else {
			total = total + edge.benefit - edge.cost
			seen[edge.start] = make(map[int]bool)
		}
		if len(seen[edge.end]) == 0 {
			seen[edge.end] = make(map[int]bool)
		}
		seen[edge.start][edge.end] = true
		seen[edge.end][edge.start] = true
	}
	fmt.Println("camino: ", path, "benef: ", total)
	return total
}

func checkNegativeCycle(e *Edge, solParcial []*Edge) bool {
	path := append(solParcial, e)
	var total int
	for index, edge := range solParcial {
		if edge.start == e.end {
			total = getPathBenefit(path[index:])
			if total < 0 {
				return true
			}
		}
	}
	return false
}
