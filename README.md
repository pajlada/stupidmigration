# stupidmigration

a stupid migration tool for Go, aimed at postgresql. only up, no down

here's a handy bash script to generate a migration
```bash
#!/bin/bash

RELATIVE_PATH=$(dirname $0)

if [ -z $1 ]; then
    echo "usage: $0 description-of-migration"
    exit 1
fi

touch "$RELATIVE_PATH/$(date --utc '+%s')-$1.sql"
```

## usage

```go
func main() {
	db, err := sql.Open("postgres", "postgres:///database_name?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	err = migrate.Migrate("../../migrations", db)
	if err != nil {
		log.Fatal("Error running migrations:", err)
	}
}
```
