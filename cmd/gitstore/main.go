package main

import (
	"fmt"
	"os"

	"github.com/Nekroze/gitstore"
)

func main() {
	var out string
	var err error
	switch os.Args[1] {
	case "read":
		out, err = gitstore.Read("go", os.Args[2])
	case "write":
		err = gitstore.Write("go", os.Args[2], os.Args[3])
	default:
		fmt.Printf("%q is not valid command.\n", os.Args[1])
		os.Exit(1)
	}

	if len(out) > 0 {
		fmt.Println(out)
	}
	if err != nil {
		fmt.Printf("%+v", err)
		os.Exit(1)
	}
}
