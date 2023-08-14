## TLDB:high performance distributed database   [[中文文档]](https://github.com/donnie4w/tldb/blob/main/README_zh.md "[中文文档]")

------------

- Tldb has the features of high availability, high performance, no data loss, excellent horizontal expansion ability and so on。
- With its own web management platform, cluster status monitoring, parameter modification, data management operations can be completed in the management platform
- Support MQ. The implementation mechanism and network characteristics of tldb itself have all the features of MQ from the bottom
- Easy to maintain. Cluster status and node status are automatically adjusted, no network islanding occurs.
- When the disk of a node is fully written or incorrectly written, the node enters the proxy mode, which does not affect the operation of the client
- Tldb supports basic operations such as creating tables, indexes, and table fields through client operations.
- Tldb supports a large number of concurrent client operations, and can well support big data writing and reading.

------------

### TLDB:application scenes

- Suitable for scenarios where the business query logic is not complicated, such as orders, logistics,IM message bodies, and wallets
- Suitable for data warehouse
- Suitable for scenarios with a large number of MQ requirements
- Suitable for scenarios that require fast data entry and reading

------------

### Problems that TLDB can solve

- Solve the problem of concurrent read and write performance of large amounts of data
- Solve the problem of large amounts of MQ information subscription publishing
- Solves problems that require rapid cluster horizontal scaling
- Solves problems where data needs to be checked back at different points in time

------------

### TLDB:technical feature

- tldb logs record data change trajectories and supports data restoration to any point in time
- The data engine uses leveldb currently;leveldb has efficient read and write performance and stability; In a specific configuration environment, it even has the advantage of millions of data seconds into the database
- Supports data compression, greatly optimizing storage space
- tldb uses the consensus mechanism, two-stage commit confirmation, and binlog to ensure data consistency among all nodes
- The node proxy mechanism ensures that abnormal nodes do not affect data operations
- The cluster uses the consensus algorithm to hash data for storage, and the number of storage node can be specified by parameters.
- Compression protocol, aggregation protocol sending and receiving, object pool optimization, tldb has excellent performance, and supports a large number of client connection concurrent operations.
- tldb supports MQ from the design of the architecture, and cluster MQ data is consistent without loss
- tldb MQ solves the problems of MQ message loss, repeated consumption, and message backlog

### TLDB:characteristics of data

- Support field index
- Field indexing is supported and the data table automatically generates a 64-bit self-growing uniquely identified ID key
- MQ data is essentially tldb data and also automatically generates ID keys

------------

### Advantages and disadvantages of TLDB compared to other distributed databases：

- tldb is easy to use, requires no installation, and has almost no maintenance costs
- tldb cluster environment is simple to use and no different from the single-node environment
- tldb data is not lost, and it is automatically split and compressed for backup
- Backup, compression, and synchronization of tldb data are automatically completed. To restore imported data, you only need to import data files on the management platform
- tldb search function is relatively weak and does not support joint indexes
- tldb data types are not as rich as relational databases

------------

### tldb：Related URL：

- Online test：http://test.tlnet.top
- Source：https://github.com/donnie4w/tldb


### download tldb

1. [version 0.0.1](http://tlnet.top/download "version 0.0.1")


### tldb database client program：

1. go   <https://github.com/donnie4w/tlcli-go>
2. java <https://github.com/donnie4w/tlcli-j>  
3. python <https://github.com/donnie4w/tlcli-py> 

### tldb orm program：

1. go https://github.com/donnie4w/tlorm-go
2. java https://github.com/donnie4w/tlorm-java


### tldb MQ client program：

1. go   <https://github.com/donnie4w/tlmq-go>
2. java <https://github.com/donnie4w/tlmq-j>  
3. python <https://github.com/donnie4w/tlmq-py> 
4. js <https://github.com/donnie4w/tlmq-js>

------------

### TLDB Startup Overview

####  Take linux as an example

###### Single-machine startup

1.  ./tldb  -clus=0 
	. When the startup parameter clus is 0, tldb is enabled in single-node mode (cluster mode by default)

2.  ./tldb  -clus=0 -mq=:5000 -admin=:4000 -dir=_data -cli=:7000
	. Binding mq port (-mq), management background port (-admin), database client port (-cli), and data file address (-dir)

---------

###### Cluster startup: The following uses starting three nodes as an example

- ./tldb -cs=:6001 -mq=:5001 -admin=:4001 -dir=_data1 -cli=:7001
- ./tldb -cs=:6002 -mq=:5002 -admin=:4002 -dir=_data2 -cli=:7002
- ./tldb -cs=:6003 -mq=:5003 -admin=:4003 -dir=_data3 -cli=:7003

###### Cluster node connection

- login in Management platform of any node， from cluster env -> Add cluster nodes Enter the cluster service addresses of other nodes：e.g. :6003 or 192.168.1.100:6000
- After adding other node addresses in the Management platform of a node, all the cluster nodes will automatically synchronize information with each other (there is no need to add cluster nodes to each other repeatedly).

------------


###### parameter description：

- -cs Cluster node interconnection address
- -mq MQ service address of the node
- -admin Management platform service address of a node
- -dir Data file address of the node
- -cli Client connection address of the node

###### format of service address

- There are four parameters that can be bound the service address: -cs, -mq, -admin, and -cli
- The address format is: domain(or IP)+":"+ port, where the domain(or IP) is the domain(or IP) bound to access, that is, if the domain(or IP) is bound, other domain(or IP) cannot be accessed. For example, -admin=db.tlnet.top:4001
- You can use ":"+port without binding the domain(or IP), for example, -admin=:4000

###### other main parameters description:

- -clitls=1 Enable the TLS service for client. The default is 0, which means disabled
- -mqtls=1 Enable the TLS service for MQ. The default is 0, which means disabled
- -admintls=1 Enable the TLS service for Management platform. The default is 0, which means disabled
- For other parameters, see Management Background or http://tldb.tlnet.top

------------

####  Other operating systems: windows, macos, freesd, solaris use tldb

######  The above operating systems use tldb no different from linux, for example:

###### The following uses windows single-node operation as an example:

	tldb.exe  -clus=0 -mq=:5000 -admin=:4000 -dir=_data -cli=:7000