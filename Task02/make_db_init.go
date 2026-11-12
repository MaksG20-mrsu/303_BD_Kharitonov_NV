package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	sqlFilename := "db_init.sql"
	sqlFile, err := os.Create(sqlFilename)
	if err != nil {
		fmt.Printf("Ошибка создания файла: %v\n", err)
		return
	}
	defer sqlFile.Close()

	// 1. Создание структуры таблиц (удаляем старые, создаем новые)
	sqlFile.WriteString("DROP TABLE IF EXISTS ratings;\n")
	sqlFile.WriteString("DROP TABLE IF EXISTS tags;\n")
	sqlFile.WriteString("DROP TABLE IF EXISTS movies;\n")
	sqlFile.WriteString("DROP TABLE IF EXISTS users;\n\n")

	sqlFile.WriteString(`CREATE TABLE movies (
    id INTEGER PRIMARY KEY,
    title TEXT,
    genres TEXT
);

CREATE TABLE ratings (
    id INTEGER PRIMARY KEY,
    user_id INTEGER,
    movie_id INTEGER,
    rating REAL
);

CREATE TABLE tags (
    id INTEGER PRIMARY KEY,
    user_id INTEGER,
    movie_id INTEGER,
    tag TEXT
);

CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    name TEXT
);

BEGIN TRANSACTION;
`)

	// 2. Функция для безопасного экранирования кавычек в текстах
	escape := func(s string) string {
		return strings.ReplaceAll(s, "'", "''")
	}

	// 3. Чтение movies.csv и запись INSERT-запросов
	fileMovies, err := os.Open("movies.csv")
	if err == nil {
		defer fileMovies.Close()
		reader := csv.NewReader(fileMovies)
		reader.Read() // пропускаем шапку таблицы
		for {
			rec, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil || len(rec) < 3 {
				continue
			}
			id := rec[0]
			title := escape(rec[1])
			genres := escape(rec[2])
			sqlFile.WriteString(fmt.Sprintf("INSERT INTO movies (id, title, genres) VALUES (%s, '%s', '%s');\n", id, title, genres))
		}
	}

	// Аналогично можно добавить чтение других файлов, если они есть (ratings.csv, tags.csv и т.д.)

	sqlFile.WriteString("COMMIT;\n")
	fmt.Println("Файл db_init.sql успешно сгенерирован вместе с данными!")
}