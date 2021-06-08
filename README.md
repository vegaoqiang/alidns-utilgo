# alidns-utilgo
阿里DNS命令行工具，支持对阿里云域名解析记录的添加，修改，删除和查询
## 安装
目前仅支持在类Unix平台运行，不支持Windows
### 下载二进制文件
```
chmod +x alidns-utilgo
```

### 自行编译
+ Linux可执行文件
```
CGO_ENABLED=0  GOOS=linux  GOARCH=amd64  go build main.go -o alidns-utilgo
```

+ macOs可执行文件
```
CGO_ENABLED=0 GOOS=darwin  GOARCH=amd64  go build main.go -o alidns-utilgo
```

## 使用方法
### -init --init
初始化账户信息，输入的access_key, access_secret和region并写入到文件，你也可以手动在当前用户home目录下创建`$HOME/.alidns-utilgo/config.json`文件并写入如下内容
```
{
    "access_key": "xxx",
    "access_secret": "xxx",
    "region": "xxx"
}
```
### -add --add
在指定的域名下添加一个新解析，搭配`-dn`一起使用，未指定`-dn`将报错
+ 在域名foo.com下添加一个bar解析，解析值为127.0.0.1
```
alidns-utilgo -add -dn bar.foo.com -v 127.0.0.1
```
> 默认的解析类型为A记录，可以通过--type=xx指定其他解析类型

### -del --del
***<font color=red>危险操作，谨慎使用</font>***<br>
删除域名下的主机记录对应的解析记录，需和`-dn`一同使用.
+ 删除foo.com域名下bar主机记录对应的所有解析记录
```
alidns-utilgo -del -dn bar.foo.com
```
+ 删除foo.com域名下bar主机记录对应的CNAME解析记录
```
alidns-utilgo -del -dn bar.foo.com -t CNAME
```
> <font color=Orange>删除主机记录时，如果未通过-t指定解析类型，将会删除该主机记录的所有解析类型。例如：bar记录同时存在MX,A,CNAME这三种解析类型，如未指定其中之一，将全部删除所有解析类型。</font>

### -update --update
更新域名解析，需要配合`-dn`和`-u`参数使用，指定的`-dn`必须能确定一个唯一的解析记录，否则报错。参数`-u`的值需要按照key=value的键值对形式提供，多个键值对以 ***`,`*** 分割,支持的key类型如下
```
value=xx // 解析记录值
rr=xx    // 主机记录
type=xx  // 解析类型  
```
+ 修改`bar.foo.com`为`abc.foo.com`
```
alidns-utilgo -update -dn bar.foo.com -u rr=abc
```
> 当`-dn=bar.foo.com`不能确定唯一解析记录时，使用-t指定解析记录类型尝试确定一个唯一解析记录。

### -list --list 
默认列出账户中所有的域名，当同时指定`-dn`，如`-dn=bar.foo.com`，将列出域名`foo.com`的`bar`解析详细信息，如果`-dn=foo.com`,则会列出`foo.com`域名下的所有解析。同时，`--list`还支持通过同时添加`--search`参数指定搜索关键字，按照`“%KeyWord%”`模式搜索，不区分大小写。`search`参数功能根据是否指定了`-dn`而不同,未同时指定`-dn`, 则`search`参数值搜索域名列表，同时指定了`-dn`,则`search`参数值在指定的`dn`下搜索所有相近的解析。注意：`-dn`的参数值应当为`foo.com`类型，如果未`bar.foo.com`,`search`参数值优先级高于`bar`,优先根据`search`参数值进行搜索`foo.com`下所有相近的解析. 当对某个域名下的所有解析进行搜索时, `search`参数值可以是域名解析的`Type`,`RR`,`Value`.
+ 列出账户中所有域名
```
alidns-utilgo --list
```

+ 在账户中搜索包含foo的域名
```
alidns-utilgo --list -s foo
```
+ 列出域名foo.com详细

```
alidns-utilgo --list -dn foo.com
```
+ 列出域名foo.com下包含bar关键字的解析详细

```
alidns-utilgo --list -dn bar.foo.com
```

+ 在域名foo.com下搜索包含bar关键字的解析

```
alidns-utilgo --list -dn foo.com -s bar
```
+ 在域名foo.com下搜索包含abc关键字的解析, 同时bar将不起作用
```
alidns-utilgo --list -dn bar.foo.com -s abc
```
> 在域名下搜索解析时，通过-s参数指定的关键字类型可以是：解析类型Type, 主机记录RR, 解析值Value