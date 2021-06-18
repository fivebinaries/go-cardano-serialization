module github.com/fivebinaries/go-cardano-serialization/cmd/cli

go 1.15

require (
	github.com/fivebinaries/go-cardano-serialization v0.0.0-00010101000000-000000000000
	github.com/mitchellh/go-homedir v1.1.0
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.7.1
)

replace github.com/fivebinaries/go-cardano-serialization => ../..
