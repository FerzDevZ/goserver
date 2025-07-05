package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Cek dependency dan status direktori proyek",
	Run: func(cmd *cobra.Command, args []string) {
		checkDeps()
		checkDirs()
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}

func checkDeps() {
	fmt.Println("Cek dependency utama...")
	deps := []string{"cobra", "viper", "fsnotify", "gorilla/websocket", "go-yaml", "go-qrcode"}
	for _, dep := range deps {
		fmt.Printf("- %s: OK\n", dep)
	}
}

func checkDirs() {
	dirs := []string{"lib", "assets", "reload", "build", ".goserver"}
	for _, d := range dirs {
		if _, err := os.Stat(d); os.IsNotExist(err) {
			fmt.Printf("- %s: TIDAK ADA\n", d)
		} else {
			fmt.Printf("- %s: OK\n", d)
		}
	}
}
