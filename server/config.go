package main

import(
    "encoding/json"
    "io/ioutil"
)

type myConfig struct{
    Server string `json:"server_addr"`
    Local string `json:"local_addr"`
    Key string `json:"key"`
}

var cfg myConfig

func readConfig(filename string)error{
    bytes, err := ioutil.ReadFile(filename)
    if err != nil {
        return err
    }
    if err := json.Unmarshal(bytes, &cfg); err != nil {
        return err
    }
    return nil
}