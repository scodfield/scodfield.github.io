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
    set @sql_str = concat('select count(*) from ', @table);
    prepare schema_table from @sql_str; execute schema_table; drop prepare schema_table;
    调用存储过程使用call proc_name
    删除存储过程使用drop [if exists] proc_name
