#include <iostream>
#include <map> 
#include <string>

using namespace std;

class Test1 {

public:
    Test1():a(1), b(0){}

    void print() { 
        cout << a << " " << b << endl;
        cout << m_double << endl;
    }

public:
    const int a;
    const int b;
    
    static constexpr double m_double = 0.00000001f;    
};

void fun() {

    Test1 t1;
    t1.print();
    
    const char* c1 = "liyang";
    unsigned char* c2 = (unsigned char*)const_cast<char*>(c1);
    cout << c2 << endl;
}
#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h>
int main () { 
    std::map<std::string, std::string> m;
    m["1"] = "11";
    m["2"] = "22";
    m["3"] = "33";
    for(auto kvp : m) {
        std::cout << kvp.first << "   " << kvp.second << std::endl;
        for(auto v:kvp.second) {
            std::cout << v << "," << std::endl;
        }
    }
    struct in_addr ia;
    ia.s_addr = htonl(3758092800);
    std::cout << inet_ntoa(ia) << std::endl;


    return 0;
}

