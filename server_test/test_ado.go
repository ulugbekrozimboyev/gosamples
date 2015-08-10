package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "code.google.com/p/odbc"
)

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("odbc", "driver={Microsoft Access Driver (*.mdb)};dbq=d:\\User.mdb")
	if err != nil {
		fmt.Println("Connecting Error")
		return
	}

	var c_id string
	rows, err := db.Query("select LOGID from CHECKINOUT")

	if err != nil {
		log.Print(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&c_id)
		if err != nil {
			log.Println("ERROR in Scan")
			log.Print(err)
		}
		//fmt.Println(c_id)
	}

	st := "CREATE TABLE device (ID INT(10), deviceTypeID INT(10),"
	st += "serialNumber VARCHAR(250), attLogStamp INT(11), operLogStamp INT(11),"
	st += "photoStamp INT(11) NULL DEFAULT NULL, transTimes VARCHAR(100) NULL DEFAULT NULL, transInterval INT(10) UNSIGNED NULL DEFAULT NULL,"
	st += "transFlag VARCHAR(100) NULL DEFAULT NULL, realtime TINYINT(1) NULL DEFAULT NULL, encrypt TINYINT(1) NULL DEFAULT NULL, delay INT(10) UNSIGNED NULL DEFAULT NULL,"
	st += "errorDelay INT(10) UNSIGNED NULL DEFAULT NULL, timeZoneAdj VARCHAR(50) NULL DEFAULT NULL, lastRequestTime DATETIME NULL DEFAULT NULL, status TINYINT(4) NOT NULL DEFAULT 0,"
	st += "PRIMARY KEY (ID), UNIQUE INDEX UQ_serialNumber (serialNumber)) COMMENT='биометрические устройства';"
	fmt.Println(st)

	res, err := db.Exec(st)
	if err != nil {
		log.Println(err)
		log.Println(res)
	}

}
