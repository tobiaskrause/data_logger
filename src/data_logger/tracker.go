package main

import (
	db "data_logger/database"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

// Tracker knows where to read the opc value and where the value should be written
type Tracker struct {
	opcUrl       string
	opcName      string
	dbConnectors []db.Database
	force        bool
	lastValue    interface{}
}

// ReadValue reads the opc value and returns the Measurement structure
func (tracker *Tracker) ReadValue() (bool, *db.Measurement) {
	contents := tracker.getJSONFromServer()
	var data map[string]string
	if contents != nil {
		err := json.Unmarshal(contents, &data)
		if err != nil {
			return false, nil
		}
		measure := db.Measurement{Name: data["name"], Timestamp: time.Now()}
		measure.SetValueFromString(data["value"])
		printf("Read data: Name: (%s) Timestamp: (%s) Value: (%f)\n", measure.Name, measure.GetTimeAsString(), measure.Value)
		return true, &measure
	}
	return false, nil
}

// WriteValue writes the given Measurement structure to all defined databases
func (tracker *Tracker) WriteValue(element *db.Measurement) {
	for index, database := range tracker.dbConnectors {
		err := database.WriteElement(*element)
		if err != nil {
			printError("Database %d: Couldn't write value to database (%s)", index, err)
		} else {
			printf("%s wrote to database %d\n", element.Name, index)
		}
	}
}

func (tracker *Tracker) getJSONFromServer() []byte {
	response, err := http.Get(tracker.opcUrl + "/" + tracker.opcName)
	if err != nil {
		printError("Cannot get JSON from server (%s)", err)
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			printError("Cannot read content from server (%s)", err)
			return nil
		}
		return contents
	}
	return nil
}
