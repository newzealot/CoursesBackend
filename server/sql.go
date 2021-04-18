package main

import (
	"log"
	"net/http"
)

// ExecQuery helps process DB.Exec
func ExecQuery(w http.ResponseWriter, q string, s ...interface{}) error {
	res, err := DB.Exec(q, s...)
	log.Printf("QUERYDONE: %s, %s\n", q, s)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}
	if rowsAffected, err := res.RowsAffected(); err != nil || rowsAffected == 0 {
		log.Printf("%s or 0 rows affected\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}
	return nil
}

// QueryAll helps process DB.Query
func QueryAll(w http.ResponseWriter, i interface{}, q string, s ...interface{}) error {
	var l int
	log.Printf("QUERY: %s, %s\n", q, s)
	results, err := DB.Query(q, s...)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}
	defer func() {
		err = results.Close()
		if err != nil {
			log.Println(err)
		} else {
			log.Println("QUERY closed")
		}
	}()
	// if Course, process
	if a, ok := i.(*[]Course); ok == true {
		for results.Next() {
			c := Course{}
			err = results.Scan(&c.Code, &c.Title, &c.Description)
			if err != nil {
				log.Println(err)
			}
			*a = append(*a, c)
		}
		l = len(*a)
	}
	// if Material, process
	if a, ok := i.(*[]Material); ok == true {
		for results.Next() {
			m := Material{}
			err = results.Scan(&m.Id, &m.Sequence, &m.FileName)
			if err != nil {
				log.Println(err)
			}
			*a = append(*a, m)
		}
		l = len(*a)
	}
	if l == 0 {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return err
	}
	return nil
}
