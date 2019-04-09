
1. centos6.9安装erlang 20.3
   安装依赖: yum install -y gcc gcc-c++ glibc-devel make automake autoconf ncurses-devel openssl-devel m4 kernel-devel
   下载源码: wget http://erlang.org/download/otp_src_20.3.tar.gz
   解压缩: tar -zxf otp_src_20.3.tar.gz
   编译安装: 
     ./otp_build autoconf
     ./configure && make && sudo make install
   erlang默认安装在/usr/local/(bin,lib/erlang)
   ./configure --prefix 指定安装路径,如 --prefix /opt/erlang/20.3
   ./configure --without-javac 不适用java编译器
2. sublime3依据文件、文件夹创建项目
3. 项目脚本gencmd,protoc,package等,shift命令左移参数,xcopy复制文件夹的文件及子文件夹(/e参数复制所有子文件夹,/s复制非空子文件夹)
4. 项目管理rebar,编译前协议tag与id(crc32循环冗余校验)映射,生成error_msg
5. game_server.spec文件必须为标准json格式,之前第一行出现多余的‘,’,导致报错,另外,xls2erl被编译成了exe文件
6. erlang escript, 第一行#!xxx 指定解释器,若以escript xxx调用escript脚本,则第一行不会起作用,
   第二行%% -*- erlang -*- 可选指令,使Emacs编辑器进入erlang源码模式
   escript脚本可以通过%%! xxx 模拟器传递参数 -smp -sname -mnesia等等
7. 关于配置文件,erlang的配置文件是为app服务的,是app的环境变量,也即xx.app文件中的env选项
   erlang有三种配置,优先级依次如下:
   a> erl -AppName1 Par1 Val1 -AppName2 Par2 Val2
   b> erl -config xxx.config 该配置文件中罗列各app的配置项, [{AppName1, [{Par1,Vla1}]}, {AppName2, [{Par2, Val2}]}].
   c> xxx.app 中的env选项
   项目中用到的其它所谓的“配置”,均由py脚本生成.config文件,再在服务启动的时候,编译生成对应的.beam文件
8. 配置导表时,若报错KeyError则输入的config name在game_server.spec里没找到json对应的key值,另外一种情况是,
   xlrd.open_workbook.sheet_by_name得到table数据时,table.nrows返回的行数超过了有效行数,比如只配了10条左右的数据,但是打印table.nrows时
   显示的是155,远远超过正常的有效行数,导致读取11行时得到的空{},d={},d["label"] 报KeyError
   若报错UTF BOM则是在修改game_server.spec文件后保存的时候默认的编码格式不是UTF-8(记事本默认编码为ANSI)
   sublime text3 -> preferences -> settings -> 在左边栏找到"default_encoding":"xxx",将该行复制到右边栏的user
   文件,值改为"UTF-8",sublime text3默认不显示当前文件的编码格式,可在左边栏找到"show_encoding"和"show_line_endings" 复制到右边栏user文件中,
   将值改为true即可,保存关闭,可在sublime的最右下角看到编码格式
   配置导表,若报错信息不全时,或出现一些不常见报错,需要调试时,可修改py文件,再由pyInstaller重新生成exe(批处理即可),进行调试
9. rebar2 get-deps的时候,重定向deps的下载路径,rebar_deps:download_source/2,  
   第一个参数AppDir是deps下载后的保存路径
10.rebar2 create-app appid=AppName 在当前目录下创建app,生成src/xx.app.src,xx_app.erl,xx_sup.erl
11.JSON (JavaScript Object Notatioin)独立于语言的轻量级的文本数据交换格式
   JSON语法:数据保存在名称/值对中;数据由逗号分隔;花括号保存对象;中括号保存数组
   JSON的值:数字(整形或浮点数);字符串(双引号包裹);逻辑值(true,false);数组([]);对象({});空(null)
   示例: {"player" : {"name" : "xxx", age : 18, "friends" : [{"name" : "yyy"}, {"name" : "zzz"}]}}
12. gen_server:call/3 timeout选项,由手册可知,在实现时,catch gen:call/4,则如果超时,返回的是捕捉到的超时信息:{'EXIT', Reason}
13. io:format函数中显示格式的控制序列为:~F.P.PadModC, F为打印宽度,P为打印精度,Pad(padding)为填充字符,Mod(modifier)为标识符,标识符对参数
    就行解释说明,比如't',表示参数为unicode字符
14. 关于玩家accid,首先是建表时在table option那可以指定初始值(table option也就是设定engine,charset等),其次是采用的mysql依赖库中
    有一个mysql:insert_id/1的接口,参数则为与mysql建立连接的进程,通过trace代码,可以发现进程在parse_ok_packet时,首先解析出来的是affected rows
    第二个解析的就是insert_id,由此可知,mysql返回的响应结果中,包含了insert的自增ID
    mysql返回是protocol发送ok报文,ok报文的结构:0标志位,affected_rows,last_insert_id,server_status,warning_count,message
    至此,解开所有疑惑
15. rebar2编译project时,显示一个'invalid syntax'的Python错误,经查是Python版本不兼容,python -V 发现是2.7,而目前用的是3.5,安装3.5+,更新环境
    变量,再次编译,同样的报错,想起来需要关闭之前打开的终端,重启之后ok
16. httpd有一个最大连接数参数:max_clients,另外再压力测试的过程中,报错:socket_closed_remotely,https://bugs.erlang.org/browse/ERL-473 
    说是httpc handle达到了keep_alive:max,暂时没看懂,明天继续
    通过查阅日志,发现在web_server的cpu负载过高的时候会返回这个错误
    查看手册,httpc_handler进程在收到tcp_closed,ssl_closed时,将session.socket = {remote_close, Socket},进而再调用deliver_answer/1时
    向request进程发送相应,Response = httpc_response:error(Request, socket_closed_remotely),然后就抛出socket_closed_remotely
    (web_server在cpu或内存等资源不足时,是否会主动关闭连接??)
    socket_closed_remotely是服务器主动关闭了连接,一般发生在短时间有大量请求时,如dos攻击,网络爬虫
    由手册可知,在调用httpc:request/N时,HTTP option中的version指明了使用的http协议版本,其默认值为:"HTTP/1.1"
