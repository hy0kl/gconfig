package gconfig

import (
	"io/ioutil"
	"log"
	"strings"

	"github.com/spf13/cast"
	"gopkg.in/yaml.v2"
)

//yaml struct
//support yaml file parse
//implemented Config interface
type YamlFile struct {
	data map[string]interface{} // Section -> key : value
}

//set function
//need clear cache after set value
func (r *YamlFile) Set(section, key string, value interface{}) {
	sMap := make(map[interface{}]interface{}, 0)
	sMap[key] = value
	r.data[section] = sMap
}

//load yaml file
func loadYamlFile(path string) (cfg Config, err error) {
	//read file
	content, ioErr := ioutil.ReadFile(path)
	if ioErr != nil {
		err = ioErr
		log.Printf("loadYamlFile error: %v", ioErr)
		return nil, err
	}

	data := make(map[string]interface{}, 0)
	//convert bytes to map
	err = yaml.Unmarshal(content, &data)
	if err != nil {
		log.Printf("loadYamlFile error: %v", err)
		return nil, err
	}
	yamlFile := new(YamlFile)
	yamlFile.data = data
	return yamlFile, nil
}

//MustValue function implemented
func (r *YamlFile) MustValue(section, key string, defaultVal ...string) string {
	defaultValue := ""
	if len(defaultVal) > 0 {
		defaultValue = defaultVal[0]
	}

	if val, ok := r.data[section]; !ok {
		return defaultValue
	} else {
		//match the key
		if data, ok := val.(map[interface{}]interface{}); ok {
			for k, v := range data {
				if cast.ToString(k) == key {
					return cast.ToString(v)
				}
			}
		} else {
			return defaultValue
		}
	}
	return defaultValue
}

//MustValueArray implemented
//split value by delim,like ",","-"
func (r *YamlFile) MustValueArray(section, key, delim string) []string {
	val := r.MustValue(section, key, "")
	if val != "" {
		return strings.Split(val, delim)
	}
	return nil
}

//GetKeyList implemented
//get all keys
func (r *YamlFile) GetKeyList(section string) []string {
	if val, err := r.GetSection(section); err != nil {
		return nil
	} else {
		data := make([]string, len(val))
		for k, _ := range val {
			data = append(data, k)
		}
		return data
	}
}

//empty function
func (r *YamlFile) GetSectionList() []string {
	return nil
}

//GetSection implemented
func (r *YamlFile) GetSection(section string) (map[string]string, error) {
	if val, ok := r.data[section]; !ok {
		return nil, nil
	} else {
		//math the type
		if data, ok := val.(map[interface{}]interface{}); ok {
			ret := make(map[string]string, len(data))
			//format map key and value
			for k, v := range data {
				ret[cast.ToString(k)] = cast.ToString(v)
			}
			return ret, nil
		}
	}
	return nil, nil
}

//GetSectionObject implemented
//object must be a pointer
func (r *YamlFile) GetSectionObject(section string, obj interface{}) error {
	//no section,use all data
	if section == "" {
		//formate map to bytes
		byt, err := yaml.Marshal(r.data)
		if err != nil {
			return err
		}
		//format bytes to object
		err = yaml.Unmarshal(byt, obj)
		if err != nil {
			return err
		}
	} else if val, ok := r.data[section]; !ok {
		//not hit value
		return nil
	} else {
		//convert value to object by section
		byt, err := yaml.Marshal(val)
		if err != nil {
			return err
		}
		err = yaml.Unmarshal(byt, obj)
		if err != nil {
			return err
		}
	}

	return nil
}
