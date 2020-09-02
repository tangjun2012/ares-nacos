package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
)

type localEnv struct {
	local *sync.Map
}

var env *localEnv

func init() {
	loadEnv()
}

func (l *localEnv) Get(key string) string {
	if l.local == nil {
		return ""
	}
	if val, ok := l.local.Load(key); ok {
		return val.(string)
	}
	return ""
}

func (l *localEnv) Set(key string, val string) {
	if l.local == nil {
		l.local = &sync.Map{}
	}
	l.local.Store(key, val)

	refreshLocalFile(l.local)
}

func (l *localEnv) Delete(key string) {
	if l.local == nil {
		return
	}
	l.local.Delete(key)
	refreshLocalFile(l.local)
}

func getLocalFile() string {
	home, _ := os.UserHomeDir()
	return home + "/" + GetString("local", "ares.env")
}

var fileLock = sync.Mutex{}

func refreshLocalFile(l *sync.Map) {
	fileLock.Lock()
	defer fileLock.Unlock()
	configFile := getLocalFile()
	os.Remove(configFile)

	m := map[string]string{}

	l.Range(func(key, value interface{}) bool {
		m[key.(string)] = value.(string)
		return true
	})
	bs, _ := json.Marshal(m)

	ioutil.WriteFile(configFile, bs, os.ModePerm)
}

func loadEnv() {
	if env != nil {
		return
	}
	configFile := getLocalFile()

	_, err := os.Stat(configFile)
	if err != nil {
		env = &localEnv{
			local: &sync.Map{},
		}
		return
	}

	bs, _ := ioutil.ReadFile(configFile)

	m := map[string]string{}
	json.Unmarshal(bs, &m)

	var local = sync.Map{}
	for key, val := range m {
		local.Store(key, val)
	}

	env = &localEnv{
		local: &local,
	}
	return
}

func GetLocal() *localEnv {
	return env
}
