package sql

import (
    "database/sql"
    "encoding/csv"
    "fmt"
    "log"
    "os"
    "strings"

    _ "github.com/go-sql-driver/mysql"
)

type MySQL struct {
    DBPtr *sql.DB
}

// OpenPostgresDB 打开数据库连接
// !在调用完数据库后一定要使用ClosePostgresDB关闭连接。
//  @param h host
//  @param u user
//  @param p password
//  @param d database name
//  @return *sql.DB db指针
//  @return error 错误
func (m *MySQL) Open(h, u, p, d string) bool {
    db, err := sql.Open("mysql",
        //"user:password@tcp(host:port)/dbname"
        fmt.Sprintf("%s:%s@tcp(%s)/%s", u, p, h, d))
    if err != nil {
        fmt.Printf("open MySQL failed err.Error(): %v\n", err.Error())
        return false
    } else {
        fmt.Println("open MySQL successfully.")
        m.DBPtr = db
        return true
    }

}

// 关闭数据库连接
// !在调用完数据库后一定要关闭连接。
func (m *MySQL) Close() {
    if m.DBPtr != nil {
        m.DBPtr.Close()
    }
    m.DBPtr = nil
    fmt.Println("mysql closed.")
}

// 执行sql语句
func (m MySQL) Execute(exp string) bool {
    _, err := m.DBPtr.Exec(exp)
    if err != nil {
        return false
    } else {
        return true
    }
}

// func (m MySQL) GetDB() *sql.DB {
//     return m.DBPtr
// }

func (m MySQL) ImportFromCVS(tableName, cvsFile string) int {
    // 打开CSV文件并创建一个*csv.Reader对象。
    file, err := os.Open(cvsFile)
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
    // fmt.Printf("rows: %v\n", rows)
    exp := m.genInsertExpModel(tableName)
    for row := 1; row < len(rows); row++ {
        tmpInsertExp := exp
        for col := 0; col < len(rows[row]); col++ {
            search := fmt.Sprintf(`##{%d}##`, col)
            tmpInsertExp = strings.ReplaceAll(tmpInsertExp, search, rows[row][col])
        }
        // fmt.Println(tmpInsertExp)
        _, err := m.DBPtr.Exec(tmpInsertExp)
        if err != nil {
            log.Fatal(err.Error())
        }
    }

    // 处理任何错误并关闭CSV文件。
    if err := file.Close(); err != nil {
        log.Fatal(err)
    }

    var count int
    err = m.DBPtr.QueryRow("SELECT COUNT(*) FROM " + tableName).Scan(&count)
    if err != nil {
        return -1
    } else {
        return count
    }
}

// getInsertExp 获取插入的语句
//  @param db 数据库
//  @param table_name 要插入表名
//  @return string 插入语句
func (m MySQL) genInsertExpModel(table_name string) string {
    query := fmt.Sprintf("SELECT COLUMN_NAME,DATA_TYPE FROM information_schema.COLUMNS WHERE TABLE_NAME = '%s'", table_name)
    rows, err := m.DBPtr.Query(query)
    if err != nil {
        return err.Error()
    }
    defer rows.Close()

    names := make([]string, 0)
    types := make([]string, 0)
    for rows.Next() {
        var columnName string
        var columnType string
        if err := rows.Scan(&columnName, &columnType); err != nil {
            return err.Error()
        }
        //* 一般情况下，如果首个为id的话，默认为自动递增的编号，在添加数据时忽略，有数据库自行添加
        if columnName == "id" {
            continue
        }
        names = append(names, columnName)
        types = append(types, columnType)
    }
    /*
       INSERT INTO mytable (file_name, flag, rating, count, description)
       VALUES ('filename1.txt', 1, 5, 10, 'This is the first file.'),
               ('filename2.txt', 0, 3, 7, 'This is the second file.'),
               ('filename3.txt', 1, 4, 12, 'This is the third file.');
        INSERT INTO mytable (column1, column2, column3) VALUES (?, ?, ?)"
    */
    exp := "INSERT INTO " + table_name + "("
    for i := 0; i < len(names); i++ {
        exp += names[i]
        if i != len(names)-1 {
            exp += ", "
        }
    }
    exp += ") VALUES ("
    for i := 0; i < len(types); i++ {
        switch types[i] {
        case "text":
            exp += fmt.Sprintf(`'##{%d}##'`, i)
        case "int":
            exp += fmt.Sprintf(`##{%d}##`, i)
        }
        if i != len(types)-1 {
            exp += ", "
        }
    }
    exp += ")"
    return exp
}
