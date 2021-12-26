package main

import (
	"fmt"
	"path"
	"path/filepath"
	"runtime"

	"github.com/fivebinaries/go-cardano-serialization/crypto"
	"github.com/fivebinaries/go-cardano-serialization/fees"
	"github.com/fivebinaries/go-cardano-serialization/lib"
	"github.com/fivebinaries/go-cardano-serialization/metadata"
	"github.com/fivebinaries/go-cardano-serialization/protocol"
	"github.com/fivebinaries/go-cardano-serialization/transactions"
	"github.com/fivebinaries/go-cardano-serialization/types"
	"github.com/fivebinaries/go-cardano-serialization/utils"
)

func main() {
	_, rnFp, _, _ := runtime.Caller(0)
	fp := filepath.Join(path.Dir(path.Dir(rnFp)), "transactions", "protocol.json")
	pr, err := protocol.LoadProtocolFromFile(fp)
	if err != nil {
		panic(err)
	}

	txBuilder := transactions.NewTransactionBuilder(
		&fees.LinearFee{
			Constant:    types.Coin(pr.TxFeeFixed),
			Coefficient: types.Coin(pr.TxFeePerByte),
		},
		1000000,
		utils.BigNum(pr.StakePoolDeposit),    // Protocol Parameter stakePoolDeposit
		utils.BigNum(pr.StakeAddressDeposit), // Protocol Parameter stakeAddressDeposit
	)

	var ttl lib.Slot = 410021

	txBuilder.Metadata = metadata.NewTransactionMetadata(
		metadata.NewGeneralTransactionMetadata(),
	)

	// add a keyhash input - for ada held in a Shelley-era normal address (Base, Enterprise, Pointer)
	prvKey, err := types.AddressFromBech32("ed25519e_sk16rl5fqqf4mg27syjzjrq8h3vq44jnnv52mvyzdttldszjj7a64xtmjwgjtfy25lu0xmv40306lj9pcqpa6slry9eh3mtlqvfjz93vuq0grl80")
	if err != nil {
		panic(err)
	}
	inputAddr, err := crypto.Ed25519KeyHashFromBytes(prvKey.ToBytes())
	if err != nil {
		panic(err)
	}

	amountIn := types.Coin(3000000)
	txBuilder.AddKeyInput(
		inputAddr,
		&types.TransactionInput{
			TransactionId: types.Hash32("8561258e210352fba2ac0488afed67b3427a27ccf1d41ec030c98a8199bc22ec"),
			Index:         0,
		},
		&types.Value{
			&amountIn,
			&types.ValueAdditionalType0{},
		},
	)

	shelleyOutputAddr, err := types.AddressFromBech32("addr_test1qpu5vlrf4xkxv2qpwngf6cjhtw542ayty80v8dyr49rf5ewvxwdrt70qlcpeeagscasafhffqsxy36t90ldv06wqrk2qum8x5w")
	if err != nil {
		panic(err)
	}

	shelleyChangeAddr, err := types.AddressFromBech32("addr_test1gz2fxv2umyhttkxyxp8x0dlpdt3k6cwng5pxj3jhsydzerspqgpsqe70et")
	if err != nil {
		panic(err)
	}

	amountOut := types.Coin(1000000)
	txBuilder.AddOutput(
		&types.TransactionOutput{
			V1: shelleyOutputAddr,
			Amount: types.Value{
				&amountOut,
				&types.ValueAdditionalType0{},
			},
		},
	)

	ttl = lib.Slot(410021)
	txBuilder.TTL = &ttl

	txBuilder.AddChangeIfNeeded(shelleyChangeAddr)

	txBody, err := txBuilder.Build()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", txBody)
}
