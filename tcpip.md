Windows下几个影响TCP并发连接的参数：
0. Ctrl+R --> regedit --> HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\services\Tcpip\Parameters
1. TcpNumConnections
   maximum number of connections of that TCP can have open simultaneously
   0 or 未设置该参数表示可以建立任意的连接(不过其它条件,如tcb,可用端口数or系统内存会限制连接上限)
2. MaxUserPort
   highest available user port when requests from the system
   typically, ephemeral ports(those used briefly) are allocated to port numbers 1024 through 5000
   服务器监听固定的端口,所以并不会存在端口不够用的情况,而客户端建立socket连接时,默认并未bind本地端口,因此需要向
   系统申请一个可用的端口,当并发连接足够多时,会出现端口不够用的场景
3. MaxFreeTcbs
   TCP control blocks(TCBs), each connection requires a control block
   TCB是一种内存驻留结构,包含套接字编码,传入和传出数据缓冲位置,已收到或未确认的字节及其它信息
   为了快速检索这些信息, windows server将TCB存储在一个哈希表中
4. MaxHashTableSize
   TCP control blocks(TCB) stored in a hash table, the value must be power of 2
   
Linux下TCP并发连接数量限制：
1. 本地端口范围
2. 单进程同时打开的文件句柄上限,默认为1024,可通过ulimit -n xxx 命令进行调整
3. 系统同时打开的文件句柄上限,可通过/etc/sysctl.conf调整
4. 系统内核的IP_TABLES防火墙对最大跟踪的TCP连接数有限制
linux TCP高并发参考:https://www.cnblogs.com/lemon-flm/p/7975812.html

Tips to remember:
1. TCP报文有两部分组成：头部和数据部分,头部的各字段体现了TCP的使用和功能,通信的另一端收到tcp报文后,去掉头部,组装接收到的消息
   头部的前20个字节是固定的,这也是头部的最小大小,头部最大60个字节,结构如下:
   a> 源端口 2byte
   b> 目的端口 2btye
      源端口告知数据来自于哪个应用程序,目的端口指明数据发送给哪个应用程序,应用程序绑定端口,结构大致如下：
      程序1 <----> Sock1(源&目的ip,源&目的port) <----> Sock2(源&目的ip,源&目的port) <----> 程序2
   c> 序列号 seq number 4byte
      tcp面向字节流,序列号就是所发送报文的首字节在字节流中的位置或序号
   d> 确认号 ack number 4byte
      确认号是接收方希望接收到的下一个报文的首字节的序号
      tcp双向通信,所以在三次握手建立连接的时候,双方同步seq number,更新各自的ack number
      这里想到个问题,在第三次握手的时候,客户端发送ACK包,其中 ACKbit=1,ack_number=服务端seq number+1,那此时客户端的seq number=??
      又或者,服务端ack number已在第一次握手时同步了客户端seq number,在连接建立之前(此时为syn_recv状态)or接收到客户端ACK包时
      并不会更改自身ack number,也即第三次握手时,客户端无需关心seq number
   e> Data offset数据偏移(首部长度) 4bit
      数据偏移指的是报文的数据部分起始处到报文起始处的距离(刚好度量了首部的长度),一个偏移量是4byte,4bit最大是15,所以TCP头部
      的最大值是60byte
   f> 保留位 6bit 
      以后使用,目前为0
   g> 标志位 6bit
      6bit标志位,每一个bit都有特定用途：
         1> 紧急 URG(Urgent) URG=1表示紧急指针(urgent pointer)有效,告诉系统该报文有紧急数据,需尽快发送,与下面的紧急指针配合使用
         2> 确认 ACK(Acknowledgment) ACK=1时确认号(ack number)有效,ACK=1的报文一般成为'确认报文'(这解决了第三次握手时客户端序列号
            的问题,第三次握手客户端发送的是'确认报文',此时ACK=1,SYN=0),连接建立后发送报文时,ACK=1
         3> 推送 PSH(Push) PSH=1表示该报文是高优先级,接收方(内核)应尽快将该报文交付给后续的应用程序,而无需等整个TCP缓冲满后再交付
         4> 复位 RST(Reset) RST=1表示tcp连接出现严重错误,需释放并重新连接,一般称该报文为'复位报文'
         5> 同步 SYN(Synchronization) SYN=1表示该报文是一个请求连接的报文,一般称为'同步报文',三次握手的第一个报文
         6> 终止 FIN(Finish) FIN=1表示该报文发送方的发送数据已发送完毕,请求释放连接,一般称为'结束报文',四次挥手断开连接时用到此标志
   h> 窗口大小 window size 2btye
      告诉对方本段tcp接收缓冲区还可以接收多少字节的数据,对方据此控制发送的速度
      窗口大小指的是从本报文段头部ack number起,可接收的字节数
   i> 校验和 TCP Checksum 2byte
      校验和由发送端填充,接收端对整个报文进行CRC计算,目的是检验报文段在传输过程中是否损坏,损坏则丢弃
   j> 紧急指针 Urgent Pointer 2byte
      仅在标志位URG=1时有效,指定了本报文段中紧急数据的字节数
      紧急数据在报文段数据部分的首部,其后是普通数据,即当URG=1时,TCP报文:头部+数据(紧急+普通)
