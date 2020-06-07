package main

import (
	"fmt"
	"strings"

	"github.com/arllanos/minesweeper-POC/services"
	"github.com/arllanos/minesweeper-POC/types"
)

func main() {
	game := &types.Game{
		Name:     "TestGame",
		Username: "ariel",
		Rows:     7,
		Cols:     7,
		Mines:    7,
	}

	services.Create(game)
	services.Start(game, "TestGame")

	for i, row := range game.Grid {
		fmt.Printf("%v = %v\n", i, strings.Split(string(row), ""))
	}

	click := &types.ClickData{
		Row:  1,
		Col:  2,
		Kind: "flag",
	}
	services.Click(game, "TestGame", click)

	click = &types.ClickData{
		Row:  1,
		Col:  1,
		Kind: "click",
	}
	services.Click(game, "TestGame", click)

	// click = &types.ClickData{
	// 	Row:  1,
	// 	Col:  2,
	// 	Kind: "flag",
	// }
	// services.Click(game, "TestGame", click)

	fmt.Println("--------------------")

	for i, row := range game.Grid {
		fmt.Printf("%v = %v\n", i, strings.Split(string(row), ""))
	}

}
