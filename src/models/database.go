package models

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func DB() {
	// Указываем параметры подключения к базе данных
	var conninfo string = "user=postgres password=postgres dbname=postgres sslmode=disable host=db"
	db, err := sql.Open("postgres", conninfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	dbName := "images"
	// Проверяем существование базы данных
	var exists bool
	err = db.QueryRow("SELECT EXISTS (SELECT FROM pg_database WHERE datname = $1)", dbName).Scan(&exists)
	if err != nil {
		log.Fatal(err)
	}
	if !exists {
		// Создаем базу данных, если она еще не существует
		_, err = db.Exec("create database " + dbName)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Подключаемся к базе данных images
	conninfo = "user=postgres password=postgres dbname=images sslmode=disable host=db"
	db, err = sql.Open("postgres", conninfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Создаем таблицу, если ее еще нет
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS images (ID TEXT, Path_or TEXT, Path_neg TEXT, Created_at TIMESTAMP)")
	if err != nil {
		panic(err)
	}

	fmt.Println("Table created successfully")
}

func DBInsert(ID string, filename string, negativeFilename string) error {
	// Указываем параметры подключения к базе данных
	var conninfo string = "user=postgres password=postgres dbname=images sslmode=disable host=db"
	db, err := sql.Open("postgres", conninfo)
	if err != nil {
		return err
	}
	defer db.Close()

	// Добавляем данные в таблицу
	_, err = db.Exec("INSERT INTO images (ID, Path_or, Path_neg, Created_at) VALUES ($1, $2, $3, $4)", ID, "uploads/"+filename, "uploads/"+negativeFilename, time.Now())
	if err != nil {
		return err
	}

	return nil
}