2. TCP连接建立（三次握手）由内核协议栈实现,连接建立后socket状态转为established,并被放入icsk_accept_queue，accept()被唤醒,返回socket
   TCP是可靠协议,可靠的保证机制包括:ack确认,超时重传,滑动窗口
   前两个很好理解,对于滑动窗口需要先了解tcp流中的数据分类,以下是发送方数据分类: 
   a> sent and acknowledged 发送且已被确认的数据,这部分数据是在窗口外的
   b> send buy not yet acknowledged 已发送但未被确认的数据,这部分属于窗口内的数据
   c> not sent, recipient ready to receive 未发送但是接收方可以接受的数据,这部分是已加载到发送缓存中,等待发送的数据,也在窗口中
   d> not sent, recipient not ready to receive 未发送且接收方也不允许发送的数据,该部分数据超出了接收方的缓存
   以下是接收方数据分类:
   e> received and ack, not send to process 已成功接收但还未被上层应用程序接收
   f> received not ack 已接收,还未恢复ack
   g> not received 有空位,还没有被接收的数据
   对于发送方来说,窗口包含两部分,一个是发送窗口(已发送,还未收到ack),一个是可用窗口(可以发送且接收方可接受),发送端窗口(两个部分)的大小根据
   接收端的接收情况,进行动态调整,接收端Ack表明已成功接收的字节序,发送端的发送窗口向右移动,至于可用窗口的大小,根据接收端报文头部的window_size
   字段进行调整,如果接收端发送的win_size=0,则发送方进入零窗口,在此期间停止发送数据
3. listen()开启监听队列,客户端SYN包到来,创建新sock,sock为状态TCP_SYN_RECV,并被存入半连接队列syn_table中
   socket编程服务器端socket流程:
   a> socket_create 创建socket,一般需要指定网络协议(AF_INET-ipv4,AF_INET6-ipv6,AF_UNIX等),套接字流类型(tcp/udp等套接字),
      协议类型(和上一个参数相关,为对应流类型的tcp/udp协议)
   b> socket_bind 绑定套接字,参数一般包括socket和ip以及port,将创建的socket套接字绑定到(ip,port)构成的这个地址,如果有其它socket需要连接
      这个socket,那么指定(ip,port)构造的地址就可以找到该socket
   c> socket_listen 监听套接字,参数就是上述创建的socket,默认创建的套接字是可以主动connect其它套接字的主动套接字,listen告诉内核某个套接字
      可以接受其它socket的连接请求,把默认的主动套接字变为被动套接字
   d> socket_accept 等待套接字的连接请求,参数为listen返回的被动套接字,accept函数从处于监听状态的socket的连接请求队列中,取出一个请求,创建
      一个新的套接字,与客户端套接字建立连接通道,如果连接成功,返回新创建的套接字fd,失败则返回invalid_socket
4. SYN攻击:客户端伪造大量IP地址,不间断的向服务器发送SYN包,塞满服务端半连接队列,导致正常的SYN请求被丢弃,SYN攻击是DDos攻击的一种,检测SYN攻击
   netstat + awk '/^tcp/' 查看SYN_RECV状态的tcp连接即可
