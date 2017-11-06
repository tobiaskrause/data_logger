Instructions
============
Small application to read opc values from a rest service (not included) and write
the values to Mongodb and ElasicSearch instances

How to build
============
Build the application with:
```
go build data_logger
```
How to testing
==============
Tests can be run with:
```
go test data_logger
```
Configuration
=======================
The projects includes a test config.
```
  "Opc" : {
    "OpcRestServerUrl" : "http://127.0.0.1/opc",
    "OpcService" : {
      "Name" : "OPCServiceName.1",
      "Values" : [
        "value1",
        "value2",
        "value3",
        "value4",
        "value5",
        "value6",
        "value7",
        "value8"
      ]
    }
  },
```
The information from Opc Section of the config file is used to build the individual
REST request.
The read a value from the REST service following URL is used:

"http://127.0.0.1/opc/OPCServiceName.1/value1"
"<OpcRestServerUrl>/<Name>/<Values[0]>"

The application can process following answer of the REST service:
```
{
  "name": "value1",
  "timestamp": "01/12/16 18:58:15",
  "value": "2.00"
}
```

Connectors configuration allow multiple MongoDB and ElasticSearch configurations.
```
  "Connectors" : [
    {
      "Name" : "MongoDb",
      "Type" : "MongoDb",
      "Host" : "127.0.0.1:27017",
      "Database" : "opc_sevice",
      "Collection" : "measurements"
    },
    {
      "Name" : "ElasticSearch",
      "Type" : "ElasticSearch",
      "Url" : "http://127.0.0.1:9200",
      "DbType" : "measurements",
      "IndexName" : "opc_service"
    }
  ]
```
Run
===
Simple call '''data_logger''' from command line to start the application. Use '''data_logger --help''' to find out about command line switches.
