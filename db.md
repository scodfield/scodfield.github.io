1. a> Redis不存在表的概念，一个Redis实例中，直接存储5大数据类型：字符串(string)、列表(list)、哈希(hash)、集合(set)、有序集合(zset)
   b> Redis expire命令用于设置key的过期时间,以秒为单位,过期后key将不再可用(被自动删除),示例:expire key_name time_in_seconds,另一个设置
      过期的命令是expireat,以Unix时间戳(unix timestamp)的格式设置key的过期时间,示例:expireat key_name time_in_unix_timestamp
   c> string类型除了最常用的get/set命令,还有一个incr命令,将key中存储的数字值加一,与之对应还有一个incrby xxx命令,key中存储的数字值增加指定的
      增量值,执行两个命令时,如果key不存在,那么key的值会先被初始化为0,再执行incr/incrby命令,如果key对应的值包含错误类型,或者字符串类型的值不能
      表示为数字,会返回一个错误,但是Redis没有数字类型,所以key存储的字符串会被解释为十进制的64位有符号整数来执行incr/incrby操作
      除了上述操作,string类型的getset,decr/decrby命令可加减key对应的数值
   d> 组合使用string类型的incr&expire命令可以实现在规定时间内的计数器,比如我们游戏中日常副本每天只能挑战三次,那么每次玩家请求挑战的时候,执行:
      newcount = eredis.incr(role_1001_daily_chapter_xx),返回最新的挑战次数newcount,如果newcount > 3,则当天挑战次数已用完,不能再请求挑战,
      至于role_1001_daily_chapter_xx这个key的过期时间,可通过eredis.expire(role_1001_daily_chapter_xx,seconds_to_next_zero)来设置为零点
      过期,这个时间可以通过明日零点的unixstamp减去当前时间得到,也可以通过eredis.expireat()命令直接设置为明日的零点时间,这样服务端程序就无需再
      每日零点刷新,即可实现记录副本的每日挑战次数
   e> Redis set集合也可以用来做计数器,比如我想统计网站的UV(unique visitor,独立访客),那么对于每个请求,parse出ip之后,通过:sadd s_name member
      命令将访客ip添加到uv_counter这个集合中,最后再通过:scard s_name取出uv_counter这个set中成员的数量,得到UV数据,但是set类型是通过hash表实现
      的,所以当统计的量比较大的时候,就非常浪费空间,而Redis提供的Hyperloglog这种数据结构就是用来解决这类统计问题的
   f> Redis Hyperloglog 是一种巧妙的近似估计海量数据基数的算法,所谓基数就是数据集中不重复元素的个数,比如数据集:{1,2,5,7,5,7,},那么基数集就是
      {1,2,5,7},基数就是4,基数估计就是在误差可接受的范围内,快速计算基数,Hyperloglog内部维护了16384个桶(bucket),用来记录每个桶内的元素数量,
      当一个元素到来时,它会散列到其中的一个桶,以一定的概率影响这个桶的计数值(每个桶有一个count,记录散列到该桶的数值的个数),将所有的桶的计数值
      进行调合均值累加,结果会非常接近真实的计数值,Hyperloglog主要提供了一下三个命令:pfadd key element 添加元素到指定的key中; pfcount key 获取
      给定key的基数估计值; pfmerge destkey sourcekye [sourcekey...] 将多个Hyperloglog key合并为到一个key中
      Hyperloglog的每个key只占用12k的内存,16384个桶(2^14),每个桶占用6bits,最大可表示的值为63,所以占用的内存:(16384 * 6)/(8 * 1024) = 12kb
