package database

import "gopkg.in/olivere/elastic.v2"

// ElasticSearch Database Implementation

//ElasticSearchDatabase holds connection information for a certain ElasticSearch instance
type ElasticSearchDatabase struct {
	Host      string
	DbType    string
	IndexName string
	client    *elastic.Client
}

// WriteElement writes data to a Elastic Serach instance
func (es *ElasticSearchDatabase) WriteElement(element Measurement) error {
	es.client = es.createConnection()
	es.createIndexIfNotExists()
	es.writeRecord(element)
	return nil
}

func (es *ElasticSearchDatabase) createConnection() *elastic.Client {
	client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL("http://192.168.1.65:9200"))
	if err != nil {
	}
	return client
}

func (es *ElasticSearchDatabase) createIndexIfNotExists() {
	if !es.hasIndex() {
		es.createIndex()
	}
}

func (es *ElasticSearchDatabase) hasIndex() bool {
	exists, err := es.client.IndexExists(es.IndexName).Do()
	if err != nil {
		exists = false
	}
	return exists
}

func (es *ElasticSearchDatabase) createIndex() bool {
	var (
		success = true
	)
	createIndex, err := es.client.CreateIndex(es.IndexName).Do()
	if err != nil {
		success = false
	}
	if !createIndex.Acknowledged {
		success = false
	}
	return success
}

func (es *ElasticSearchDatabase) writeRecord(element Measurement) {
	_, err := es.client.Index().
		Index(es.IndexName).
		Type(es.DbType).
		BodyJson(element).
		Do()
	if err != nil {
		return
	}
}
