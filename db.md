1. Redis不存在表的概念，一个Redis实例中，直接存储5大数据类型：字符串(string)、列表(list)、哈希(hash)、集合(set)、有序集合(zset)
2. Redis并没有数字类型，转换成了字符串，包括有序集合(zset)中的score
3. Redis保证单个命令的原子性，而Redis的事务并没有作任何原子性的保证，且事务中各个指令互不影响，既不会回滚，也不会终止后续指令的执行，感觉Redis事务更像批处理
4. Redis的save/bgsave命令备份当前数据库数据到dump.rdb文件，恢复操作是将该文件copy到Redis安装目录，直接启动服务(如果需要多个Redis服务怎么办。。。)
5. Redis分区是将数据放到多个实例中(好奇Redis的集群和分区的区别)
6. Redis集群可以在运行时进行动态的伸缩性调整，灵活性比较高
7. 非关系型数据库方便处理相互之间没有耦合性的非结构化数据
8. MongoDB的读写性能(只能达到1kqps??)
9. 分区和索引哪个优先执行？
	查询条件子句中含有分区条件，则先过滤分区，再在对应分区执行索引操作，否则在所有分区执行索引操作。因此，在实践中尽量避免分区列和索引列的不匹配，或者条件子句中包含过滤分区的条件
10. 如何实现读写分离？
	目前mysql读写分离通常采用的做法是主节点处理写操作，并异步复制更新数据到从节点，而从节点处理select等查询操作，总结就是主从复制+读写分离。读写分离的好处：
	1) 增加冗余备份,提高可用性
	2) 一主多从或多主多从,减轻单一节点压力,实现负载均衡,提升系统性能
	主从结构存在的问题是主从之间异步复制导致的数据不一致,既主从间的数据延迟问题
11. 数据量很大？ rownumber ranknumber？？
12. 两个索引是不是一定用得上？
	不一定，包含但不限于下列情况索引将会失效：
	1) 查询子句的条件字段使用了函数
	2) 查询子句的条件字段运用了数学运算
	3) like子句的首字符为通配符
	4) 带or的子句含有非索引字段
	5) 索引字段隐式转换
	6) mysql估计使用索引比全表扫描慢
	7) order by子句,排序条件不是查询表达式时才会使用索引
	8) join子句,外键和主键的数据类型不同
	9) 复合索引中单独引用非第一个索引字段
	10) not in和<>操作
	11) 索引字段为字符串,但使用时未加引号
13. mysql分区类型：range分区(partition by range())、list分区(partition by list(),与range分区类似,区别是list分区的列的值是某列表中的一个,而range的值是连续的)、hash分区(partition by hash(),确保要分区的数据平均分布到各个分区)、key分区(partition by key(),与hash分区类似,区别是hash函数为mysql提供)
14. mysql创建分区时自动对分区表从0开始编号
15. mysql delimiter设置分隔符,默认为分号,可以通过delimiter重新设置,与存储过程并无必然联系,一般在执行含分号的多条语句时使用,比如自定义函数、存储过程或者触发器(慎用,太消耗资源)
16. mysql存储过程是一组完成特定功能的sql语句块,经过预编译保存在进程字典中,常用语批量处理一些重复性高的操作;
    mysql5.7中information_schema.routines表查看存储过程信息
    create procedure proc_name(parameters)
    begin
    	some valid sql statements
    end
    存储过程默认绑定当前数据库,若需指定某数据库则在存储过程名前加数据库名即可
    存储过程参数分为输入参数IN,输出参数OUT,输入输出参数INOUT,默认为IN类型,声明方式为:[IN|OUT|INOUT] para_name datatype
    常用的数据类型datatype有:int float varchar()
    sql语句块请参考sql教程,常用的有变量定义赋值相关的set,declare,控制语句if - then - else,case及循环语句,其它的还有字符串类的concat,substring,length,数学类的abs,floor,format,rand,round,日期时间类的current_date,current_time,current_timestamp,date,year,month等,具体的函数可以在需要的时候查询具体的参数和用法
    存储过程还可以使用prepare,execute,deallocate|drop预处理语句,用法如下:
	set @table = 'tables';
    set @sql_str = concat('select count( * ) from ', @table);
    prepare schema_table from @sql_str; execute schema_table; drop prepare schema_table;
    调用存储过程使用call proc_name
    删除存储过程使用drop [if exists] proc_name
