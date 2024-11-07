package main

import (
	"fmt"

	"github.com/aminGhafoory/daq/internal/database"
	"go.bug.st/serial"
)

type DataLoggerService struct {
	DB         *database.Queries
	SerialPort *serial.Port
}

func (dls *DataLoggerService) fetchData() {
	//dls.SerialPort.fetch

}

func writeToDB() {

	//dls.serialPort.write

}

func parseData() {
	//prepare Data
}

func DataLogger() {

	for {
		//fetchData()
		//writeToDB()
		fmt.Println("??") //TODO!

	}

}
