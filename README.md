# fiber-boilerplate
A simple boilerplate for [Fiber](https://github.com/gofiber/fiber)

[![Go Report Card](https://goreportcard.com/badge/github.com/fabienbellanger/fiber-boilerplate)](https://goreportcard.com/report/github.com/fabienbellanger/fiber-boilerplate)
[![MIT license](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)

## Sommaire
-  [Makefile commands](#Makefile-commands)
-  [Routes](#Routes)
    -  [Web](#Web)
    -  [API](#API)
-  [API documentation](#API-documentation)
    -  [Swagger](#Swagger)
-  [Golang web server in production](#golang-web-server-in-production)
-  [Mesure et performance](#mesure-et-performance)
    -  [pprof](#pprof)
    -  [trace](#trace)
    -  [cover](#cover)
-  [TODO](#TODO)


## Makefile commands

| Makefile command | Go command | Description |
|---|---|---|
| `make update` | `go get -u && go mod tidy` | Update Go dependencies |
| `make serve` | `go run cmd/main.go` | Start the Web server |
| `make serve-race` | `go run --race cmd/main.go` | Start the Web server with data races option |
| `make serve-pkger` | `pkger && go run cmd/main.go` | Run Pkger and start the Web server |
| `make build` | `go build -o fiber-boilerplate -v cmd/main.go` | Build application with pkger |
| `make test` | `go test -cover -v ./...` | Launch unit tests |


## Routes
[Liste des routes](ROUTES.md)


## API documentation

### Swagger
I must install [Swag](https://github.com/swaggo/swag/cmd/swag):
```bash
go get -u github.com/swaggo/swag/cmd/swag
```
Then run:
```bash
swag init -g cmd/main.go
```
The documentation is available in your browser at `http://localhost:3000/swagger/index.html`.


## Golang web server in production
-  [Systemd](https://jonathanmh.com/deploying-go-apps-systemd-10-minutes-without-docker/)
-  [ProxyPass](https://evanbyrne.com/blog/go-production-server-ubuntu-nginx)
-  [How to Deploy App Using Docker](https://medium.com/@habibridho/docker-as-deployment-tools-5a6de294a5ff)

### Creating a Service for Systemd
```bash
touch /lib/systemd/system/<service name>.service
```

Edit file:
```
[Unit]
Description=<service description>
After=network.target

[Service]
Type=simple
Restart=always
RestartSec=5s
ExecStart=<path to exec with arguments>

[Install]
WantedBy=multi-user.target
```

| Commande | Description |
|---|---|
| `systemctl start <service name>.service` | To launch |
| `systemctl enable <service name>.service` | To enable on boot |
| `systemctl disable <service name>.service` | To disable on boot |
| `systemctl status <service name>.service` | To show status |
| `systemctl stop <service name>.service` | To stop |


## Benchmark
Use [Drill](https://github.com/fcsonline/drill)
```bash
$ drill --benchmark drill.yml --stats --quiet
```


## Mesure et performance
Go met à disposition de puissants outils pour mesurer les performances des programmes :
-  pprof (graph, flamegraph, peek)
-  trace
-  cover

=> Lien vers une vidéo intéressante [Mesure et optimisation de la performance en Go](https://www.youtube.com/watch?v=jd47gDK-yDc)

### pprof
Lancer :
```bash
curl http://localhost:8888/debug/pprof/heap?seconds=10 > <fichier à analyser>
```
Puis :
```bash
go tool pprof -http :7000 <fichier à analyser> # Interface web
go tool pprof --nodefraction=0 -http :7000 <fichier à analyser> # Interface web avec tous les noeuds
go tool pprof <fichier à analyser> # Ligne de commande
```

### trace
Lancer :
```bash
go test <package path> -trace=<fichier à analyser>
curl localhost:<port>/debug/pprof/trace?seconds=10 > <fichier à analyser>
```
Puis :
```bash
go tool trace <fichier à analyser>
```

### cover
Lancer :
```bash
go test <package path> -covermode=count -coverprofile=./<fichier à analyser>
```
Puis :
```bash
go tool cover -html=<fichier à analyser>
```

## TODO
-  [x] Utiliser Zap
-  [ ] **Attention** : Le middleware websocket de Fiber génère une data race avec le hub ! Voir si cela sera corrigé à l'avenir (lever une issue sur Github ?)
-  [ ] Mettre en place la stack Prometheus + Grafana pour la télémétrie
