package main

import (
	"fmt"
	"os"
)

func topoOrder(recipes map[string]Recipe, start string) ([]string, error) {
	var topo []string
	const temporary = 1
	const permanent = 2
	mark := make(map[string]int)

	var visit func(s string) error
	visit = func(s string) error {
		m, _ := mark[s]
		if m == permanent {
			return nil
		} else if m == temporary {
			return fmt.Errorf("Not a DAG")
		}
		mark[s] = temporary
		if recipe, ok := recipes[s]; ok {
			for _, input := range recipe.Inputs {
				if err := visit(input.Name); err != nil {
					return err
				}
			}
		}
		mark[s] = permanent
		topo = append(topo, s)
		return nil
	}

	err := visit(start)
	if err != nil {
		return nil, err
	}

	reverse(topo)
	return topo, nil
}

func reverse(a []string) {
	for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
		a[left], a[right] = a[right], a[left]
	}
}

func calcOreNeeded(topoSorted []string, recipes map[string]Recipe, fuelCount int) int {
	need := make(map[string]int)
	need["ORE"] = 0
	for name := range recipes {
		need[name] = 0
	}

	need["FUEL"] = fuelCount
	for _, name := range topoSorted {
		r := recipes[name]
		i := need[name]
		if i%r.Output.Amount > 0 {
			i = i/r.Output.Amount + 1
		} else {
			i = i / r.Output.Amount
		}
		for _, input := range r.Inputs {
			need[input.Name] += input.Amount * i
		}
	}
	return need["ORE"]
}

func main() {
	recipeList, err := readRecipes(os.Stdin)
	if err != nil {
		fmt.Printf("Could not parse input: %v\n", err)
		os.Exit(1)
	}

	recipes := make(map[string]Recipe)
	recipes["ORE"] = Recipe{ItemWithAmount{"ORE", 1}, nil}
	for _, r := range recipeList {
		recipes[r.Output.Name] = r
	}

	topoSorted, err := topoOrder(recipes, "FUEL")

	minOreNeeded := calcOreNeeded(topoSorted, recipes, 1)
	fmt.Printf("Part 1: %v\n", minOreNeeded)

	const ORE_MAX = 1000000000000
	min := 1
	max := 2
	for {
		on := calcOreNeeded(topoSorted, recipes, max)
		if on > ORE_MAX {
			break
		} else {
			min = max
			max *= 2
		}
	}
	for {
		pivot := min + ((max - min) / 2)
		on := calcOreNeeded(topoSorted, recipes, pivot)
		fmt.Printf("%v %v -> %v: %v\n", min, max, pivot, on)
		if on > ORE_MAX {
			max = pivot
		} else {
			min = pivot
		}
		if min == max || min == max-1 {
			break
		}
	}
	fmt.Printf("Part 2: %v", min)
}
