package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/ulugbekrozimboyev/gosamples/daemons/db"
)

var NO_RESULT string = "No record found."

func connectDB() {
	// db configs, TODO read from config file
	var user, password, datb, host string = "root", "", "biotrackonline", "localhost"

	fmt.Println("Connection...")

	// connection to database
	/*for connection := db.Connect(datb, host, user, password); connection == false; {
		// if connection failed will try after 5 second
		time.Sleep(5 * time.Second)
		fmt.Println("Re-connect...")
		connection = db.Connect(datb, host, user, password)
	}*/

	if db.Connect(datb, host, user, password) == false {
		fmt.Println("Connection failed")
	}

	// successfully connaction
	fmt.Println("Connected!")
}

func main() {

	connectDB()
	var wg sync.WaitGroup
	for {

		for i := 0; i < 10; i++ {
			log.Print("Read new one---------------")

			tempInoutItems := db.ReadTempInout(1)

			if len(tempInoutItems) == 0 {

				fmt.Println(NO_RESULT)

				switch NO_RESULT {
				case "No record found.":
					NO_RESULT = "No record found.."
				case "No record found..":
					NO_RESULT = "No record found..."
				default:
					NO_RESULT = "No record found."
				}
				time.Sleep(1 * time.Second)
				continue
			}

			for _, item := range tempInoutItems {
				wg.Add(1)
				db.UpdateTempInoutItem(item)
				move(item, &wg)

			}
		}

		//wg.Wait()

	}

}

/*
	move to `inout` table from `temp_inout`

	Do following things:
	// not now 1. update `temp_inout` , set isBusy = 1
	2. write to `inout`
		a. get employee
		b. insert to `inout`
	3. delete item from `temp_inout`
*/
func move(item db.TempInoutItem, wg *sync.WaitGroup) bool {
	fmt.Println("running")
	defer wg.Done()

	// get device by serial number of device
	device, hasDevice := db.GetDeviceById(item.DeviceSN)
	fmt.Println(hasDevice)

	if !hasDevice {
		fmt.Println("Device not found")
		msg := "Device not found"
		db.AddLogToStampInoutError(item, msg)
		return false
	}

	fmt.Println(device)

	// get employee Id
	employeeID := db.GetEmployeeID(device.CompanyID, item.Pin)
	fmt.Println(employeeID)
	if employeeID == -1 {
		// employee not found
		msg := "Emplyee not found"
		db.AddLogToStampInoutError(item, msg)
		return false
	}

	// get status and event code
	eventCode, status := getInoutEventCodeAndStatus(device, item)

	fmt.Printf("eventCode = %d & status = %d\n", eventCode, status)

	// write inout
	result := db.InsertInoutRecord(item, device, employeeID, eventCode, status)

	return result
}

func getInoutEventCodeAndStatus(device db.Device, tempItem db.TempInoutItem) (int, int) {
	event_code := 0
	status := 1

	if device.InMode == 1 && device.OutMode == 0 {

		if tempItem.Status == device.InCode || tempItem.Status == device.InLunchCode {
			event_code = 0
		}
		status = 1
	} else if device.InMode == 0 && device.OutMode == 1 {

		if tempItem.Status == device.OutLunchCode || tempItem.Status == device.OutCode {
			event_code = 1
		}
		status = 0
	}

	if device.InMode == 1 && device.OutMode == 1 {

		if tempItem.Status == device.InLunchCode || tempItem.Status == device.InCode {
			status = 1
		} else {
			status = 0
		}
	}

	return event_code, status
}
