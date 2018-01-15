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

func TestDemRun(t *testing.T) {
	Run(8899)
}

func TestSupRun(t *testing.T) {
	Run(8081)
}
