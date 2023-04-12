# f2db
将CVS文件中的数据导入到数据库中。

## Build

```
go build main.go 
```

## Usage

```
Usage of params:
    f2db [-f|-t|-c|-i|-a]

    Usage:
      [1]从文件中导入,追加到数据库中 f2db -t table_name -f example.csv -a
      [2]从文件中导入,并且清除原有数据 f2db -t table_name -f example.csv -c
      [3]用文件中数据初始化数据库 f2db -i
         按照json文件中对应文件名进行匹配。

  -a    [append] 在表后进行追加。
  -c    [clean] 清空表中的数据
  -f string
        [file] 导入的文件。支持cvs
  -i    [init] 初始化数据库。
  -t string
        [table] 远程数据库表名
```




