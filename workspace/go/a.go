package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os/exec"
	"strconv"
	"strings"
)

func fun() {
	var arr_t [2]int = [2]int{1, 2}
	var slice_t []int = arr_t[:]
	fmt.Println("arr_t element:")
	for _, v := range arr_t {
		fmt.Print(v, " ")
	}
	fmt.Println("\nslice_t element:")
	slice_t_2 := append(slice_t, append(slice_t, 10, 11)...)
	for _, v := range slice_t {
		fmt.Print(v, " ")
	}
	fmt.Println()
	for _, v := range slice_t_2 {
		fmt.Print(v, " ")
	}
}

var str2 string = "<mid=e35203d896ba6dfec014d72e384eb658><m2=6f2a4203c7656284fdcb1b18429857d785cb3aaa2b76><v=2><product=360safe><combo=dau><ipartner=h_home><pa=x86><pid=h_home><sysver=5.1.2600><version=11.0.0.2001><log=http://s.360.cn/safe/traynew.html?mid:e35203d896ba6dfec014d72e384eb658&m2:6f2a4203c7656284fdcb1b18429857d785cb3aaa2b76&partner:h_home&ipartner:h_home&instdays:1052&ver:11.0.0.2001&sysver:5.1(Build2600)&syssp:3&sever:8.1.1.254&iever:8.0.6001.18702&platform:0&processors:2&computertype:2&mfg:LENOVO&hddt:12041&hddv:298&sysdisk:97&ram:2042&ovinst:1&ty:1><ip=60.221.226.161><time=20170510 04:55:18.045589><s=kill28.bjcc><np=TCP>"

func fun3() {
	var arr []string = make([]string, 0)
	arr = []string{"m2", "mid"}

	str := str2
	for _, v := range arr {
		if index := strings.Index(str, "<"+v+"="); index != -1 {
			fmt.Println("index:", index)
			str_before := str[:index]
			str = str[index:]

			index2 := strings.Index(str, ">")
			fmt.Println("index2:", index2)
			str = str_before + str[index2+1:]
		}
	}
	fmt.Printf("str2:%s\nstr:%s\n", str2, str)
}

func fun4() {
	str := "aa bb"
	arr := bytes.Split([]byte(str), []byte(" "))
	fmt.Printf("var:%v", string(arr[0]))
}

func getMemLoad() {
	out, err := exec.Command("free").Output()
	if err != nil {
		fmt.Printf("err: %s\n", err)
		return
	}
	//fmt.Printf("free out: %s\n", string(out))

	cmd := exec.Command("grep", "Mem")
	cmd.Stdin = strings.NewReader(string(out))
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err = cmd.Run()
	//fmt.Printf("grep out: %s\n", stdout.String())

	arr := strings.Fields(stdout.String())
	//fmt.Printf("arr:%v\n", arr)
	free, _ := strconv.Atoi(arr[3])
	buffer, _ := strconv.Atoi(arr[5])
	cache, _ := strconv.Atoi(arr[6])
	free = free + buffer + cache
	all, _ := strconv.Atoi(arr[1])
	fmt.Printf("all:%d, free:%d\n", all, free)
}

func fun5() {
	c, err := ioutil.ReadFile("/proc/loadavg")
	if err != nil {
		fmt.Printf("err:%v", err)
		return
	}
	cpu_load := strings.Fields(string(c))
	fmt.Printf("cpu_load:%v\n", cpu_load)

	getMemLoad()
	a, _ := math.Modf(float64(111763688 / float64(131869972) * 100))
	fmt.Println(a)
}

func main() {
	str := "a\nbb\nccc\nd\neeee\nff\n"
	byteReader := bytes.NewReader([]byte(str))
	var scanner *bufio.Scanner
	//var gzReader *gzip.Reader

	again := false
	line_num := 1
	buffer := make([]byte, 1, 1)
	line_max_size := len(buffer)
Loop:
	//if isGz {
	//	if !again {
	//		gzReader, err = gzip.NewReader(byteReader)
	//		if err != nil {
	//			log.Printf(p.sub, p.id, "ProxyProcessor.Process gzip.NewReader err: %v", err)
	//			return errors.New("gzip reader create fail")
	//		}
	//		defer gzReader.Close()
	//	}
	//	scanner = bufio.NewScanner(gzReader)
	//} else {
	//	scanner = bufio.NewScanner(byteReader)
	//}
	scanner = bufio.NewScanner(byteReader)

	//scanner.Split()
	if !again {
		scanner.Buffer(buffer, len(buffer))
	} else {
		line_max_size *= 2
		buffer := make([]byte, line_max_size, line_max_size)
		scanner.Buffer(buffer, line_max_size)
	}

	i := 1
	for scanner.Scan() {
		if again && i < line_num {
			i++
			continue
		}

		lineByte := scanner.Bytes()
		lineCount++
		line_num++
		// if len(bytes.TrimSpace(lineByte)) <= 0 {
		//  continue
		// }
		line := string(lineByte)
		log.Printf("line:%s", line)
	}

	if err := scanner.Err(); err != nil {
		log.Printf(p.sub, p.id, "ProxyProcessor.Process scanner.Err err: %v", err)

		if strings.Contains(err.Error(), "token too long") {
			again = true
			goto Loop
		}

		// if err := p.checker.SetFail(hour+p.checkerSuffix, readableId); err != nil {
		// logger.ErrorSIf(p.sub, p.id, "ProxyProcessor.Process SetFail err: %v, id: %s", err, readableId)
		// }
		return err
	}
}
