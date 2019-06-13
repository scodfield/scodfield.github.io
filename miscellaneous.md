1. RPC
   RPC框架通过提供一种透明的调用机制,让使用者不必显式区分本地调用和远程调用,从而让构建分布式计算or应用更加方便
   RPC调用分为同步和异步调用两种,同步调用指客户端发起调用后,等待执行结果的返回,异步调用指客户端不用等待执行返回,可以通过回调等方式获取结果
2. 正则表达式描述了一种字符串匹配的模式
   正则中的非打印字符:\f 换页符; \n 换行符; \r 回车符; \s 任意空白字符,包括空格,制表符,换页符等; \S 任意非空白符,等价于[^\f\n\r\t\v];
   \t 制表符; \v 垂直制表符
   正则中的特殊字符:$ 字符串的结尾位置; ^ 字符串的开始位置; () 子表达式的开始和结束位置,在中括号表达式中表示不在该集合中;
   * 匹配前面的子表达式零次或多次; + 匹配前面的子表达式一次或多次; ? 匹配前面的子表达式零次或一次; . 匹配除换行符\n以为的任意一个字符; 
   [] 中括号表达式; | 两项之间的一个
   正则中的限定符,表示表达式中的一个组件,必须要出现多少次才能满足匹配,除了*/+/?之外,还包括{n},{n,},{n,m}
   注: * + 限定符都是贪婪的,它们会尽可能多的匹配字符,只有在它们后面加上'?',就可以实现非贪婪或者最小匹配
   正则中的定位符,用来描述字符串或单词的边界, ^/$ 表示字符串的开始和结束; \b 描述单词的前/后边界; \B 非单词边界
3. 需要在另外一台电脑上checkout一些文件,输入地址&账号&密码后,TortoiseSVN提示:
   Unable to connect to a repository atu URL "xxxxxxx" Access to '/svn/xxx' forbidden
   在网上查了一下,SVN --> Settings --> Saved Data --> 几个相关的clear按钮点一下,之后重新checkout,bingo...
   参考:https://blog.csdn.net/wx_lanyu/article/details/84207303
4. 记录两个sublime的快捷键
    一个是常用的Ctrl+D,选中文本,当有多个的时候,或者想只选中其中的几个的时候,常用的是该快捷键,如果想对所有文本进行操作,一个一个选显然太low,
    可替代的就是先选中文本,再Alt+F3,即可选中全部文本;
    另一个是列操作,常用的是Ctrl+Alt+上下箭头,后边遇到个情况,需要多从mysql中导出的1w多条记录进行列操作,这时候一直按箭头显然也不现实,
    可替代的是先Ctrl+A选中所有的记录,再Ctrl+Shift+L ('L'大小写均可)
5. 从系统架构来看,目前的服务器可分为三大体系结构,分别是对称多处理器(SMP),非一致存储访问(NUMA)以及海量并行处理(MPP)
   SMP(Symmetric Multi Processor),所谓的对称是指多个cpu对称工作,没有主次和从属关系,因此SMP也被称为一直存储访问结构(UMA,Uniform Memory Access)
   对称多处理器系统内有很多紧耦合多处理器,所有的CPU共享总线,内存,I/O等系统资源,因此该系统最大的特点就是共享所有的资源
   cpu之间没有区别,平等的访问内存,外设,共用一个操作系统,操作系统管理一个队列,每个处理器依次处理队列中的进程,如果两个处理器同时访问一个资源
   (同一段内存地址),由硬/软件的锁机制处理资源争用问题,对SMP的扩展方式包括增加内存和cpu,更换更快的cpu,扩充I/O;
   NUMA(Non-uniform Memory Access)非一致存储访问架构,是一种为多处理器电脑设计的内存架构,内存访问时间取决于处理的内存位置,NUMA服务器的基本特征是
   具有多个cpu模块,每个cpu模块又有多个cpu组成,并且具有独立的本地内存和I/O槽口,由于节点之间可以通过互联模块(crossbar switch)进行连接和信息交互,
   因此每个cpu都可以访问整个系统内存,显然访问本地内存的速度远远高于访问其它节点内存的速度,这也就是非一致的由来,NUMA的主要缺陷也在于此,当添加多个
   cpu模块时系统性能无法线性增加;
   MMP(Massive Parallel Processing)与NUMA不同,MPP提供了另外一种扩展系统的方式,其基本结构是由多个SMP服务器(称之为节点),通过节点互联网络连接而成,
   每个SMP节点只访问自己的本地资源,是一种完全无共享的结构,在MPP系统中,每个SMP节点也可以运行自己的内存,总线,操作系统和数据库等,
   与NUMA不同的是,MPP不存在访问其它节点内存的问题,节点之间的信息交互是通过互联网络实现的,这个过程一般称之为数据重分配(Data Redistribution)
   参考:http://www.elecfans.com/baike/computer/fuwuqi/20171023568144.html
6. Protocol buff编码格式及sxxx的zigzag编码:
   https://www.cnblogs.com/cobbliu/archive/2013/03/02/2940074.html
   https://izualzhy.cn/protobuf-encode-varint-and-zigzag#4-zigzag%E7%BC%96%E7%A0%81
   protocol buff的数据类型:
   VARIINT 可变长度整型,该类型数据使用varint编码对所传入的数据进行压缩存储,int32/64,uint32/64,sint32/64,bool,enum属该类型
   FIXED32/64 固定长度整型,不会对传入的数据进行varint压缩,只存储原始数据,fixed32,sfixed32,float属FIXED32,fixed64,sfixed64,double属FIXED64
   LENGTH_DELIMITED 长度界定型数据,主要针对string,bytes,embedded messages,packed repeated field,简言之就是针对string类型,repeated和嵌套
   类型,对这些类型数据进行编码时需要保存它们的长度信息
   START_GROUP 组的开始标志,组也可以是repeated或嵌套类型
   END_GROUP 组的结束标志,其余同上
