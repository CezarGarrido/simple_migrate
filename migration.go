package simple_migrate

import (
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Migration struct {
	Id          int64
	Description string
	Created_at  *time.Time
}

func NewMigration(db *sql.DB) {
	comandos := flag.String("migration", "", "")
	createCommand := flag.String("migration:create", "", "")
	flag.Parse()

	migrate := Migration{}
	switch *comandos {
	case "init":
		InitTable(db)
	case "up":
		migrate.MigrationUp(db)
	case "down":
		migrate.MigrationDown(db)
	case "list":
		migrate.MigrationList(db)
	}
	if createCommand != nil {
		fmt.Println("Creating migration", *createCommand)
		err := createFileMigration(*createCommand)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}

func InitTable(db *sql.DB) error {
	sql := `CREATE TABLE IF NOT EXISTS migrations ( 
		id INT(11) NOT NULL AUTO_INCREMENT, 
		description TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    PRIMARY KEY (id)
		) ENGINE = InnoDB;`

	if _, err := db.Exec(sql); err != nil {
		return err
	}
	return nil
}

func createFileMigration(name string) error {
	f, err := os.Create("./migrations/up/" + time.Now().Format("20060102150405") + "_" + name + ".up.sql")
	if err != nil {
		return err
	}
	defer f.Close()

	return nil
}

func (this Migration) UpFiles(dir string) (files []string, err error) {
	files, err = filepath.Glob(filepath.Join(dir, "*.up.sql"))
	return
}

func (this Migration) DownFiles(dir string) (files []string, err error) {
	files, err = filepath.Glob(filepath.Join(dir, "*.down.sql"))
	return
}

func (this Migration) MigrationUp(db *sql.DB) {
	files, _ := this.UpFiles("./migrations/up/")
	fmt.Println("Performing UP migrations.")
	for _, n := range files {
		fmt.Printf("%s\n", n)
		if checkExist(db, n) {
			fmt.Println("Migration already exists.")
			continue
		}
		query, err := ioutil.ReadFile(n)
		if err != nil {
			panic(err)
		}
		stringSlice := strings.Split(string(query), ";")
		fmt.Println(stringSlice)
		for _, comando := range stringSlice {
			comando = strings.TrimSpace(comando)
			comando = strings.Trim(comando, " ")
			if len(comando) == 0 {
				continue
			}
			if _, err := db.Exec(comando); err != nil {
				panic(err.Error())
			}

		}
		registerMigration(n, db)
	}
	fmt.Println("Migrations UPs executed successfully.")
}

func (this Migration) MigrationDown(db *sql.DB) {
	files, _ := this.DownFiles("./migrations/down/")
	fmt.Println("Performing DOWN migrations.")
	for _, n := range files {
		fmt.Printf("%s", n)
		if checkExist(db, n) {
			fmt.Println("Migration already exists.")
			continue
		}
		query, err := ioutil.ReadFile(n)
		if err != nil {
			panic(err)
		}
		stringSlice := strings.Split(string(query), "")
		fmt.Println(stringSlice)
		for _, comand := range stringSlice {

			comand = strings.TrimSpace(comand)
			comand = strings.Trim(comand, " ")
			if len(comand) == 0 {
				continue
			}
			if _, err := db.Exec(string(query)); err != nil {
				panic(err.Error())
			}
		}
		registerMigration(n, db)
	}
	fmt.Println("Migrations DOWNs executed successfully.")
}

func registerMigration(migration string, db *sql.DB) {
	data := regexp.MustCompile(`\d{4}\d{2}\d{2}\d{2}\d{2}\d{2}`)
	submatchall := data.FindAllString(migration, -1)
	re := regexp.MustCompile("_(.*?).up")
	descricaoMigration := re.FindStringSubmatch(migration)
	for _, element := range submatchall {
		dataRegistro := str2Date(element)
		db, err := db.Prepare("INSERT INTO migrations(description, created_at) VALUES(?,?)")
		if err != nil {
			panic(err.Error())
		}
		if _, err := db.Exec(descricaoMigration[1], dataRegistro); err != nil {
			panic(err.Error())
		}
	}
}

func checkExist(db *sql.DB, migration string) bool {
	re := regexp.MustCompile("_(.*?).up")
	descricaoMigration := re.FindStringSubmatch(migration)
	var count bool
	_ = db.QueryRow("SELECT (SELECT COUNT(*) FROM migrations WHERE description = ?) > 0", descricaoMigration[1]).Scan(&count)

	return count
}

func (this Migration) MigrationList(db *sql.DB) {

	fmt.Println("Listing Migrations")
	results, err := db.Query("SELECT id, description, created_at FROM migrations")
	if err != nil {
		panic(err.Error())
	}
	var migration Migration
	for results.Next() {
		err = results.Scan(&migration.Id, &migration.Description, &migration.Created_at)
		if err != nil {
			panic(err.Error())
		}
		const padding = 3
		w := tabwriter.NewWriter(os.Stdout, 22, 0, padding, '-', tabwriter.Debug)
		fmt.Fprintln(w, strconv.FormatInt(migration.Id, 10)+"\t"+migration.Description+"\t"+migration.Created_at.Format("2006/01/06"))
		w.Flush()
	}
}

func str2Date(data string) (ret time.Time) {
	ret, err := time.Parse("20060102150405", data)
	if err != nil {
		fmt.Println(err.Error())
	}
	return
}
