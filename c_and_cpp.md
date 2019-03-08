1. c++中的虚函数用来实现多态,多态简单来说就是用父类型指针指向子类对象,进而调用实际的子类成员函数,虚函数的动态绑定是一种泛型技术(试图以
   不变的代码实现不同的算法),每个包含虚函数的类都有一个虚函数表,表中主要是该类的虚函数的地址,也即虚函数表是一个函数指针数组
2. GCC和G++的区别,手册:http://gcc.gnu.org/onlinedocs/gcc/G_002b_002b-and-GCC.html#G_002b_002b-and-GCC
   GCC最开始是GNU C Compiler,是开源免费的C语言编译器,后来又加入了对C++,Pascal,Objective-C等其它语言的支持,名字就变成了 GNU Compiler Collection
   G++是GCC编译器集合中,专门用来处理C++语言的
   GCC提供了C语言预处理器:C Preprocessor 简称CPP
   传统编译器的工作原理基本上是三段式的,分为前端(Front end),优化器(Optimizer),后端(Back end),前端负责解析源代码,语义检查,生成抽象语法树(Abstract
   Sytax Tree, AST),优化器对中间代码进行优化,后端负责生成机器代码,这一过程后端会最大化利用目标机器的特殊指令,以提高代码的性能
   GCC实现了很多前端,支持多种语言,但它是一个完整的可执行文件,没有给其它语言的开发提供重用中间代码的接口,也就是说GCC太重,非模块化
   Low Level Virtual Machine(LLVM) 是一个定位较底层的虚拟机,它的有点正是GCC的缺点,用来解决编译器代码重用的问题,LLVM与其它编译器最大的区别是
   它不仅仅是Compiler Collection,也是Libraries Collection,也就是说LLVM是一个编译器,也是一个SDK
   Clang是一个支持C,C++,Objective-C和Objective-C++等编程语言的编译器前端,Clang采用LLVM作为其后端,Clang主要有C++编写,Clang相比GCC
   编译速度快,内存占用小,基于库的模块化设计,易于IDE集成,出错提示更友好,而GCC的优势在于,除了支持C系列,还支持Java,Fortran,Go等语言,其次是使用更广泛
   支持更多平台
