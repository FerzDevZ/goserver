package server

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/websocket"
	"github.com/skip2/go-qrcode"
)

type LiveServer struct {
	Port      int
	Host      string
	StaticDir string
	AssetsDir string
	WatchDir  string
	Prod      bool
	NoInject  bool
	Clients   map[*websocket.Conn]bool
	Mu        sync.Mutex
	Watcher   *fsnotify.Watcher
}

func NewLiveServer(port int, host, staticDir, assetsDir, watchDir string, prod, noInject bool) *LiveServer {
	return &LiveServer{
		Port:      port,
		Host:      host,
		StaticDir: staticDir,
		AssetsDir: assetsDir,
		WatchDir:  watchDir,
		Prod:      prod,
		NoInject:  noInject,
		Clients:   make(map[*websocket.Conn]bool),
	}
}

func (s *LiveServer) Start() error {
	if !s.Prod {
		w, err := fsnotify.NewWatcher()
		if err != nil {
			return err
		}
		s.Watcher = w
		go s.watchFiles()
		if err := s.addWatchRecursive(s.WatchDir); err != nil {
			return err
		}
	}
	http.HandleFunc("/ws", s.wsHandler)
	http.HandleFunc("/__dashboard__", s.dashboardHandler)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(s.AssetsDir))))
	http.HandleFunc("/", s.serveIndex)
	url := fmt.Sprintf("http://%s:%d", s.Host, s.Port)
	fmt.Println("Server running at:", url)
	printQRCode(url)
	return http.ListenAndServe(fmt.Sprintf("%s:%d", s.Host, s.Port), nil)
}

func (s *LiveServer) serveIndex(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join(s.StaticDir, r.URL.Path)
	if info, err := os.Stat(path); err == nil && !info.IsDir() {
		if !s.Prod && !s.NoInject && strings.HasSuffix(path, "index.html") {
			injectReloadJS(w, path)
			return
		}
		http.ServeFile(w, r, path)
		return
	}
	// fallback ke index.html
	indexPath := filepath.Join(s.StaticDir, "index.html")
	if !s.Prod && !s.NoInject {
		injectReloadJS(w, indexPath)
		return
	}
	http.ServeFile(w, r, indexPath)
}

func injectReloadJS(w http.ResponseWriter, file string) {
	b, err := os.ReadFile(file)
	if err != nil {
		http.Error(w, "Not found", 404)
		return
	}
	if !bytes.Contains(b, []byte("reload.js")) {
		b = bytes.Replace(b, []byte("</body>"), []byte("<script src=\"/reload/reload.js\"></script></body>"), 1)
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(b)
}

func (s *LiveServer) wsHandler(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	s.Mu.Lock()
	s.Clients[conn] = true
	s.Mu.Unlock()
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
	s.Mu.Lock()
	delete(s.Clients, conn)
	s.Mu.Unlock()
	conn.Close()
}

func (s *LiveServer) broadcastReload() {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	for c := range s.Clients {
		_ = c.WriteMessage(websocket.TextMessage, []byte("reload"))
	}
}

func (s *LiveServer) watchFiles() {
	for {
		select {
		case event, ok := <-s.Watcher.Events:
			if !ok { return }
			if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Remove|fsnotify.Rename) != 0 {
				s.broadcastReload()
			}
		case err, ok := <-s.Watcher.Errors:
			if !ok { return }
			fmt.Println("Watcher error:", err)
		}
	}
}

func (s *LiveServer) addWatchRecursive(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil { return err }
		if info.IsDir() {
			return s.Watcher.Add(path)
		}
		return nil
	})
}

func printQRCode(url string) {
	var buf bytes.Buffer
	_ = qrcode.WriteColor(url, qrcode.Medium, 256, &buf, qrcode.Black, qrcode.White)
	fmt.Println("Scan QR code to open on mobile:")
	io.Copy(os.Stdout, &buf)
}

func (s *LiveServer) dashboardHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`<html><body><h2>goserver Dashboard</h2><ul><li>Active WebSocket clients: ` + fmt.Sprint(len(s.Clients)) + `</li></ul></body></html>`))
}

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil { return "localhost" }
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			return ipnet.IP.String()
		}
	}
	return "localhost"
}
