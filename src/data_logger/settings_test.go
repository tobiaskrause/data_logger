package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	configFile = "../../config/test_settings.json"
)

func TestSettingsLoad(t *testing.T) {
	NewSettings(configFile)
}

func TestSettingsOpcService(t *testing.T) {
	assert := assert.New(t)
	settings, err := NewSettings(configFile)
	assert.Nil(err)
	opcService, err := settings.OpcService()
	assert.Nil(err)
	assert.Equal("OPCServiceName.1", opcService.Name, "OPC Service Name")
	assert.Equal("value1", opcService.Values[0], "First Value")

}

func TestSettingsOpcUrl(t *testing.T) {
	assert := assert.New(t)
	settings, err := NewSettings(configFile)
	assert.Nil(err)
	opcURL := settings.OpcUrl()
	assert.Nil(err)
	assert.Equal("http://127.0.0.1:5000/opc", opcURL, "OPC Url")
}

func TestSettingsConnectors(t *testing.T) {
	assert := assert.New(t)
	settings, err := NewSettings(configFile)
	assert.Nil(err)
	connectors, err := settings.Connectors()
	assert.Nil(err)
	assert.Equal(2, len(connectors), "Size Of Connectors")
	for _, connector := range connectors {
		switch connector.ConnectorType() {
		case MongoDb:
			mongodbConnector := connector.(*MongoDbConnection)
			assert.Equal("MongoDb", mongodbConnector.Name, "Connector Name")
		case ElasticSearch:
			elasicSearchConnector := connector.(*ElasticSearchConnection)
			assert.Equal("ElasticSearch", elasicSearchConnector.Name)
		default:
			t.Fail()
		}
	}
}
