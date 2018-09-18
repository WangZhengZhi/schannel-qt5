package config

import (
	"testing"

	"encoding/json"
	"os"
)

func TestConfigPath(t *testing.T) {
	err := os.Setenv("HOME", "/home/testing")
	if err != nil {
		t.Errorf("set $HOME ERROR: %v\n", err)
	}

	if path, err := ConfigPath(); err != nil {
		t.Errorf("ConfigPath return an ERROR: %v\n", err)
	} else if path != "/home/testing/.local/share/schanclient.json" {
		t.Errorf("wrong path: %v", path)
	}
}

func TestMarshalUserConf(t *testing.T) {
	u := new(UserConfig)
	u.SSRNodeConfigPath.Data = "/tmp/testing/t.json"
	u.SSRBin.Data = "/tmp/testing/a.out"
	u.LogFile.Data = "/tmp/a.log"

	data, err := json.MarshalIndent(u, "", "\t")
	if err != nil {
		t.Error(err)
	}
	t.Log(string(data))
}

func TestUnmarshalUserConf(t *testing.T) {
	u := new(UserConfig)
	data := `{"ssr_config_path":"/tmp/testing/t.json","ssr_bin":"/tmp/testing/a.out","log_file":"/tmp/a.log"}`
	if err := json.Unmarshal([]byte(data), u); err != nil {
		t.Error(err)
	}
	t.Log(*u)
}