2. Redis应用:
   a> 位图bitmap,位图并不是一种特殊的数据结构,本质上是二进制字符串,也可以看作是byte数组(go的底层文件/socket存储就是字节数组),redis string类型
      中有几个操作bit位的命令,get/set可以直接获取和设置整个位图的内容,也有单独的bit位操作命令getbit/setbit等,这些命令将byte数组看成bit数组来
      处理,常用命令如下:setbit key offset value 对key所存储的字符串,设定or清除指定偏移bit上的值,setbit online 1001 1; getbit key offset
      对key所存储的字符串,获取指定偏移bit上的值,getbit online 1001; bitcount key [start] [end] 计算key存储的字符串中,被设置为1的bit位的
      数量,set online no_one_online, bitcount online; bitop operation destkey key [otherkeys...] 对一个或多个bitmap进行位操作,并将结果
      保存到destkey上,位操作operation,可以是AND,OR,NOT,XOR这四种操作中的任意一种
      位图基于bit位,非常节省空间,设置的时候时间复杂度为O(1),读取的时候时间复杂度为O(N),二进制计算速度非常快,基于上述优点,位图适用于各类的统计,
      比如游戏中玩家是否在线,统计当前在线玩家人数,记录并统计玩家当月的签到情况等等
   b> 布隆过滤器bloom filter,用于判断一个元素是否在某个集合,布隆过滤器实际上是一个很长的二进制bit数组和一系列的hash函数,有点是空间效率和查询
      时间都比较高,缺点是有一定的误判率及删除困难,bit数组初始全部为0,当一个元素被加入到集合当中时,这个元素被k个hash函数计算出k个值,这些值是
      在bit数组中的offset,然后将bit数组中k个映射值对应的bit位置为1,至于判断一个元素是否存在,同样经过上述k个hash函数的计算,如果计算出的k个位置
      上全为1,则判断这个元素很大可能已存在,如果有一个位置为0,则该元素肯定不存在,由此可见,布隆过滤器有一定的误判率,且由于误判率的存在,导致难以
      删除指定元素,因为无法确切保证该元素是否存在,布隆过滤器的特点使得它适用于大数据规模下不需要精确过滤的场景,比如检查垃圾邮件地址,爬虫url地址
      去重,解决缓存穿透等,关于缓存穿透问题,比如恶意构造出一系列还未注册的玩家id,请求玩家的数据,此时缓存系统中并未有该玩家的数据,所以会导致对数据
      库的大量读写请求,那么就可以考虑构建一个所有已注册玩家ID的布隆过滤器,一旦发现请求还未注册玩家的信息,直接返回相关提示,而不用执行先访问缓存
      再访问数据库等一系列操作,大大减少缓存穿透对数据库的访问压力
      redis4.0提供了布隆过滤器插件,下载编译安装Rebloom插件,redis添加参数:rebloom_module="/path/to/rebloom.so" 启动参数: 
      --loadmodule $rebloom_module, 主要的命令:bf.add 添加元素; bf.exists 查询元素是否存在; bf.madd 一次添加多个元素; bf.mexists 一次
      查询多个元素是否存在,redis中有两个值决定布隆过滤器的准确性,error_rate 允许的错误率,这个值越低,所需的bit数组越大,占用空间也越大; 
      initial_size 布隆过滤器可以存储的元素个数,当实际存储的元素个数超过这个值后,过滤器的准确率会下降,可以使用bf.reserve命令来设置这两个值,
      bf.reserve bf_key 0.001 1000, 一般都要在调用bf.add命令前,显式调用bf.reserve来创建一个布隆过滤器,如果bf_key已存在,再调用bf.reserve
      时会报错,如果不调用bf.reserve,默认的的error_rate是0.01,initial_size是100
      注:由以上描述可知,布隆过滤器实际上比没有存储具体的元素,其次它不支持计数
      延伸阅读,布谷鸟过滤器:https://juejin.im/post/5cfb9c74e51d455d6d5357db
   c> redis实现分布式锁,当多个进程不再一个系统当中时,就需要分布式锁控制多个进程对资源的访问,redis实现分布式锁主要基于以下几个命令:sennx,
      expire,del, setnx key value 只在key不存在时,才会将值设为value,如果key已存在,会返回0, expire key seconds 对key设置过期时间(秒)
      del key 删除key, 该类型分布式锁的原理是将某个key视为互斥量(go.sync.Mutex),如果能成功执行setnx key value则获得锁,通过expire给该互斥
      key设置过期时间,防止逻辑处理过程中出现异常导致无法释放锁,del key命令删除key,类似于最后释放锁,在用redis实现分布式锁的过程中,有几个点
      需要注意:首先是要保证只有在成功拿到锁之后,再执行key的过期时间设置; 其次是每个进程在setnx key是value值最好是随机值,在del key前,先get
      一下,只在当前key的值和自己保存的随机值一致时,才执行del操作,这样做的目的是防止逻辑处理时间过长,导致key已失效,且下一个进程已获得该key,
      那么若直接执行del操作,则会误删下一个进程的锁请求
3. Redis保证单个命令的原子性，而Redis的事务并没有作任何原子性的保证，且事务中各个指令互不影响，既不会回滚，也不会终止后续指令的执行，感觉Redis
   事务更像批处理
