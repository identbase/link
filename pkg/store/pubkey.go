package store

import (
	"fmt"
)

/*
PublicKey is used for various purposes. An identity server has some
long-term public-private keypairs. These are named in a scheme
algorithm:identifier, e.g. ed25519:0. When signing an association, the standard
Signing JSON algorithm applies.

The identity server may also keep track of some short-term public-private
keypairs, which may have different usage and lifetime characteristics than the
service's long-term keys. */
type PublicKey struct {
	Algorithm  string `json:"algorithm"`
	Identifier string `json:"identifier"`
	Content    string `json:"public_key"`
}

/*
Key is the key used to save and lookup a particular item. */
func (p *PublicKey) Key() string {
	return fmt.Sprintf("%s:%s", p.Algorithm, p.Identifier)
}
