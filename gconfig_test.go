package gconfig

import (
	"testing"
)

func TestGetConf(t *testing.T) {
	//dir, _ := os.Getwd()
	dir := GetAbsPath()
	SetConfigFile(dir + "/conf/conf.ini")

	if GetConf("goconfig", "hosts") == "127.0.0.1 127.0.0.2 127.0.0.3" {
		t.Log("pass")
	} else {
		t.Error("fail")
	}

	if len(GetConfArr("goconfig", "hosts")) == 3 {
		t.Log("pass")
	} else {
		t.Error("fail")
	}

	if GetConfStringMap("goconfigStringMap")["name"] == "goconfig" {
		t.Log("pass")
	} else {
		t.Error("fail")
	}
}
