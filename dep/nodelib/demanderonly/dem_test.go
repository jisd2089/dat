package demanderonly

/**
    Author: luzequan
    Created: 2018-01-02 19:06:23
*/
import (
	"testing"
	"fmt"
)

func TestDem(t *testing.T) {
	testRule := DEM.RuleTree.Trunk["ruleTest"].ParseFunc
	fmt.Println("**********", testRule)
}
