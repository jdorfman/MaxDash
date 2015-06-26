package main

import (
    "os"
    "errors"
    "reflect"
    "runtime"
    "encoding/json"
)

type Config struct {
    Cpus    string      `json:"cpus"`
    DbAddr  string      `json:"dbaddr"`
}

func (c *Config) validate() bool {
    setCPUs(c.Cpus)

    return true
}

// get data from config.json
func getConfig(confFile string) (*Config, error) {
    file, err := os.Open(confFile)
    if err != nil {
        return nil, errors.New("Error: Config file not found. Please check " + confFile)
    }
    defer file.Close()

    var config Config
    jsonDecoder := json.NewDecoder(file)
    if err = jsonDecoder.Decode(&config); err != nil {
        return nil, errors.New("Unable to decode config file. " + err.Error())
    }

    // validate before allowing to continue
    if !config.validate() {
        return nil, errors.New("Validation of config.json failed, exiting.")
    }

    return &config, nil
}

func setCPUs(cpuValue interface{}) error {
    switch reflect.ValueOf(cpuValue).Kind().String() {
    case "string":
        if reflect.ValueOf(cpuValue).String() == "max" {
            runtime.GOMAXPROCS(runtime.NumCPU())
            return nil
        } else {
            return errors.New("Unrecognized string in cpus field of config.json")
        }
    case "float64": 
        runtime.GOMAXPROCS(int(reflect.ValueOf(cpuValue).Float()))
        return nil
    }
    return errors.New("Unrecognized value in cpus field of config.json")
}