17. 需求需要去掉lists最后一个元素,使用lists:droplast/1,需要列表非空
18. remsh登录后台节点时,报错:Protocal 'inet_tcp':register/listen error: eaddrinuse,查了一下是epmd进程的问题
    epmd(erlang port mapper daemon) 用于erlang节点和ip及端口的映射,一般erl启动时,若启动参数包含-name -sname时就会自动启动该进程
    不过按网上的说法:lsof -i:4369 之后并未发现占用端口的其它进程,只好搬出杀手锏重启大法,问题解决
19. epmd进程启动时,为本地erlang节点分配一个动态端口,通过启动参数-kernel inet_dist_listen_min inet_dist_listen_max可以指定端口范围
    如: erl -sname thd_node1 -kernel inet_dist_listen_min 5860 inet_dist_listen_max 5870 
    极端情况是min=max 则只能启动一个节点,若强行启动第二个节点则会报上面eaddrinuse的错误,原因当然是无法在epmd注册,没有可用端口分配给新节点了
    erlang集群中,节点之间的通信用的就是epmd分配的端口
    可用通过application:set_env(kernel,inet_dist_listen_min/max,xxx)动态改变端口范围,若动态改变则需要重启epmd进程,因为即便erlang节点关闭
    epmd进程仍然会存在,重启使参数生效
    epmd进程默认使用的端口是4369
20. NIF(Native Implemented Functions)可以使我们用C来实现相同的逻辑,但运行速度比纯erlang更快,与C语言程序速度很接近
    nif和port driver都可以用来扩张erlang,但nif更简洁有效
    C编译生成的动态库(.so)在erlang调用C模块时,动态加载到进程空间中,调用nif无需上下文切换,但是安全性不高,nif的crash会导致erlang进程的crash
    nif运行于调度线程,在nif返回前,调度线程不能做其它的事,因此对于执行时间较长的nif,将会影响虚拟机的公平调度
    erlang建议的nif执行时间不超过1ms,对于执行时间较长的可考虑使用脏调度器
    脏调度器:http://www.cnblogs.com/zhengsyao/p/dirty_scheduler_otp_17rc1.html
21. windows编译erlang nif需要erlang的头文件和静态库:$ERL_ROOT\erts9.3\include\, $ERL_ROOT\erts9.3\lib\erts_MD.lib
    翻阅手册,windows下的编译命令: cl -LD -MD -Fe nif_test.dll nif_test.c
    linux下的编译命令: gcc -fPIC -shared -o niftest.so niftest.c -I $ERL_ROOT/usr/include/
    在vps上测试,bingo...
22. 调用lua进行战斗验证的时候,装完lua,试了一下,退出的时候发现常规的q,exit不管用了,搜了一下,unix下用:Ctrl-D,windows下用:Ctrl-Z
23. C与lua通过一个虚拟的栈进行交互,第一个压入栈的索引为1,依次递增,直至栈顶,也可用通过负数来访问栈顶元素,-1为栈顶,依次类推
    lua_to*(lua_State* L, int index)返回返回栈中index索引的值,并将该值转换为星号代表的类型,如:lua_tostring(L,-1) 
    上述函数只是返回对应索引的值,并不会pop,如需弹出元素可使用:lua_pop(L,1)
24. C调用lua函数,可通过:luaL_loadfile(L,filename)加载函数
    加载之后,通过:lua_getglobal(L,"fun_name")将函数压入栈中,接着将函数所需的参数按照定义的顺序依次压入栈中,函数为lua_push*,星号为对应类型
    如:lua_pushinteger(L,3) lua_pushstring(L, "hehehe") lua_pushboolean(L,1)
    参数压入之后,通过:lua_pcall完成对函数的调用,lua_pcall有四个参数,第一个为lua_State,第二个为传递给函数的参数个数,第三个为期望返回的结果个数
    第四个为错误处理的函数索引(在栈中的索引,为0表示没有错误处理函数,若有则需要先把其压入栈中),例:lua_pcall(L,1,1,0)
    lua函数运行错误的话,lua_pcall返回一个非零值,并在栈中压入一条错误信息
    则可通过: if( lua_pcall(L,0,0,0) != 0) printf("error:%s\n", lua_tostring(L,-1)); 打印错误信息
    函数执行完毕,栈中数据被弹出栈,返回值按顺序入栈,即最后一个返回值再栈顶,此时可以通过:lua_to*(L,-1) lua_pop(L,1) 依次获取返回值
25. nif执行时报错,崩掉的是整个节点,也没有crash dump,更没有日志
    nif默认不会被抢占,也不会消耗reduction计数,因此耗时的nif会导致系统延迟
26. lua: /usr/local/include erlang:/opt/erlang20/lib/erlang
27. C编译成动态库之后,在实际运行的时候给变量赋值,此时才会加载path指定的lua文件,而且如果需要更改lua文件,更改之后,调用C中的load函数即可,
    一般load函数包括: L = luaL_newstate(); luaL_openlibs(L); lua状态机变量L,一般可定义为一个全局变量:static lua_State* L = NULL;
28. 调用前端写的脚本时,lua5.1 vs lua5.3的一些坑
    attempt to call a nil value (global 'xxx'):
    loadstring -- load
    unpack -- table.unpack
    attempt to index a functioni value (global 'iparis'):
    原来以为又是一个淘汰的函数,最后看了下,原来是前端大哥脚本写错了(ipairs(xxx)写成了ipairs[xxx])
29. 调用前端lua脚本时,节点直接就崩了,也没留下任何日志或消息,因此需要将lua报错信息返回回来,此处用的lua_pcall/4
    可在头文件或C文件定义一个traceback/1:
    static int traceback(lua_State* L) 
      { const char* msg = lua_tostring(L,-1); if(msg){ luaL_traceback(L,L,msg,1); } else { lua_pushliteral(L,"no message") } return 1 }
    traceback函数应该先压入栈中,比如在调用calc函数之前:
    lua_pushcfunction(L,traceback); 
    lua_getglobal(L,"calc"); lua_pushstring(L,buf);
    if(0 != lua_pcall(L,1,LUA_MULTRET,1)) {
      // trace_back
    }
    const char* calc_result = lua_tostring(L,-1);
    memset(buf,0,1024);
    strncpy(buf,calc_result,strlen(calc_result));
    lua_pop(L,1);
    return enif_make_string(env,buf,ERL_NIF_LATIN1);
