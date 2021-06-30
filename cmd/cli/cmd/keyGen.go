package cmd

import (
	"crypto/rand"
	"fmt"
	"github.com/fivebinaries/go-cardano-serialization/bip32"
	"github.com/spf13/cobra"
)

var normalKey, extendedKey, byronKey bool
var verificationKeyFile, signingKeyFile string

// keyGenCmd represents the keyGen command
var keyGenCmd = &cobra.Command{
	Use: "key-gen [--normal-key | --extended-key | --byron-key]\n" +
		"\t\t--verification-key-file FILE \n" +
		"\t\t--signing-key-file FILE",
	Short: "Create an address key pair.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !(normalKey || extendedKey || byronKey) {
			normalKey = true
		} else if !((normalKey && !(extendedKey || byronKey)) || (extendedKey && !(normalKey || byronKey)) ||
			(byronKey && !(normalKey || extendedKey))) {
			return cmd.Usage()
		}

		if len(verificationKeyFile) == 0 || len(signingKeyFile) == 0 {
			return cmd.Usage()
		}

		if normalKey {
			fmt.Println("Generating normal key")
		} else if extendedKey {
			fmt.Println("Generating extended key")
			seed := make([]byte, 32)
			_, err := rand.Read(seed)
			if err != nil {
				fmt.Println("Failed to generate random seed:", err)
				return err
			}

			privateKey, err := bip32.NewXPrv(seed)
			if err != nil {
				fmt.Println("Failed to generate new private key", err)
				return err
			}

			_, err = privateKey.Derive(0x1)
			if err != nil {
				fmt.Println("Failed to derive new key", err)
				return err
			}
		} else if byronKey {
			fmt.Println("Generating byron key")
		}
		return nil
	},
}

func init() {
	addressCmd.AddCommand(keyGenCmd)

	keyGenCmd.Flags().BoolVar(&normalKey, "normal-key", false, "Use a normal Shelley-era key.")
	keyGenCmd.Flags().BoolVar(&extendedKey, "extended-key", false, "Use an extended ed25519 Shelley-era key.")
	keyGenCmd.Flags().BoolVar(&byronKey, "byron-key", false, "Use a Byron-era key.")
	keyGenCmd.Flags().StringVar(&verificationKeyFile, "verification-key-file", "", "Output filepath of the verification key.")
	keyGenCmd.Flags().StringVar(&signingKeyFile, "signing-key-file", "", "Output filepath of the signing key.")
}