17. mysql有哪些索引类型？
	mysql索引由数据库引擎实现,所以不同的引擎实现的索引类型及实现方式都有所区别,索引是表空间的一部分,索引信息存储在information_schema.innodb_sys_indexes
	以下内容参考自:https://segmentfault.com/q/1010000003832312
	从数据结构角度分为BTree索引、hash索引、空间索引
	从索引存储角度分为聚簇索引和非聚簇索引
	从逻辑角度分为主键索引(primary key)、普通/单列索引(index)、唯一索引(unique)、全文索引(fulltext)、组合索引
	hash索引以hash形式组织索引结构,每个键对应一个值,散列分布,所以单个查询速度很快,不支持范围查找和排序
	BTree或B+Tree索引,应用广泛,以平衡树的形式来组织,适合排序、范围查找,关于BTree和B+Tree的详细介绍可参考:http://blog.codinglabs.org/articles/theory-of-mysql-index.html
18. memcached只支持字符串类型,而Redis提供五种数据类型,且每种数据类型都有专属命令
19. Redis在提高可用性上采用主从异步复制结构,且并不保证数据的强一致性,也即是主节点在给客户端回复之后,才会向从节点发送写操作(可能会在未来提供同步写方法),另一中可能导致不一致的是集群出现网络分区,部分节点被孤立
20. Redis只有集群模式运行的节点才能组成集群,普通节点无法组成集群,集群模式运行需要开启redis.conf中的 cluster-enabled yes
21. 明确缓存的应用场景:一般用在不经常变化或者变化之后一定时间内影响不会太大的数据,那些实时性、一致性要求高的地方比如12306出票不能用缓存,毕竟缓存的目的是为了提高效率
22. Redis为什么这么快?
23. 为什么说Redis是单线程?
24. Redis提供了发布订阅功能,不过存在两个问题,首先是Redis系统稳定性,若订阅者处理消息过慢,会导致消息堆积,输出缓冲区变大,将导致Redis速度变慢,甚至崩溃;其次是网络传输可靠性,若订阅过程中发生断线,则会丢失断线期间的消息
25. Redis的基本事务暂未提供事务执行过程中的回滚操作
26. Redis提供的两种持久化方法：快照,类似于定时写,可以将某一时刻的数据写入硬盘;只追加文件,将每次执行的写命令复制到硬盘
27. 写文件到硬盘的过程:内存-->缓冲区-->硬盘,Redis同步AOF文件appendfsync选项为no时,并不会显式执行同步操作,而是由操作系统决定何时同步,AOF持久化的问题在于不断增大的文件体积
28. Redis复制的启动过程
	主服务器:
		接收从服务器的sync请求;
		执行bgsave,创建快照,缓冲区记录新的写命令;
		发送快照文件,继续记录新的写命令;
		发送缓冲区里的写命令;
		上述执行完,后续每执行一个写命令,向从服务器同步相同的写命令
	从服务器:
		发送sync请求;
		根据现有配置决定如何反馈客户端(直接返回错误或者根据现有数据相应);
		丢弃旧数据(如果有的话),载入主服务器发来的快照文件;
		完成快照的解释操作,开始接受请求;
		执行发送过来的写命令(包括主服务器缓冲区及新的)
29. Redis的watch属于乐观锁,也即watch会在发现其它客户端修改数据的情况下才会通知执行watch的客户端,
    watch并不会阻止其它客户端对数据进行修改
30. 自己构建Redis锁时除了获取、释放锁之外,还需考虑的几个问题：
	a) 若取得锁的进程执行时间过长,触发锁超时,锁被自动释放,进程执行完后可能会错误的释放掉其它进程持有的锁
	b) 持有锁的进程崩溃,若其它进程无法检测到崩溃进程及其状态,则只能等待锁超时
	c) 锁超时被自动释放,若有多个进程同时申请锁,且都得到锁(在高负载的情况下,很容易出现)
