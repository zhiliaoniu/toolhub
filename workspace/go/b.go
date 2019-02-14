package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	_ "runtime/pprof"
	_ "strconv"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	go server()
	go printNum()
	var i = 1
	for {
		// will block here, and never go out
		i++
		time.Sleep(time.Second * 1)
		//fmt.Println("i:" + strconv.Itoa(i))
	}
	fmt.Println("for loop end")
	time.Sleep(time.Second * 1)
}

func printNum() {
	i := 0
	for {
		fmt.Println(i)
		i++
		time.Sleep(time.Second * 1)
	}
}

func HelloServer(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "hello, world!\n")
}

func server() {
	Mux := http.NewServeMux()
	Mux.HandleFunc("/abc", HelloServer)
	Mux.HandleFunc("/abc1j", HelloServer)
	Mux.HandleFunc("/abc2", HelloServer)
	Mux.HandleFunc("/abcadfaf", HelloServer)
	Mux.HandleFunc("/abcwewe", HelloServer)
	Mux.HandleFunc("/abasdfac", HelloServer)
	fmt.Printf("Mux:%v\n", Mux)
	//srv := &http.Server{Handler: Mux}
	//http.HandleFunc("/", HelloServer)
	//err := http.ListenAndServe(":12345", nil)
	http.HandleFunc("/abc", HelloServer)
	http.HandleFunc("/abc1j", HelloServer)
	http.HandleFunc("/abc2", HelloServer)
	http.HandleFunc("/abcadfaf", HelloServer)
	http.HandleFunc("/abcwewe", HelloServer)
	http.HandleFunc("/abasdfac", HelloServer)
	fmt.Printf("DefaultServeMux:%v\n", http.DefaultServeMux)
	srv := &http.Server{Handler: http.DefaultServeMux}
	listener, err := net.Listen("tcp", ":12345")
	go func() {
		err = srv.Serve(listener)
		if err != nil {
			fmt.Fprintf(os.Stderr, "server exit. err:%v", err)
		}
	}()
	//err = srv.Serve(listener)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
