package agollo

import (
	"testing"
	"drcs/dep/agollo/test"
)

func TestStart(t *testing.T) {
	go runMockConfigServer(onlyNormalConfigResponse)
	go runMockNotifyServer(onlyNormalResponse)
	defer closeMockConfigServer()

	Start()

	value := getValue("key1")
	test.Equal(t,"value1",value)
}
