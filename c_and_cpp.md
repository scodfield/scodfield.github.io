1. c++中的虚函数用来实现多态,多态简单来说就是用父类型指针指向子类对象,进而调用实际的子类成员函数,虚函数的动态绑定是一种泛型技术(试图以
   不变的代码实现不同的算法),每个包含虚函数的类都有一个虚函数表,表中主要是该类的虚函数的地址,也即虚函数表是一个函数指针数组
2. GCC和G++的区别,手册:http://gcc.gnu.org/onlinedocs/gcc/G_002b_002b-and-GCC.html#G_002b_002b-and-GCC
   GCC最开始是GNU C Compiler,是开源免费的C语言编译器,后来又加入了对C++,Pascal,Objective-C等其它语言的支持,名字就变成了 GNU Compiler Collection
   g++是GCC编译器集合中的C++编译器, gcc则是GCC中的C编译器 
   GCC提供了C语言预处理器:C Preprocessor 简称CPP
   典型的编译过程:调用预处理器,比如CPP;调用实际的编译器,比如cc,cc1;调用汇编器(assembler),比如as;调用链接器,比如ld
   编译器是可以更换的,所以gcc调用的是C Compiler,而g++调用的是C++ Compiler
   gcc与g++的区别包括:对于.c/.cpp文件,gcc分别当做C和cpp文件编译(c/c++语法强度不一样,c++更严谨一些),而g++则统一当做cpp文件处理;
   编译时g++会调用gcc,但是gcc不能自动连接c++的标准库文件(STL),所以gcc编译c++文件时,需要加上-lstdc++参数,而g++可以自动连接库文件
   传统编译器的工作原理基本上是三段式的,分为前端(Front end),优化器(Optimizer),后端(Back end),前端负责解析源代码,语义检查,生成抽象语法树(Abstract
   Sytax Tree, AST),优化器对中间代码进行优化,后端负责生成机器代码,这一过程后端会最大化利用目标机器的特殊指令,以提高代码的性能
   GCC实现了很多前端,支持多种语言,但它是一个完整的可执行文件,没有给其它语言的开发提供重用中间代码的接口,也就是说GCC太重,非模块化
   Low Level Virtual Machine(LLVM) 是一个定位较底层的虚拟机,它的有点正是GCC的缺点,用来解决编译器代码重用的问题,LLVM与其它编译器最大的区别是
   它不仅仅是Compiler Collection,也是Libraries Collection,也就是说LLVM是一个编译器,也是一个SDK
   Clang是一个支持C,C++,Objective-C和Objective-C++等编程语言的编译器前端,Clang采用LLVM作为其后端,Clang主要有C++编写,Clang相比GCC
   编译速度快,内存占用小,基于库的模块化设计,易于IDE集成,出错提示更友好,而GCC的优势在于,除了支持C系列,还支持Java,Fortran,Go等语言,其次是使用更广泛
   支持更多平台
3.  gcc常用编译参数:
    -c 只激活预处理,编译和汇编,生成.o的目标文件
    -S 只激活预处理和编译,生成.s的汇编代码
    -E 只激活预处理,不会生成文件,不过可以通过重定向输出到另一个文件,如: gcc -E fight.c > fight_pre.txt 
    -o 指定输出目标,默认为.out
    -L 指定编译时路径,参考第60条
    -shared 生成共享目标文件(动态库)
    -share 尽量使用动态库
    -static 禁止使用动态库
    -g 指示编译器在编译的时候,产生调试信息
    -OX X = 0,1,2,3 编译器优化选项的4个级别,0 没有优化,1 为缺省, 3 优化级别最高
    -M 生成文件的关联信息,包含目标文件所有依赖
    -MM 同-M, 但它会忽略由#include<file>造成的依赖
    -MD 同-M, 生成的关联信息将会输出到.d文件里
    -MMD 同-MM, 输出到.d文件
    编译执行的4个步骤:
    预处理 预处理器cpp; 
    编译 将预处理后的文件转换成汇编代码,生成.s文件 编译器egcs; 
    汇编 将汇编文件转换为目标代码(机器码),生成.o文件 汇编器as;
    链接 连接目标代码,生成可执行程序 链接器ld