5. 在进行本地压测的时候,出现一个情形,首次启动客户机没有任何问题,加机器人再次启动客户机时,日志报:connrefused
   第一反应是上次的客户机tcp连接没有释放,由于客户机tcp连接设置了保活机制(keepalive),而默认的保活检测时间过长(默认是7200s),以及保活重发次数
   另外一个会影响的参数是系统释放连接资源之前的等待时间(TcpTimedWaitDelay),系统不释放资源导致本地动态端口,TCB,TCB hashtable等系统资源不足,
   从而拒绝连接,注册表调整参数
   a. Ctrl+R --> regedit --> HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\services\Tcpip\Parameters 
   b. 右键 --> 新建 --> DWORD --> 输入键值 --> 选中键值&右键 --> 修改 --> 基数栏选择十进制&填入具体参数值
   c. MaxUserPort 65534
   d. MaxFreeTcbs 16000
   e. MaxHashTableSize 65536
   f. TcpTimedWaitDelay 30 (s)
   g. KeepAliveTime 120000 (ms)
   h. KeepAliveInterval 1000 (系统未收到响应而重发保活信号的间隔,ms)
   调整参数,关闭注册表,重启电脑即可
6. windows下查看tcp、udp及端口等统计情况:netstat -an | find "ESTABLISHED" /c 统计活跃状态的tcp连接,状态与linux类似,包括
   LISTENING,CLOSE_WAIT,ESTABLISHED,TIME_WAIT
7. netstat参数 
   -a 显示所有连接和监听端口
   -n 以数字形式显示地址和端口号
   -o 显示与每个连接相关的所属进程
   -p proto 显示proto指定的协议连接,包括TCP,UDP,TCPv6,UDPv6
   比如显示所有tcp连接: netstat -an -p tcp
8. windows下与linux类似的为find,常用参数如下: 
   /v 指定不包含指定所有行
   /c 对指定的行进行技术
   /i 不区分大小写
   "" 指定要搜索的字符串
   比如显示端口为20001的所有连接: netstat -ano | find "20001" ,类似的还有一个findstr,不需要(""),如: netstat -ano | findstr 20001
9. HTTP/1.0 短连接,每次request都会建立一个单独的连接,因此请求较多时,连接的建立和释放会占用大量的系统资源
   HTTP/1.1 支持长连接,管线处理,在一个连接上可以传送多个请求和响应,并且客户端可以在上一个请求未返回前再次发送请求,不过服务器则需要保证
   按客户端请求的顺序,返回响应,1.1还新增了一些请求/响应头域来扩展功能,比如:status code, request method(options,put,delete..),host域
   HTTP状态码, 1xx 消息,服务端已收到请求,需要请求者继续执行操作; 2xx 成功,请求收到并被处理; 3xx 重定向,需要进一步操作完成请求; 4xx 客户端
   错误,语法错误或找不到请求的资源; 5xx 服务端错误,服务器在处理请求的过程中出现错误
   常见状态码:
    200 ok 请求成功; 305 use proxy 请求的资源必须通过代理访问; 400 bad request 客户端语法or请求参数有误; 401 unauthorized 请求需要用户
    认证; 403 forbidden 服务器理解请求,但拒绝执行; 404 not found 服务器未找到请求的资源; 500 internal server error 服务器内部错误; 501 
    not implemented 服务器不支持请求的功能; 502 bad gateway 网关or代理服务器尝试发送请求时,从上游服务器收到一个无效的响应
10. 服务器在httpc:request/4时报错:inet eaddrinuse, 同时:erl -sname test  也会报同样的错误,POSIX Error Codes说是address被占用
    通过对比其他资料,这个address应该就是需要占用的端口,httpc申请的是动态端口,那么很有可能是动态端口不够用
    /proc/sys/net/ipv4/ip_local_port_range 查看本地TCP/UDP端口范围
    可由: echo '32768 60999' > /proc/sys/net/ipv4/ip_local_port_range  更改端口范围
11. TCP/UDP端口类型:
    0 - 1023 固定端口,与常见应用紧密绑定的端口
    1024 - 49151 注册端口,与应用松散绑定
    49452 - 65535 动态端口,可与任意应用绑定
