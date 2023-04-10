package psql

import (
    "fmt"
    "testing"
)

func TestAnalysis_Kernel_CVE(t *testing.T) {
    db, err := OpenPostgresDB("172.16.5.114", "ly", "123456", "lydb")
    if err != nil {
        fmt.Printf("err.Error(): %v\n", err.Error())
        return
    }
    defer ClosePostgresDB(db)

    fmt.Println("OPen OK")
    ImportFromCVS(db, "file_suffix_table", "/Users/neumann/Desktop/excavator/file_suffix_table.csv")

}
