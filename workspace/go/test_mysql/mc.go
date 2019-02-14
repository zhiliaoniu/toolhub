package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func main() {
	//fmt.Printf("%05d\n", 5)

	//return
	db, err := sql.Open("mysql", "jxz_db@jxz_bd:xywaD3kfz@tcp(221.228.79.244:8066)/budao?interpolateParams=true") //interpolateParams=true
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	fmt.Println("ping Ok.")

	//my test
	r, err := db.Query("select videourl from video_0 where vid = ?", 1)
	//execSql := fmt.Sprintf("select videourl from video_0 where vid = %d", 1)
	//r, err := db.Query(execSql)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer r.Close()
	var url string
	for r.Next() {
		err := r.Scan(&url)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(url)
	}
	err = r.Err()
	if err != nil {
		log.Fatal(err)
	}
	//-----------------------------------------------------
	// Prepare statement for inserting data
	fmt.Println("======b")
	stmtIns, err := db.Prepare("select videourl from video_0 where vid = ?") // ? = placeholder
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

	// Insert square numbers for 0-24 in the database
	fmt.Println("======a")
	r, err = stmtIns.Query(1) // Insert tuples (i, i^2)
	fmt.Println("======3")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer r.Close()
	for r.Next() {
		err := r.Scan(&url)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(url)
	}
	err = r.Err()
	if err != nil {
		log.Fatal(err)
	}
	return
}