4. Redis的save/bgsave命令备份当前数据库数据到dump.rdb文件，恢复操作是将该文件copy到Redis安装目录，直接启动服务(如果需要多个Redis服务怎么办。。。)
5. Redis分区是将数据放到多个实例中(好奇Redis的集群和分区的区别),Redis集群可以在运行时进行动态的伸缩性调整，灵活性比较高
6. Redis Hash是一个string类型的field和value的映射表,适于存储对象,同时Redis set是一个string类型的无序集合,且set是通过哈希表实现的,那么也少不了
   映射,则Redis的哈希用的是哪个哈希算法,有无用到一致性哈希,一致性哈希又用到了哪些地方?
   redis底层数据结构:
   typedef strutc redisobj {
       unsigned type;
       unsigned encoding;
       void * ptr;
   }
   其中type:string,list,hash,set,zset; encoding:REDIS_ENCODING_INT/EMBSTR/RAW/ZIPLIST/LINKEDLIST/HT/INTSET/SKIPLIST
7. 非关系型数据库方便处理相互之间没有耦合性的非结构化数据
8. a> MongoDB的读写性能(只能达到1kqps??)
   b> MongoDB进程启动后,除了和普通进程一样加载binary,依赖的各种library到内存之外,作为一个DBMS,还需要管理客户端连接,处理请求,维护数据库
      元数据和存储引擎等很多工作,这些工作都设计内存的分配和释放,默认情况下MongoDB使用Google的tcmalloc作为内存分配器,内存占用的大头主要
      包括存储引擎,客户端连接及请求的处理; 
      MongoDB3.2.0以后,默认使用的存储引擎是WiredTiger,可通过cacheSizeGB字段来配置WiredTiger引擎使用内存的上限,默认配置为系统可用内存
      的60%,为了控制内存的使用,WiredTiger在内存使用接近一定阈值的时候开始做淘汰eviction,目前有4个参数来支持WiredTiger引擎的eviction策略
      eviction_target 默认为80(百分比),当cache used超过eviction_target(阈值计算为:0.8 * cacheSizeGB),后台evic线程开始淘汰clean page;
      eviction_trigger 默认为95(同上),当cache used超过eviction_trigger(阈值计算同上),用户线程也开始淘汰clean page;
      eviction_dirty_target 默认为5(同上),当cache dirty超过eviction_dirty_target,后台evic线程开始淘汰dirty page;
      eviction_dirty_trigger 默认为20(同上),当cache dirty超过eviction_dirty_trigger,用户线程也开始淘汰dirty page;
      需要注意的是当used >= 95%或者dirty >= 20%,并一直持续,说明内存淘汰压力比较大,用户的请求线程会被阻塞,参与到page淘汰,这时候客户端请
      求时延增加,可以考虑增加内存,或更好更快的磁盘以提高IO能力;
      TCP连接及请求处理,tcp协议栈除了为连接维护socket元数据之外,每个连接会有一个read buffer和write buffer,除了协议栈,mongod进程针对每
      个连接起一个线程,专门负责处理这个连接上的请求,该线程的线程栈通常只有十几k,但mongod为线程栈配置的上限是1MB,可通过proc查看线程的具体
      开销,上述协议栈的read/write buff可通过ss -m命令查看每个连接的buff信息(具体命令参考common_commands.md)
9. 分区和索引哪个优先执行？
	查询条件子句中含有分区条件，则先过滤分区，再在对应分区执行索引操作，否则在所有分区执行索引操作。因此，在实践中尽量避免分区列和索引列的不匹配，或者条件子句中包含过滤分区的条件
10. 如何实现读写分离？
	目前mysql读写分离通常采用的做法是主节点处理写操作，并异步复制更新数据到从节点，而从节点处理select等查询操作，总结就是主从复制+读写分离。读写分离的好处：
	1) 增加冗余备份,提高可用性
	2) 一主多从或多主多从,减轻单一节点压力,实现负载均衡,提升系统性能
	主从结构存在的问题是主从之间异步复制导致的数据不一致,既主从间的数据延迟问题
