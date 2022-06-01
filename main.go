// Database connection

package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/robfig/cron/v3"
)

func main() {
	msc, _ := time.LoadLocation("Europe/Moscow")
	c := cron.New(cron.WithLocation(msc))

	c.AddFunc("@every 1m", func() {
		deleteFromParentsTimer()
	})

	c.Start()

	for {
		time.Sleep(time.Second * 1)
	}

	// result, err := db.Exec("insert into parents (surname, name, middlename, phone, status, processing_inf, mail, date_reg) values ($1, $2, $3, $4, $5, $6, $7, $8)", "Kozlicina", "Olga", "Anatolevna", "8999232341", "Mother", true, "Mail.com", time.Now())
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(result.RowsAffected()) // количество добавленных строк
}

func deleteFromParentsTimer() {
	connStr := "user=postgres password=6858 host=127.0.0.1 port=5432 dbname=test_1 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	type dateReg struct {
		id        int
		dateR     time.Time
		statusPay bool
	}

	rows, err := db.Query("select id, date_reg, status_pay from parents")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	dateRegs := []dateReg{}

	for rows.Next() {
		d := dateReg{}
		err := rows.Scan(&d.id, &d.dateR, &d.statusPay)
		if err != nil {
			log.Fatal(err)
		}
		dateRegs = append(dateRegs, d)
	}

	for _, d := range dateRegs {
		if !d.statusPay {
			differDate := time.Since(d.dateR)
			fmt.Println(differDate)

			if differDate.Hours() > 720 {
				result, err := db.Exec("delete from parents where id = $1", d.id)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println(result.RowsAffected())
			} else {
				fmt.Println("payment deadline has not expired")
			}
		}
	}
}
