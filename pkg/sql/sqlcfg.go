package sql

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
)

type SQLCfg struct {
    Platform string      `json:"platform"`
    Host     string      `json:"host"`
    User     string      `json:"user"`
    Password string      `json:"password"`
    DBName   string      `json:"dbname"`
    Tables   []TableInfo `json:"tables"`
}

type TableInfo struct {
    TableName string `json:"name"`
    TableFile string `json:"file"`
    CreateExp string `json:"create"`
}

func LoadConfigJsonFile(filepath string) SQLCfg {
    var sqlcfg SQLCfg
    fileBytes, err := ioutil.ReadFile(filepath)
    if err != nil {
        fmt.Printf("err.Error(): %v\n", err.Error())
    }

    // 将JSON内容解析为结构体
    err = json.Unmarshal(fileBytes, &sqlcfg)
    if err != nil {
        fmt.Printf("err.Error(): %v\n", err.Error())
    }
    return sqlcfg
}