11. a> mysql5.7.8引入json类型,且使用的是内部的二进制格式而非字符串,支持对文档的快速查询
    b> mysql8.0支持窗口函数(window functions),按官网定义窗口函数对查询结果做聚合类操作(A window function performs an aggregate-like 
       operation on a set of query rows),聚合操作一般是对组数据(group by)操作,输出单个值,比如求每个student的平均分:
       select id,name,avg(score) as avg_score from scores group by stu_id order by avg_score,窗口函数后一般跟over子句(over clause),
       over子句格式为OVER(window_spec),聚合函数窗口函数后都可跟over子句,比如: avg(score) as avg_score OVER(partition by stu_id) from scores
       order by avg_score, window_spec包括4个部分,常用的是下面两个:partition_clause,partition by xxx 指定如何对query rows分组; 
       order_clause order by xxx asc|desc 指定如何排序及排序方向;
       窗口函数中和排序相关的有以下四种:
       row_number() 连续排序,相同的值序号不一样,select row_number() over(order by s.score) as rank,s.stu_id,s.name,s.score from scores s
       rank() 跳跃排序,相同的值归为一组且序号一样,select rank() over(同上) 同上....
       dense_rank() 连续排序,相同的值归为一组且序号一样,select dense_rank() over(同上) 同上....
       ntile(group_num) 将所有记录分为group_num个组,每组中各个元素的序号都一样,select ntile(4) over(同上) 同上....,将学生按成绩分成分四档
    c> mysql是一个客户端/服务器架构的软件,mysql服务器可同时与多个客户端连接,每个连接称之为一个会话(session),不同的会话可能同时发送一系列sql
       请求,服务器把每个会话发送的sql语句放到不同的事务里进行处理,这样就可能出现不同的事务同时访问相同的记录,而理论上事务应该满足ACID,其中I就是
       Isolation 隔离性,也就是说不同的事务应该彼此互不干扰,那么最简单的同时也是隔离级别最高的方法就是串行化(serializable),如果一个事务在对某个
       记录访问时,其它事务排队,当该事务提交后,其它事务才可以继续访问该记录,这种方法的缺点是牺牲了系统的并发处理能力
       在引入事务隔离之前,先普及一下并发事务可能出现的几个问题:
       首先是脏读,事务A读取了事务B更新但未提交的数据(数据处于更改的中间状态),事务B回滚,那么此时事务A读取到的数据就是脏数据;
       其次是不可重复读,事务A多次读取某个数据,事务B在A读取的过程中,对数据更新并提交,导致A在多次读取时,数据的结果不一致;
       最后是幻读,事务A多次做同一个涉及某个范围的查询,事务B对A范围内的数据做了Insert/Delete操作,导致A的查询结果数量多了或少了,就像出现幻觉
       mysql有四种事务隔离级别,除了上述所说的串行化之外,还有以下三种:
       read uncommitted 读未提交,如果一个事务读能取到另一个未提交事务修改过的数据,这种隔离级别就称为读未提交;
       read committed 读已提交,如果一个事务只能读取到另一个已提交事务修改过的数据,并且其它事务每对该数据进行一次修改并提交后,该事务都能查询
       到最新的数据,那么这种隔离级别就称为读已提交
       repeatable read 可重复读,如果一个事务只能读取到另一个已提交事务修改过的数据,但是第一次读取某条记录后,即使再有其它事务对该记录修改并
       提交,该事务之后再读取时,读到的仍是第一次读到的值,而不是每次都读取到最新的数据,这种隔离级别称为可重复读,InnoDB引擎默认的采用该隔离级别,
       可重复读使用了MVCC(多版本并发控制)机制,每当有事务对记录进行insert/delete/update操作时,都会更改记录隐藏的版本号,在可重复读隔离级别下
       select操作并不会记录更新版本号,而是去读事务第一次select时记录的版本号(类似于快照),需要注意的是,可重复读隔离级别下,select读取的是旧的
       版本号,insert/update/delete读取的是最新的版本号,比如:A事务select 库存 = 10,这个时候B事务生成一个订单,减去1,并提交,再轮到A事务生成
       订单,同时库存-1,执行完B事务及A事务的扣减操作后,再次查询select 库存 = 8,也就是说在可重复读隔离级别下,同一个事务中如果执行了inset/upd
       ate/delete操作,此时会更新到最新的版本号,select读的也是最新的版本号
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
16. mysql存储过程是一组完成特定功能的sql语句块,经过预编译保存在进程字典中,常用于批量处理一些重复性高的操作;
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
22. Redis为什么这么快? 主要有几个方面:
    a> 内存操作,相比读磁盘速度快的多
    b> 单线程,避免竞争锁
    c> IO多路复用,同时监控多个客户端连接
    d> 零拷贝,直接从操作系统的页缓存(page cache)拷贝到网卡,减少拷贝次数
