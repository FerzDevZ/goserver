package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "goserver",
	Short: "goserver - Static web server with live reload and modern project structure",
	Long:  `goserver adalah CLI tool untuk membangun dan menyajikan aplikasi web statis dengan live reload, mirip Flutter CLI.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
