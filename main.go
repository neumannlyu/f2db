package main

import (
    ms "f2db/pkg/sql"
    "flag"
    "fmt"
    "os"
)

var (
    _file_name       string
    _table_name      string
    _is_clean        bool
    _is_init_db      bool
    _is_table_append bool
)

func main() {
    // 初始化检查
    if !parseCommandLineArguments() {
        return
    }

    // 加载配置
    sqlcfg := ms.LoadConfigJsonFile("db.json")
    // check
    if len(sqlcfg.Platform) == 0 {
        fmt.Println("load db.json error.")
        return
    }

    var isql ms.ISQL
    switch sqlcfg.Platform {
    case "postgresql":
        isql = &ms.PostgresSQL{
            DBPtr: nil,
        }
        // 打开数据库
        if !isql.Open(sqlcfg.Host, sqlcfg.User, sqlcfg.Password, sqlcfg.DBName) {
            return
        }
    case "mysql":
        isql = &ms.MySQL{
            DBPtr: nil,
        }
        // 打开数据库
        if !isql.Open(sqlcfg.Host, sqlcfg.User, sqlcfg.Password, sqlcfg.DBName) {
            return
        }
    default:
        fmt.Println("No Suppoted Platform!!!")
        return
    }

    // 清除数据表中所有的数据
    if _is_clean {
        for _, table := range sqlcfg.Tables {
            if isql.Execute(fmt.Sprintf("DELETE FROM %s;", table.TableName)) {
                fmt.Println("清空 " + _table_name + " 成功。")
            } else {
                fmt.Println("清空 " + _table_name + " 失败。")
            }
        }
    }

    // 初始化表。如果执行初始化，会先移除所有的表，然后根据表达式创建新的表格
    if _is_init_db {
        for _, table := range sqlcfg.Tables {
            // 移除表
            if isql.Execute(fmt.Sprintf("DROP TABLE IF EXISTS %s;", table.TableName)) {
                fmt.Println("移除 " + table.TableName + " 成功。")
            } else {
                fmt.Println("移除 " + table.TableName + " 失败。")
            }
            // 创建表
            if isql.Execute(table.CreateExp) {
                fmt.Println("创建 " + table.TableName + " 成功。")
            } else {
                fmt.Println("创建 " + table.TableName + " 失败。")
            }

            // 插入数据
            count := isql.ImportFromCVS(table.TableName, table.TableFile)
            fmt.Printf(" %s >> %s table successfully. Total %d records。\n", table.TableFile, table.TableName, count)
        }
    }

    if _is_table_append {
        for _, table := range sqlcfg.Tables {
            // 插入数据
            count := isql.ImportFromCVS(table.TableName, table.TableFile)
            fmt.Printf(" %s >> %s table successfully. Total %d records。\n", table.TableFile, table.TableName, count)
        }
    }

    // 清理资源
    isql.Close()
}

// parseCommandLineArguments 初始化。分析参数
//  @return bool 是否正常初始化，检查参数是否通过。
func parseCommandLineArguments() bool {
    // 1. 解析参数
    pf := flag.String("f", "", "[file] 导入的文件。支持cvs")
    pt := flag.String("t", "", "[table] 远程数据库表名")
    pc := flag.Bool("c", false, "[clean] 清空表中的数据")
    pi := flag.Bool("i", false, "[init] 初始化数据库。")
    pa := flag.Bool("a", false, "[append] 在表后进行追加。") //_is_table_append
    flag.Usage = func() {
        fmt.Fprintf(os.Stderr, "Usage of params:\n")
        fmt.Fprintf(os.Stderr, "    f2db [-f|-t|-c|-i|-a]\n\n")
        fmt.Fprintf(os.Stderr, "    Usage:\n")
        fmt.Fprintf(os.Stderr, "      [1]从文件中导入,追加到数据库中 f2db -t example_table -f example.csv -a\n")
        fmt.Fprintf(os.Stderr, "      [2]从文件中导入,并且清除原有数据 f2db -t example_table -f example.csv -c\n")
        fmt.Fprintf(os.Stderr, "      [3]用文件中数据初始化数据库 f2db -i\n")
        fmt.Fprintf(os.Stderr, "         按照Json文件中对应规则进行数据填充\n")
        flag.PrintDefaults()
    }
    flag.Parse()
    if len(os.Args) == 1 {
        flag.Usage()
        return false
    }
    _file_name = *pf
    _table_name = *pt
    _is_clean = *pc
    _is_init_db = *pi
    _is_table_append = *pa
    if !_is_init_db && _file_name == "" {
        fmt.Println("-f 需要指定导入文件。")
        flag.Usage()
        return false
    }
    return true
}
