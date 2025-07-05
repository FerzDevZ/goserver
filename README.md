# goserver

CLI tool untuk membangun dan menyajikan aplikasi web statis (HTML/CSS/JS) dengan live reload, watcher, WebSocket, QR code, dashboard, dan build system. Struktur proyek mirip Flutter CLI.

## Fitur Utama
- Static server (lib/), auto-inject reload.js
- Live reload via WebSocket
- File watcher (lib/, assets/)
- Build system (custom script)
- Dashboard dev info
- QR code URL lokal
- Modular CLI: create, serve, build, doctor, config
- Konfigurasi via goserver.yaml
- Mode production (build/ tanpa reload.js/ws)

## Instalasi
```
go install github.com/FerzDevZ/goserver@latest
```

## Perintah CLI
- `goserver create myapp` — generate struktur proyek
- `goserver serve` — jalankan server dengan live reload
- `goserver build` — jalankan build script
- `goserver doctor` — cek dependency & struktur
- `goserver config` — baca/tulis goserver.yaml

## Struktur Proyek
```
myapp/
├── goserver.yaml
├── assets/
├── build/
├── lib/
│   ├── index.html
│   ├── style.css
│   └── main.js
├── reload/
│   └── reload.js
└── .goserver/
```

## License
MIT
# goserver
