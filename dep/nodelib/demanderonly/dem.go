package demanderonly

/**
    Author: luzequan
    Created: 2017-12-28 17:22:00
*/
import (
	"dat/core/interaction/request"
	. "dat/core/dataflow"
	"fmt"
)

func init() {
	DEM.Register()
}

var DEM = &DataFlow{
	Name:        "demtest",
	Description: "demtest",
	// Pausetime:    300,
	// Keyin:        KEYIN,
	// Limit:        LIMIT,
	EnableCookie: false,
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {
			fmt.Println(ctx)
			ctx.AddQueue(&request.DataRequest{
				Url:  "http://www.inderscience.com/info/inarticletoc.php?jcode=ijguc&year=2016&vol=7&issue=1",
				Rule: "ruleTest",
			})
		},

		Trunk: map[string]*Rule{
			"ruleTest": {
				ParseFunc: func(ctx *Context) {
					fmt.Println(ctx)
				},
			},
		},
	},
}