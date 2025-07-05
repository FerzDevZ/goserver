package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Jalankan build script dari goserver.yaml",
	Run: func(cmd *cobra.Command, args []string) {
		viper.SetConfigName("goserver")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		if err := viper.ReadInConfig(); err != nil {
			fmt.Println("Config error:", err)
			os.Exit(1)
		}
		buildEnabled := viper.GetBool("build.enabled")
		if !buildEnabled {
			fmt.Println("Build tidak diaktifkan di goserver.yaml")
			return
		}
		script := viper.GetString("build.script")
		if script == "" {
			fmt.Println("Script build tidak ditemukan di config")
			return
		}
		fmt.Println("Menjalankan build:", script)
		cmdExec := exec.Command("bash", "-c", script)
		cmdExec.Stdout = os.Stdout
		cmdExec.Stderr = os.Stderr
		cmdExec.Dir = "."
		if err := cmdExec.Run(); err != nil {
			fmt.Println("Build gagal:", err)
			os.Exit(1)
		}
		// Copy hasil build ke build/
		fmt.Println("Build selesai. Output ke folder build/")
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
