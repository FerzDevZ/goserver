package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"io/ioutil"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Buat struktur proyek baru mirip Flutter",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		if err := createProject(name); err != nil {
			fmt.Println("Gagal membuat proyek:", err)
			os.Exit(1)
		}
		fmt.Println("Proyek berhasil dibuat:", name)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}

func createProject(name string) error {
	dirs := []string{"assets", "build", "lib", "reload", ".goserver"}
	for _, d := range dirs {
		if err := os.MkdirAll(filepath.Join(name, d), 0755); err != nil {
			return err
		}
	}
	// Copy template files
	files := map[string]string{
		"lib/index.html": "templates/index.html",
		"lib/style.css": "templates/style.css",
		"lib/main.js": "templates/main.js",
		"reload/reload.js": "templates/reload.js",
	}
	for dst, src := range files {
		content, err := ioutil.ReadFile(src)
		if err != nil {
			return err
		}
		if err := ioutil.WriteFile(filepath.Join(name, dst), content, 0644); err != nil {
			return err
		}
	}
	// Create goserver.yaml
	cfg := []byte("name: " + name + "\ndescription: A static web project with live reload\nversion: 1.0.0\nserver:\n  port: 8080\n  host: localhost\n  watch_dir: lib\n  static_dir: lib\n  assets_dir: assets\nbuild:\n  enabled: true\n  script: \"npm run build\"\n")
	return ioutil.WriteFile(filepath.Join(name, "goserver.yaml"), cfg, 0644)
}