23. 为什么说Redis是单线程?
    https://mp.weixin.qq.com/s?__biz=MjM5ODI5Njc2MA==&mid=2655824091&idx=1&sn=440fc9ed2587803bb65c8ccc276f9b48&chksm=bd74e50c8a036c1a3a1148295b2d4a5128e49db8847c9d3269ea93b5426db3fcd80daaf73dad&scene=21#wechat_redirect
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
32. a> mongod启动: mongod --dbpath xxx --logpath yyy --fork --auth --bind_ip zzz
    b> mogodb支持多种类型的索引,包括单字段索引,复合索引,多key索引和文本索引,索引可以提高查询的效率,但mongodb的索引不能被以下
       查询利用:正则表达式及非操作,如$nin,$not等; 算术运算符,如$mod等; $where子句(这几个限制倒和mysql差不多,除了where子句),
       按手册,MongoDB索引还有以下限制:额外开销,每个索引占据一定的存储空间,在插入,更新和删除操作时也需要对索引进行操作;内存(RAM)
       使用,索引是存储在内存中的,如果索引大小超过内存限制,MongoDB会删除一些索引,这将导致性能下降;一些范围限制,集合中索引不能
       超过64个,索引名的长度不能超过125字节,一个复合索引最多可以有31个字段
    c> single field index 单字段索引,db.t_name.ensureIndex({key:1/-1}),对t_name表的key字段创建索引,1代表升序,-1代表降序
       '_ id'是创建表的时候自动创建的索引,此索引不能被删除,当系统已有大量数据时,创建索引是个非常耗时的操作,因此可以将该操作
       放到后台进行,只需指定ensureIndex()的可选参数background,如:db.t_name.ensureIndex({key1:1}, {background:true}),常用的
       可选参数除了这个用于后台创建索引的background,还有unique,用法同background
       更多ensureIndex的可选参数及用法可参考手册:https://www.mongodb.org.cn/tutorial/18.html
    d> compound index 复合索引,ensureIndex()方法中可以同时使用多个字段创建复合索引,如:db.blog.ensureIndex({"title":1,"desc":-1},
       {unique:true}),上述语句创建一个复合索引,同时还是一个唯一索引,针对多个字段的复合索引,先按第一个字段排序,第一个字段相同的文档
       按第二个字段排序,依次类推
    c> multikey index 多key索引,当建立索引的字段为数组时,创建出的索引称为多key索引,比如有以下结构:
       {"title" : "some tips about mongodb", "tags": ["grammer", "index", "sharding"]}, 我们针对"tags"字段创建数组索引:
       db.blog.ensureIndex({"tags":1}), 当使用以下语句检索时,加上explain()命令查看是否使用了索引:
       db.blog.find({"tags":"index"}).explain()
    d> 还可以针对子文档字段简历索引,比如有以下结构:
       {"author" : {"id":1001,"name":"thd","rank":100}, "title" : "some tips about mongo", "tags" : ["grammer","index","sharding"]}
       我们可以针对author这个子文档的三个字段简历索引:db.blog.ensureIndex({"author.id":1,"author.name":1,"author.rank":1}), 创建
       之后,我们可以针对子文档字段来检索数据:db.blog.find({"author.rank":{$lt:1000}}), 查询排名前一千的玩家blog,针对该类型索引的查询
       表达式必须遵循创建索引时字段的顺序,如果顺序不同,则不会使用索引,比如:db.blog.find({"author.name":"xte","author.id":1001})
    e> 更多MongoDB索引,可参考手册:https://docs.mongodb.com/manual/indexes/
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
    IT又一次直接重启内网,按照上述步骤之后,提示:child process failed,exist with error number 100, 参考下面的资料,以repair的方式重启,
    还是提示100,在系统上找了找mongod.log启动日志,在最新的启动日志发现如下提示:
    exception in initAndListen: 72 Requested option conflicts with current storage engine option for directoryPerDB;
    you requested false but the current server storage is already set to true and cannot be changed, terminating
    ok, 在/etc/mongod.conf文件的storage项下面加入:directoryPerDB true, 保存退出,再次以repair方式启动:mongod -f /etc/mongod.conf 
    --repair, 参照下面的资料,再执行:mongod -f /etc/mongod.conf, 一切正常,按照上面正常退出mongod的方法,杀掉mongod进程之后,
    再以mongod --dbpath xxx --logpath yyy --fork --auth --directoryperdb 命令启动,启动游戏服务器,连接成功
    注:如果以上解决100的方法还不行,可以考虑终极方法,删除整个mongo的数据和日志,重启
    参考:https://blog.csdn.net/sinat_30397435/article/details/50774175
