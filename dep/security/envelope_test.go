package security

import (
	"drcs/dep/cncrypt"
	"testing"
)

func TestEnvelope(t *testing.T) {
	priKey := "81EB26E941BB5AF16DF116495F90695272AE2CD63D6C4AE1678418BE48230029"
	_privateKey = priKey
	cncrypt.Init(priKey)

	pubKey := "03160e12897df4edb61dd812feb96748fbd3ccf4ffe26aa6f6db9540af49c94232"
	content := "1"

	envelope := PackEnvelope(pubKey, content)
	// fmt.Println(envelope.EncryptedKey)
	// fmt.Println(envelope.Content)

	data, err := UnpackEnvelope(envelope)
	if err != nil {
		t.Fatal(err)
	}
	if data != content {
		t.Fatalf("%s != %s", data, content)
	}
}
