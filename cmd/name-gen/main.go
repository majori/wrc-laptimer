package main

import (
	"fmt"
	"os"

	"github.com/majori/wrc-laptimer/pkg/username"
)

func main() {
	seed := os.Args[1]
	fmt.Println(username.GenerateFromSeed(seed))
}
