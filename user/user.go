package user

import (
	"database/sql"

	_ "github.com/lib/pq"
)

var UserArr = []User{}

type User struct {
	Id int64 `json:"id"`
}

func UpdateUserArr(db *sql.DB) {
	rows, err := db.Query("select id from users")
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var id int64
		rows.Scan(&id)
		UserArr = append(UserArr, User{Id: id})
	}
}
