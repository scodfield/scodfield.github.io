
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
8. 配置导表时,若报错KeyError则输入的config name在game_server.spec里没找到json对应的key值
   若报错UTF BOM则是在修改game_server.spec文件后保存的时候默认的编码格式不是UTF-8(记事本默认编码为ANSI)
   sublime text3 -> preferences -> settings -> 在左边栏找到"default_encoding":"xxx",将该行复制到右边栏的user
   文件,值改为"UTF-8",sublime text3默认不显示当前文件的编码格式,可在左边栏找到"show_encoding"和"show_line_endings" 复制到右边栏user文件中,
   将值改为true即可,保存关闭,可在sublime的最右下角看到编码格式
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
