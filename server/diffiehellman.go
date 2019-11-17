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

// GeneratepModulusValue : Escolhe um valor primo para `pModulusValue`
func (dh *DiffieHellman) GeneratepModulusValue(ceilingValue int) {
	fmt.Print("\n-> Generating Modulus Value... ")
	dh.pModulusValue = ceilingValue
	for !isPrime(dh.pModulusValue) {
		dh.pModulusValue--
	}
	fmt.Println("OK")

}

// GenerategBaseValue : Escolhe um valor primo para `gBaseValue` a partir de `pModulusValue`
func (dh *DiffieHellman) GenerategBaseValue() {
	fmt.Print("-> Generating Base Value... ")
	dh.gBaseValue = RAND.Intn(dh.pModulusValue)
	for !isPrime(dh.gBaseValue) {
		dh.gBaseValue = RAND.Intn(dh.pModulusValue)
	}
	fmt.Println("OK")

}

// GeneratePrivateValue : Gera valor secreto do servidor
func (dh *DiffieHellman) GeneratePrivateValue() {
	fmt.Print("-> Generating Private Value... ")
	dh.privateValue = RAND.Intn(dh.pModulusValue)
	fmt.Println("OK")
}

// GeneratePublicValue : Gera valor público do servidor
func (dh *DiffieHellman) GeneratePublicValue() {
	fmt.Print("-> Generating Public Value... ")
	dh.publicValue = int(
		math.Mod(math.Pow(float64(dh.gBaseValue), float64(dh.privateValue)), float64(dh.pModulusValue)))
	fmt.Println("OK")
}

// GenerateSharedPrivateKey : Gera chave privada compartilhada
func (dh *DiffieHellman) GenerateSharedPrivateKey(sharedPublicValue int) {
	fmt.Print("-> Generating Shared Private key... ")
	res := int(
		math.Mod(math.Pow(float64(sharedPublicValue), float64(dh.privateValue)), float64(dh.pModulusValue)))

	dh.sharedPrivateKey = getHash([]byte(string(res)))
	fmt.Println("OK")
}