12. 做个试验,在浏览器直接输入:server_ip:port/anything, 查看日志
    11:03:34.564 gw进程init ok
    11:03:34.564 gw进程do_terminate,报head_error,这个是因为在解析http请求时候,<<Len:32>> = <<"GET ">>,解析协议头部发送的body长度时
    Len解析出来的是1195725856,过滤异常包,直接断掉连接即可
    11:03:34.567 gw进程init ok (浏览器在上一次断开后,紧接着就再一次的发起连接请求)
    11:03:35.597 gw进程报相同的错误
    11:05:15.616 11:05:15.617 发起两次请求,#Port<0.7882> #Port<0.7883>,有意思的是7882有打印报错信息,7883没有打印任何信息,且与11:10
    发起的连接相比,#Port<>不在连续
    11:10:16.597 11:10:16.598 11:10:18.629 发起三次请求,#Port<0.9078> #Port<0.9079> #Port<0.9080>,三次请求连续,且均有打印错误信息
    后续直到12:04,浏览器未再发起连接请求,试验用的是chrome浏览器
    12:05:27.454 在firefox上试验了一下,发送的请求次数更多(10次), 12:05:27.454 -- 12:05:36.562, #Port<0.9747> -- #Port<0.9758>
    不过Firefox直到12:15:41都没再发起请求,估计各个浏览器在实现上不太一样
13. tcp挥手,无所谓客户端和服务端,都可以主动发起FIN,关闭的原则是先关读通道,再关写通道,比如典型的客户端发起关闭连接,服务器收到FIN报文后,关闭
    服务器的读通道,客户端收到ACK后关闭自己的写通道,后两次挥手重复操作
    被动关闭方在发送ACK后,进入CLOSE_WAIT状态,该状态下,被动方要做的就是检查是否还有要发送给对方的数据,若没有,那就可以关闭这个socket了
    为什么三次挥手,其中的一个原因是防止已失效的连接请求报文又发送到了服务端
14. a> OSI七层网络模型从下到上为:物理层,数据链路层,网络层,传输层,会话层,表示层,应用层,第二层数据链路层上的数据称为帧(Frame),第三层
       网络层上的数据称为包(Packet),第四层传输层上的数据称为段(Segment);
       表示层的功能是对数据进行格式化,代码转换,数据加密等;
       会话层提供会话控制,提供的服务包括认证(Authentication),权限(Permission),会话恢复(Session restoration)
       传输层提供端到端的接口,网络层为数据包选择路由,数据链路层传输有地址的帧及错误检测,物理层以二进制形式在物理媒体上传输数据;
    b> TCP/IP协议族采用的是五层模型,将表示层和会话层放到应用层中,不过在使用过程中,出现了一个提供安全加密服务的层(SSL/TLS),有了安全层
       之后各应用层协议都可以加上一个S(Security),如HTTPS就是原本的HTTP协议有了SSL/TLS的保护;
    c> 数据链路层对数据帧的长度有一个限制,也就是链路层所能承受的最大数据长度,这个值称为最大传输单元(Max Transmission Unit,MTU),
       数据链路层协议中的以太网Ethernet协议和IEEE802.3协议,都定义了MTU,分别是1500和1492字节,需要注意的是,MTU定义的是MAC帧中数据区
       的长度,并不包括MAC帧的首尾长度(共18B),也就是说MTU限制的是IP数据报的长度,如果IP层的数据报长度大于MTU,则需要对数据报分片;
       由于IP报文头部占用20字节,所以IP数据报的数据区的最大长度是1480B,而tcp的头部占用20B,所以传输tcp报文段数据区的最大长度是1460B,
       UDP头部占用8B,所以UDB报文数据区的最大长度为1472B; 不过从协议结构上看,IP报文头部中有2个字节表明报文总长度,所以IP报文的最大长度
       是65535B;
       链路层的MAC帧定义中,14字节的头部,4字节的尾部; IP协议头20字节; TCP协议头20字节; udp协议头8字节
