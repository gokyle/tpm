package tpm

import "fmt"
import "testing"

var tpm *TPMContext

func TestInitTPM(t *testing.T) {
	var err error
	tpm, err = NewTPMContext()
	if err != nil {
		t.Fatalf("%v", err)
	}
}

func TestRand(t *testing.T) {
	rdata, err := tpm.Random(32)
	if err != nil {
		t.Fatalf("%v", err)
	}
	fmt.Printf("TPM random: %x\n", rdata)
}

func TestShutdownTPM(t *testing.T) {
	err := tpm.Destroy()
	if err != nil {
		t.Fatalf("%v", err)
	}
}
