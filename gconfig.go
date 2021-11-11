package gconfig

import (
	"fmt"
	"log"
	"os"
	"sync"
)

//Config interface definition
//developers can implement this interface to combine with other config library
type Config interface {
	//get value by section and key,if not exist,return the defaultVal
	MustValue(section, key string, defaultVal ...string) string
	//get value by sction and key and split value by delim,return array
	MustValueArray(section, key, delim string) []string
	//get section value,return array
	GetKeyList(section string) []string
	//get section value,return map
	GetSection(section string) (map[string]string, error)
	//get all section list
	GetSectionList() []string
	//get object value by section
	GetSectionObject(section string, obj interface{}) error
	//set value by section and key when need
	Set(section, key string, value interface{})
}

var (
	// global cache
	// only load the config file once
	// after config init complete, all config data get from cache
	configCache = struct {
		sync.RWMutex
		cache map[string]Config
	}{cache: make(map[string]Config, 0)}

	//global config
	gCfg Config

	specifiedConfig string

	// default config file
	configFilename = "conf/conf.ini"
)

func InitConfig() {
	//check if config has inited
	if gCfg != nil {
		return
	}

	var cfgFile string
	var err error

	//set the default path
	if specifiedConfig != "" {
		cfgFile = specifiedConfig
	} else {
		cfgFile = configFilename
	}

	log.Printf("config file: %s", cfgFile)
	_, err = os.Stat(cfgFile)
	if err != nil {
		panic(fmt.Sprintf(`can open config file: %s, err: %v`, cfgFile, err))
	}

	// load config from path
	if gCfg, err = load(cfgFile); err != nil {
		gCfg = nil
		log.Printf("config error: %v", err)
	}
}

func SetConfigFile(filename string) {
	specifiedConfig = filename
	gCfg = nil
}

func Set(section, key string, value interface{}) {
	InitConfig()
	if gCfg == nil {
		log.Printf("config error: NOT_FOUND[sec:%s,key:%s]", section, key)
		return
	}
	gCfg.Set(section, key, value)
}

func ClearConfigCache() {
	configCache.Lock()
	configCache.cache = make(map[string]Config, 0)
	gCfg = nil
	configCache.Unlock()
}

func GetConf(sec, key string) string {
	//init
	InitConfig()
	if gCfg == nil {
		log.Printf("config error: NOT_FOUND[sec:%s,key:%s]", sec, key)
		return ""
	}
	//if value not existed return ""
	return gCfg.MustValue(sec, key, "")
}

func GetConfDefault(sec, key, def string) string {
	//init
	InitConfig()
	if gCfg == nil {
		log.Printf("config error: NOT_FOUND[sec:%s,key:%s]", sec, key)
		return ""
	}
	//if value not existed return def
	return gCfg.MustValue(sec, key, def)
}

func GetConfArr(sec, key string) []string {
	//init
	InitConfig()
	if gCfg == nil {
		log.Printf("config error: NOT_FOUND[sec:%s,key:%s]", sec, key)
		return []string{}
	}

	//if value not existed return " "
	return gCfg.MustValueArray(sec, key, " ")
}

func GetConfStringMap(sec string) (ret map[string]string) {
	//init
	InitConfig()
	if gCfg == nil {
		log.Printf("config error: NOT_FOUND[sec:%s]", sec)
		return nil
	}

	var err error
	//if value not existed return empty map
	if ret, err = gCfg.GetSection(sec); err != nil {
		log.Printf("Conf,err:%v", err)
		ret = make(map[string]string, 0)
	}

	return
}

func GetConfArrayMap(sec string) (ret map[string][]string) {
	//init
	InitConfig()
	if gCfg == nil {
		log.Printf("config error: NOT_FOUND[sec:%s]", sec)
		return nil
	}
	ret = make(map[string][]string, 0)
	//get all keys
	confList := gCfg.GetKeyList(sec)
	//get all config by range keys
	for _, k := range confList {
		ret[k] = gCfg.MustValueArray(sec, k, " ")
	}

	return
}

func ConfMapToStruct(sec string, v interface{}) error {
	//init
	InitConfig()
	if gCfg == nil {
		log.Printf("config error: NOT_FOUND[sec:%s]", sec)
		return nil
	}

	return gCfg.GetSectionObject(sec, v)
}

func load(cfgFile string) (cfg Config, err error) {
	//default load module
	fileType := "ini"

	if cfgFile[len(cfgFile)-4:] == "yaml" {
		fileType = "yaml"
		//if the path suffix is not ".ini", completed path by append ".ini"
	} else if cfgFile[len(cfgFile)-4:] != ".ini" {
		cfgFile = cfgFile + ".ini"
	}

	//config cache
	var ok bool
	configCache.RLock()
	//read cache first
	cfg, ok = configCache.cache[cfgFile]
	configCache.RUnlock()
	//no cache
	if !ok {
		//load file and create cache
		configCache.Lock()
		if fileType == "yaml" {
			if cfg, err = loadYamlFile(cfgFile); err == nil {
				configCache.cache[cfgFile] = cfg
			}
		} else {
			if cfg, err = loadIniFile(cfgFile); err == nil {
				configCache.cache[cfgFile] = cfg
			}
		}

		configCache.Unlock()
	}
	return
}
