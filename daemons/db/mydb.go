package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

const layout = "2006-01-02 15:04:05"

func Connect(mdb, host, usr, pwd string) bool {
	var err error
	//fmt.Println(mdb, host, usr, pwd)

	s := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", usr, pwd, host, mdb)
	db, err = sql.Open("mysql", s)
	if err != nil {
		log.Print("error Connect")
		log.Print(err.Error())
		return false
	}
	err = db.Ping()

	return true

}

func ReadTempInout(limit int) []TempInoutItem {
	var tempInoutItems []TempInoutItem

	var sql string = "SELECT * FROM temp_inout WHERE isBusy = ? LIMIT ?"

	rows, err := db.Query(sql, 0, limit)
	if err != nil {
		log.Print("error ReadTempInout")
		return tempInoutItems
	}
	// close cursor after reading
	defer rows.Close()

	for rows.Next() {
		var item TempInoutItem
		var atStr string
		err := rows.Scan(&item.DeviceSN, &item.Pin, &item.time, &item.Status, &item.verify, &atStr, &item.IsBusy)
		if err != nil {
			log.Print("error ReadTempInout scan")
			log.Print(err.Error())
		}

		item.at, _ = time.Parse(layout, atStr)

		tempInoutItems = append(tempInoutItems, item)
	}

	return tempInoutItems
}

/*
	this function will update only one column: `isBusy`
*/

func UpdateTempInoutItem(v TempInoutItem) bool {
	sql := "UPDATE temp_inout SET isBusy = ? WHERE deviceSN = ? AND pin = ? AND `time` = ? AND `at` = ?"
	res, err := db.Exec(sql, 1, v.DeviceSN, v.Pin, v.time, v.at)
	if err != nil {
		log.Print("error update tempinout ")
		log.Print(err.Error())
		log.Print(res)
		return false
	}

	return true
}

/*
	get device by serial number

	returns device if found
	in another case returns empty device
*/
func GetDeviceById(serialNumber string) (Device, bool) {
	var device Device
	sql := "SELECT ID, deviceTypeId, companyID, locationID, inMode, outMode, inCode, outCode, inLunchCode, outLunchCode FROM device WHERE serialNumber=? LIMIT 1;"

	deviceItem, err := db.Query(sql, serialNumber)

	if err != nil {
		log.Print("error GetDeviceById ")
		log.Print(err.Error())
		return device, false
	}
	defer deviceItem.Close()

	for deviceItem.Next() {

		err := deviceItem.Scan(&device.ID, &device.deviceTypeID, &device.CompanyID, &device.locationID, &device.InMode, &device.OutMode, &device.InCode, &device.OutCode, &device.InLunchCode, &device.OutLunchCode)
		if err != nil {
			log.Print("error GetDeviceById scan")
			log.Print(err.Error())
			return device, false
		}

		device.serialNumber = serialNumber

		return device, true
	}

	return device, false
}

/*
	write logs to database, exactly to `stampinout_error` table
*/
func AddLogToStampInoutError(v TempInoutItem, msg string) {
	sql := "INSERT INTO stampinout_error(deviceSN, pin, `time`, `status`, `verify`, `at`, msg) VALUES( ? , ?, ?, ?, ?, ?, ?)"
	res, err := db.Exec(sql, v.DeviceSN, v.Pin, v.time, v.Status, v.verify, v.at, msg)
	if err != nil {
		log.Printf("error %s , ", msg)
		log.Print(err.Error())
		log.Print(res)
	}

	DeleteTempInout(v)
}

/*
	Delete item from `temp_inout`
*/
func DeleteTempInout(v TempInoutItem) {
	sql := "DELETE FROM temp_inout WHERE deviceSN=? AND pin=? AND `time` = ? AND `status`=? AND `verify`=?"
	res, err := db.Exec(sql, v.DeviceSN, v.Pin, v.time, v.Status, v.verify)
	if err != nil {
		log.Print("error DeleteTempInout ")
		log.Print(err.Error())
		log.Print(res)
	}
	return
}

/*
	this function return employee id by pincode at one company
	return -1 when function could't find emplyee
*/

func GetEmployeeID(companyID, pinCode int) int {
	sql := "SELECT ID FROM employee WHERE companyID= ? AND pinCode= ? AND `state`='active' AND isActive=1 LIMIT 1"
	rows, err := db.Query(sql, companyID, pinCode)
	if err != nil {
		log.Print("error GetEmployeeID ")
		log.Print(err.Error())
		return -1
	}
	defer rows.Close()

	for rows.Next() {
		var employeeID int

		err := rows.Scan(&employeeID)
		if err != nil {
			log.Print("error GetEmployeeID scan")
			log.Print(err.Error())
			return -1
		}

		return employeeID
	}

	return -1
}

