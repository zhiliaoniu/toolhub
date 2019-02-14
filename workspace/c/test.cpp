#include <stdio.h>
#include <iostream>
#include <time.h>
#include <stdlib.h>
#include <string>
#include <string.h>
#include <set>
#include <vector>
#include <map>
using namespace std;

void fun1() {

    cout << "sizeof(int32_t):" << sizeof(int32_t) << endl;
    cout << "sizeof(time_t):" << sizeof(time_t) << endl;

    time_t t = atoll("1423223433");
    tm tmp;
    tm* t2 = localtime_r((time_t*)&t, &tmp);
    cout << "t:" << t << endl;
    cout << "t2:" << t2 << endl;

    char t3[50];
    strftime(t3, sizeof(t3), "%Y-%m-%d %H:%M:%S", t2);

    cout << strlen(t3) << endl;
    cout << "t3:" << t3 << endl;


    cout << "----------char test---------" << endl;
    
    char c = char(0xA1);
    char buf[5] = {0};
    printf("first: %02hhx\n", c);

    if(c == 0xA1) cout<< "c == 0xA1" << endl;
    else cout << "not equal " << endl;
    cout << 0xA1 << endl;
    cout << char(161) << endl;
    cout << "0xA1:" << c << endl;
}

bool fun2 () {
    std::set<std::string> ss;
    ss.insert("abc.b.c.a");

    cout << *(ss.find("abc.b.c.a")) << endl;

    return true;
}

class A {
    public:
        A() {
            cout << "In A " << endl;
        }
        void funa(){ cout << "In A a:" << this << endl;}
        
};

class B {
    public:
        B(A* aa):a(aa){
            cout << "In B a:" << a << endl;
        }
        B() {//a = NULL;
            s = "a";
        }
        void funb(){ 
            a->funa();
        }
        A* a;
        std::string s;
        int i;
        double d;
};   

int fun2(int& lhs) {
    int tmp = 1;
    lhs += tmp;
    return lhs;
}

int main (int argc, char* argv[]) {

    //fun2();
    //A a;

    //B b(&a);
    //b.funb();


    size_t size = 0;
    for (int i = 0; environ[i]; i++) {
        size += strlen(environ[i]) + 1;
        std::cout << i << ":" << environ[i] << std::endl;
    }

    strncpy(argv[0], "fuckfuckfuck", 12);
    argv[0][12]='\0';

    int count = 0;
    for(int i = 0; i < 1000; ++i) {
        int result = fun2(count);
        sleep(1);
        std::cout << result << std::endl;
    }

    //static const std::string s = "yangshengzhi";
    //std::map<std::string, std::string> ssmap;
    //for (int i = 0; i < 10000000; ++i) {
    //    ssmap["yangshengzhi"] = "yangshengzhi";
    //}

    return 0;
}
