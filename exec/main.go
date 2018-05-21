package main

/**
    Author: luzequan
    Created: 2018-01-02 11:21:18
*/
import (
	"path/filepath"
	"os"
	"fmt"
	"drcs/exec/web"
)

func main() {

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
	}
	fmt.Println("current dir: " + dir)

	web.Run()
}