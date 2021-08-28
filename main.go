package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

const (
  host     = "localhost"
  user     = "postgres"
  port = 5432
  password = "postgres"
  dbname   = "postgres"
)

func uploadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Uploading file \n")
	
	// parsing the input, type multipart/form-data
	r.ParseMultipartForm(10 << 20)

	//retriving files from posted upload 
	file, handler, err := r.FormFile("myFile")

	if err != nil {
		fmt.Println("Error reading the file")
		fmt.Println(err)
		return
	}
	defer file.Close() 

	c, _ := os.Create(handler.Filename)

	// Copy the uploaded file to the created file on the filesystem
	if _, err := io.Copy(c, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	f, _ := os.Open(handler.Filename)
	fileScanner := bufio.NewScanner(f)
	lineCount := 0
	for fileScanner.Scan() {
		lineCount++
	}
	
	fmt.Fprintf(w, "File Size: %+v\n", handler.Size/1024)
	fmt.Fprintf(w, "MIME Header: %+v\n", handler.Header)
	fmt.Fprintf(w, "File Name: %+v\n", handler.Filename)
	fmt.Fprintf(w, "Total Number of lines: %+v\n", lineCount)
	x := handler.Filename
	y := handler.Size/1024

	insert_data_to_db(x,y,lineCount)
}

func setupRoutes(){
	http.HandleFunc("/upload", uploadFile)
	http.ListenAndServe(":8080", nil)
}

func insert_data_to_db(fName string, fSize int64, lCount int){
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
    "password=%s dbname=%s sslmode=disable",
    host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
  	panic(err)	
	}
	defer db.Close()

	sqlstatemnet := `INSERT INTO TEST_1234 (FILE_NAME,FILE_SIZE,TOTAL_NUMBER_LINES)
	VALUES ($1,$2,$3)
	RETURNING id`

	id := 0
	err = db.QueryRow(sqlstatemnet, fName, fSize, lCount).Scan(&id)

	if err != nil {
		panic(err)	
	}
	fmt.Println("New record ID is:", id)
}


// func create_table(){
// 	sqlQuery := `CREATE TABLE test_1234 (
// 		id SERIAL NOT NULL PRIMARY KEY,
// 		file_name TEXT,
// 		file_size INT,
//		total_number_lines INT,
// 		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
// 	  );`
// }


func init() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
    "password=%s dbname=%s sslmode=disable",
    host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
  	panic(err)	
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected to PostgreSQL!")
}

func main(){
	fmt.Println("Go file upload file")
	setupRoutes()
}