35. mongo首次运行时并未开启验证,登录后需要首先创建一个超级管理员用户,用于管理其他用户
    由手册可知,mongo3.0以后,在admin中创建(As of MongoDB 3.0, with the localhost exception, you can only create users 
    on the admin database),创建命令为db.createUser({user:"xxx", pwd:"yyy", roles:[{},...]})
    createUser的roles选项,指定了被创建用户的"角色",包括内置角色和自定义角色,mongo采用基于角色授权的方式管理对数据库的相关操作
    为了便于管理,首个超级管理员账号,可以选择内置角色中属于superuser roles的root角色,root拥有所有权限(provides full privileges
    on all resources) localhost exception which allows you to create a user administrator in the admin database
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
       加索引,在where,order by等条件所涉及的列上加索引,减少全表扫描
       避免在where子句里使用!=,<>操作符或函数,操作符或函数可能导致引擎放弃使用索引而进行权标扫描
       设计表时,字段尽量使用数字类型,减少字符型字段,引擎在查询和连接字符时,会逐个比较字符串中的每一个字符,如确实需要使用字符,可考虑find_in_set
       避免频繁创建和删除临时表,减少系统表资源的消耗
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
    为了安全起见,mysql的root用户默认是不开远程访问权限的,在mysql系统库的user表,select user,host from user; 可以查看所有用户的的用户名及访问地址
    此时若要修改某个用户的访问地址,update user set host ='%'/'ip_addr' where user = 'xxx'; flush privileges; mysql远程访问除了设置用户的访问
    地址之外,还需要mysql位于公网,及设置mysql所在服务器防火墙
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
    CAP&Mysql cluster:http://messagepassing.blogspot.com/2012/03/cap-theorem-and-mysql-cluster.html
    关于CAP的一些争论和质疑:https://blog.csdn.net/chen77716/article/details/30635543
64. mongodb数据库备份:mongodump -h db_host[:port] -uxxx -pyyy -d db_name -o /path/to/dump , 该命令将对应的数据库实例备份到指定的本地路径
    数据恢复:mongorestore -h db_host[:port] -d db_name [--dir]/path/to/restore
    mongodump备份的数据格式为:collection_name1.bson,collection_name1.metadata.json,...
    mongo在当前路径恢复数据,mongorestore -h xxx db_name:port [-uyyy -pzzz] ./mongo_dump_dir/ 如果远程服务器开启了验证,需要加上-u/p参数 
65. mongo复制集和分片 
    mongo的复制指的是数据在多个mongod实例之间的同步,它提供数据的冗余备份,提高数据的可用性和安全性,具体部署可参考上面的记录
    mongo的分片指的是将大量的数据集合分成一块块的小集合,并将这些小集合分散到多个mongod实例上
    复制集和分片都有集群的概率,也都涉及到多个mongod实例,二者的区别在于:复制集提供了数据复制的技术,而分片提供的是数据分割的技术
    分片在允许MongoDB存储海量数据的同时,提供可接受的读写吞吐量,其主要有以下三个组件构成:
    Shard 实际存储数据,一个Shard在生产环境中,一般可由多个分布于不同机器节点组成的一个Replica set承担,防止单点故障
    Config Server 存储整个ClusterMetadata,包括chunk信息
    Query Router 前端路由,客户端由此接入,使得集群对前端应用透明
    分片实例:
    首先启动shard server, mongod --port xxx --dbpath yyy --logpath zzz --logappend --fork, 改变port和dbpath启动多个mongod实例
    其次启动config server, 启动命令同上,指定不同的port,dbpath,logpath即可
    然后启动route process, mongos --port xxx --configdb config_server_ip:port --fork --logpath yyy --chunkSize zzz
    最后配置sharding, mongo shell登录到mongos,添加shard节点
	shell> mongo router_ip:port admin, 登录mongos,并转入admin库(mongo router_ip:port; use admin)
	mongos> db.runCommand({addshard:"shard_server1:port1"})
	......
	mongos> db.runCommand({addshard:"shard_servern:portn"})
	mongos> db.runCommand({enablesharding:"shard_db_name"}) # 设置分片存储的数据库
	mongos> db.runCommand({shardcollection:"shard_collection_name",key:{doc_index_field:key_index}}) # 设置分片存储的数据库
    mongos的启动参数:
    --chunkSize 指定chunk大小,单位是MB,默认200MB
    --maxConns 最大并发连接数,默认1000000
    MongoDB database commands:
    db.runCommand({addshard:"replica_set/hostname:port"}); 向分片集群添加shard_server
    db.runCommand({enablesharding:"db_name"}); 激活指定数据库分片
    db.runCommand({shardcollection:"db_name.collection_name",key:<shardkey>}); 允许指定collection分片
    shardkey决定了documents在shards上的分布,shardkey的结构是:{document_indexed_field:key_index},第一个字段是文档的索引字段或组合索引
    key_index有三个取值:1 indexed field的前向遍历分布; -1 indexed field的后向遍历分布; hashed 指定的hash key分布
