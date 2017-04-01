package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Global variables
var mejorSol []*Edge
var beneficioDisponible int
var solParcial []*Edge

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	if len(os.Args) <= 2 {
		fmt.Println("Para ejecutar ./main <nombre-archivo> <valor-optimo>")
		return
	}

	beginning := time.Now()
	args := os.Args
	file, _ := os.Open(args[1])
	// file, _ := os.Open("./instanciasPRPP/CHRISTOFIDES/P01NoRPP")
	//file, _ := os.Open("./instanciasPRPP/RANDOM/R9NoRPP")
	// file, _ := os.Open("./instanciasPRPP/DEGREE/D2NoRPP")
	//file, _ := os.Open("./instanciasPRPP/GRID/G16NoRPP")
	lineScanner := bufio.NewScanner(file)
	line := 0
	g := NewGraph(1)
	maxBenefit, b, c := 0, 0, 0
	for lineScanner.Scan() {
		contents := strings.Fields(lineScanner.Text())
		if line == 0 {
			number, _ := strconv.ParseInt(contents[len(contents)-1], 0, 0)
			g = NewGraph(int(number))
		}
		if _, err := strconv.Atoi(contents[0]); err == nil {
			startNode, _ := strconv.ParseInt(contents[0], 0, 0)
			endNode, _ := strconv.ParseInt(contents[1], 0, 0)
			cost, _ := strconv.ParseInt(contents[2], 0, 0)
			benefit, _ := strconv.ParseInt(contents[3], 0, 0)
			g.AddEdge(int(startNode), int(endNode), int(cost), int(benefit))
			b = g.Benefit(int(startNode), int(endNode))
			c = g.Cost(int(startNode), int(endNode))
			if b-c >= 0 {
				maxBenefit = maxBenefit + b - c
			}
		}
		line++
	}
	var path []*Edge
	path = getCycleGRASP(g)             // Get Greedy initial Path
	path = removeNegativeCycle(g, path) // Remove Negative Cycle
	mejorSol = make([]*Edge, 0)         // Global variable bestPath
	// Copy path to new array
	for _, elem := range path {
		mejorSol = append(mejorSol, elem)
	}
	beneficioDisponible = maxBenefit // Global variable maxBenefit
	_ = time.AfterFunc(time.Duration(120)*time.Minute, func() {
		fmt.Println("Archivo: ", args[1])
		fmt.Println("Tiempo limite excedido")
		buf := bufio.NewReader(os.Stdin)
		fmt.Print("Desea continuar la busqueda? Y/n")
		sentence, err := buf.ReadBytes('\n')
		if err != nil {
			fmt.Println(err)
		} else if sentence == 'N' {
			os.Exit(2)
		}
	})
	g.branchAndBound(1)
	ending := time.Since(beginning)
	branchValue := getPathBenefit(mejorSol)

	salida, err := os.Create(args[1] + "-salida.txt")
	check(err)
	defer salida.Close()
	stringValue := strconv.Itoa(getPathBenefit(mejorSol))
	stringPath := []string{}
	stringTime := ending.String()
	_, err = salida.WriteString(stringValue)
	check(err)
	_, err = salida.WriteString("\n")
	check(err)
	stringPath = append(stringPath, strconv.Itoa(1))
	for _, edge := range mejorSol {
		number := edge.end
		text := strconv.Itoa(number)
		stringPath = append(stringPath, text)
	}
	result := strings.Join(stringPath, " ")
	result = "d " + result + " d"
	_, err = salida.WriteString(result)
	check(err)
	_, err = salida.WriteString("\n")
	check(err)
	optimalValue, _ := strconv.ParseInt(args[2], 0, 0)
	var heuristicDeviation float64
	if optimalValue != 0 {
		heuristicDeviation = float64(100 * (float64(optimalValue) - float64(branchValue)) / float64(optimalValue))
	}
	_, err = salida.WriteString(strconv.FormatFloat(heuristicDeviation, 'f', 2, 32))
	check(err)
	_, err = salida.WriteString("%\n")
	check(err)
	_, err = salida.WriteString(stringTime)
	check(err)
	salida.Sync()
	fmt.Println("Archivo: ", args[1])
	fmt.Println("Ciclo Branch and bound: ", mejorSol)
	fmt.Println("Total: ", branchValue)
	fmt.Println("Tiempo: ", ending)
	fmt.Printf("Desviacion: %.2f %%", heuristicDeviation)
}
