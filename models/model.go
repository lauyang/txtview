package models

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

// sqlite数据库
type TxtView struct {
	Id      int64
	Title   string
	Content string
	Del     bool
}

func InsertTxtView(title, content string) (int64, error) {
	// 数据库path是相对路径（相对于main.go）
	db, err := sql.Open("sqlite3", "./db/txtview.db")
	// 函数代码执行完后关闭数据库
	defer db.Close()
	if err != nil {
		return -1, err
	}

	stmt, err := db.Prepare("INSERT INTO txtview(title, content, del) VALUES(?, ?, ?)")
	defer stmt.Close()
	if err != nil {
		return -1, err
	}

	res, err := stmt.Exec(title, content, false)
	if err != nil {
		return -1, err
	}
	return res.LastInsertId()
}

func QueryAll() ([]TxtView, error) {
	db, err := sql.Open("sqlite3", "./db/txtview.db")
	defer db.Close()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("SELECT * FROM txtview WHERE del = FALSE")
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	var txtViews []TxtView

	for rows.Next() {
		var id int64
		var title string
		var content string
		var del bool
		err = rows.Scan(&id, &title, &content, &del)
		if err != nil {
			return nil, err
		}
		todo := TxtView{id, title, content, del}
		txtViews = append(txtViews, todo)
	}
	return txtViews, nil
}

func DelTxtView(textViewId int64, del bool) (int64, error) {
	db, err := sql.Open("sqlite3", "./db/txtview.db")
	defer db.Close()
	if err != nil {
		return 0, err
	}
	stmt, err := db.Prepare("UPDATE txtview SET del=? WHERE id=?")
	defer stmt.Close()
	if err != nil {
		return 0, nil
	}

	res, err := stmt.Exec(del, textViewId)
	if err != nil {
		return 0, nil
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return 0, nil
	}
	return affect, nil
}

func GetTxtView(todoId int64) (string, string, error) {
	db, err := sql.Open("sqlite3", "./db/txtview.db")
	defer db.Close()
	if err != nil {
		return "", "", err
	}
	// 只查询一行数据
	row := db.QueryRow("SELECT title, content FROM txtview WHERE id=?", todoId)
	var title, content string
	e := row.Scan(&title, &content)
	if e != nil {
		return "", "", e
	}
	return title, content, nil
}

func EditTodo(id int64, title string) (int64, error) {
	db, err := sql.Open("sqlite3", "./db/txtview.db")
	defer db.Close()
	if err != nil {
		return 0, err
	}
	stmt, err := db.Prepare("UPDATE txtview SET title=? WHERE id=?")
	defer stmt.Close()
	if err != nil {
		return 0, nil
	}

	res, err := stmt.Exec(title, id)
	if err != nil {
		return 0, nil
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return 0, nil
	}
	return affect, nil
}
