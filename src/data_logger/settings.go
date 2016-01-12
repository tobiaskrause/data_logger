package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
)

// OpcService config structure
type OpcService struct {
	Name   string
	Values []string
}

func setField(objectValue reflect.Value, name string, value interface{}) error {
	structFieldValue := objectValue.FieldByName(name)
	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}
	fieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if val.Kind() == reflect.Slice {
		fieldTypeSlice := reflect.MakeSlice(fieldType, 0, 0).Interface()
		for _, v := range val.Interface().([]interface{}) {
			switch fieldType.Elem().Kind() {
			case reflect.String:
				fieldTypeSlice = append(fieldTypeSlice.([]string), v.(string))
			default:
				return fmt.Errorf("Cannot recognize slice type (%s)", fieldType.String())
			}
		}
		val = reflect.ValueOf(fieldTypeSlice)
	} else if fieldType != val.Type() {
		return fmt.Errorf("Provided value type (%s) did not match with object field type (%s)",
			val.Type().String(), fieldType.String())
	}
	structFieldValue.Set(val)
	return nil
}

func fillStruct(object interface{}, m map[string]interface{}) error {
	for key, value := range m {
		objectValue := reflect.ValueOf(object).Elem()
		err := setField(objectValue, key, value)
		if err != nil {
			return err
		}
	}
	return nil
}

// Connector interface
type Connector interface {
	ConnectorType() string
	ConnectorName() string
}

// Const definition of the ConnectorTypes
const (
	MongoDb       = "MongoDb"
	ElasticSearch = "ElasticSearch"
)

// MongoDbConnection holfs information of the connection to a Mongo server
type MongoDbConnection struct {
	Name       string
	Type       string
	Host       string
	Database   string
	Collection string
}

// ConnectorType returns the MongoDbConnectionType
// This is a string to identify the connection
func (connector *MongoDbConnection) ConnectorType() string {
	return connector.Type
}

// ConnectorName returns the custom name of the connector
func (connector *MongoDbConnection) ConnectorName() string {
	return connector.Name
}

// ElasticSearchConnection holds information of the connection to a ElasticSerach server
type ElasticSearchConnection struct {
	Name      string
	Type      string
	Url       string
	DbType    string
	IndexName string
}

// ConnectorType return the ElasticSearchConnectionType
// This is a string to identify the connection
func (connector *ElasticSearchConnection) ConnectorType() string {
	return connector.Type
}

// ConnectorName returns the custom name of the connector
func (connector *ElasticSearchConnection) ConnectorName() string {
	return connector.Name
}

// Settings holds the data of a config file
type Settings struct {
	FileName  string                 // Path to the config file
	configMap map[string]interface{} // Raw data of the config file
}

// NewSettings reads the config file of the given path
func NewSettings(aFileName string) (*Settings, error) {
	settings := Settings{FileName: aFileName}
	err := settings.load()
	return &settings, err
}

func (settings *Settings) load() error {
	dat, err := ioutil.ReadFile(settings.FileName)
	if err != nil {
		return err
	}
	err = json.Unmarshal(dat, &settings.configMap)
	if err != nil {
		return err
	}
	return nil
}

func (settings *Settings) opc() map[string]interface{} {
	return settings.configMap["Opc"].(map[string]interface{})
}

// OpcService returns the OpcService structure.
func (settings *Settings) OpcService() (OpcService, error) {
	service := &OpcService{}
	err := fillStruct(service, settings.opc()["OpcService"].(map[string]interface{}))
	return *service, err
}

// OpcUrl returns the REST URL of the OPC server
func (settings *Settings) OpcUrl() string {
	opcUrl := settings.opc()["OpcRestServerUrl"]
	return opcUrl.(string)
}

// Connectors returns an array of all connectors defined by the config file
func (settings *Settings) Connectors() ([]Connector, error) {
	connectors := []Connector{}
	connectorsMap := settings.configMap["Connectors"].([]interface{})
	var err error
	for _, connectorMap := range connectorsMap {
		switch connectorMap.(map[string]interface{})["Type"] {
		case MongoDb:
			tmpConnector := &MongoDbConnection{}
			err = fillStruct(tmpConnector, connectorMap.(map[string]interface{}))
			connectors = append(connectors, tmpConnector)
		case ElasticSearch:
			tmpConnector := &ElasticSearchConnection{}
			err = fillStruct(tmpConnector, connectorMap.(map[string]interface{}))
			connectors = append(connectors, tmpConnector)
		}
	}
	return connectors, err
}
