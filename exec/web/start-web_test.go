package web

/**
    Author: luzequan
    Created: 2018-01-04 13:42:46
*/
import (
	"testing"
	"dat/core"

	_ "dat/dep/nodelib/demanderonly"
)

func init() {
	assetnode.AssetNodeEntity.Init()
}

func TestRun(t *testing.T) {
	Run()
}
