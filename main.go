package main

import (
    "database/sql"
    "f2db/pkg/psql"
    "flag"
    "fmt"
    "log"
    "os"
)

var (
    _file_name       string
    _host_ip         string
    _user_name       string
    _user_password   string
    _database_name   string
    _table_name      string
    _is_drop         bool
    _is_init_db      bool
    _is_table_append bool
)
var _db_tables = []string{"file_suffix_table", "kernel_vulnerability_table", "known_file_table"}
var _db_create_table_exps = []string{
    `CREATE TABLE file_suffix_table(id SERIAL PRIMARY KEY NOT NULL, suffix_name TEXT NOT NULL, importance INT NOT NULL, description TEXT);`,
    `CREATE TABLE kernel_vulnerability_table(id SERIAL PRIMARY KEY NOT NULL, vul_id TEXT NOT NULL, affected_kernel_ver TEXT NOT NULL, vul_type TEXT, vul_description TEXT, severity INT, fix_suggestion TEXT);`,
    `CREATE TABLE known_file_table(id SERIAL PRIMARY KEY NOT NULL,file_name TEXT NOT NULL,idx INT NOT NULL, category INT,importance INT NOT NULL, count INT,md5 TEXT,description TEXT);`,
}

func main() {
    // 初始化检查
    if !_init() {
        return
    }

    // 打开数据库连接
    db, err := psql.OpenPostgresDB(_host_ip, _user_name, _user_password, _database_name)
    if err != nil {
        log.Fatalln(err.Error())
    }
    defer psql.ClosePostgresDB(db)

    _init_db(db)

    // 插入数据
    insertData(db)
}

// _init 初始化。分析参数
//  @return bool 是否正常初始化，检查参数是否通过。
func _init() bool {
    // 1. 解析参数
    pf := flag.String("f", "", "[file] 导入的文件。支持cvs")
    ph := flag.String("h", "", "[host] 远程数据库主机的IP")
    pu := flag.String("u", "", "[user] 远程数据库连接的用户名")
    pp := flag.String("p", "", "[password] 远程数据库连接的密码")
    pd := flag.String("d", "", "[database] 远程数据库名")
    pt := flag.String("t", "", "[table] 远程数据库表名")
    pc := flag.Bool("c", false, "[clean] 清空表中的数据")
    pi := flag.Bool("i", false, "[init] 初始化数据库。")
    pa := flag.Bool("a", false, "[append] 在表后进行追加。") //_is_table_append
    flag.Usage = func() {
        fmt.Fprintf(os.Stderr, "Usage of params:\n")
        fmt.Fprintf(os.Stderr, "    f2db <-h -u -p -d> [-f|-t|-c|-i|-a]\n\n")
        fmt.Fprintf(os.Stderr, "    Usage:\n")
        fmt.Fprintf(os.Stderr, "      [1]从文件中导入,追加到数据库中 -h 127.0.0.1 -u usr -p 123456 -d db -t example_table -f example.csv -a\n")
        fmt.Fprintf(os.Stderr, "      [2]从文件中导入,并且清除原有数据 -h 127.0.0.1 -u usr -p 123456 -d db -t example_table -f example.csv -c\n")
        fmt.Fprintf(os.Stderr, "      [3]用文件中数据初始化数据库 -h 127.0.0.1 -u usr -p 123456 -d db -i\n       按照文件与表对应规则：")
        fmt.Fprintf(os.Stderr, "         按照文件与表对应规则进行数据填充：\n")
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
    _host_ip = *ph
    _user_name = *pu
    _user_password = *pp
    _database_name = *pd
    _table_name = *pt
    _is_drop = *pc
    _is_init_db = *pi
    _is_table_append = *pa
    if _host_ip == "" || _user_name == "" || _user_password == "" || _database_name == "" {
        fmt.Println("-h -u -p -d 字段都不允许为空。")
        flag.Usage()
        return false
    }

    if !_is_init_db && _file_name == "" {
        fmt.Println("-f 需要指定导入文件。")
        flag.Usage()
        return false
    }
    return true
}

// _init_db 初始化数据库
//  @param db
func _init_db(db *sql.DB) {
    // 判断是否需要初始化数据库
    if _is_init_db {
        for i := 0; i < len(_db_tables); i++ {
            // 移除表
            table := _db_tables[i]
            _, err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s;", table))
            if err != nil {
                fmt.Println("移除 " + table + " 成功。")
            } else {
                fmt.Println("移除 " + table + " 失败。")
            }

            // 创建表
            _, err = db.Exec(_db_create_table_exps[i])
            if err != nil {
                fmt.Println("创建 " + table + " 成功。")
            } else {
                fmt.Println("创建 " + table + " 失败。")
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
func insertData(db *sql.DB) {
    if _is_init_db {
        fmt.Printf("导入%s文件成功。%s表中现在记录数为 %d。\n", "kft.csv", "known_file_table", psql.ImportFromCVS(db, "known_file_table", "kft.csv"))
        fmt.Printf("导入%s文件成功。%s表中现在记录数为 %d。\n", "kft.csv", "file_suffix_table", psql.ImportFromCVS(db, "file_suffix_table", "fst.csv"))
        fmt.Printf("导入%s文件成功。%s表中现在记录数为 %d。\n", "kft.csv", "kernel_vulnerablity_table", psql.ImportFromCVS(db, "kernel_vulnerablity_table", "kvt.csv"))
    }

    if _is_drop {
        fmt.Printf("导入%s文件成功。%s表中现在记录数为 %d。\n", _file_name, _table_name, psql.ImportFromCVS(db, _table_name, _file_name))
    }

    if _is_table_append {
        fmt.Printf("导入%s文件成功。%s表中现在记录数为 %d。\n", _file_name, _table_name, psql.ImportFromCVS(db, _table_name, _file_name))
    }
}
