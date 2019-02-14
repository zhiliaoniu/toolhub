#include <iostream>
//#include <shared_ptr.h>
#include <vector>
#include <memory>
#include <atomic>

using namespace std;
struct A {
    //A() {cout << "construct" << endl;}
    //A(const A& aa) {a = aa.a; cout << "copy construct" << endl;}
    //A& operator = (const A& aa) {a = aa.a; cout << "fuzhi construct" << endl;}
    //~A() {cout << "destruct" << endl;}
    int a;
    string b;
    //vector<shared_ptr<string>> vs;
};

atomic<shared_ptr<A> > a1;
//atomic<string*> c1;
//shared_ptr<A> a3;

void func() {
    /*/
    shared_ptr<A> a2 = make_shared<A>();
    a2->a = 4;
    a2->b = "mmm";

    //vector<shared_ptr<string>>& vs = a2->vs;
    //vs.push_back(make_shared<string>("abc"));

    //a3 = a2;
    //atomic_store(&a1, a2);
    a1.store(a2);
    //cout << *(a2->vs[0])<<endl;
    cout << (a2->b)<<endl;
    */

    //A c2;
    //c2.a = 9;
    //c1.store(c2);
}

int main() 
{
    func();
    //cout << a3->a << endl;
    //cout << a3->vs.size() << endl;
    //cout << *(a3->vs[0]) << endl;
    //cout << "--------------" << endl;

    /*
    shared_ptr<A> bb;
    //bb = atomic_load(&a1);
    bb = a1.load();
    cout << bb->a << endl;
    cout << bb->b << endl;
    //cout << bb->vs.size() << endl;
    //cout << *(bb->vs[0]) << endl;
    */
    //A c3 = c1.load();
    //cout << c3.a << endl;
   // cout << c3.b << endl;
}
