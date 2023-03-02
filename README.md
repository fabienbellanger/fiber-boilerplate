# fiber-boilerplate

A simple boilerplate for [Fiber](https://github.com/gofiber/fiber)

[![Go Report Card](https://goreportcard.com/badge/github.com/fabienbellanger/fiber-boilerplate)](https://goreportcard.com/report/github.com/fabienbellanger/fiber-boilerplate)
[![MIT license](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=square)](https://pkg.go.dev/github.com/fabienbellanger/fiber-boilerplate)

## Sommaire

- [Commands list](#commands-list)
- [Makefile commands](#makefile-commands)
- [Routes](#routes)
- [Swagger](#swagger)
- [Golang web server in production](#golang-web-server-in-production)
- [Go documentation](#go-documentation)
- [Mesure et performance](#mesure-et-performance)
  - [pprof](#pprof)
  - [trace](#trace)
  - [cover](#cover)
- [TODO](#todo)

## Commands list

| Command               | Description                 |
| --------------------- | --------------------------- |
| `<binary> run`        | Start server                |
| `<binary> logs -s`    | Server logs reader          |
| `<binary> logs -d`    | Database (GORM) logs reader |
| `<binary> register`   | Create a new user           |

## Makefile commands

| Makefile command    | Go command                                    | Description                                 |
| ------------------- | --------------------------------------------- | ------------------------------------------- |
| `make update`       | `go get -u && go mod tidy`                    | Update Go dependencies                      |
| `make serve`        | `go run cmd/main.go`                          | Start the Web server                        |
| `make serve-race`   | `go run --race cmd/main.go`                   | Start the Web server with data races option |
| `make build`        | `go build -o go-url-shortener -v cmd/main.go` | Build application                           |
| `make test`         | `go test -cover ./...`                        | Launch unit tests                           |
| `make test-verbose` | `go test -cover -v ./...`                     | Launch unit tests in verbose mode           |
| `make logs`         | `go run cmd/main.go logs -s`                  | Start server logs reader                    |

## Routes

[Liste des routes](ROUTES.md)

## Swagger

TODO

## Golang web server in production

- [Systemd](https://jonathanmh.com/deploying-go-apps-systemd-10-minutes-without-docker/)
- [ProxyPass](https://evanbyrne.com/blog/go-production-server-ubuntu-nginx)
- [How to Deploy App Using Docker](https://medium.com/@habibridho/docker-as-deployment-tools-5a6de294a5ff)

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

| Commande                                   | Description        |
| ------------------------------------------ | ------------------ |
| `systemctl start <service name>.service`   | To launch          |
| `systemctl enable <service name>.service`  | To enable on boot  |
| `systemctl disable <service name>.service` | To disable on boot |
| `systemctl status <service name>.service`  | To show status     |
| `systemctl stop <service name>.service`    | To stop            |

## Benchmark

Use [Drill](https://github.com/fcsonline/drill)

```bash
$ drill --benchmark drill.yml --stats --quiet
```

## Go documentation

Installer `godoc` (pas dans le répertoire du projet) :

```bash
go get -u golang.org/x/tools/...
```

Puis lancer :

```bash
godoc -http=localhost:6060 -play=true -index
```

## Mesure et performance

Go met à disposition de puissants outils pour mesurer les performances des programmes :

- pprof (graph, flamegraph, peek)
- trace
- cover

=> Lien vers une vidéo intéressante [Mesure et optimisation de la performance en Go](https://www.youtube.com/watch?v=jd47gDK-yDc)

### pprof

Lancer :

```bash
curl http://localhost:<port>/debug/pprof/heap?seconds=10 > <fichier à analyser>
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

- [x] Utiliser Zap
- [x] Add OpenAPI documentation [RapiDoc](https://rapidocweb.com/examples.html)
- [x] Validation des données avec github.com/go-playground/validator
- [x] Si la connexion à la base de données est coupée, cela retourne une 401 au lieu d'une 500.
- [x] Supprimer Pkger (plus maintenu)
- [ ] **Attention** : Le middleware websocket de Fiber génère une data race avec le hub ! Voir si cela sera corrigé à l'avenir (lever une issue sur Github ?)
- [ ] Add Docker support
  - [ ] Try OpenTelemetry [middleware](https://github.com/gofiber/contrib/tree/main/otelfiber)
  - [ ] Mettre en place la stack Prometheus + Grafana pour la télémétrie
  - [ ] Add Prometheus metrics ([Example](https://github.com/stefanprodan/dockprom))
  - [ ] Create a first user to use API
- [ ] Try test suite [ginkgo](https://github.com/onsi/ginkgo) and [gomega](https://github.com/onsi/gomega)
