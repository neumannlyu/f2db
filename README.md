# f2db
将CVS文件中的数据导入到数据库中。

用法：

```
f2db <-h -u -p -d> [-f|-t|-c|-i|-a]

    Usage:
      [1]从文件中导入,追加到数据库中 -h 127.0.0.1 -u usr -p 123456 -d db -t example_table -f example.csv -a
      [2]从文件中导入,并且清除原有数据 -h 127.0.0.1 -u usr -p 123456 -d db -t example_table -f example.csv -c
      [3]用文件中数据初始化数据库 -h 127.0.0.1 -u usr -p 123456 -d db -i
       按照文件与表对应规则：         按照文件与表对应规则进行数据填充：
         kft.csv <> known_file_table
         fst.csv <> file_suffix_table
         kvt.csv <> kernel_vulnerablity_table
  -a    [append] 在表后进行追加。
  -c    [clean] 清空表中的数据
  -d string
        [database] 远程数据库名
  -f string
        [file] 导入的文件。支持cvs
  -h string
        [host] 远程数据库主机的IP
  -i    [init] 初始化数据库。
  -p string
        [password] 远程数据库连接的密码
  -t string
        [table] 远程数据库表名
  -u string
        [user] 远程数据库连接的用户名
```