31. watch命令锁住Redis的整个键(与mysql锁表相同),因此频繁操作时影响执行的性能,但这并不是说粗粒度锁一无是处,
    因为使用多个细粒度锁时,也会引发死锁风险
32. mongod启动
    mongod --dbpath xxx --logpath yyy --fork --auth --bind_ip zzz
33. mongo连接远程数据库
	mongo ip
	mongo ip:port
	mongo ip:port/db_name
	mongo ip:port/db_name -u xxx -p yyy
    示例: mongo 192.168.1.100:27027/thd_game
34. 十一放假期间公司调整网络,回来发现mongo已关闭,重启,报错 error number 1,搜了一下发现原因是mongod非正常关闭导致的
    删掉$dbpath下的mongod.lock及$logpath下的日志文件,重启mongod即可
    查询mongo文档(https://docs.mongodb.com/manual/tutorial/manage-mongodb-processes/),正确关闭mongo有如下几种方式:
      use shutdownServer() 登入mongo shell,键入 
        use admin
        db.shutdownServer()
      use --shutdown linux command line,键入
        mongod --shutdown
      use CTRL-C 以交互模式(interactive mode)运行mongod实例时适用此方法
      use kill linux command line,键入 kill mongod_process_id 或者 kill -2 mongod_processed_id (决不能使用kill -9) 
35. mongo首次运行时并未开启验证,登录后需要首先创建一个超级管理员用户,用于管理其他用户
    由手册可知,mongo3.0以后,在admin中创建(As of MongoDB 3.0, with the localhost exception, you can only create users 
    on the admin database),创建命令为db.createUser({user:"xxx", pwd:"yyy", roles:[{},...]})
    createUser的roles选项,指定了被创建用户的"角色",包括内置角色和自定义角色,mongo采用基于角色授权的方式管理对数据库的相关操作
    为了便于管理,首个超级管理员账号,可以选择内置角色中属于superuser roles的root角色,root拥有所有权限(provides full privileges on all resources)
    localhost exception which allows you to create a user administrator in the admin database
    管理员的创建流程可参考mongodb doc的Enable Auth(https://docs.mongodb.com/manual/tutorial/enable-authentication/)
36. 为project创建database和user时,可选择内置角色中属于Database Administration Roles的dbAdmin,dbOwner
    可由db.createRole()自定义角色类型,db.createRole({role,writeConcern})
    role结构为:{role:"role_name", privileges:[],roles:[]}
    privileges字段包含对角色的各种授权,授权结构为:{resource:{}, actions:["action1",...]}
    privileges.resource结构为:{ db: <database>, collection: <collection> } 或者 { cluster : true }, 指定了database,collection及集群相关
    privileges.actions包含授权的查询和写入等操作,包括:find,insert,remove,update等
    roles字段,指定从其他角色及其数据库继承的授权
    createUser示例如下:
    	use admin
	db.CreateUser(
		{role: "thd_adm",
		 privileges: [{{ resource: { db: "thd_game", collection: "" }, actions: [ "find", "update", "insert", "remove" ]}}],
		 roles: [{role:"root", db:"admin"}]
		 }, 
		 { w: "majority" , wtimeout: 5000 }
	)
37. 创建User时的database为用户的authentication database, 两种方式实现验证登录:
	mongo命令行 --auth -u xxx -p yyy --authenticationDatabase zzz
	mongo直接连接mongd/mongos --> use user_authentication_database  db.auth(user_name,pwd)
38. 搭建外网数据库环境时用到的几个命令
	安装mysql客户端: yum install mysql
    	mysql导入客户端本地文件: mysql -h xxx -uyyy -pzzz db_name < local_sql_path
	source命令同样可以导入文件,不过需要先登录到数据库终端,source命令只能导入服务端本地文件 
	mysql shell中LOAD DATA INFILE读取文件:LOAD DATA LOCAL INFIEL 'xxx.txt' INTO TABLE yyy;
	LOAD DATA的local字段可以指定从客户机路径导入文件,如果没有指定则从服务器路径导入
39. centos yum 安装mongodb
	/etc/yum.repos.d/下,创建mongodb的repo文件 mongodb-org-3.6.repo,文件内容为:
	[mongodb-org-3.6]
	name=MongoDB Repository
	baseurl=https://repo.mongodb.org/yum/redhat/$releasever/mongodb-org/3.6/x86_64/
	gpgcheck=0
	enabled=1 
	保存退出,执行: yum install -y mongodb-org 
40. mysql使用group_concat出现下面的结果:
	@str_sql := CONCAT('DROP TABLE IF EXISTS ', GROUP_CONCAT(`table_name`))
	NULL
	语句是:
	SELECT @str_sql := CONCAT('DROP TABLE IF EXISTS ', GROUP_CONCAT(`table_name`))
    	FROM `information_schema`.`tables`
    	WHERE `table_schema` = database() AND `table_name` LIKE pattern;
	经过试验,NULL是执行的结果,首次初始化时,没有符合条件的表,所以返回的是NULL,再次执行,则会返回@str_sql,值类似于:
	@str_sql := CONCAT('DROP TABLE IF EXISTS ', GROUP_CONCAT(`table_name`))
	DROP TABLE IF EXISTS xxx_pattern,yyy_pattern,zzz_pattern
  注:Unless otherwise stated, group functions ignore NULL values. 也就是说group_concat会自动忽略NULL值
41. mysql varchar类型执行+-等数学运算时,会自动隐式转换,原则是:字符开头的一律为0,数字开头的直接截取到第一个不是字符的位置
42. 日志报错:duplicate entry 'xxx' for key 'PRIMARY',主键冲突,再次插入数据前,先删除已有的数据,否则用update
43. mysql> delete table 并不会清除主键的auto_increment计数,drop table 可以
	delete table 不能批量删除表,比如我用concat和group_concat拼接具有相同前缀的表删除语句,提示syntax error
44. 查看auto_increment: select table_shcema,table_name,auto_increment from information_schema.tables where table_schema = database() and table_name like 'xxx';
    可在当前database执行: drop database database(); 执行后再: select database(); 结果为NULL 
45. 设置主键自增初始值可以在建表时通过:AUTO_INCREMENT=xxxx制定,如: create table tb_name (column_name,column_type) AUTO_INCREMENT=100000
    也可以通过: alter table tb_name AUTO_INCREMENT=10000 来设定,不过这种方式要确保新的值要比表中已有的值大
46. mysql insert/update语句执行后,返回自增ID
	mysql shell登录服务器,执行insert/update后,可通过:select LAST_INSERT_ID(); 或者: select max(auto_increment_column) from tb_name;
	对于active connections,mysql服务器为每个连接单独维护last_insert_id值,各个连接之间互不影响
	当insert多条数据时,返回的插入的第一条数据的自增ID,而不是max(auto_increment_column)
47. 贴一个淘宝写的mysql源码分析--连接与认证过程的网址,以备后续参考学习:http://mysql.taobao.org/monthly/2018/08/07/ 
48. Mysql的一些约束规范
    数据表的字符集一般都是UTF8,如果需要存储emoji表情,需要使用UTF8mb4,mysql5.5.3以后支持
    单表数据量控制在1亿以下
    所有字段均定义为NOT NULL,除非要存NULL,两个缺点:使用Null时InnoDB引擎需要额外一个字节存储,浪费空间;默认Null值过多,影响优化器选择执行计划
    使用varchar存储变长字符串,varchar(N) N指的是字符个数,不是字节数
    使用decimal存储精确浮点数
    少用blob,test
    不在数据库存储图片和文件,一般存地址
    单个索引的字段数不超过5,单表索引数不超过5
    避免索引的隐式转换和冗余索引
    索引可以加速读,也会引入额外的写入和锁,降低写入能力
    避免大表的join,优化器对join优化策略比较简单
    避免在数据库中进行大量的数学运算
48. DDL(Data Definition Languages)数据定义语言,主要用来定义或改变表的结构,包括create,drop,alter,truncate
    DML(Data Manipulation Languages)数据操作语言,用于表中数据的添加,删除,更新和查询,包括insert,delete,update,select
    DCL(Data Control Languages)数据控制语言,用于授予或收回对数据的访问权限,包括grant(授权),revoke(取消授权)
    mysql执行DDL时需要锁表,锁表期间数据无法写入,online DDL主要采用Facebook OSC和5.6 OSC,5.6 OSC未解决DDL时从库的延时问题,而facebook OSC
    则采用触发器+change log的方式,腾讯互娱的online DDL是通过修改InnoDB存储格式来实现的(https://segmentfault.com/a/1190000004946420)
49. 查看mysql的数据存储位置
    mysql shell: show global variables like "%datadir%";
    shell: cat /etc/my.conf or cat /etc/init.d/mysqld 查看datadir选项
    mysql shell: show global variables like "%dir%"; 可查看更多和路径有关的配置值
50. 存储统计信息时,mysql报:data truncated for column 'xxx' at row 1
    通常出现这个报错的原因有:乱码,超过字段定义的长度,存在非法字符(类型不一致or编码问题)
    修改表字段类型or名: alter table t_name change old_column new_column new_data_type
    修改表字段类型: alter table t_name modif column column_name new_data_type
    修改数据库编码: alter database db_name character set = character_set_name
    修改表编码: alter table t_name convert to character set character_set_name
51. mysql的优化从哪些方面入手:
    a) 硬件太老:CPU,RAM,Disk,其它还包括网卡性能,机房的网络
	不同mysql版本对cpu的利用,5.1及之前的版本最多只能用4个core,5.5可以用24个,5.6可以用48+个
        在高并发环境,基本是靠提高内存缓存来减少对磁盘的访问,对于一般业务可按总数据的15%-20%(经验值)来规划缓存的热点数据
        核心库查询的平均响应时间不能超过30ms,如果超过则可能已达到承载极限,需要扩容,可对query的响应时间进行长期监控
        mysql的归档日志(binlog),重做日志(redolog),回滚日志(undolog),中继日志(relaylog)等文件为顺序读写,可以放到机械硬盘上
        数据文件(datafile),数据库表结构数据文件(ibdata file,.ibd)等为随机读写,可以放到SSD上
    b) 数据库设计不好,表结构,索引,单表大小,高级特性
       数据库设计不要过多使用触发器,函数,存储过程,mysql5.6之前的版本子查询的性能很差,如果需要在生产环境使用子查询,选5.6+版本
    c) 程序写的烂,
       应用程序尽量用连接池,特别是大型高并发应用程序,应用连接池减少连接的创建开销
       复杂语句,针对业务进行优化业务及查询语句,或者mysql主从架构,只从从库查询, 
       无效逻辑,全表扫描,可以在查询语句加where等限制条件,或者建立sql审查,经过DBA审核后,才能发布上线,对于大批量的更新操作,
       将任务拆分成小任务,分批更新
