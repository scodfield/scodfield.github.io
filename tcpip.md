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
