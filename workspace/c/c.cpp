#include <iostream>
#include <stdlib.h>
#include <sys/time.h>
#include <time.h>
#include <thread>
#include <mutex>
#include <future>
#include <chrono>

using namespace std;

void add_fun(std::promise<int> p) {
    std::this_thread::sleep_for(std::chrono::seconds(1));
    std::cout << "sleep end" << std::endl;
    p.set_value(2);
}

int main () {
    int t11 = time(NULL);
    std::promise<int> p;
    std::future<int> f = p.get_future();

    std::thread t1(add_fun, std::move(p));
    std::cout << "thread after" << std::endl;

    f.wait();
    int ret = f.get();
    std::cout << "ret:" << ret << std::endl;

    t1.join();

    int t12 = time(NULL);
    std::cout << "cost time: " << t12 - t11 << std::endl;
    return 0;
}
