package sqlconnect

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDb(dbname string) (*sql.DB, error) {
	fmt.Println("Connecting to MariaDB...")

	connectionString := "root:zzy15198@tcp(127.0.0.1:3306)/" + dbname
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		// panic(err)
		return nil, err
	}
	fmt.Println("Connected to MariaDB successfully")
	return db, nil
}
