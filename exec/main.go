package main

/**
    Author: luzequan
    Created: 2018-01-02 11:21:18
*/
import (
	"drcs/exec/web"

	_ "drcs/dep/nodelib/dep"
	_ "drcs/dep/nodelib/batchdistribution"


)

func main() {
	web.Run()
}