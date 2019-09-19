# simple_migrate 

A simple golang migration manager, database/sql


## Getting Started

### Installing

```
   go get github.com/CezarGarrido/simple_migrate

```
## Usage

Example mysql: go run -migration=up

```go
package main

import (
	"database/sql"
	"fmt"
	"time"
	migrate "github.com/CezarGarrido/simple_migrate"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dbSource := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		"username",
		"password",
		"hostname",
		"port",
		"dbname",
	)
	db, err := sql.Open("mysql", dbSource)
	if err != nil {
		panic(err)
	}
  migrate.NewMigration(db)

  //code
}

```
## Authors
Cezar Garrido Britez  
[@CezarCgb18](https://twitter.com/CezarCgb18)