52. 一篇对比mongodb和mysql读写性能的文章:https://blog.csdn.net/clh604/article/details/19608869
    mongodb充分利用内存资源,读性能比mysql快了一个量级
    在不指定"_id"的情况下,mongodb的写性能是mysql的一倍,但在指定"_id"的情况下,mongodb的性能低于mysql指定primary key,且波动很大
    mongodb指定"_id"的情况下,每次插入数据,都要检查该"_id"是否可用,当数据量较大时(1亿条),判重开销将拖慢整个数据库的插入速度
    mongodb不指定"_id"时,系统自动计算唯一ID,"_id"为ObjectID类型,ObjectID为12byte,每个byte包含2位16进制字符,所以"_id"是一个24位的字符串
    ObjectID 12字节包含以下部分:0-3byte是时间戳(精确到秒一级),4-6byte是机器码,7-8byte是生成ObjectID的进程PID,9-11byte为计数器,
    一个mongod进程维护一个全局唯一的计数器,保证同一秒的ObjectID唯一
53. mongodb data/文件夹下的 mongod.lock文件,如果下次启动的时候还存在的话,需要删掉才能启动成功
54. mongodb的主从部署:
	./mongod --dbpath /data/master --port 1000 --master
	./mongod --dbpath /data/slave --port 1001 --slave --source localhost:1000
    主从部署,--master指定主节点,--slave 表明slave节点,--source 指定主节点的地址和端口
    mongodb的replica Set(复制集)部署:
	./mongod --dbpath /data/set1 --port 1000 --replSet set_name
	./mongod --dbpath /data/set2 --port 1001 --replSet set_name
	./mongod --dbpaht /data/set3 --port 1002 --replSet set_name
   Replica Set需要指定相同的复制集名称(--replSetC参数),mongod实例可在同一台设备,也可是多台设备
   单节点、主从、及复制集只需要启动mongod即可,mongod用来分片存储数据,mongos则用于shard集群中的路由处理
