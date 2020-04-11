package models

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

// sqlite数据库
type TxtView struct {
	Id         int64
	Title      string
	Content    string
	Del        bool
	UnLockTime int64
	Lock       string
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
		var id, unLockTime int64
		var title, content, lock string
		var del bool
		err = rows.Scan(&id, &title, &content, &del, &unLockTime)
		if err != nil {
			return nil, err
		}
		if time.Now().Unix() > unLockTime {
			lock = "未锁定"
		} else {
			lock = "锁定"
		}
		todo := TxtView{id, title, content, del, unLockTime, lock}
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

func GetTxtView(todoId int64) (string, string, int64, error) {
	db, err := sql.Open("sqlite3", "./db/txtview.db")
	defer db.Close()
	if err != nil {
		return "", "", 0, err
	}
	// 只查询一行数据
	row := db.QueryRow("SELECT title, content, unLockTime FROM txtview WHERE id=?", todoId)
	var title, content string
	var unLockTime int64
	e := row.Scan(&title, &content, &unLockTime)
	if e != nil {
		return "", "", 0, e
	}
	return title, content, unLockTime, nil
}

func IsExistTxtViewByTitle(titleParam string) (bool, error) {
	db, err := sql.Open("sqlite3", "./db/txtview.db")
	defer db.Close()
	if err != nil {
		return true, err
	}
	// 只查询一行数据
	row := db.QueryRow("SELECT title FROM txtview WHERE title=? and del=0 LIMIT 1", titleParam)
	var title string
	e := row.Scan(&title)
	if e != nil {
		if e.Error() == sql.ErrNoRows.Error() {
			return false, nil
		}
		return true, e
	}
	if title != "" {
		return true, nil
	}
	return false, nil
}

func UpdateTxtViewUnLockTime(txtViewId, unLockTime int64) (int64, error) {
	db, err := sql.Open("sqlite3", "./db/txtview.db")
	defer db.Close()
	if err != nil {
		return 0, err
	}
	stmt, err := db.Prepare("UPDATE txtview SET unlocktime=? WHERE id=?")
	defer stmt.Close()
	if err != nil {
		return 0, nil
	}

	res, err := stmt.Exec(unLockTime, txtViewId)
	if err != nil {
		return 0, nil
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return 0, nil
	}
	return affect, nil
}

func EditTodo(id, unLockTime int64, title, content string) (int64, error) {
	db, err := sql.Open("sqlite3", "./db/txtview.db")
	defer db.Close()
	if err != nil {
		return 0, err
	}
	stmt, err := db.Prepare("UPDATE txtview SET title=?, content=?, unlocktime=? WHERE id=?")
	defer stmt.Close()
	if err != nil {
		return 0, nil
	}

	res, err := stmt.Exec(title, content, unLockTime, id)
	if err != nil {
		return 0, nil
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return 0, nil
	}
	return affect, nil
}
