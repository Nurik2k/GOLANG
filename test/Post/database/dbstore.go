package database

import (
	"context"
	"database/sql"
	"fmt"
)

type User struct {
	Id        string `json:"id"`
	Login     string `json:"login"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	Name      string `json:"name"`
	LastName  string `json:"last_name"`
	Birthday  string `json:"birthday"`
}

type dbStore struct {
	db *sql.DB
}

type DbInterface interface {
	SignIn(login, password string) (isSigned bool, err error)
	Get(limit, offset int) ([]User, error)
	GetById(id string) (user User, err error)
	Create(user *User) error
	Edit(user *User) error
	Delete(id string) error
}

func (db *dbStore) SignIn(login, password string) (isSigned bool, err error) {
	ctx := context.Background()

	err = db.db.PingContext(ctx)
	if err != nil {
		return false, err
	}

	tsql := "SELECT Password FROM GoUser WHERE Login = @Login"

	row := db.db.QueryRowContext(ctx, tsql, sql.Named("Login", login))
	if err != nil {
		return false, err
	}

	var Password string

	err = row.Scan(&Password)
	if err != nil {
		return false, err
	}

	if err := row.Err(); err != nil {
		return false, err
	}

	if Password != password {
		return false, nil
	}

	return true, nil
}

func (db *dbStore) Get(limit, offset int) ([]User, error) {
	ctx := context.Background()

	//limit offset
	tsql := "SELECT [Id], [First_Name], [Name], [Last_Name], [Birthday], [Login], [Password] FROM [Users].[dbo].[GoUser] " +
		"ORDER BY Id OFFSET @Offset ROWS FETCH NEXT @Limit ROWS ONLY"

	rows, err := db.db.QueryContext(ctx, tsql, sql.Named("Offset", offset), sql.Named("Limit", limit))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//make
	users := []User{}

	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.FirstName, &user.Name, &user.LastName, &user.Birthday, &user.Login, &user.Password)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (db *dbStore) GetById(id string) (user User, err error) {
	//var t string
	//var u User
	//var u2 *User // == []int map[string]int pointer = nil

	ctx := context.Background()

	//user = new(User)
	//user = &User{}

	err = db.db.PingContext(ctx)
	if err != nil {
		return user, err
	}

	if db.db == nil {
		return user, err
	}

	tsql := "SELECT Id, Login, Password, First_Name, Name, Last_Name, Birthday FROM GoUser WHERE Id = @Id"

	//Посмотреть возвращение только одной операции row
	row := db.db.QueryRowContext(ctx, tsql, sql.Named("Id", id))

	err = row.Scan(&user.Id, &user.Login, &user.Password, &user.FirstName, &user.Name, &user.LastName, &user.Birthday)
	if err != nil {
		return user, err
	}

	return user, nil
}

// return err
func (db *dbStore) Create(user *User) error {
	ctx := context.Background()

	err := db.db.PingContext(ctx)
	if err != nil {
		return err
	}

	if db.db == nil {
		return err
	}

	//Запрос return @Id. user.Id = @Id
	tsql := "INSERT INTO GoUser(Login, Password, First_Name, Name, Last_Name, Birthday) OUTPUT INSERTED.Id VALUES(@Login, @Password, @First_Name, @Name, @Last_Name, @Birthday);"

	stmt, err := db.db.Prepare(tsql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(
		ctx,
		sql.Named("Login", user.Login),
		sql.Named("Password", user.Password),
		sql.Named("First_Name", user.FirstName),
		sql.Named("Name", user.Name),
		sql.Named("Last_Name", user.LastName),
		sql.Named("Birthday", user.Birthday),
	).Scan(&user.Id)

	if err != nil {
		return err
	}

	return nil
}

func (db *dbStore) Edit(user *User) error {
	ctx := context.Background()

	err := db.db.PingContext(ctx)
	if err != nil {
		return err
	}

	tsql := fmt.Sprintf("UPDATE GoUser SET Login = @Login, Password = @Password, First_Name = @First_Name, Name = @Name, Last_name = @Last_name, Birthday = @Birthday WHERE Id = @Id")

	stmt, err := db.db.Prepare(tsql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		sql.Named("Id", user.Id),
		sql.Named("Login", user.Login),
		sql.Named("Password", user.Password),
		sql.Named("First_Name", user.FirstName),
		sql.Named("Name", user.Name),
		sql.Named("Last_Name", user.LastName),
		sql.Named("Birthday", user.Birthday),
	)
	if err != nil {
		return err
	}

	return nil
}

func (db *dbStore) Delete(id string) error {
	ctx := context.Background()

	err := db.db.PingContext(ctx)
	if err != nil {
		return err
	}

	tsql := fmt.Sprintf("DELETE FROM GoUser WHERE Id = @Id;")

	_, err = db.db.ExecContext(ctx, tsql, sql.Named("Id", id))
	if err != nil {
		return err
	}

	return nil
}

func NewDb() (DbInterface, error) {
	db, err := sql.Open("sqlserver", "Server=localhost;Database=Users;User Id=sa;Password=yourStrong(!)Password;port=1433;MultipleActiveResultSets=true;TrustServerCertificate=true;")
	if err != nil {
		return nil, err
	}

	store := &dbStore{
		db: db,
	}

	err = store.db.Ping()
	if err != nil {
		return nil, err
	}

	return store, nil
}
