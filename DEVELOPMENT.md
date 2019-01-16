## Développement

### Installation

Prérequis :
- [go](https://golang.org/) (testé avec 1.11)
- [dep](https://github.com/golang/dep)
- [docker-compose](https://docs.docker.com/compose/install/) (testé avec 1.19.0)

Il faut d'abord installer les dépendances.
```
dep ensure
```

Puis démarrer la base de données.
```
docker-compose up -d
```

Enfin lancer l'application.
```
go run main.go
```

### Organisation

- database/ : couche d'accès à la base de données
- httpserver/ : serveur HTTP
- service/ : services métiers
- tests/e2e/ : tests end-to-end
- vendor/ : contient les dépendances
- main.go : point d'entrée
