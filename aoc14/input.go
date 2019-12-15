package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type ItemWithAmount struct {
	Name   string
	Amount int
}

type Recipe struct {
	Output ItemWithAmount
	Inputs []ItemWithAmount
}

func readIntoItemWithAmount(s string, item *ItemWithAmount) error {
	s = strings.TrimSpace(s)
	if _, err := fmt.Sscanf(s, "%d %s", &item.Amount, &item.Name); err != nil {
		return err
	}
	return nil
}

func readRecipes(input io.Reader) (recipes []Recipe, err error) {
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()
		inOut := strings.Split(line, " => ")
		if len(inOut) != 2 {
			return nil, fmt.Errorf("bad line: %v", line)
		}
		var recipe Recipe
		if err := readIntoItemWithAmount(inOut[1], &recipe.Output); err != nil {
			return nil, err
		}
		inputs := strings.Split(inOut[0], ", ")
		var items []ItemWithAmount
		for _, s := range inputs {
			var item ItemWithAmount
			if err := readIntoItemWithAmount(s, &item); err != nil {
				return nil, err
			}
			items = append(items, item)
		}
		recipe.Inputs = items
		recipes = append(recipes, recipe)
	}
	return
}
