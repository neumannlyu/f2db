package sql

type ISQL interface {
    // 打开数据库连接
    // !在调用完数据库后一定要使用Close关闭连接。
    //  @param h host
    //  @param u user
    //  @param p password
    //  @param d database name
    //  @return *sql.DB db指针
    //  @return error 错误
    Open(h, u, p, d string) bool

    // 关闭数据库连接
    // !在调用完数据库后一定要关闭连接。
    Close()

    // GetDB() *sql.DB
    // 执行表达式
    Execute(exp string) bool

    // 从CVS文件中导入数据
    ImportFromCVS(tableName, cvsFile string) int
}