55. mysql对字符集的支持分为4个层次,分别为服务器,数据库,表和连接
    常用的几个字符集:
	character-set-server/default-character-set 服务器字符集,默认情况下使用的
	character-set-database 数据库字符集
        character-set-table 表字符集,若创建数据和表示未指定字符集,则默认采用的就是服务器字符集
	character-set-client 客户端字符集,当客户端向服务器发送请求时,请求以该字符集进行编码
	character-set-results 结果字符集,服务器向客户端返回结果时,以该字符集进行编码,客户端如果没有定义该字符集,
	则默认采用character-set-client字符集
	chracter-set-connection 数据库连接字符集
    处理中文时,可将character-set-server和character-set-client设置为GB2312,若需要处理多种语言,可将二者设置为utf8
56. mysql的默认字符集
    编译mysql时,指定了一个默认的字符集latin1
    安装mysql时,可以在配置文件(my.ini)中指定默认字符集,若未指定,则继承自编译时指定的字符集
    启动mysqld时,可以在命令行中指定一个默认字符集,没指定则继承自配置文件中的字符集,此时character-set-server被设置为此字符集
    创建一个新数据库时,若未指定,则character-set-database继承自character-set-server
    创建表时,若未指定,则表的默认字符集继承自character-set_database
    创建或修改column时,可通过character set xxx指定该column使用的字符集,若未指定,则使用表默认的字符集
