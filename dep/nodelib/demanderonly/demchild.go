package demanderonly

/**
    Author: luzequan
    Created: 2018-01-16 15:37:57
*/
import (
	. "dat/core/databox"
	"fmt"
)

func init() {
	DEMCHILD.Register()
}

var DEMCHILD = &DataBox{
	Name:        "demchild",
	Description: "demchild",
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {
			fmt.Println("demchild start...")
		},
		Trunk: map[string]*Rule{
			"split": {
				ParseFunc: func(ctx *Context) {
					fmt.Println("start ...")

				},
			},
		},
	},
}
