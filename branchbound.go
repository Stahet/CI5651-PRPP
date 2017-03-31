package main

import (
	"fmt"
	"math"
	"sort"
)

func (g *Graph) estaEnSolucionParcial(e *Edge, solParcial []*Edge) bool {
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
	fmt.Println("SolParcialBenefit", beneficioSolParcial, "beneficioMejor:", getPathBenefit(mejorSol))
	maxBeneficio := beneficioDisponible - int(math.Max(0, float64(beneficioE))) + beneficioSolParcial
	if maxBeneficio <= getPathBenefit(mejorSol) {
		return false
	}
	return true
}

func (g *Graph) branchAndBound(e int, solParcial []*Edge, mejorSol []*Edge, beneficioDisponible int) ([]*Edge, []*Edge, int) {
	if e == 1 {
		if getPathBenefit(solParcial) > getPathBenefit(mejorSol) {
			mejorSol = solParcial
		}
	}
	//fmt.Println(mejorSol)
	sucesores := g.obtenerListaSucesores(e)
	estaSolParcial, cumpleAco, cicloNeg := true, true, true
	fmt.Print("solParcial: ", solParcial, "benef:", getPathBenefit(solParcial))
	fmt.Println("sucesores: ", sucesores)
	for _, edge := range sucesores {
		estaSolParcial = g.estaEnSolucionParcial(edge, solParcial)
		cumpleAco = g.cumpleAcotamiento(edge, solParcial, mejorSol, beneficioDisponible)
		cicloNeg = g.checkNegativeCycle(edge, solParcial)
		fmt.Println("nodo:", e, "lado: ", edge, "| netB:", g.NetBenefit(edge.start, edge.end), "| estaSolPar:", estaSolParcial, "| cumpleAc:", cumpleAco, "| NegCycle:", cicloNeg, "| benefDisponible:", beneficioDisponible, "| cond:", !estaSolParcial && cumpleAco && !cicloNeg)

		if !estaSolParcial && cumpleAco && !cicloNeg {
			fmt.Println()
			solParcial = append(solParcial, edge)
			beneficioDisponible = beneficioDisponible - int(math.Max(0, float64(g.NetBenefit(edge.start, edge.end))))
			g.AddOcurr(edge.start, edge.end)
			solParcial, mejorSol, beneficioDisponible = g.branchAndBound(edge.end, solParcial, mejorSol, beneficioDisponible)
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
	return solParcial, mejorSol, beneficioDisponible
}
