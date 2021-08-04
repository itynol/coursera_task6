package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// тут вы пишете код
// обращаю ваше внимание - в этом задании запрещены глобальные переменные

type HandlersSuite struct {
	DB *sql.DB
}

func (h *HandlersSuite) selectFromTable (w http.ResponseWriter, r *http.Request, tableName string) {
	rows, err := h.DB.Query(fmt.Sprintf("SELECT * FROM %s;", tableName))
	if err != nil || rows == nil {
		if rows == nil {
			log.Println("No rows in result set")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	col, _ := rows.Columns()
	colNum := len(col)
	items := make([]interface{}, colNum)
	for i := 0; i < colNum; i++ {
		items[i] = new(sql.RawBytes)
	}
	for rows.Next() {
		var str []string
		err = rows.Scan(items...)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		for i := 0; i < colNum; i++ {
			irr := items[i].(*sql.RawBytes)
			str = append(str, string(*irr))
		}
		io.WriteString(w, strings.Join(str, " ") + "\n")
	}
	return
}

func (h *HandlersSuite) mainPage(w http.ResponseWriter, r *http.Request) {
	tableName := r.URL.Path[1:]
	rows := new(sql.Rows)
	var err error
	if tableName != "" {
		h.selectFromTable(w, r, tableName)
	} else {
		rows, err = h.DB.Query("SHOW TABLES;")
	}
	if err != nil || rows == nil {
		if rows == nil {
			log.Println("No rows in result set")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	for rows.Next() {
		var res string
		err = rows.Scan(&res)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		io.WriteString(w, res + "\n")
	}
	w.WriteHeader(http.StatusOK)
}

func NewDbExplorer(db *sql.DB) (http.Handler, error) {
	adm := http.NewServeMux()
	fmt.Println("NewDbExplorer")
	hS := &HandlersSuite{
		DB: db,
	}
	adm.HandleFunc("/", hS.mainPage)
	return adm, nil
}
