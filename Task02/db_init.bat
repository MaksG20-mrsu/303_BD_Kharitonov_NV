#!/bin/bash
go run make_db_init.go
sqlite3 movies_rating.db ".read db_init.sql"