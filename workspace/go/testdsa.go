package main

import (
	"crypto/dsa"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

func main() {

	var key = []byte(`-----BEGIN ENCRYPTED PRIVATE KEY-----
MIIBuDCCASwGByqGSM44BAEwggEfAoGBAP1/U4EddRIpUt9KnC7s5Of2EbdSPO9E
AMMeP4C2USZpRV1AIlH7WT2NWPq/xfW6MPbLm1Vs14E7gB00b/JmYLdrmVClpJ+f
6AR7ECLCT7up1/63xhv4O1fnxqimFQ8E+4P208UewwI1VBNaFpEy9nXzrith1yrv
8iIDGZ3RSAHHAhUAl2BQjxUjC8yykrmCouuEC/BYHPUCgYEA9+GghdabPd7LvKtc
NrhXuXmUr7v6OuqC+VdMCz0HgmdRWVeOutRZT+ZxBxCBgLRJFnEj6EwoFhO3zwky
jMim4TwWeotUfI0o4KOuHiuzpnWRbqN/C/ohNWLx+2J6ASQ7zKTxvqhRkImog9/h
WuWfBpKLZl6Ae1UlZAFMO/7PSSoDgYUAAoGBAP1R1jLPc1kikRwexRvKZhmR01hx
FTCYrRaDX8/g+gmQAWWHf0fOrAi0R7dr6BRlT3unfNMgAi8U2+Iet7vpSz1EgG4Z
XRc4XSK704jhMV0FPF98OFKFDBWlxJsNnt/MwKiwIA9KHbC89OzJGSap02Mqfa0f
8LzMUkP848EZDJkD
-----END ENCRYPTED PRIVATE KEY-----`)

	block, _ := pem.Decode(key)
	if block == nil {
		fmt.Errorf("expected block to be non-nil", block)
		panic("block is nil")
	}
	fmt.Println("block:", block)
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		fmt.Errorf("could not unmarshall data: `%s`", err)
	}
	fmt.Println("pub:", pub)
	switch pub := pub.(type) {
	case *rsa.PublicKey:
		fmt.Println("pub is of type RSA:", pub)
	case *dsa.PublicKey:
		fmt.Println("pub is of type DSA:", pub)
	case *ecdsa.PublicKey:
		fmt.Println("pub is of type ECDSA:", pub)
	default:
		panic("unknown type of public key")
	}

	fmt.Printf("done")
}