15. Linux环境网络IO的同步,异步,阻塞,非阻塞,复用
    对于一个网络IO,比如读取read来说,它涉及两个对象(调用IO的用户process/thread,内核kernel),两个阶段(等待数据准备,数据从内核copy到进程)
    blocking IO:用户进程调用recvfrom系统调用,在kernel准备数据及copy数据阶段,用户进程都是阻塞的,直到kernel完成copy,返回ok,用户进程解除block状态;
    non-blocking IO:用户进程调用recvfrom,kernel未准备好数据时,并不会阻塞用户进程,而是直接返回一个error,用户进程接收到error,得知数据还没准备好
    于是再次发起read操作,调用recvfrom,如此循环,直到kernel准备好数据后再次受到用户进程的系统调用,此时kernel copy数据到用户缓冲区,返回ok,非阻塞IO
    需要在准备数据阶段,不断的发起询问,所谓非阻塞指的就是准备数据阶段,每次询问都会立即得到一个结果;
    IO复用(IO multiplexing),或者称之为event driven IO:常见的实现为select/epoll,基本原理是select/epoll方法会不断的轮询所负责的socket,当某个
    socket的数据到达,就通知用户进程,调用流程如下:当用户进程调用select,整个进程block,kernel监视select负责的socket,当其中任意一个socket的数据准备
    好了,select就返回,此时用户进程发起read操作,kernel将数据copy到用户进程,select用到了两个system call(select,recvfrom),select的优势不在于单个
    处理的更快,而是可以同时处理更多connection;
    异步IO:用户进程发起异步读取操作之后,去做其它事,kernel接收到异步读之后,立刻返回,所以并不会阻塞用户进程,kernel等待数据准备,并把数据copy到用户
    缓存之后,kernel给用户进程发送一个signal,告诉它read操作已完成
    POSIX同步与异步的定义:IO操作时进程阻塞为同步,非阻塞为异步
    则阻塞,非阻塞与复用均为同步IO,非阻塞IO虽然在数据未准备好时并未阻塞,但当kernel准备好数据之后,用户进程再次调用recvfrom copy数据时会被阻塞
