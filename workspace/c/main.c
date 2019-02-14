#include <stdio.h>
#include <algorithm> 
#include <iostream>
#include <vector>
#include <map>

#include "main.h"

//void fun(std::string&& s) {
//    std::cout << "s:" << s << std::endl;
//    std::string s2 = std::move(s);
//    std::cout << "s2:" << s2 << std::endl;
//    return;
//}
//

void fun(int a[]) {

    std::cout << "len:" << sizeof(a) << std::endl;
}

int main () {

    //printf("Hello World");
    //std::string bar = "bar-string";
    //std::vector<std::string> myvector;

    ////myvector.push_back (std::move(bar)); 
    ////std::cout << "vec[0]: " << myvector[0] << std::endl;
    //fun(std::move(bar));
    //std::cout << "bar: " << bar << std::endl;

    //typedef std::map<int, std::shared_ptr<std::string> > ISMap;
    //ISMap m;
    //m[1]=std::shared_ptr<string>(new std::string("abc"));



    int a[6];
    std::cout << ABSL_ARRAYSIZE(a) << std::endl;

    std::cout << sizeof(a) << std::endl;
    int (&b)[6] = a;

    fun(b);

    return 0;

}