30. 在调试lua的时候,为了尽量避免lua文件语法错误,可以在调用之前,在cmd/shell里通过:luac xxx.lua 来检查相关的语法错误
31. C与lua互调时,每次都会得到一个新的虚拟栈,栈里是传递给对方的函数及参数,调用返回时,将栈清空,将返回值放入栈中
    所以可以看到,在调试lua脚本时,加载lua的load函数与调用计算的calc函数,每次都会将traceback函数首先压入栈中
32. erlang进程字典:无锁,hash实现,内容参与gc,根据霸爷的测试,进程字典比ets表快了一个量级,数据为进程私有,不能跨进程,进程销毁时回收
    ets:读写锁,可跨进程读写,只在归属进程销毁时回收,ets数据不在进程中,查询时会复制一份到进程中,进程字典则不存在复制操作
    mnesia:底层由ets/dets实现,跨节点读写,节点间同步由mnesia维护
33. 进程初始大小参数
    min_heap_size 最小堆大小
    min_bin_vheap_size 进程最小二进制虚拟堆大小
    vm启动时可通过+hxxx选项或spawn进程时通过spawn_opt指定数值来定制进程创建的heap初始值
    erlang gc会动态调整堆大小,如果进程空间不足,则会不断的申请和回收内存,所以对于需要较大内存的进程,可指定较大的初始化内存
34. erl_process.h process结构的定义(line966-1131)
    其中进程字典: ProcDict* dictionary; 进程字典可能是NULL
    翻阅erl_process_dict.h, typedef struct proc_dict {unsigned int sizeMask; ...; Eterm data[1];} ProcDict;
    那么很明显进程字典中的数据存放在data数组中,翻阅erl_process_dict.c,先看一下put这个bif函数的对应实现:pd_hash_put
    最关键的当然是ARRAY_PUT(p->dictionary, hval, tpl); 其中
    hval = pd_hash_value(p->dictionary, id); id也就是key,返回的是对应的hash值
    tpl = TUPLE2(hp, id, value); id,value分别对应key,value,返回的是{key,value}的二元tuple
    继续翻阅可知ARRAY_PUT与ARRAY_GET同为数据访问的宏定义(Array access macro),如下:
    #define ARRAY_GET(PDict, Index) (ASSERT((Index) < (PDict)->arraySize), (PDict)->data[Index])
    #define ARRAY_PUT(PDict, Index, Val) (ASSERT((Index) < (PDict)->arraySize), (PDict)->data[Index] = (Val))
    可知,data是一个以key的hash值为下标的Eterm数组,回过来在看pd_hash_value这个宏定义
    #define pd_hash_value(Pdict, Key) pd_hash_value_to_ix(Pdict, MAKE_HASH((Key)))
    #define MAKE_HASH(Term) ((is_small(Term)) ? (Uint32) unsigned_val(Term) :  \
     ((is_atom(Term)) ?  (Uint32) atom_val(Term) : make_internal_hash(Term, 0)))
    MAKE_HASH返回一个Uint32型整数
    pd_hash_value_to_ix(ProcDict* pdict, Uint32 hx)则返回:hx & pdict->sizeMask 作为数组索引,此处有疑问,取余运算后,如何保证索引唯一
    创建进程字典的时候,调用ensure_array_size(ProcDict** ppdict, unsigned int size),data数组的大小则由next_array_size(size)决定,
    返回的是大于等于size的最小的(10 * 2^N),数组每个值填充NIL,也就是说data数组中存在冗余数据,并不是存N个数据,就申请N个空间
 35. 记一个深坑
    先看下两个函数的输出
    erlang:term_to_binary('XYZ_TEST').   ==> <<131,100,0,8,88,89,90,95,84,69,83,84>>
    erlang:atom_to_binary('XYZ_TEST').   ==> <<"XYZ_TEST">>  // io:format("~w",[<<"XYZ_TEST">>])  ==> <<88,89,90,95,84,69,83,84>>
    可以看到term_to_binary/1函数转换的binary,多输出了四个字符(<<131,100,0,8>>),并且有一个131,该值超出了ASCII范围(0-127)
    还有一个坑,在5.5,5.7上用相同的init.sql分别初始化数据库,分别插入上述两个binary串,5.5全部返回ok
    5.7在第一个binary串报错:Error 1366 (HY000) Incorrect string value:'xxx' for column:'yyy' at row 1,第二个binary串ok
36. 关于erlang的调度
    erlang进程在每一个cpu核心上运行一个线程(平衡cpu负载),每个线程运行一个调度器,+sbt参数可以将调度器和核心绑定,可运行的erlang 
    process放入调度器的运行队列,获得时间片后开始运行
    调度器通过复杂的过程在调度器之间迁移一些进程,平衡多个调度器的负载
    抢占指的是调度器能够强制剥夺任务的执行,erlang进程和端口都有一个reduction计数,任何操作都要消耗reduction,包括函数的调用,BIF的调用,进程中
    堆的垃圾回收,存取ETS及发送消息,开发者确保每一步操作都会消耗reduction,reduction消耗完,进程被调度器放回运行队列,接着调用队列中下一个进程,
    因此erlang可以说是真正实现了抢占式多任务并能做到软实时的少数语言之一,相比于吞吐量,erlang更看重的是低延迟,当erlang系统负载较高时,可对耗时
    的任务实现自动降级(计算过程中,更多的被抢占),参考来源:http://jlouisramblings.blogspot.com/2013/01/how-erlang-does-scheduling.html
37. SMP (Symmetrical Multi Processor)对称多处理器
    没有smp支持时,VM在主线程只会运行一个scheduler,调度器从运行队列中取出可运行的进程和IO任务,此时无需对数据进行加锁
    有smp支持的VM可运行多个调度器,并且可通过"+S"参数指定调度器的数量,smp的启动和关闭可通过"-smp [enable|disable|auto]"来指定