16. select/poll适用于所有的Unix系统,epoll则是Linux,从根本上说,poll/select这两个系统调用使用的是相同的代码,工作方式基本相同
    select定义:https://github.com/torvalds/linux/blob/v4.10/fs/select.c#L634-L656
    do_select定义:https://github.com/torvalds/linux/blob/v4.10/fs/select.c#L404-L542
    poll定义:https://github.com/torvalds/linux/blob/v4.10/fs/select.c#L1005-L1055
    do_poll定义:https://github.com/torvalds/linux/blob/v4.10/fs/select.c#L795-L879
    select和poll的一些区别:
    select()基于位图(bitmap),poll()基于fd数组,因此select()的一个缺陷就是它的大小的固定的(FD_SETSIZE,1024),即便可以通过某些方式绕过
    poll返回的结果类型更多,select只会返回读、写和报错,第二个区别是fd较少时,poll的效率比select高
    至于原因嘛,从源码上就可以看出来,do_select(select.c#L440)便利fd时,从0开始知道找到fd(fd的本质是个索引值),而do_poll(select.c#L818)则是
    遍历fd数组,如果当前只有4个fd,则poll只需遍历4次,而select需要从0遍历到max_select_fd
17. epoll调用包括epoll_create,epoll_ctl,epoll_wait,epoll_create开启epolling,内核返回一个ID,epoll_ctl告诉内核要监听的fd,调整fd set               (declare_interest()),epoll可以同时监听多种不同类型的fd,包括但不限于pipes,FIFOs,sockets,POSIX message queue,devices等,
    epoll_waite等待内核返回可用的fd,也就是获取事件(get_next_event()),
    select/poll都是无状态的,所以每次调用的时候都会提供整个fd set,优化就是内核维护状态相关的fd set,避免每次都返回整个fd set,
    Linux和BSD各自的实现是epoll/kqueue
    kqueue()函数类似于epoll_create(),kevent()则集成了epoll_ctl()和epoll_wait()
    从性能上说epoll的一个缺陷是,单次系统调用无法更新多个fd状态,比如有100个fd需要更新状态,那么epoll不得不调用100次epoll_ctl(),那么过度调用系统调用
    将会导致系统性能降级,相反在一次kevent()调用中,可以指定对多个fd进行状态更新
    epoll的另外一个限制就是,它基于文件描述符工作,但是时钟,信号,信号量,进程,(linux中的)网络设备等不是文件,无法对这些非文件类型使用基于select/poll/
    epoll的事件复用技术,Linux提供了许多补充性质的系统调用,比如signalfd(),eventfd(),timerfd_create()来转换非文件类型的文件描述符,然后就可以使用
    epoll,只是不是那么优雅,而kqueue中kevent结构体支持多种非文件事件,例如,程序可以获得一个子进程退出事件,通过设置filter=EVFILT_PROC,ident=pid
    fflags=NOTE_EXIT,发现一篇很清晰的对比epoll_vs_kqueue的文章:http://people.eecs.berkeley.edu/~sangjin/2012/12/21/epoll-vs-kqueue.html
18. C10K问题就是Client 10000,说的是同时连接到服务器的客户端数量超过10000,即使硬件性能足够,依然无法正常提供服务,简而言之就是单机1w个并发连接问题
    C10K问题受到创建进程数,内存空间等的限制,即便我们使用64位创建进程,提高进程创建的上限,使用虚拟内存,扩大内存的使用空间,然而问题依然存在,进程和
    线程的创建都需要消耗一定的内存,每生成一个栈空间,都会产生内存开销,当使用内存超过物理内存时,一部分数据会被持久化到磁盘上,导致性能的大幅下降
    C10K问题的突破是单个线程或进程管理多个客户端请求,通过异步编程和事件触发机制,IO的非阻塞,IO多路复用等来提高性能,底层解决方案是epoll,kqueue,
    libevent等,应用层面的解决方案有OpenRestry,Golang,Node.js
    参考来源:https://medium.com/@chijianqiang/%E7%A8%8B%E5%BA%8F%E5%91%98%E6%80%8E%E4%B9%88%E4%BC%9A%E4%B8%8D%E7%9F%A5%E9%81%93-c10k-%E9%97%AE%E9%A2%98%E5%91%A2-d024cb7880f3
    https://my.oschina.net/xianggao/blog/664275
19. vps上的梯子,一旦请求连接就会出现十多个处于SYN_RECV状态的连接,暂时无解,不过刚好是个机会了解下tcp相关参数的配置
    tcp ipv4的参数位置: /proc/sys/net/ipv4/
    针对SYN_RECV有三个相关的参数:
    tcp_syn_retries:integer 默认为5,对于新建连接,内核要发送多少个SYN请求,才会决定放弃,对通信良好的网络可调整为2
    tcp_synack_retries:integer 默认为5,对于客户端的连接请求SYN,服务端内核会发送SYN+ACK数据报,该值决定了内核放弃连接之前所发送的SYN+ACK报次数
    tcp_syncookies:integer 默认为1,表示开启syn cookie功能,tcp_syncookies可有效防范SYN Flood攻击,原理是在收到客户端的SYN,并返回SYN+ACK包时
    不分配一个专门的数据区,而是根据SYN包计算一个cookie值,收到客户端ACK包时,由cookie值检查该ACK包是否合法,如果合法再分配专门的数据区处理tcp连接
20. 客户端反映两次请求的文件长度不一致,查了一下发现由于是基于轮询的负载,两次请求返回不同的登录服地址,登录服httpd服务对应的文件的大小不一致
    最后将连接设置为keep-alive,保证同一连接,多次请求返回同一登录服地址
    不过,由此倒引出了http响应中的content-length字段,在content-encoding,gzip,chunked等不同情境下,content-length字段大小不一致
    content-length 是http消息实体的传输长度,注意区分消息实体长度和消息实体传输长度,在服务器开gzip的情况下,消息实体长度是压缩前的长度,消息实体
    传输长度是压缩后的长度,在实际的交互中,客户端获取消息长度的规则:
    如果content-length存在且有效的情况下,则必须和消息实体的长度一致;
    如果存在transfer-encoding,则在header中不能有content-length字段,如果有也会忽略
    如果采用短连接,可以直接通过关闭服务器连接来确定消息的传输长度
    参考:https://segmentfault.com/a/1190000006194778
21. tcp的粘包问题,首先看下所谓的粘包是指的什么问题,粘包是指发送方发送的若干包数据到达接收方时粘成了一个包,从接收缓冲区看,后一个包的数据的头紧挨着
    前一个包数据的尾,出现粘包的直接原因有可能来自发送方,也有可能来自接收方
    a> 发送方原因,tcp默认采用Nagle算法,而nagle算法主要做两件事,其一,只有上一个数据包得到确认,才会发送下一个包;其二,收集多个小包,在收到一个确认
       时,一起发送多个小包组成的大包,此时就会出现粘包,可见Nagle算法的主要作用是减少网络传输中的报文段数量
    b> 接收方原因,tcp接收到数据包之后,把数据包保存在缓存中,然后由应用程序读取缓存中的包,此时如果tcp接收包到缓存中的速度大于应用程序的读取速度,那么
       多个包就会被缓存,应用程序读取的时候就会读到多个首尾相连粘在一起的包
    tcp协议传输的对象是字节流,也就是一串连续的字节,上层应用处理的是一个一个逻辑上的消息或者请求,但是tcp传输的字节流中并没有消息的边界,同时tcp报文
    的头部也没有表示数据长度的字段,所以上述两个原因是tcp传输过程中有可能出现粘包问题的根源
    所以从本质上看,tcp和粘包没有关系,tcp保证的是字节流传输的正确性,至于如何用这些数据那是上层应用的事,所以解决问题的关键还是在于上层应用如何从字节
    流中提取逻辑上的一个个消息,常见的解决粘包问题的方法如下:
    a> 发送方设置tcp_nodelay选项来关闭Nagle算法,只能解决发送方的问题
    b> 应用层解决,包括以下几种方法:格式化数据,每条数据以固定的格式开始&结尾(开始符&结束符),但是要保证发送的数据中不能包含开始&结束符,这种方法有一
       定的局限性; 发送长度,每一条数据都包含数据长度,比如,数据的前4个字节是本数据包的实际长度(目前本项目用的就是这种方法),这样可以通过读取指定长度
       有效的区分各个消息
    UDP就不存在粘包问题,udp基于报文传输,一次传输一个报文,接收方一次只接收一个独立的报文,同时udp头部有2byte用来表示数据报文的长度
22. tcp的流量控制和拥塞控制
    流量控制是控制发送方的发送速度,利用滑动窗口机制得以实现,发送方的发送窗口不会大于接收方给的接收窗口,流量控制是点到点的控制
    拥塞控制是防止过多的数据被注入到网络中,避免网络中的路由或链路出现过载,拥塞问题是一个全局问题,涉及到主机,路由器以及所有降低网络传输性能的其它因素
    拥塞控制的4种算法:慢启动(slow start),拥塞避免(congestion avoidance),快重传(fast restrangmit)和快恢复(fast recovery)
    发送方会维护一个拥塞窗口(cnwd,congestion window)的状态变量,拥塞窗口的大小取决于网络的拥塞程度,并动态的变化
    慢启动算法的思路是:新建立的连接不能一开始就发送大量的数据包,而是根据网络情况逐渐增加每次发送的量,初始时cwnd=1个最大报文段(MSS)大小,每当一个报文
    被确认就翻番(1->2->4->8->...),为了防止cwnd增长过大引起网络拥塞,还需设置一个慢启动门限(ssthresh)状态变量,当cwnd < ssthresh时,使用慢启动算法,
    cwnd > ssthresh时,改用拥塞避免算法
    拥塞避免算法的思路是:让拥塞窗口cwnd缓慢增长,每经过一个RTT,cwnd加1,而不是加倍,这样cwnd按线性规律缓慢增长
    无论是在慢启动阶段还是拥塞避免阶段,只要出现网络拥塞(收不到确认报文),就把慢开始门限ssthresh设置为出现拥塞时发送窗口(send window)大小的一半,然后
    把cwnd设置为1,执行慢开始算法,目的是迅速减少主机发送到网络中的报文数量,使得发生拥塞的路由器有足够时间把队列中挤压的报文处理完毕
    tcp连接有时会因为等待重传超时而空闲较长的时间,慢启动和拥塞避免无法很好的解决这个问题,因此提出了快重传和快恢复算法
    快重传要求发送方在接收到一个失序的报文段时就立即发送重复确认报文,而不是等到自己发送数据时顺带发送确认,快重传算法规定如果发送方连续接收到三个重复
    确认报文,就立即发送接收方尚未收到的报文段,而不继续等待设置的重传计时器到期,比如有7个报文段p1-p7,接收方在收到p1,p2之后收到了p4,此时接收方快重
    传算法发送失序报文p2的第一个重复确认报文,如果后续接收到的报文还不是p2,而是p5,则发送失序报文p2的第二个确认报文,以此类推,当发送第三个重复确认报文
    后,发送方立即重传p2
    与快重传算法配合的还有快恢复算法,指的是如果触发快重传,则此时将ssthresh门限减半,当时接下来并不执行慢启动算法,考虑到如果能连续收到多个重复确认
    报文,所以此时发送方并不认为网络出现拥塞,所以并不执行慢启动算法,而是将cwnd=此时减半后的ssthresh(不把cwnd=1),然后在新的ssthresh和cwnd的基础上
    执行拥塞避免算法,这个过程就是快恢复
    注:拥塞窗口大小以字节为单位,与发送窗口相同
