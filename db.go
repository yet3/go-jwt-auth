package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	instance *sql.DB
}

func NewDB() *DB {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
    PRAGMA journal_mode = WAL;
    PRAGMA foreign_keys = ON;
    
    CREATE TABLE IF NOT EXISTS users (
      id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
      email TEXT NOT NULL,
      password TEXT NOT NULL,

      accessToken TEXT,
      refreshToken TEXT,

      UNIQUE(email)
    );

    CREATE TABLE IF NOT EXISTS user_permissions (
      id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
      userId INTEGER NOT NULL,

      name TEXT NOT NULL,
      value TEXT NOT NULL,

      FOREIGN KEY (userId) REFERENCES users (id)
    );
  `)

	if err != nil {
		log.Fatal(err)
	}

	return &DB{
		instance: db,
	}
}

func (db *DB) Close() {
	if db.instance == nil {
		return
	}

	db.instance.Close()
}
