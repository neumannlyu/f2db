package sql

import (
    "database/sql"
    "encoding/csv"
    "fmt"
    "log"
    "os"
    "strings"

    _ "github.com/lib/pq"
)

type PostgresSQL struct {
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
func (pg *PostgresSQL) Open(h, u, p, d string) bool {
    db, err := sql.Open("postgres",
        fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
            h, u, p, d))
    if err != nil {
        fmt.Printf("open PostgresSQL failed err.Error(): %v\n", err.Error())
        return false
    } else {
        fmt.Println("open PostgresSQL successfully.")
        pg.DBPtr = db
        return true
    }

}

// 关闭数据库连接
// !在调用完数据库后一定要关闭连接。
func (pg *PostgresSQL) Close() {
    if pg.DBPtr != nil {
        pg.DBPtr.Close()
    }
    pg.DBPtr = nil
    fmt.Println("postgres closed.")
}

// 执行sql语句
func (pg PostgresSQL) Execute(exp string) bool {
    _, err := pg.DBPtr.Exec(exp)
    if err != nil {
        return false
    } else {
        return true
    }
}

// func (pg PostgresSQL) GetDB() *sql.DB {
//     return pg.DBPtr
// }

func (pg PostgresSQL) ImportFromCVS(tableName, cvsFile string) int {
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
    exp := pg.genInsertExpModel(tableName)
    for row := 1; row < len(rows); row++ {
        tmpInsertExp := exp
        for col := 0; col < len(rows[row]); col++ {
            // ! 为了兼容mysql。
            // ! mysql \会被解析为转译，但是在postgres会直接当作字符
            // ! 将所有的\\替换为\
            rows[row][col] = strings.ReplaceAll(rows[row][col], "\\\\", "\\")
            search := fmt.Sprintf(`##{%d}##`, col)
            tmpInsertExp = strings.ReplaceAll(
                tmpInsertExp, search, rows[row][col])
        }
        // fmt.Println(tmpInsertExp)
        _, err := pg.DBPtr.Exec(tmpInsertExp)
        if err != nil {
            log.Fatal(err.Error())
        }
    }

    // 处理任何错误并关闭CSV文件。
    if err := file.Close(); err != nil {
        log.Fatal(err)
    }

    var count int
    err = pg.DBPtr.QueryRow("SELECT COUNT(*) FROM " + tableName).Scan(&count)
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
func (pg PostgresSQL) genInsertExpModel(table_name string) string {
    colsname := []string{}
    colstype := []string{}
    rows, err := pg.DBPtr.Query(
        "SELECT column_name, data_type FROM information_schema.columns" +
            " WHERE table_name ='" + table_name + "';")
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

        //* 一般情况下，如果首个为id的话，默认为自动递增的编号，
        // 在添加数据时忽略，有数据库自行添加
        if columnName == "id" {
            continue
        }

        colsname = append(colsname, columnName)
        colstype = append(colstype, dataType)
    }
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
            exp += fmt.Sprintf(`##{%d}##`, i)
        case "text":
            exp += fmt.Sprintf(`'##{%d}##'`, i)
        default:
            exp += fmt.Sprintf(`##{%d}##`, i)
        }
        if i != len(colstype)-1 {
            exp += ", "
        }
    }
    exp += ");"
    return exp
}
