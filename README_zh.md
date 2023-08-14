## TLDB  高性能分布式数据库	 [[English document]](https://github.com/donnie4w/tldb/blob/main/README.md "[English document]")

------------

- tldb 具备高可用、高性能，数据不丢失，极好的水平扩展能力等特性。
- 自带web管理后台，集群状态监测，参数修改，数据管理操作等等均可在后台完成
- 支持MQ。tldb本身的实现机制与网络特性从底层具备了MQ所有特性。
- 极易维护。集群状态与节点状态自动调整，不出现网络孤岛现象。
- 节点磁盘写满或写入错误时,节点进入代理模式,不影响客户端的操作
- tldb数据通过客户端操作，支持建立表,索引,表字段等基础操作。
- tldb支持大量的客户端并发操作，可以很好应对大数据写入与读取。

------------

### TLDB 适用场景

- 适用业务查询逻辑简单的场景,如 订单,物流,IM消息体,钱包等业务场景
- 适用于数据仓库
- 适用大量MQ要求的场景
- 适用大量数据库客户端并发读写的场景
- 适用需要数据快速入库及读取的场景

------------

### TLDB 能解决的问题

- 解决大量数据并发读写性能问题
- 解决大量MQ信息订阅发布的问题
- 解决需要快速集群水平扩展的问题
- 解决需要在不同时间点回查数据的问题

------------

### TLDB 技术特点

- tldb日志记录数据变化轨迹，支持将数据还原到任意时间点
- 数据引擎目前使用leveldb，leveldb本身具备高效的读写性能与稳定性；在特定配置环境下，甚至具备百万数据秒级入库的优异性能
- 数据引擎数据压缩与tldb提供的数据压缩，极大优化了存储空间。
- tldb通过共识机制，两阶段提交确认和binlog日志 保证了所有节点数据的一致性。
- 通过节点代理机制，保证了异常节点不影响数据操作
- 集群通过共识算法对数据进行散列存储，存储份数支持由参数指定。
- 压缩协议，聚合协议收发，对象池优化，使tldb有优异性能，并支持大量的客户端连接并发操作。
- tldb从架构底层支持MQ，集群MQ数据一致，不丢失
- tldb MQ解决MQ消息丢失、重复消费，消息积压等问题

### TLDB 数据特点

- 支持字段索引
- 数据表自动生成64位自增长唯一标识的ID键
- MQ数据本质是数据表数据，也自动生成ID键

------------

### TLDB相较其他分布式数据库的优缺点：

- tldb极易使用，无需安装，几乎没有维护成本
- tldb集群环境使用简单，与单机没有区别
- tldb数据不丢失，自动切分压缩备份
- tldb数据的备份压缩同步都自动完成，恢复导入数据也只需要在后台导入数据文件既可
- tldb搜索功能相对较弱，不支持联合索引
- tldb数据类型没有关系型数据库丰富

------------


### tldb相关网址：

1. 后台测试链接：[http://dbtest.tlnet.top](http://dbtest.tlnet.top/ "http://dbtest.tlnet.top")
2. 源码链接：[https://github.com/donnie4w/tldb](https://github.com/donnie4w/tldb "https://github.com/donnie4w/tldb")
3. 官方网站 ： [http://tldb.tlnet.top](http://tldb.tlnet.top "http://tldb.tlnet.top")

### tldb 客户端：

1. go   <https://github.com/donnie4w/tlcli-go>
2. java <https://github.com/donnie4w/tlcli-j>  
3. python <https://github.com/donnie4w/tlcli-py> 

### tldb MQ 客户端：

1. go   <https://github.com/donnie4w/tlmq-go>
2. java <https://github.com/donnie4w/tlmq-j>  
3. python <https://github.com/donnie4w/tlmq-py> 
4. js <https://github.com/donnie4w/tlmq-js>

------------

### TLDB 启动

####  以linux为例

###### 单机启动

1.  ./tldb  -clus=0 
	. 直接运行时启动默认端口，启动参数clus等于0时，tldb启动单机模式，默认是集群模式

2.  ./tldb  -clus=0 -mq=:5000 -admin=:4000 -dir=_data -cli=:7000
	. 可以指定mq端口(-mq)，管理后台端口(-admin),数据库客户端端口(-cli),数据文件地址(-dir)

---------

###### 集群启动：以启动3个节点为例

	节点 1
	./tldb  -cs=":6001" -mq=":5001" -admin=":4001" -dir="_data1" -cli=":7001"
	节点 2
	./tldb  -cs=":6002" -mq=":5002" -admin=":4002" -dir="_data2" -cli=":7002"
	节点 3
	./tldb  -cs=":6003" -mq=":5003" -admin=":4003" -dir="_data3" -cli=":7003"

###### 集群节点相关联

1. 打开其中任意一个节点的管理后台，注册登录后 由 集群环境 -> 增加集群节点并连接 中添加其他节点的集群服务地址 如 :6003 或 192.168.1.100:6000
2. 在一个节点后台添加其他节点地址之后，不同集群节点之间就会自动相互同步信息，(无需重复相互添加集群节点)


------------


###### 启动参数说明：

- -cs 节点集群服务地址，及节点互联地址(单机运行无需设置)
- -mq 节点MQ服务地址
- -admin 节点管理后台服务地址
- -dir 节点的数据文件地址
- -cli 节点客户端链接地址

###### 服务地址说明：

- 有4个参数是指定服务地址的，分别为-cs，-mq，-admin，-cli
- 地址格式为：域名(或IP)+":"+端口，其中域名(或IP)为绑定访问的域名(或IP)，即如果指定了域名(或IP)时，其他的域名(或IP)则无法访问. 如：-admin=db.tlnet.top:4001
- 可以不绑定域名(或IP)，直接使用 ":"+端口 如：-admin=:4000

###### 其他主要参数说明：

- -clitls=1 		节点客户端服务开启TLS服务地址，默认0
- -mqtls=1 		节点MQ服务开启TLS服务地址，默认0
- -admintls=1 		节点管理后台开启TLS服务地址，默认0
- 其他参数请参考管理后台或<http://tldb.tlnet.top>


------------

####  其他操作系统：windows ，macos ，freesd ，solaris 使用tldb

###### 以上操作系统使用tldb与linux无区别，如：

###### windows 单机运行为例：

	tldb.exe  -clus=0 -mq=:5000 -admin=:4000 -dir=_data -cli=:7000