38. 补一篇erlang调度原理的中文翻译博客:http://www.cnblogs.com/zhengsyao/p/how_erlang_does_scheduling_translation.html
    服务器在玩家登录的大概前60s,cpu一直飙高,查了下改了两个和调度器有关的启动参数:
    +sbwt none|very_short|short|medium|long|very_long  none则关闭虚拟机调度器的spinlock(自旋锁?),可以降低cpu,之所以降低cpu是因为
    调度器在run out of work之后,进入sleep状态之前,有一段时间的busy wait时间,在此状态仍在会占用cpu,设置为none则会直接进入sleep;
    +swt low|medium|high|very_high 调度器唤醒灵敏度, high,very_high也可降低cpu,不过very_high有可能导致调度器睡死,业务大量堆积时无法唤醒
    参数调整之后,效果不明显,再查有资料说可能是当有大量进程时,调度器消耗过大,要验证需要得到调度器的实际cpu使用率,可用recon:scheduler_usage(time)
    压测时遇到一个坑,每次运行到某条协议时,cpu总是居高不下(即便并发量只有1k的情况下),原以为是虚高,用scheduler_usage发现利用率为1.0,
    简直吐血,后来用etop启动监控,发现进程卡在某个函数,然后用recon:info(Pid,location)定位调用栈,分析了一下,发现notify的split函数每次发送时
    都要循环列表,得到不同状态的子列表,一般来说这个地方不会耗费太多时间,不过瞄了一眼测试使用的数据,列表元素没有低于300,甚至还有3600的
    清理缓存,停服,清数据库,重新开测,登录阶段(1k,10s内)cpu在30%左右,剩余时间基本在25%左右波动,bingo
    调度器相关的参考资料汇总:http://blog.yufeng.info/archives/2963
    http://highscalability.com/blog/2014/2/26/the-whatsapp-architecture-facebook-bought-for-19-billion.html
    http://www.cnblogs.com/lulu/p/3978378.html
    http://www.cnblogs.com/lulu/p/4032365.html
    https://blog.csdn.net/erlib/article/details/40948557
39. spawn(fun() -> etop:start([{sort,memory},{lines,20}]) end).  etop会阻塞进程
40. 压测时,日志进程堆积了大量的消息,大多数进程占用的内存都在32M+(最大45M),此时的标准是2min,216条,Reds指标(reductions计数)都在14w+(最小148014)
    应该进一步调整标准,增加向日志节点写的频率,reductions初始只有2k,10w+的Reds表明进程在执行过程中会被频繁抢占,执行效率不高
41. 压测时,网关进程在读取协议报文长度时,读取到了<"GET ">>,按数值匹配时,值为5亿+,接收缓冲直接爆掉,考虑对解析的报文长度值进行判断,超过阈值,
    可以将玩家主动踢下线,在对长度进行过滤的基础上,对prim_inet:async_recv进行catch,以防缓冲区空间不足
    关于<<"GET ">>,在浏览器中输入服务器地址和端口,回车,复现报错,因为浏览器和服务器建立连接后,发送了一个http请求,http请求报文的
    请求行为:请求方法|空格|URL|空格|协议版本|回车|换行,若为get请求,则请求行的前4个字节为:"GET ",
42. 数据库的读写性能终究还是有上限,此时需要考虑如何冷启动,比如加缓存,减少对数据的二次读写,或者提前加载部分数据
43. prim_inet:async_recv在异步接收的时候,若超时,要考虑处理包未完全接收的情况,尤其是弱网情况下,async_recv(Socket,Length,TimeOut),分别表示
    Socket,要接收的字节数及超时时间,-1 = infinity,单位是
44. prim_inet底层调用的C文件是erts/emulator/drivers/common/inet_drv.c
    对于async_recv,对应的就是tcp_inet_ctl(inet_drv.c:9933)函数的case TCP_REQ_RECV语句块(inet_drv.c:10180)
    在语句块的开始,tcp描述符变量(tcp_descriptor*)desc,首先判断其状态(连接重置econnreset,关闭closed,未连接enotconn)
    当期望读取的字节大于TCP_MAX_PACKET_SIZE时,返回enomem
    若一切顺利,将会调用sock_select函数接收数据,并返回ctl_reply(INET_REP_OK, tbuf, 2, rbuf, rsize)
    异常时返回的是ctl_error(EReason, rbuf, rsize)
    由语句块第一部分对desc状态的检查可知,若erlang的gw进程没有调用发送或接收函数,则底层socket的状态并不会主动返回
    也就是存在一种情况,若对方突然断开连接或断网(未完成tcp的挥手),则本方的底层socket在断开瞬间就改变了状态,但prim_inet及其上层进程并未收到
    状态改变,除非prim_inet及其上层进程主动调用接收或发送函数,才会触发对状态改变的检查,并返回相应的状态
45. 梳理了一下建号流程,发现访问了两次myslq(读写各一次),使用mysql的自增+锁保证ID的唯一
    查询账号表,统计了登录前7分钟每一分钟的注册个数,最高是964(16.06人/s),即便是不考虑百万级的并发,这效率也低的可怜
    优化方案:放弃由mysql保证的ID自增+唯一,在登录服维护一个自增的ID,每次玩家登录,首先访问mnesia,若已注册,则直接返回服务器地址即可
    若为信号,则访问维护ID的常驻进程,获取自己的ID,对维护ID的常驻进程而言,接受玩家注册进程的call消息,需要考虑的是进程mailbox是否能撑得住峰值请求
    目前想到的解决办法,一是提升配置,在创建进程时预先为该进程预留较大的堆内存,二是将该请求作为handle_call第一匹配消息,同时确保消除无效的消息,减少
    消息匹配主循环中无效消息的数量,减少mailbox中消息的堆积
46. 想在start werl时,默认执行 -s xxx yyy 自动启动服务器,结果是不行,不知道有没有其他方法可以实现
47. RFC(Request for Comments)请求意见稿,是由互联网工程任务组(IETF)发布的一组备忘录,是用来记录互联网规范,协议和过程的标准文件,基本的互联网
    通信协议都有在RFC文档内详细说明,常见互联网协议的RFC编号:IP 791, TCP 793, UDP 768, HTTP1.1 2616, FTP 959
48. xmerl处理xml文档,输出vsn和encoding可以通过export/export_simple的第三个参数(子元素列表)实现,具体是prolog属性
    xmerl:export_simple(Content,CallBack,[{prolog, "<?xml version=\"1.0\" encoding=\"utf-8\"?>"}]).
    xmerl_scan:file/1返回的xmlElement并未包含prolog(如果原xml文件有的话)
    发现xmer:export/export_simple在导出的时候,并未包含xmlComment,即便原始的xml文件包含这些内容(比如和客户端有关的并不会每次都用到的一些热更信息)
