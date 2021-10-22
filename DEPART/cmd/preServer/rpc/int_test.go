package rpc

import (
	"math/big"
	"testing"
)

func TestChainIndexerSingle(t *testing.T) {
	b := big.NewInt(1231232312312)

	str := b.String()

	c := big.NewInt(0)
	c.SetString(str, 10)

	println(b.Int64())
	println(c.Int64())
}