66. mysql之where,having,两者功能类似,均可用于筛选数据,其后的表达式也都一样
    区别在于:where 针对表中的列发挥作用,查询数据; having 针对查询结果中的列起作用,筛选数据
    select goods_id,goods_name,cur_price - original_price as price from goods having price > 200;
67. mysql创建表的三种方式
    典型create命令:create [temporary] table [if not exists] tb_name create_defination
    create table like 参照已有表的定义:create [temporary] table [if not exists] tb_name like old_tb_name
    根据select的结果集创建表:create [temporary] table [if not exists] tb_name as query_expression; (注:该方式创建的新表,没有主键和索引)
    例:create table goods_as as select id,name,price from goods; 
68. mysql count(*)与count(column_name)及sum()的一些小区别
    count(*) 计算所有的行数; count(column_name) 计算字段非NULL的行数
    sum(column_name) 计算列名对应的值(非NULL值)的和, sum(表达式) 如果记录满足表达式,则加1
69. Navicat导入文件时,若第一行不是字段行,则不能有NULL字段
    mysql合并查询的多行数据:group_concat, select field_1,group_concat(field_2) [as field2],group_concat(field_3) [as field3],field_4 
    from tb_name [group by field_6]; 上述语句会将所有结果合并到一行,若需要按某个字段合并,可加上group by field_name;
    mysql备份数据的一种方式:
    select */fields_list from tb_name [where/order/group...] into outfile "/path/to/backup" fields terminated by "xx" 
    lines terminated by "yy"; 
    该备份方式数据存储在mysql server上,若为"./xxx.txt",则在/data/mysql/路径下
