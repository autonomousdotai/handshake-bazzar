package configs

import (
	"os"
	"strconv"
)

var ServicePort, _ = strconv.Atoi(os.Getenv("SERVICE_PORT"))
var DB = os.Getenv("DB")
var CdnDomain = os.Getenv("CDN_DOMAIN")
var CdnHttps, _ = strconv.ParseBool(os.Getenv("SERVICE_PORT"))
var SolrServiceUrl = os.Getenv("SOLR_SERVICE_URL")
var StorageServiceUrl = os.Getenv("STORAGE_SERVICE_URL")
