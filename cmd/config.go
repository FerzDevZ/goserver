package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Baca/tulis konfigurasi goserver.yaml",
	Run: func(cmd *cobra.Command, args []string) {
		viper.SetConfigName("goserver")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		if err := viper.ReadInConfig(); err != nil {
			fmt.Println("Config error:", err)
			os.Exit(1)
		}
		if len(args) == 2 && args[0] == "--set" {
			kv := args[1]
			var key, value string
			_, err := fmt.Sscanf(kv, "%[^=]=%s", &key, &value)
			if err != nil {
				fmt.Println("Format: --set key=value")
				return
			}
			viper.Set(key, value)
			if err := viper.WriteConfig(); err != nil {
				fmt.Println("Gagal update config:", err)
				return
			}
			fmt.Println("Config diperbarui:", key, "=", value)
		} else {
			fmt.Println("Konfigurasi saat ini:")
			for _, k := range viper.AllKeys() {
				fmt.Printf("%s: %v\n", k, viper.Get(k))
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
