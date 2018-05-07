package security

import (
	"drcs/dep/cncrypt"
	"testing"
)

func TestSignature(t *testing.T) {
	priKey := "81EB26E941BB5AF16DF116495F90695272AE2CD63D6C4AE1678418BE48230029"
	_privateKey = priKey
	cncrypt.Init(priKey)

	pubKey := "03160e12897df4edb61dd812feb96748fbd3ccf4ffe26aa6f6db9540af49c94232"
	content := []byte{'1'}

	signature, _ := Signature(content)
	result := VerifySignature(pubKey, content, signature)
	if result != true {
		t.Fatal("verify failed")
	}
}
