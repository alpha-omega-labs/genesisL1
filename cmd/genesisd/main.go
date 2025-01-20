package main

import (
	"os"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/crypto-org-chain/cronos/app"
	"github.com/crypto-org-chain/cronos/cmd/genesisd/cmd"

	// [ADDED CODE]
	"fmt"
	"os/exec"
	"github.com/spf13/viper"
)

func main() {
	rootCmd, _ := cmd.NewRootCmd()

	// [ADDED CODE]: Launch sidecar in background if --bio is set
	go func() {
		if viper.GetBool("bio") {
			sideCmd := exec.Command("python", "sidecar.py")
			sideCmd.Stdout = os.Stdout
			sideCmd.Stderr = os.Stderr
			if err := sideCmd.Start(); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to start Biopython sidecar: %v\n", err)
			} else {
				fmt.Println("Biopython sidecar started on port 8000.")
			}
		}
	}()

	if err := svrcmd.Execute(rootCmd, cmd.EnvPrefix, app.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
