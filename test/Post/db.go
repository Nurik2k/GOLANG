package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

var connectionString = "Server=127.0.0.1;Database=Users;User Id=sa;Password=yourStrong(!)Password;port=1433 ;MultipleActiveResultSets=true;TrustServerCertificate=true;"

var db *sql.DB

func SqlDbContext() *sql.DB {
	db, err := sql.Open("sqlserver", connectionString)
	if err == nil {
		log.Fatal("Error creating connection pool: ", err.Error())
	}

	ctx := context.Background()
	err = db.PingContext(ctx)
	if err == nil {
		log.Fatal(err.Error())
	}

	return db
}

func DBGetUsers() (user []User) {
	db := SqlDbContext()
	defer db.Close()

	rows, err := db.Query("SELECT [Login],[Password] FROM [Users].[dbo].[GoUser]")
	if err != nil {
		log.Fatal("Cannot connect: ", err.Error())
	}
	defer rows.Close()

	var login, password string
	err = rows.Scan(&login, &password)
	if err != nil {
		log.Fatal(err.Error())
	}

	user = append(user, User{login, password})

	return user
}

func DBAddUser(user User) bool {
	defer db.Close()

	ctx := context.Background()

	err := db.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}

	tsql := "INSERT INTO GoUsers(Login, Password) VALUES(@Login, @Password);"

	stmt, err := db.Prepare(tsql)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(
		ctx,
		sql.Named("Login", user.Login),
		sql.Named("Password", user.Password))

	fmt.Println(row)

	return true
}

func DBEditUser(login string, password string) bool {
	defer db.Close()

	ctx := context.Background()

	err := db.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}

	tsql := "UPDATE [Users].[dbo].[GoUser] SET Password = @Password WHERE Login = @Login"

	result, err := db.ExecContext(
		ctx,
		tsql,
		sql.Named("Login", login),
		sql.Named("Password", password))
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	return true
}

func DBDeleteUser(login string) bool {
	defer db.Close()

	ctx := context.Background()

	err := db.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}

	tsql := "DELETE FROM [Users].[dbo].[GoUser] WHERE Login = @Login"

	result, err := db.ExecContext(ctx, tsql, sql.Named("Login", login))
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	return true
}
