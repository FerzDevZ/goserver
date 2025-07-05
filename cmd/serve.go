package cmd

import (
	"fmt"
	"github.com/FerzDevZ/goserver/internal/server"
	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/websocket"
	"github.com/skip2/go-qrcode"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var (
	servePort    int
	serveHost    string
	serveOpen    bool
	serveNoInject bool
	serveWatchDir string
	serveLogLevel string
	serveProd    bool
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Jalankan static server dengan live reload",
	Run: func(cmd *cobra.Command, args []string) {
		viper.SetConfigName("goserver")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		_ = viper.ReadInConfig()
		port := servePort
		if port == 0 {
			port = viper.GetInt("server.port")
			if port == 0 { port = 8080 }
		}
		host := serveHost
		if host == "" {
			host = viper.GetString("server.host")
			if host == "" { host = "localhost" }
		}
		staticDir := viper.GetString("server.static_dir")
		if staticDir == "" { staticDir = "lib" }
		assetsDir := viper.GetString("server.assets_dir")
		if assetsDir == "" { assetsDir = "assets" }
		watchDir := serveWatchDir
		if watchDir == "" { watchDir = viper.GetString("server.watch_dir") }
		if watchDir == "" { watchDir = "lib" }
		prod := serveProd
		if !prod { prod = cmd.Flag("prod").Changed }
		noInject := serveNoInject
		ls := server.NewLiveServer(port, host, staticDir, assetsDir, watchDir, prod, noInject)
		if err := ls.Start(); err != nil {
			fmt.Println("Server error:", err)
		}
	},
}

func init() {
	serveCmd.Flags().IntVar(&servePort, "port", 0, "Port server")
	serveCmd.Flags().StringVar(&serveHost, "host", "", "Host address")
	serveCmd.Flags().BoolVar(&serveOpen, "open", false, "Auto open di browser")
	serveCmd.Flags().BoolVar(&serveNoInject, "no-inject", false, "Jangan sisipkan reload script")
	serveCmd.Flags().StringVar(&serveWatchDir, "watch-dir", "", "Direktori yang dipantau")
	serveCmd.Flags().StringVar(&serveLogLevel, "log-level", "info", "Log level")
	serveCmd.Flags().BoolVar(&serveProd, "prod", false, "Mode production")
	rootCmd.AddCommand(serveCmd)
}

// Implementasi fungsi startWatcher, startWebSocketServer, injectReload, openBrowser, printQRCode, dashboardHandler, getLocalIP
