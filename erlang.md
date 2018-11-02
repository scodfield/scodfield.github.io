
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
