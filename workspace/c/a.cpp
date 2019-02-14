#include <iostream>
#include <vector>
#include <unordered_set>
#include <algorithm>

//int f(char* a) {
//    std::cout << a << std::endl;
//    int i = 1;
//    return 3==2?3:4;
//}
enum e {
    LL = 1<<30
};

std::unordered_set<std::string> audio_exts{"wav", "mp3", "m4a", "ogg", "amr", "aac"};

int main() {
    //char a[10];
    //int i = 1000000000;
    //while (i--) {
    //    f(a);
    //}
    //std::cout << f(a) << std::endl;
    //int b = -1620649632;
    //unsigned int c = b;
    //std::cout << c << std::endl;
    //
    std::string ext_name = "WAV";
    std::cout << "ext_name: " << ext_name << std::endl;
    std::transform(ext_name.begin(), ext_name.end(), ext_name.begin(), ::tolower);
    std::cout << "ext_name: " << ext_name << std::endl;

    std::vector<int> v = {1, 3, 5, 7, 9, 8,6, 4, 2, 0};
    for (auto& i : v) {
        std::cout << i << std::endl;
        i += 1;
    }
    std::cout << "-----------" << std::endl;
    for (auto& i : v) {
        std::cout << i << std::endl;
        i += 1;
    }

    int r = int(LL) & int(1);
    std::cout << r << std::endl;
    return 0;
}