70. 梳理了一下初始化流程,感觉还可以再精简一下mongo读取,每个玩家可以再减少两次读取操作,那么在玩家大规模登录的时候(比如10w),节省的读取次数还是很可观的
    首先一个要确认的就是单个文档的大小,毕竟要把其它功能模块合并到已有的模块中去,而一般在连接mongo时,若safe=False,则客户端在向数据库发送插入,删除
    等操作时,是不需要等待数据操作结果的(成功or失败),若单个文档超过大小上限,那么客户端程序并不会报错,解决方法是开启安全验证,连接时safe=True
    由手册可知,MongoDB目前单个文档的大小上限是16M(https://docs.mongodb.com/manual/core/document/ Document Limitations)
71. mysql int类型
    数据类型   字节  有/无符号最小-最大值
    tinyint   1byte -128-127/0-255 
    smallint  2byte -32768-32767/0-65535
    mediumint 3byte -8388608-8388607/0-16777215
    int       4byte -2147483648-2147483647/0-4294967295
    bigint    8byte -9223372036854775808-9223372036854775807/0-18446744073709551615
    int(M) int类型后面的这个数值M,表示的是最大显示宽度,最大有效显示宽度是255
    M的值和int(M)所占用的存储空间没有任何关系(和int/tinyint有关系),M=10表示告诉数据库该字段存储的数据的宽度为10,如果存储的数据不是10位数,只有
    该值仍在int类型的有效范围内,mysql也能正常存储
72. mongo嵌套查询
    db.collection_name.find({'_id':unique_id},{'_id同层属性1.嵌套属性1':1, '_id同层属性1.嵌套属性2':0}
    1 - 表示查询结果中保留该字段;  0 - 表示查询结果中去除该字段
73. mysql有四种类型日志
    error log 错误日志,记录mysqld的一些错误
    general query log 一般查询日志,记录mysqld正在做的事情,比如客户端的连接和断开,客户端每条sql statment,详细记录了客户端传给服务器的每条查询
    该日志非常影响性能
    slow query log 慢查询日志,记录一些查询比较慢的sql statement,常用于开发者调优
    Binary log 二进制日志,记录一些事件,包括数据库的改动,建表,数据改动等,也包括一些潜在的改动,比如:delete from t where id = xxx;记录所有改动
    潜在改动的sql statement,以二进制的形式保存在磁盘中,bin log 的作用:可用于查看数据库的变更历史(任何时间点的所有sql操作);数据库的增量备份和恢复;
    (增量备份和基于时间点的恢复);mysql复制(主主,主从复制)
    mysql默认关闭bin log,可通过修改mysql配置文件打开,linux下配置文件为my.cnf,一般在/etc/my.cnf,windows下是my.ini 或者 my-default.ini
    开启Binlog,需要修改 log_bin[=base_name] base_name是生成的Binlog文件的前缀,没有的话,用的是pid-file选项的值
    max_binlog_size 每个Binlog文件的大小,最小是4m,最大与默认是1G
    log_bin_index[=file_name] mysqld为了追踪已使用过的Binlog文件,会创建一个Binlog索引文件,默认是Binlog的basename.index
    expire_logs_day  Binlog过期清理时间
    binlog_cache_size Binlog缓存大小
    max_binlog_cache_size Binlog最大缓存大小
    binlog_format=format Binlog文件的模式:
    row level 不记录每条sql语句的上下文信息,记录的是每一行数据被修改的情况,然后在slave端对相同的数据进行修改,不会出现某些情况下的存储过程,
    function,trigger的调用和触发无法被复制的问题,缺点是会产生大量的日志,尤其是alter table的时候
    statement level 每一条会修改数据的sql语句会记录到Binlog文件中,优点是不需要记录每一条语句和每一行数据的变化,减少Binlog日志量,节约磁盘IO
    slave在复制的时候sql进程会解析成和原来在master端执行过的相同的sql再次执行,确定是容易出现主从不一致,如执行sleep(),last_insert_id()等函数
    mixed level 混合模式,结合row level和statement level的优点,一般的操作使用statement模式保存,对于statement模式无法复制的操作,使用row模式
    Binlog相关设置更改后,重启mysqld,在mysql shell执行:show variables like '%log_bin%'; show variables like '%binlog%';
    mysqlbinlog命令可以将Binlog日志转换成mysql语句,默认情况下Binlog日志是二进制文件
    mysqlbinlog的参数: -d 指定数据库的Binlog; -h -u -p 指定hostname,用户和密码; -o 指定跳过前N条记录; -r 输出结果到指定文件
    例 解析login数据的Binlog,并写入login.sql文件:mysqlbinlog -d login mysql-bin.000001 -r login.sql
    详细参数,可mysqlbinlog --help
74. MongoDB注册windows服务,以管理员身份运行cmd,运行以下命令注册windows服务:
    mongod --dbpath "d:\mongodb\data" --logpath "d:\mongodb\logs\error.log" --install  --auth --directoryperdb --serviceName MongoDB
    --serviceDisplayName MongoDB 
    服务创建完成后,默认是停止状态,需要手动打开服务,接着运行:net start MongoDB (如果还不生效,打开任务管理器,找到服务,手动打开)
74. mysql清除表数据,truncate vs delete
    truncate 整体删除,速度比delete快,不写服务器log,truncate不激活触发器(trigger),会重置identity(标识列,自增字段)
    delete 逐条删除,写服务器log,不会重置identity
    如果只删除部分数据,只能使用delete配合where条件语句
75. mysql的几个模糊查询:like + 通配符; field in (xxx); where find_in_set(filed,strlist); where field REGEXP 'reg_exp'
76. mysql与mongo或者关系型数据库与非关系型数据库的缺点:
    a> 事务，mongo4.0引入了事务,但还是有局限性
    b> 关系完整性(外键),如果数据之间有关系,mongo需要借助应用程序
    c> 数据结构,mongo在结构方面提供了极大的灵活性,但是为了保证数据的可用性，应用程序需要做大量的维护
