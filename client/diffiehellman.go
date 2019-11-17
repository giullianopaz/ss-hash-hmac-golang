package main

import (
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"math"
)

func isPrime(value int) bool {
	if value <= 1 {
		return false
	}

	for i := 2; i < value; i++ {
		if value%i == 0 {
			return false
		}
	}
	return true
}

func getHash(text []byte) string {
	hasher := sha512.New()
	hasher.Write(text)
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

// DiffieHellman :
type DiffieHellman struct {
	pModulusValue    int
	gBaseValue       int
	privateValue     int
	publicValue      int
	sharedPrivateKey string
}

// SetpModulusValue : Recebe o valor primo para `pModulusValue`
func (dh *DiffieHellman) SetpModulusValue(value int) {
	fmt.Print("\n-> Setting Modulus Value... ")
	dh.pModulusValue = value
	fmt.Println("OK")

}

// SetgBaseValue : Recebe o valor primo para `gBaseValue`
func (dh *DiffieHellman) SetgBaseValue(value int) {
	fmt.Print("-> Setting Base Value... ")
	dh.gBaseValue = value
	fmt.Println("OK")

}

// GeneratePrivateValue : Gera valor secreto do servidor
func (dh *DiffieHellman) GeneratePrivateValue() {
	fmt.Print("-> Generating Private Value... ")
	dh.privateValue = 1 + RAND.Intn(dh.pModulusValue)
	fmt.Println("OK")
}

// GeneratePublicValue : Gera valor público do servidor
func (dh *DiffieHellman) GeneratePublicValue() {
	fmt.Print("-> Generating Public Value... ")

	fmt.Printf("\nmath.Mod(math.Pow(%f, %f), %f))\n",
		float64(dh.gBaseValue), float64(dh.privateValue), float64(dh.pModulusValue))

	dh.publicValue = int(
		math.Mod(math.Pow(float64(dh.gBaseValue), float64(dh.privateValue)), float64(dh.pModulusValue)))
	fmt.Println("OK")
}

// GenerateSharedPrivateKey : Gera chave privada compartilhada
func (dh *DiffieHellman) GenerateSharedPrivateKey(sharedPublicValue int) {
	fmt.Print("-> Generating Shared Private key... ")

	fmt.Printf("\nmath.Mod(math.Pow(%f, %f), %f))\n",
		float64(sharedPublicValue), float64(dh.privateValue), float64(dh.pModulusValue))

	res := int(
		math.Mod(math.Pow(float64(sharedPublicValue), float64(dh.privateValue)), float64(dh.pModulusValue)))

	dh.sharedPrivateKey = getHash([]byte(string(res)))
	// dh.sharedPrivateKey = res
	fmt.Println("OK")
}