57. mysql表的字符集和字符序
    字符集(character set)定义了字符及字符的编码
    字符序(collation)定义了字符的比较规则
    mysql乱码与客户端,数据库连接,数据库,查询结果的字符集有关,插入数据时,客户端,数据库连接,数据库的字符集需保持一致,因为在这三个地方要进行
    编码转换,同插入数据,查询数据时,返回结果,数据库连接,客户端字符集需保持一致
    mysql shell: show character set;  or use information_shcema; select * from character_sets; // 显示mysql支持的字符集
    mysql shell: show collation where charset = 'xxx'; or use information_schema; select * from collations 
	where character_set_name = 'xxx'; // 显示字符集支持的字符序
    mysql shell: show variables like 'character_set_server'; show variables like 'collation_server'; // 查看mysql server的字符集,字符序
    mysql server字符集&字符序可在启动mysqld时通过命令行指定,或修改配置文件,或在运行时,mysql shell: set character_set_server = 'utf8';
    或者在编译时,cmake . -DDEFAULT_CHARSET=utf8mb4 -DDEFAULT_COLLATION=utf8mb4_general_ci
    创建或修改数据库时,通过CHARACTER SET及COLLATE字段指定字符集和字符序,例:
	create/alter database db_name [default] character set xxx [default] collate yyy;
    查看数据库的字符集和字符序
    mysql shell: use db_name; select @@character_set_database, @@collation_database; 
    也可通过information_schema及建库语句查看,建库语句为:show create database db_name;
    数据库表的字符集与字符序与数据库类似,需要注意的是:如果建表时指定了charset_name与collation_name则使用指定的字符集
    若只指定了charset_name,则collation_name采用charset_name对应的默认字符序
    若只指定了collation_name,则charset_name采用collation_name关联的字符集
    若两者均为指定,则采用数据库的字符集与字符序
