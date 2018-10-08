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
32. mongo连接远程数据库
	mongo ip
	mongo ip:port
	mongo ip:port/db_name
	mongo ip:port/db_name -u xxx -p yyy
    示例: mongo 192.168.1.100:27027/thd_game
