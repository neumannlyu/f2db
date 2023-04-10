package main

import (
    "database/sql"
    "encoding/json"
    "f2db/pkg/psql"
    "flag"
    "fmt"
    "io/ioutil"
    "os"
)

var (
    _file_name       string
    _table_name      string
    _is_drop         bool
    _is_init_db      bool
    _is_table_append bool
)

func main() {
    // 初始化检查
    if !_init() {
        return
    }

    var dbstruct psql.StDBJson
    fileBytes, err := ioutil.ReadFile("db.json")
    if err != nil {
        fmt.Printf("err.Error(): %v\n", err.Error())
    }

    // 将JSON内容解析为结构体
    err = json.Unmarshal(fileBytes, &dbstruct)
    if err != nil {
        fmt.Printf("err.Error(): %v\n", err.Error())
    }

    // 打开数据库连接
    db, err := psql.OpenPostgresDB(dbstruct.Host, dbstruct.User, dbstruct.Password, dbstruct.DBName)
    if err != nil {
        fmt.Printf("OpenPostgresDB err.Error(): %v\n", err.Error())
    }
    defer psql.ClosePostgresDB(db)

    _init_db(db, dbstruct.Tables)

    // 插入数据
    insertData(db, dbstruct.Tables)
}

// _init 初始化。分析参数
//  @return bool 是否正常初始化，检查参数是否通过。
func _init() bool {
    // 1. 解析参数
    pf := flag.String("f", "", "[file] 导入的文件。支持cvs")
    pt := flag.String("t", "", "[table] 远程数据库表名")
    pc := flag.Bool("c", false, "[clean] 清空表中的数据")
    pi := flag.Bool("i", false, "[init] 初始化数据库。")
    pa := flag.Bool("a", false, "[append] 在表后进行追加。") //_is_table_append
    flag.Usage = func() {
        fmt.Fprintf(os.Stderr, "Usage of params:\n")
        fmt.Fprintf(os.Stderr, "    f2db <-h -u -p -d> [-f|-t|-c|-i|-a]\n\n")
        fmt.Fprintf(os.Stderr, "    Usage:\n")
        fmt.Fprintf(os.Stderr, "      [1]从文件中导入,追加到数据库中 f2db -t example_table -f example.csv -a\n")
        fmt.Fprintf(os.Stderr, "      [2]从文件中导入,并且清除原有数据 f2db -t example_table -f example.csv -c\n")
        fmt.Fprintf(os.Stderr, "      [3]用文件中数据初始化数据库 f2db -i\n")
        fmt.Fprintf(os.Stderr, "         按照文件与表对应规则进行数据填充(或者json文件中修改)：\n")
        fmt.Fprintf(os.Stderr, "         kft.csv <> known_file_table\n")
        fmt.Fprintf(os.Stderr, "         fst.csv <> file_suffix_table\n")
        fmt.Fprintf(os.Stderr, "         kvt.csv <> kernel_vulnerablity_table\n")
        flag.PrintDefaults()
    }
    flag.Parse()
    if len(os.Args) == 1 {
        flag.Usage()
        return false
    }
    _file_name = *pf
    _table_name = *pt
    _is_drop = *pc
    _is_init_db = *pi
    _is_table_append = *pa
    if !_is_init_db && _file_name == "" {
        fmt.Println("-f 需要指定导入文件。")
        flag.Usage()
        return false
    }
    return true
}

// _init_db 初始化数据库
//  @param db
func _init_db(db *sql.DB, tables []psql.StTable) {
    // 判断是否需要初始化数据库
    if _is_init_db {
        for _, table := range tables {
            // 移除表
            _, err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s;", table.TableName))
            if err != nil {
                fmt.Println("移除 " + table.TableName + " 失败。" + err.Error())
            } else {
                fmt.Println("移除 " + table.TableName + " 成功")
            }

            // 创建表
            _, err = db.Exec(table.CreateExp)
            if err != nil {
                fmt.Println("创建 " + table.TableName + " 失败。" + err.Error())
            } else {
                fmt.Println("创建 " + table.TableName + " 成功")
            }
        }
    }

    if _is_drop {
        _, err := db.Exec(fmt.Sprintf("DELETE FROM %s;", _table_name))
        if err != nil {
            fmt.Println("清空 " + _table_name + " 成功。")
        } else {
            fmt.Println("清空 " + _table_name + " 失败。")
        }
    }
}

// insertData 从csv文件中导入数据
//  @param db
func insertData(db *sql.DB, tables []psql.StTable) {
    if _is_init_db {
        for _, table := range tables {
            fmt.Printf("导入%s文件成功。%s表中现在记录数为 %d。\n", table.TableFile, table.TableName, psql.ImportFromCVS(db, table.TableName, table.TableFile))
        }
    }

    if _is_drop {
        fmt.Printf("导入%s文件成功。%s表中现在记录数为 %d。\n", _file_name, _table_name, psql.ImportFromCVS(db, _table_name, _file_name))
    }

    if _is_table_append {
        fmt.Printf("导入%s文件成功。%s表中现在记录数为 %d。\n", _file_name, _table_name, psql.ImportFromCVS(db, _table_name, _file_name))
    }
}