func IsWritenInout(newInoutItem InoutItem) bool {
	sql := "SELECT ID FROM `inout` WHERE companyID=? AND employeeID=? AND deviceSN=? AND pinCode=? AND `dt`=? AND `dt_year`=Year(?) AND `dt_month`=Month(?) AND `eventCode`=? AND `status`=?"

	rows, err := db.Query(sql, newInoutItem.CompanyID, newInoutItem.employeeID, newInoutItem.deviceSN, newInoutItem.pinCode, newInoutItem.dt, newInoutItem.dt, newInoutItem.dt, newInoutItem.eventCode, newInoutItem.status)
	if err != nil {
		log.Print("error IsWritenInout ")
		log.Print(err.Error())
		return false
	}
	defer rows.Close()

	for rows.Next() {
		fmt.Println("record found")
		return true
	}

	return false
}

/*
	Insert to `inout`
*/
func InsertInoutRecord(tempInout TempInoutItem, device Device, employeeID int, eventCode int, status int) bool {
	var newInoutItem InoutItem
	var result = false

	newInoutItem.CompanyID = device.CompanyID
	newInoutItem.employeeID = employeeID
	newInoutItem.deviceSN = tempInout.DeviceSN
	newInoutItem.pinCode = tempInout.Pin
	newInoutItem.dt = tempInout.time
	newInoutItem.eventCode = eventCode
	newInoutItem.status = status

	// check this was record writen already
	if IsWritenInout(newInoutItem) == false {

		fmt.Println("record not found, writing new one")
		sql := "INSERT INTO `inout`(companyID, employeeID, deviceSN, pinCode, `dt`, `dt_year`, `dt_month`, eventCode, `status`) VALUES(?, ?, ?, ?, ?, Year(?), Month(?), ?, ? )"

		res, err := db.Exec(sql, newInoutItem.CompanyID, newInoutItem.employeeID, newInoutItem.deviceSN, newInoutItem.pinCode, newInoutItem.dt, newInoutItem.dt, newInoutItem.dt, newInoutItem.eventCode, newInoutItem.status)
		if err != nil {
			fmt.Println("-------writing error--------")
			log.Print(err.Error())
			log.Print(res)
			return false
		}

		result = true

		fmt.Println("-------writing finished--------")
	}

	InsertToStampInout(tempInout, device.CompanyID, employeeID, eventCode)
	DeleteTempInout(tempInout)

	return result
}

func InsertToStampInout(item TempInoutItem, companyID, employeeID, eventCode int) {
	sql := "INSERT INTO `stampinout`(companyID, employeeID, deviceSN, pin, `time`, `status`, `verify`, `at`) VALUES(?, ?, ?, ?, ?, ?, ?, ?)"

	res, err := db.Exec(sql, companyID, employeeID, item.DeviceSN, item.Pin, item.time, eventCode, item.verify, item.at)
	if err != nil {
		log.Print("error InsertToStampInout ")
		log.Print(err.Error())
		log.Print(res)
	}

}

// a record from table `temp_inout`
type TempInoutItem struct {
	DeviceSN string
	Pin      int
	time     string
	Status   int
	verify   uint8
	at       time.Time
	IsBusy   uint8
}

// a record from table `inout`
type InoutItem struct {
	ID         int
	CompanyID  int
	employeeID int
	deviceSN   string
	pinCode    int
	dt         string
	eventCode  int
	status     int
}

// device
type Device struct {
	ID           int
	deviceTypeID int
	CompanyID    int
	locationID   int
	serialNumber string
	InMode       int
	OutMode      int
	InCode       int
	OutCode      int
	InLunchCode  int
	OutLunchCode int
}

// read TempInoutItem
func (v TempInoutItem) read(id int) TempInoutItem {
	var item TempInoutItem
	return item
}

// read TempInoutItem
func (v TempInoutItem) delete(id int) bool {
	return false
}

// string function to print understanable string for users
func (v TempInoutItem) String() string {
	return fmt.Sprintf("deviceSN: {%v}, pin: {%v}, time: {%v}, status: {%v}, verify: {%v}, at: {%v}, isBusy: {%v}", v.DeviceSN, v.Pin, v.time, v.Status, v.verify, v.at, v.IsBusy)
}

/*
	* InoutItem
	write InoutItem
*/

func (v InoutItem) write() {
	// TODO write
}

/*
	* InoutItem
	string function to print understanable string for users
*/
func (v InoutItem) String() string {
	return fmt.Sprintf("ID: {%v}, companyID : {%v},	employeeID: {%v}, deviceSN: {%v}, pinCode: {%v}, dt: {%v}, eventCode: {%v}, status: {%v}", v.ID, v.CompanyID, v.employeeID, v.deviceSN, v.pinCode, v.dt, v.eventCode, v.status)
}

/*
	* Device
	string function to print understanable string for users
*/
func (v Device) String() string {
	return fmt.Sprintf("ID: {%v}, deviceTypeID: {%v}, companyID : {%v}, locationID: {%v}, serialNumber: {%v}", v.ID, v.deviceTypeID, v.CompanyID, v.locationID, v.serialNumber)
}
