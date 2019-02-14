#include <iostream>
#include <thread>
#include <algorithm>
#include <ctime>

using namespace std;

void fun(int num) 
{
    for(int i = 0; i < 2; ++i) {
        cout << i << endl;
    }
    cout << "num:" << num << endl;
}

void fun2(){
    int a[10]={0};
    srand(time(NULL));

    generate(a,a+10,[]()->int { return rand() % 100; });

    cout<<"before sort: "<<endl;
    for_each(a, a+10, [&](int i){ cout<< i <<" "; });
    cout<<endl;

    cout<<"After sort"<<endl;
    sort(a,a+10);
    for_each(a, a+10, [&](int i){ cout<< i <<" "; });
    cout<<endl;

    return;
}

int main () {
    int n = 10, n2 = 5;
    std::thread t{[&](int num)->int{
        cout << "in thread num:"<< num << endl;
        n = 11;
        n2 = 11;
        return num;
    }, 666
    };

    t.join();

    cout << "n:" << n << "  n2:" << n2 << endl;

    //fun2();
    return 0;
}


