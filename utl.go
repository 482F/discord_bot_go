package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

func sleep(ms time.Duration) {
	time.Sleep(time.Millisecond * ms)
	return
}

func checkErr(e error) bool {
	if e != nil {
		fmt.Println(os.Stderr, e)
		os.Exit(1)
	}
	return true
}

func readTxtFile(filepath string) (string, error) {
	bt, err := ioutil.ReadFile(filepath)
	var str string = string(bt)
	str = str[:len(str)-1] //改行が入ってしまうので取り除く
	return str, err
}

func arrayHasString(ss []string, target string) bool {
	for _, s := range ss {
		if s == target {
			return true
		}
	}
	return false
}
