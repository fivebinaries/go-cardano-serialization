// A basic example of how to initialize an address from bech32 or base58(byron)
// address string

package main

import (
	"log"

	"github.com/fivebinaries/go-cardano-serialization/address"
)

func main() {
	addr, err := address.NewAddress("addr1vxvfjdcr6nd3y65mkqx5gfvgwjarvkt527h35wq3g7a447chmv4zj")
	//addr, err := address.NewAddress("addr_test1qqe92py4mf3ffrtmjuwjpzu6jwlw0zmr50h8ey67qcehlmty5kcrvg2ds9fkpg32t535l9v6lkgaj5cunufgvz5f7snql2fawd")

	if err != nil {
		log.Fatal(err)
	}

	switch addr.(type) {
	case *address.BaseAddress:
		log.Println("Base Addr:", addr.String())
	case *address.EnterpriseAddress:
		log.Println("Enterprise Address:", addr.String())
	default:
		log.Println("Not baseAddress")
	}
}
