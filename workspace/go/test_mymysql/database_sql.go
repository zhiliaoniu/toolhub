// Example demonstrates using of mymysql/godrv driver.
package main

import (
	"database/sql"
	_ "fmt"
	"log"

	_ "github.com/ziutek/mymysql/godrv" // Go driver for database/sql package
)

func main() {
	db, err := sql.Open("mymysql", "tcp:221.228.79.244:8066*budao/jxz_db@jxz_bd/xywaD3kfz")
	//db, err := sql.Open("mymysql", "tcp:221.228.79.2:6305*budao/jxz_db@jxz_bd/YES")
	if err != nil {
		log.Fatal(err)
	}

	//id := 1
	//var query = "SELECT email from users WHERE id = ?"

	rows, err := db.Query("select * from squareNum where number = ?", 1)
	if err != nil {
		log.Fatal(err)
	}
	var (
		number int
		squareNumber int
	)
	for rows.Next() {
		err := rows.Scan(&number, &squareNumber)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(number, squareNumber)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	rows.Close()
	return;

	// Prepare statement for inserting data
	stmtIns, err := db.Prepare("INSERT INTO squareNum VALUES(?, ?)") // ? = placeholder
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates


	// Insert square numbers for 0-24 in the database
	for i := 0; i < 25; i++ {
		_, err = stmtIns.Exec(i, (i * i)) // Insert tuples (i, i^2)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
	}



	return

	/*
	var email string
	for rows.Next() {
		if err := rows.Scan(&email); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Email address: %s\n", email)

	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	*/
}
