package generator

import (
	"encoding/json"
	log "github.com/inconshreveable/log15"
	"io/ioutil"
)

type ParamDefinition struct {
	Name string `json:"name"`
	Type      string `json:"type"`
	TypeAlias string `json:"typeAlias,omitempty"`
	Rule      string `json:"rule,omitempty"`
	Optional  bool `json:"optional,omitempty"`
	Describe  string `json:"describe"`
	SubParam  []ParamDefinition `json:"subParam,omitempty"`
}

type InterfaceDefinition struct {
	Name  	string `json:"name"`
	Brief           string `json:"brief"`
	Describe        string `json:"describe,omitempty"`
	InputPreProcess string `json:"inputPreProcess,omitempty"`
	InputParamList []ParamDefinition `json:"inputParamList"`
}

type ServiceDefinition struct {
	Name string `json:"name"`
	HostList      map[string]string `json:"hostList"`
	InterfaceList []InterfaceDefinition `json:"interfaceList"`
}

func LoadServiceDefinitionFromFile(name string) (*ServiceDefinition, error) {
	dat, err := ioutil.ReadFile(name)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	svc := &ServiceDefinition{}
	err = json.Unmarshal(dat, svc)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return svc, nil
}
