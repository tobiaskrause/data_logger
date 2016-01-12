package main

import (
	"container/list"
	db "data_logger/database"
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	options Options
)

// Options of the command line
type Options struct {
	outputOnConsole *bool
	configFile      *string
}

func main() {
	addSignalHandlers()
	options = parseCommandLine()
	settings, err := NewSettings(*options.configFile)
	if err != nil {
		printError("ERROR: Cannot load config from %s (%s)\n", *options.configFile, err)
		return
	}
	// Create database connections
	connectors, err := settings.Connectors()
	if err != nil {
		printError("Cannot get connectors from config (%s)\n", err)
		return
	}

	print("Connectors\n")
	print("==========\n")
	databases := make([]db.Database, len(connectors))
	for index, connector := range connectors {
		databases[index] = connector2DbObj(connector)
		printf("Name: %s Type: %s\n", connector.ConnectorName(), connector.ConnectorType())
	}
	printf("\n")

	// Add each opc value to the timer
	opcService, err := settings.OpcService()
	if err != nil {
		printError("Cannot get OpcService from config (%s)", err)
		return
	}
	print("OPC Service Values\n")
	print("==================\n")
	opcURL := settings.OpcUrl()

	serviceURL := getServiceURL(opcURL, opcService.Name)
	opcTimer := OPCTimer{time.Minute * 1, time.Minute * 15, list.New()}
	for _, value := range opcService.Values {
		tracker := Tracker{opcUrl: serviceURL, opcName: value, dbConnectors: databases, force: true}
		opcTimer.AddTracker(&tracker)
		printf("Name: %s\n", value)
	}

	// Start timer cycle
	var wg sync.WaitGroup
	print("Starting timer.\n")
	opcTimer.Run(&wg)
	wg.Wait()
	print("Finished.\n")
}

func addSignalHandlers() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		os.Exit(0)
	}()
}

func printf(format string, a ...interface{}) {
	if *options.outputOnConsole {
		log.Printf(format, a...)
	}
}

func printError(format string, a ...interface{}) {
	errorStr := "ERROR: " + format
	log.Printf(errorStr, a...)
}

func getServiceURL(url string, serviceName string) string {
	return url + "/" + serviceName
}

func connector2DbObj(connector Connector) db.Database {
	switch connector.ConnectorType() {
	case MongoDb:
		mdbConnector := connector.(*MongoDbConnection)
		return &db.MgoDatabase{Host: mdbConnector.Host, Db: mdbConnector.Database, Collection: mdbConnector.Collection}
	case ElasticSearch:
		eSConnector := connector.(*ElasticSearchConnection)
		return &db.ElasticSearchDatabase{Host: eSConnector.Url, DbType: eSConnector.Type, IndexName: eSConnector.IndexName}
	default:
		return nil
	}
}

func parseCommandLine() Options {
	var options Options
	options.outputOnConsole = flag.Bool("console", false, "Print log messages to stdout")
	options.configFile = flag.String("config", "settings.json", "Path to json config file")
	flag.Parse()
	return options
}
