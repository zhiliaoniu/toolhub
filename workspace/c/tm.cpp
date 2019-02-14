#include <iostream>
#include <fstream>

using namespace std;

void swap(int& a, int& b){
    int tmp = a;
    a = b;
    b = tmp;
}

void f2(int* arr, int b, int e){
    if(b >= e){
        return;
    }
    int f = b, l = e;
    int mid = arr[b];
    while(f < l){
        while(f < l && arr[l] >= mid){
            --l;
        }
        arr[f] = arr[l];
        while(f < l && arr[f] <= mid){
            ++f;
        }
        arr[l] = arr[f];
    }
    arr[f] = mid;
    f2(arr, b, f-1);
    f2(arr, f+1, e);
}

int main () {
    string a;
    fstream f("./c");
    int num = 0;
    while(getline(f, a)) {
        int n = 0;
        for (int i = 0; i < a.length(); ++i) {
            if(a[i] >= 97 && a[i] <= 122) {a[i] = a[i] - 97 + 10;}
            if(a[i] >= 48 && a[i] <= 57) {a[i] = a[i] - 48;}
            n += a[i];
        }
        if(n%2==0){
            cout << "n:" << n << endl;
            num++;
        }
    }
    cout << " num:" << num << endl;

    return 0;
}



