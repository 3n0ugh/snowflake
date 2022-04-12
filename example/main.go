package main

import (
	"fmt"

	"github.com/3n0ugh/snowflake"
)

func main() {
	n, err := snowflake.NewNode(30, 3)
	if err != nil {
		fmt.Println(err)
	}

	id, err := n.Generate()
	if err != nil {
		fmt.Println(err)
	}

	id2, err := n.Generate()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("ID: %d\n", id)
	fmt.Printf("String: %s\n", id.String())
	fmt.Printf("Uint64: %d\n", id.UInt64())

	fmt.Printf("DecomposeID: %v\n", snowflake.DecomposeID(id))
	fmt.Printf("DecomposeID: %v\n", snowflake.DecomposeID(id2))
}