49. mnesia集群中没有主节点概念,集群节点依次启动之后,对于相同的表而言,读(where_to_read)是本地,写(where_to_write)是所有节点
    dirty_write的SyncMode是async_dirty,向所有involved node发送async_dirty消息,
    由手册知,此处的involved_node = mnesia:table_info(m_table,where_to_read)
50. 我们来看一下mnesia:create_table的调用过程
    mnesia:create_table --> mnesia_schema:create_table --> mnesia_schema:schema_transaction(fun() -> do_multi_create_table(TabDef) end)
    do_multi_create_table返回创建表的操作列表:[{op, create_table, vsn_cs2list(Cs)}],操作列表被插入tidstore,而schema_transaction调用的则是
    mnesia:transaction(Fun) --> mnesia:transaction/6 --> mnesia_tm:transaction/6 --> mnesia_tm:execute_inner/9 --> 
    mnesia_tm:execute_transaction/5 --> mnesia_tm:apply_fun/3 --> mnesia_tm:t_commit(Type), 这个type就是mnesia:transaction/1中的Kind,
    默认是async,mnesia_tm:arrange(Tid,Store,Type)返回#prep{},arrange函数调用do_arrange/5函数便利Store表
    其中,后续用到的#prep.protocol字段,当Store中的Key有op/restore_op时,protocol=asym_trans,arrange函数在调用do_arrange后返回{N,Prep}
    N是Store的大小 --> mnesia_tm:multi_commit/5,第一个参数就是arrange返回的Prep#prep.protocol,multi_commit匹配到asym_trans
    multi_commit中,首先通过ask_commit/5向所有的disc/ram_nodes请求提交,而后调用recv_call/4接收involved nodes的vote结果,
    ask_commit/5的第三个参数是投票结果,默认是do_commit,involved nodes返回vote_yes时并不会更改vote_result
    当有一个node返回的是vote_no时,vote_result={do_abort, Reason}
    vote_result = do_commit时,先调用tell_participants(Pids, {Tid, pre_commit}),向involved nodes具体的Pid发送准备提交的消息
    之后调用rec_acc_pre_commit(Pids, Tid, Store, {C,Local},do_commit, DumperMode, [], [])
    rec_acc_pre_commit接收返回的pre_commit消息,除非返回do_abort/mnesia_down
    rec_acc_pre_commit的Res=do_commit时,先tell_participants(GoodPids, {Tid, committed}),后do_commit/3提交操作
    do_commit/3 --> do_update/4 --> mnesia_dumper:update(Tid, C#commit.schema_ops, DumperMode),do_update_op/3
    mnesia_dumper:update/3 --> mnesia_controller:release_schema_commit_lock() --> cast({release_schema_commit_lock, self()})
    --> opt_start_worker --> opt_start_loader --> load_and_reply --> load_table_fun --> mnesia_loader:disc/net_load_table
    bingo...
51. erlang socket的三种消息接收模式active,passive,active once,可以通过gen_tcp:connect,gen_tcp:listen来设置
    active 主动模式,系统底层接收到数据后,主动向对应的erlang控制进程发消息,该模式不能进行流量控制,若客户端疯狂发包,将会造成控制进程消息堆积
    passive 被动模式,erlang控制进程显式调用gen_tcp,gen_udp的recv/2,3函数,可以控制流量 
    active once/N 混合模式,这种模式的主动仅针对前1/N条消息,之后就进入passive模式,必须显式调用inet:setopts/2重新设置模式,才能接受新的消息
52. erlang:pot_control/port_command --> erts_internal:port_control/port_command --> erl_bif_port.c 
    对应的函数分别为:erts_internal_port_control/command_3 
    其中erts_internal_port_control_3 --> io.c erts_port_control/5, erts_internal_prot_command_3 --> io.c erts_port_output/6
    erts_port_xxx函数的声明在erl_port.h,而erts_port_command/5最后调用的也是erts_port_output/6
    erts_port_control和erts_port_output又是调用的erts_schedule_proc2port_signal/8
    erts_schedule_porc2port_signal/8和erts_schedule_port2port_signal/4调用的是erl_port_task.c erts_port_task_schedule,进行port task调度
53. 需要拿手机号搞些东西,第一个问题就是匹配手机号,erlang的re模块
    re:compile/1,2 第一个参数就是正则表达式,语法也是标准的正则表达式语法
    匹配手机号: {ok, MP} = re:compile("^1[0-9]{10}"). re:run("123xxxxyyyy",MP). 
54. 同集群内mnesia节点之间同步数据,同一份数据copy到不同的节点,如果存玩家的public_info,需要考虑冷热数据的区分,当然第一步是预估数据量的大小
    erlang:system_info(wordsize). 返回当前系统一个word所占的字节
    mnesia:table_info(m_tab,size). 返回m_tab表存储的元素个数
    mnesia:table_info(mtab,memory). 返回m_tab表占用的word,由wordsize可计算出占用的字节,取一下平均,可预估最后占用的内存大小
    大概算了一下,以目前存储的数据,100w条记录大概占用604M
55. erlang节点的一些限制
    进程数量限制,erlang:system_info(process_limit). 可通过启动参数"+P"修改
    端口数量限制,erlang:system_info(port_limit). 可通过环境变量ERL_MAX_PORTS或启动参数"+Q"修改
    ets表数量限制,erlang:system_info(ets_limit). 可通过启动参数"+e"修改
    atom数量限制,erlang:system_info(atom_limit). 可通过启动参数"+t"修改
    同时打开的文件和socket数量取决于能使用的端口数量及操作系统限制
56. supervisor监督树结构中,若type=worker类型的子进程再生成一个supervisor的子进程,该子supervisor进程仍然在整个监督树中,也就是说supervisor的
    父进程不一定非得是supervisor进程
57. try catch可以嵌套,不过try语句块避免使用尾递归,因为erlang虚拟机始终保持对try块的引用,以防出现异常,所以try块调用尾递归会造成大量内存消耗
58. gcc -fPIC -shared -o xxx.so xxx.c -I"XXXX" -lluajit-5.1时,提示:Failed to load NIF library libluajit-5.1.so.2 cannot open shared         object file no such file or directory, ls /usr/local/include 可知libluajit-5.1.so.2是一个软连接文件
    源文件为同目录下的libluajit-5.1.so.2.1.0,但是为什么会提示打不开该文件?
    通过-L参数直接指定libluajit-5.1.so.2.1.0,提示:Failed to load ... undefined symbol lua_settop
    若指定同目录下的libluajit-5.1.a,报错同上
    gcc -l 指定要链接的库名,完整的库名包括开头的lib和".so",".a"的结尾
    gcc -L 添加链接库的搜索路径
    gcc产生警告信息的编译选项,大多以-W开头,常用的是-Wall
    预处理、编译、汇编、连接,连接库时,默认情况下编译器会优先加载动态库,如需要强制加载静态库,可以通过-static选项,如:gcc xxx -static -llib_name yyy
    静态库搜索顺序:-L参数,环境变量LIBRARY_PATH,默认路径(参考动态库的默认路径)
    动态库的搜索顺序:-L参数,系统变量LD_LIBRARY_PATH,/etc/ld.so.conf文件(修改文件后,需执行ldconfig),gcc安装时配置的路径(gcc --print-search-dir | grep libraries),这个也是默认路径,一般是/usr/lib,/lib
59. rebar可通过port_specs,port_env这两个配置项编译nif模块,具体的编译连接参数可参考rebar_port_compile.erl,一般是配置下CFLAGS,CXXFLAGS,LDFLAGS
    如果考虑到跨平台,可在rebar.config.script脚本中,通过os:type()决定相应平台的配置,达到动态更改rebar.config的目的
    关于nif的编译配置,还可以翻阅rebar.config.sample文件
60. 关于第58条,"提示打不开该文件?", 编译时通过-L指定lua动态库所在路径,不过运行时并未生效,解决方法之一是修改/etc/ld.so.conf文件
    在文件末尾加上和-L参数一样的路径(/usr/local/lib),保存,退出,执行ldconfig命令,再次编译,在erl shell里加载,bingo...
    -L "链接" 的时候,去搜索的路径,它只是指定了程序在编译链接时库的路径,并不影响程序 "执行" 时库的路径,程序执行时,系统还是会在默认路径下查找库
    如果找不到,还是会报cannot open shared object file,此时可以通过修改LD_LIBRARY_PATH环境变量(无需root权限)或者/etc/ld.so.config文件
    error loading module 'xxx' from yyy/zzz.lua:1: unexpected symbol near 'aaa', 如果是第一行就报错,则zzz.lua文件的格式不是utf8
    有可能是utf-8 with BOM或者其他格式,在编辑器中save with encoding保存即可(又一个前端的坑...)
61. 典型编译命令: gcc -fPIC -shared -o xxx.so xxx.c -I. -I/path/to/erl/include -I/path/to/lua/include -g -Wall -Werror -O3 
    -fno-strict-aliasing -lstdc++ -L/paht/to/other/lib -lother_lib_name
    编译参数释义参考
62. 准备给项目加上Nginx,用于客户端的热更新,Nginx的几个常用命令(手册):
    若指定--prefix=xxx,则安装后,将该目录加入PATH,直接调用nginx,即可启动服务
    nginx -s reload 重新加载nginx.conf文件
    nginx -s stop 快速关服
    nginx -s quit 优雅关服
    nginx.conf常用配置参数说明
    user nobody # 运行用户
    worker_processes 1 # 工作进程数量,通常和cpu数量相等,worker processes的最优值取决于(不限于)cpu数量,存储数据的硬盘数量和负载模式
    worker_rlimit_nofile 100000 # worker进程的最大打开文件数限制,如果没有设置,则为操作系统限制
    error_log logs/error.log # 全局错误日志
    pid logs/nginx.pid # PID文件,nginx启动后,有一个master process, master的pid保存在nginx.pid文件中,可由ps命令查看master和worker进程信息
    events {} # 事件模块,nginx中处理网络连接的配置
    events.accept_mutex [on|off] # 默认为on,使用连接互斥锁进行顺序的accept()系统调用,即网络连接的序列化,防止发生惊群现象(一个网络连接到来,
    同时唤醒多个睡眠进程,但是只有一个进程能获得连接,这样会影响系统性能)
    events.accept_mutex_delay Nms # 默认500ms,如果一个进程没有互斥锁,它将至少在这个值的时间后被回收
    events.debug_connection [ip|CIDR] # 0.3.54版本后,这个参数支持CIDR地址池格式,这个参数可以指定只记录由某个客户端IP产生的debug信息
    events.multi_accept [on|off] # 默认为off,nginx接到一个新连接通知后调用accept()来接受尽量多的连接(单个进程是否同时接受多个连接)
    events.use [kqueue|rtsig|epoll|/dev/poll|select|poll|eventport] # 指定事件驱动模型
    events.worker_connections # 单个worker process进程的最大并发连接数
    http {} # http模块
    http.sendfile on|off # 设置为on,则启用linux上的sendfile系统调用,减少了内核态和用户态之间的两次内存拷贝,在磁盘中读取文件后,直接在内核态
    发送到网卡,提高发送文件的效率
    http.keepalive_timeout Ns # server端对连接的保持时间,默认75s
    http.send_timeout Ns # 客户端在Ns内未接收nginx发送的数据包,则nginx关闭该连接
    http.sendfile_max_chunk xx # 每个进程每次调用时传输的最大大小,默认为0,即不设上限
    http.client_max_body_size XX # 客户端发送较大http包体的数据时,nginx不需要接收到完整的包体,就可以告诉用户请求过大,不被接受
    http.include /path/to/file # include只是一个包含另一个文件的命令, include /usr/local/nginx/conf/mime.types; 引入网络资源的媒体类型
    http.gzip on # nginx采用gzip的压缩的形式发送数据,可减少发送的数据量
    http.error_page error_code url # 错误页,例:error_page 404 www.google.com;
    http.server.listen # 侦听端口or地址:端口,listen 8080;127.0.0.1:8088;*.8090
    http.server.server_name # 设置服务器名,nginx解析http请求的host头,和server模块进行匹配
    http.server.location # 匹配URL,执行不同的应用配置
    location有精确匹配,前缀字符串匹配,正则表达式匹配和通用匹配四种方式,优先级依次如下:
    location = /xx/yy {} % 精确匹配,优先级最高
    location ^~ /xxx/ {} % 提高匹配优先级的前缀字符串匹配,优先级次之,匹配之后不再进行正则表达式匹配,如果没有^~ 即便前缀匹配到了,仍要进行正则匹配
    location ~[*] reg_exp {} % 区分/不区分大小写的正则表达式匹配,~区分大小写,~* 不区分大小写
    location /xxx/ {} % 前缀字符串匹配,为提升优先级,只有在正则不匹配时,才会采用该匹配
    location / {} % 通用匹配,匹配所有请求
    upstream是nginx的负载均衡模块,常见的配置方式如下:
    轮询, upstream xxx {server name_or_ip1:port1; server name_or_ip2:port2;} 按请求的时间顺序,依次分配不同的server;
    权重, upstream xxx {server name_or_ip1:port1 weight=w1; server name_or_ip2:port2 weight=w2;} weight字段为轮询的几率,值越大,几率越高;
    ip_hash, upstream xxx {ip_hash; server name_or_ip1:port1; server name_or_ip2:port2;} 按访问ip的hash结果分配server,则同意客户端访问
    同一个服务器,可解决session问题
    fair, upstream xxx {server name_or_ip1:port1; server name_or_ip2:port2; fair;} 按后端服务器的响应时间来分配,时间短的优先分配,该方式
    需要安装upstream fair balancer模块
    url_hash, upstream xxx {server name_or_ip1:port1; server name_or_ip2:port2; hash $request_uri; hash_method yyy;} 
    按方位url的hash结果分配server,在upstream模块中加入hash语句,hash_method指定了hash算法(crc32),需要安装upstream hash模块
    upstream的使用方法是配合http.server.location.proxy_pass字段,proxy_pass = "http://" + upstream_name即可(proxy_pass=http://upstream_1)
    nginx的几个常见配置项:$remote_addr,$http_x_forwarded_for 记录客户端ip地址;$remote_user 记录客户端用户名称;$time_local 访问时间和时区;
    $request 记录请求的url和http协议;$status 记录请求状态;$http_user_agent 记录客户端浏览器的相关信息;
    $http_referer 记录从哪个页面链接访问过来的; $body_bytes_sent 记录发送给客户端文件的主体内容大小
    nginx源码理解:https://www.kancloud.cn/digest/understandingnginx/202587
    nginx动态代理:https://segmentfault.com/a/1190000007059973
63. lua源码安装的时候,/usr/local/lib默认只生成了liblua.a的静态库文件,可以通过修改两个Makefile文件,实现在编译安装的时候同时生成.a和.so文件
    参考文章:https://blog.csdn.net/yzf279533105/article/details/77586747
64. nginx缓存设置相关参数
    proxy_cache_path /path/to/cache # 定义缓存文件的存储路径
    levels=1:2 # 设置缓存文件的保存方式,未设置则直接保存到缓存路径,1:2表示缓存文件将根据其key的md5值保存在缓存路径的子目录中,最多创建3层
    keys_zone=cache_name:meta_size # 定义缓冲区名称和缓存key和其它元数据的存储空间
    max_size=xxx # 最大磁盘缓存空间
    inactive=xxx # 设置缓存时间,60m 则60分钟没有访问就删除
    proxy_cache_key "$scheme$proxy_host$uri$is_args$args" # 区分缓存文件的key,前项示例为默认值
    location.proxy_cache cache_name # location块内,指定缓存区域
    清除nginx缓存的插件:ngx_cache_purge
65. 整理一下线上定位问题的大致思路:优先看内存和进程消息队列,其次可看相关的日志,最后是一些统计数据(binary,atom大小,网关/玩家进程大小)
    对于MMO/RPG,卡顿的时候,可以看下地图进程的帧率,也就是每秒对地图数据进行多少次更新
    压测或者上线后,另一个指标:网关进程和玩家进程的CPU与内存占比,据此可进一步优化
66. httpc:request/4收到socket_closed_remotely,做了个测试
    test_pressure(Times,Rem) -> timer:tc(fun() -> lists:foreach(Seq) -> _Return = httpc:request(get,{Url,[]},[],[]), 
      case Seq rem Rem of true -> timer:sleep(1); _ -> next end, lists:seq(1,Times) end).
    总的思路是尽可能的提高并发性,如:Times=10000,Rem=100,在该条件下用时5.04s -- 5.21s之间,大概每秒不到2k次请求
    在该状态下:netstat -anlp | grep :port | wc -l 统计连接数量,去掉统计,可看到大量连接处于'TIME_WAIT'状态
    当:Times=10000,Rem=1000时,第一次用时4.99,紧接着再次用上述参数调用,返回socket_closed_remotely,同时可以看到服务端上次调用的连接还未释放
    连接数由101-->193,此时服务端连接数只有50,其他参数如cpu,memory及节点的port_count等均未有异常,莫非云平台有防DDOS机制?
67. 项目中一些值得注意的地方
    玩家public_info的大小:750条记录,89030Word=712240B,0.927k/p,也就是说单个玩家大概有1k
    玩家gateway进程大小:88616B=86.53kb,单个玩家大概106k-1088k不等,也就是说网关进程最多的时候也会有1m大小
    之前大概估算过,单个玩家进程大概在2m左右,那么如果网关按1m估算,则此时网关与进程的比值将会是1/2,若按100k算,比值为1/20
    目前网关进程的大小即便是在这么小的取样范围内,最高都达到了1m,那么很有必要对网关进行gc,尽量将网关进程控制在100k左右
    对mnesia做了个小测试,将1-1000000以{{seq,Num}, Num}的格式,插入到mnesia表中,用时4.656s,memory为91.575m
    再将上述100w条数据删除,用时3.812s,重复上述操作,插入和删除分别用时4.625s,3.547
68. 可通过在DNS服务器上添加主机A记录的方式,实现负载均衡
    DNS域名解析负载均衡有点是:实现简单,方便;省去维护负载均衡服务器的麻烦;DNS还支持基于地理位置的域名解析,返回距离用户最近的服务器地址,加速访问
    DNS的主要缺陷是:修改记录后,有生效时间;无法检测服务器状态,可能会一直返回down掉的节点ip;无法区分服务器的处理能力
69. 隔壁项目出现回档,原因是大量玩家同时登陆,而回写的随机时间范围较小,导致比较多的玩家在同一时间回写mysql(单点回写峰值),
    回写失败后的处理也不太合理,采用的是将玩家踢下线,把数据写到mnesia,而mnesia的操作是脏读,脏写,(cpu和内存没有问题),
    以下为猜测:玩家被踢下线后,估计是又点了登录,又走了一遍初始化流程,从mysql载入数据,并放到mnesia中,由于是脏操作,估计把踢下线之前的数据给覆盖了
    问题的关键在于,如何将大量在线玩家的回写时间,尽可能的均匀分散,升级数据库性能,适当扩大随机回写的范围(5min-->10min),另外还可以考虑在玩家进程
    回写的时候,catch一下异常,若回写失败,则放弃此次回写,再随机一个新的时间,新的时间可适当扩大/缩小范围
    大量玩家同时回写的另外一种情景:网络出现闪断,玩家同时下线,同时调用回写函数
70. rebar2编译NIF的配置选项:port_specs,port_env
71. 资源监控工具:https://www.cnblogs.com/arnoldlu/p/9462221.html
72. 隔壁项目节点崩了,记一下erl_crash:
    =erl_crash_dump:0.3
    Sat Mar  2 17:38:00 2019
    Slogan: Kernel pid terminated (application_controller) ({application_start_failure,kernel,{{shutdown,{failed_to_start_child,
    net_sup,{shutdown,{failed_to_start_child,net_kernel,{'EXIT',nodistribution}}}}},{k
    System version: Erlang/OTP 18 [erts-7.3] [source] [64-bit] [smp:8:8] [async-threads:10] [hipe] [kernel-poll:false]
    Compiled: Wed Sep 14 10:27:19 2016
    net_kernel nodistribution 这个还是第一次见,搜了一下,有一篇参考,不过貌似别人是因为错误的iptables rule引起的
    先放下参考,后续研究下:http://erlang.2086793.n4.nabble.com/net-kernel-fails-to-start-nodistribution-td4714181.html
73. 任何线上数据库的操作,都要进行备份,另外为了防止误操作(rm -rf),一定要进行严格的权限管理
    可采取的方案包括:首先在搭建环境的时候,使用root账户,将mysql和mongo的数据目录放在根目录的/data下;
    其次新建普通账户,服务器搭建,程序文件更新均使用普通账户
74. erlang:process_info(Pid,Item_or_Itemlist).
    第一个参数Pid,必须在本地,原因嘛就是: erl_bif_info.c(line.1006) process_info_2函数有一个判断,is_not_internal_pid(pid)
    当Pid为外部节点进程时,返回:BIF_ERROR(BIF_P, BADARG);
75. 遇到一个奇怪的问题,和mnesia:dirty_write/1有关,问题描述如下:
    mnesia:dirty_write({t1,k1,v1}); mnesia:dirty_write({t1,k2,v1}); mnesia:dirty_write({t1,k1,v2}); mnesia:dirty_write({t1,k3,v2});
    第1/2条语句执行之后,第3/4条语句执行之前,一般至少间隔1-2s,更长时间也会有几十s(87s)会依次执行下面两条查询语句:
    mnesia:dirty_read({t1,k3); mnesia:dirty_read({t1,k1}), 若两条查询都为空,才会执行3/4的insert函数
    现在的问题是,在上述流程走完之后,只剩下后边3个插入的数据:{t1,k2,v1}, {t1,k1,v2}, {t1,k3,v2}, 缺失第一条插入的数据{t1,k1,v1}
    按照手册 mnesia:dirty_write/1 --> mnesia_tm:dirty/2 protocol默认为async_dirty,在提交数据之前,dirty函数调用mnesia_tm:prepare_items
    返回的是一个名为prep的record,主要用的是#prep.records这个数组,包含的是需要同步数据的节点,然后调用mnesia_tm:async_send_dirty/4
    async_send_dirty/4函数首先会判断第一个send_node是不是本地&where_to_read_node,如果满足这两个条件则调用mnesia_tm:do_dirty/2函数
    do_dirty/2函数调用mnesia_log:log/1;mnesia_tm:do_commit/3,提交本地node
    若不同时满足上述两个条件,则进入第二个判断,send_node是否等于where_to_read_node,如果是则向read_node发送同步消息:{sync_dirty, _, _,_}
    对其余节点发送异步消息:{async_dirty, _, _,_}
    还有就是若要同步数据的node == node() == where_to_read_node,mnesia_tm:async_dirty/2的返回结果就是mnesia_tm:do_commit/3的返回结果
    do_commit/3 依次调用:do_snmp/2; 三次do_update/4,参数区别是ram_copies,disc_copies,disc_only_copies; do_update_ext/3;
    返回结果则是上述函数调用的返回结果,因为本地表为ram_copies,且为dirty_write,重点关注一下do_update/4 --> do_update_op/3(#line.1814)
    do_update_op/3 --> mnesia_lib:db_put/3 --> mnesia_lib:db_put(ram_copies, Tab, Val) -> ?ets_insert(Tab, Val), ok;
    mnesia.hrl --> ?ets_insert/2宏定义调用的是ets:insert/2(#line.27)
    由上述流程可知,mnesia:dirty_write/1,在写入本地节点的时候,最后调用的是ets:insert/2函数,那么问题来了,为什么在间隔超过20s以上的情况下
    3/4两条插入语句还能执行,或者为什么,插入语句没有读到已插入的第一条数据?
76. 接75,上述问题也有可能是dirty_read的时候返回undefined,导致顺利执行3/4两条插入语句
    mnesia:dirty_read/1 --> mnesia:dirty_rpc/4 参数是(Tab, mnesia_lib, db_get, [Tab,Key]) --> mnesia_lib:db_get/2
    mnesia_lib:db_get/2 --> db_get(ram_copies, Tab, Key) -> ?ets_lookup(Tab, Key);
    mnesia.hrl --> ?ets_insert/2宏定义调用的是ets:lookup/2(#line.25)
    dirty_read/1的实现更加的简洁明了,貌似不太可能出问题,那么是dirty_write/1在某些情况下不稳定,会丢失数据?
    要么mnesia的实现机制有隐藏的bug,在某些情况下会导致read/write操作不稳定?
