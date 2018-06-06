package configs

import (
	"encoding/json"
	"log"
	"os"
)

var AppConf = AppConfig{}

type AppConfig struct {
	ServicePort          int        `json:"service_port"`
	DbUrl                string     `json:"db_url"`
	CdnUrl               string     `json:"cdn_url"`
	DispatcherServiceUrl string     `json:"dispatcher_service_url"`
	StorageServiceUrl    string     `json:"storage_service_url"`
	SolrServiceUrl       string     `json:"solr_service_url"`
	PubsubConf           PubsubConf `json:"pubsub_conf"`
}

type PubsubConf struct {
	CredsFile    string `json:"creds_file"`
	ProjectId    string `json:"project_id"`
	Topic        string `json:"topic"`
	Subscription string `json:"subscription"`
}

func Initialize(confFile string) {
	file, err := os.Open(confFile)
	if err != nil {
		log.Println(err)
	}
	decoder := json.NewDecoder(file)
	conf := AppConfig{}
	err = decoder.Decode(&conf)
	if err != nil {
		log.Println(err)
	}
	AppConf = conf
}
