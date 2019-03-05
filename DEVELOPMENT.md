## Développement

### Installation

Prérequis :
- [go](https://golang.org/) (testé avec 1.11)
- [dep](https://github.com/golang/dep)
- [mockery](https://github.com/vektra/mockery)
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
- domain/ : services métier
- models/ : entités métier
- app/ : config et bootstrap application
- tests/e2e/ : tests end-to-end
- vendor/ : contient les dépendances
- main.go : point d'entrée


### Mise à jour des mocks

```
go generate ./...
```

### Tests

Lancement d'un test :
```
go test -p 1 -v tests/e2e/inspections_test.go
```

Lancement de tous les tests :
```
go test -p 1 -v ./...
```

- `DEBUG_HTTP=1` permet d'afficher les requêtes et réponses HTTP et leur body (JSON).
- `FILHARMONIC_DATABASE_SQL_LOG=1` permet d'afficher les requêtes SQL.


### Benchmarks

On utilise [k6](https://k6.io/), qui permet d'écrire des scénarios en javascript.

Pour lancer un test pendant 10 secondes avec 4 utilisateurs en parallèle.
```sh
docker run -it --rm -v $PWD/tests/benchmark/:/benchmark --network host loadimpact/k6 run -u 4 -d 10s /benchmark/creation_inspection.js
```

Pour afficher les requêtes HTTP et leur contenu :
```sh
docker run -it --rm -v $PWD/tests/benchmark/:/benchmark --network host loadimpact/k6 run --http-debug=full /benchmark/creation_inspection.js
```


### Modifications du modèle

Pour vérifier que le schéma des migrations correpond bien à celui des structures, lancer le script suivant :
```
./database/scripts/check_migrations.sh
```

Afin que le script de vérification fonctionne, il faut veiller à ce que chaque nouveau champ ajouté à un modèle soit défini après tous les autres champs car PostgreSQL ne permet pas de [changer la position des colonne](https://wiki.postgresql.org/wiki/Alter_column_position) sans recréer la table.