58. 记录下导出数据库、表的用法,命令行下:mysqldump -hxxxx -uyyy -pzzz db_name > xxx/yyy/db_name.sql 导出数据库
    mysqldump -hxxxx -uyyy -pzzz db_name t_name > xxx/yyy/t_name.sql
    导入数据,命令行下:mysql -hxxxx -uroot -pyyy db_name < xxx/yyy/db_name.sql
    mysql shell下source导入,use db_name; source local/path/to/db.sql;
    mysql shell下select导出,select * from t_name into outfile 'xxx/t_name.txt';
    select into oufile导出时,若mysql server的 --secure-file-priv选项有效,则导出操作失败
    select导出的另一个可能遇到的问题是报:Access denied for user 'root'@'%' (using password:YES),查了一下说是root用户没有FILE权限
    不过这个问题可以绕过,直接在本地shell: echo "select * from t_name" | mysql -hxxx -uroot -pyyy db_name > path/to/local/t_name.txt
59. echo+管道的方式导出表数据到本地文件,如果跨越多个库可以采用追加的方式"... >> path/to/local/t_name.txt"
    这种方式的缺点其一是需要多次手动执行,其二就是结果中会有重复的字段名(字段名表1数据字段名表2数据),格式上有些尴尬
60. 常见join操作,用来统计些数据 
    select count(distinct t1.field_name) from t_name1 t1 left/right/inner join t_name1/2 t2 on t1.field_name = t2.field_name where xxx
    left join 以左表数据为基准,同理right join,inner join则是取交集
    关于select xxx as yyy,除了可以给字段取别名,比如:select user_id as id, user_name as name from user; 之外,还可以给select出来的临时表命名
    比如:select id,name from (select user_id as id, user_name as name from user) as temp_user where xxx
61. 上回遇到一个没有FILE权限的提示,此处记录一下mysql权限相关的命令
    授权, mysql> grant privilege on *.* to *@* identified by "password" with grant option;
	privilege可以是所有权限(all privileges),也可以是select,insert,update,drop,alter等
        on x.y x指代数据库,y指代数据库中的表,*.*表示当前server的任何数据库及表
        to x@y x指代用户名,y指代限制的登录主机地址,可以是IP,IP段,域名或者任何地址('%'),'root'@'%'表示允许root用户在任意地址远程登录
        identified by "password" 指定登录密码
        with grant option 表示该用户可以将自己拥有的权限授权给其它用户
    刷新权限, mysql> flush privileges;
    查看当前用户权限, mysql> show grants;  
    查看某用户权限, mysql> show grants for 'xx'@'yyyy';
    回收权限, mysql> revoke privilege on *.* from 'xx'@'yyyy';
62. 导入.frm,.myd,.myi文件
    .frm 描述了表的结构, .myd 为表的数据记录, .myi 则是表的索引数据
    mysql> create database new_db;
    在mysql的data文件夹下打开new_db文件夹,将上述三个文件copy到里面
    mysql> show tables [desc t_name]; 查看表信息
63. 分布式存储系统的一个重要的三选二理论:C(Consistency,一致性)A(Availability,可用性)P(Partition Tolerant,分区容错性)
    一致性,分布式系统中所有数据备份,在同一时刻是否具有相同的值,也即是所有节点的数据时刻保持一致
    可用性,指对任意非失败节点的请求都能在有限时间得到响应
    分区容错,指的是允许节点之间丢失任意多的消息,网络分区状态一般包括节点网络不同,节点繁忙失去响应,单机房故障等
    In order to model partition tolerance, the network will be allowed to lose arbitrarily many messages sent from one node to another.
    上述对P的定义可知,在分布式系统中P是一个必选项,因为现实中,我们面对的是不可靠的网络和可能宕机的机器
