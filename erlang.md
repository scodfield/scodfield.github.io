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
