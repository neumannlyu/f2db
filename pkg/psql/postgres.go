package psql

import (
    "database/sql"
    "encoding/csv"
    "fmt"
    "log"
    "os"
    "strings"

    _ "github.com/lib/pq"
)

// OpenPostgresDB 打开数据库连接
// !在调用完数据库后一定要使用ClosePostgresDB关闭连接。
//  @param h host
//  @param u user
//  @param p password
//  @param d database name
//  @return *sql.DB db指针
//  @return error 错误
func OpenPostgresDB(h, u, p, d string) (*sql.DB, error) {
    return sql.Open("postgres", fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", h, u, p, d))
}

// 关闭数据库连接
// !在调用完数据库后一定要关闭连接。
func ClosePostgresDB(db *sql.DB) {
    db.Close()
}

// 从cvs文件中导入到数据库中
// todo 以后写个工具来导入数据库，顺便写个文档来总结go怎么和数据库相连接
//  @param file_path
// @return count 导入后表中的记录数
func ImportFromCVS(db *sql.DB, table_name, file_path string) (count int) {
    // 打开CSV文件并创建一个*csv.Reader对象。
    file, err := os.Open(file_path)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()
    reader := csv.NewReader(file)

    // 设置*csv.Reader对象的属性，例如字段分隔符和引用字符等。
    reader.Comma = ','       // 使用逗号作为字段分隔符
    reader.LazyQuotes = true // 允许不规则的引用字符

    // 通过调用*csv.Reader对象的方法从CSV文件中读取数据行。
    rows, err := reader.ReadAll()
    if err != nil {
        log.Fatal(err)
    }
    insertExp := getInsertExp(db, table_name)
    for row := 1; row < len(rows); row++ {
        tmpInsertExp := insertExp
        for col := 0; col < len(rows[row]); col++ {
            search := fmt.Sprintf(`%%$%d`, col)
            tmpInsertExp = strings.ReplaceAll(tmpInsertExp, search, rows[row][col])
        }
        _, err := db.Exec(tmpInsertExp)
        if err != nil {
            log.Fatal(err.Error())
        }
    }

    // 处理任何错误并关闭CSV文件。
    if err := file.Close(); err != nil {
        log.Fatal(err)
    }

    err = db.QueryRow("SELECT COUNT(*) FROM " + table_name).Scan(&count)
    if err != nil {
        return -1
    } else {
        return
    }
}

// getInsertExp 获取插入的语句
//  @param db 数据库
//  @param table_name 要插入表名
//  @return string 插入语句
func getInsertExp(db *sql.DB, table_name string) string {
    colsname := []string{}
    colstype := []string{}

    fmt.Println("SELECT column_name, data_type FROM information_schema.columns WHERE table_name ='" + table_name + "';")
    rows, err := db.Query("SELECT column_name, data_type FROM information_schema.columns WHERE table_name ='" + table_name + "';")
    fmt.Printf("rows: %v\n", rows)
    if err != nil {
        panic(err)
    }
    defer rows.Close()

    for rows.Next() {
        var columnName string
        var dataType string
        err = rows.Scan(&columnName, &dataType)
        if err != nil {
            panic(err)
        }
        fmt.Println(columnName, dataType)
        // 一般情况下，如果首个为id的话，默认为自动递增的编号，在添加数据时忽略，有数据库自行添加
        if columnName == "id" {
            continue
        }
        colsname = append(colsname, columnName)
        colstype = append(colstype, dataType)
    }

    // fmt.Sprintf("INSERT INTO known_file_table(file_name, idx, category, importance,count,md5,description) VALUES ( '%s', %s,%s,%s,%s,'%s','%s');",
    //     rows[i][0], rows[i][1], rows[i][2], rows[i][3], rows[i][4], rows[i][5], rows[i][6])
    exp := "INSERT INTO " + table_name + "("
    for i := 0; i < len(colsname); i++ {
        exp += colsname[i]
        if i != len(colsname)-1 {
            exp += ", "
        }
    }
    exp += ") VALUES ("
    for i := 0; i < len(colstype); i++ {
        switch colstype[i] {
        case "integer":
            exp += fmt.Sprintf(`%%$%d`, i)
        case "text":
            exp += fmt.Sprintf(`'%%$%d'`, i)
        default:
            exp += fmt.Sprintf(`%%$%d`, i)
        }
        if i != len(colstype)-1 {
            exp += ", "
        }
    }
    exp += ");"
    return exp
}
