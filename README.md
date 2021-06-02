### -add --add
在指定的域名下添加一个新解析，搭配`-dn`一起使用，未指定`-dn`将报错
+ 在域名foo.com下添加一个bar解析，解析值为127.0.0.1
```
alidns-utilgo -add -dn bar.foo.com -v 127.0.0.1
```
> 默认的解析类型为A记录，可以通过--type=xx指定其他解析类型

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