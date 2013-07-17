package main

import (
  "path/filepath"
	"os"
	"archive/tar"
	"flag"
	"fmt"
	"io"
	"strconv"
	"bufio"
	"strings"
	"compress/gzip"
)

func visit(path string, f os.FileInfo, err error) error {
//	_, filename:= filepath.Split(path)
	all_data := 0
	suffix := filepath.Ext(path)
	if suffix == ".tar" {
		f, err := os.Open(path)
		if err != nil {
			panic(nil)
		}
		defer f.Close()
		tr := tar.NewReader(f)
		for {
			h, err := tr.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				panic(err)
			}
			fmt.Println(h.Name)

			fw, err := os.OpenFile("C:/output/" + h.Name, os.O_CREATE | os.O_WRONLY, 777)
			if err != nil {
				panic(err)
			}
			defer fw.Close()
  
			_, err = io.Copy(fw, tr)
			if err != nil {
				panic(err)
			}

			err, flow := ungzfile("C:/output/" + h.Name)
			all_data = all_data + flow
			fmt.Println(flow)
		}
	}
	return nil
	
} 


func ungzfile(path string) (err error, flow int) {
	gz_data := 0
	fr, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fr.Close()
  
    // gzip read
	gr, err := gzip.NewReader(fr)
	if err != nil {
		panic(err)
	}
	defer gr.Close()
  
	// tar read
	tr := tar.NewReader(gr)
  
	for {
		h, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
  
		fmt.Println(h.Name)
		fw, err := os.OpenFile("C:/output/log/" + h.Name, os.O_CREATE | os.O_WRONLY, 777)
		if err != nil {
			panic(err)
		}
		defer fw.Close()
  
		_, err = io.Copy(fw, tr)
		if err != nil {
			panic(err)
		}
		err, num := deallog("C:/output/log/" + h.Name)
		gz_data = gz_data + num
	}
	return nil, gz_data
}


func deallog(path string) (err error, num int) {
	data := 0
	fr, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fr.Close()
  	r := bufio.NewReader(fr)
	line, isPrefix, err := r.ReadLine()
	for err == nil && !isPrefix {
		s := string(line)
		words := strings.Split(s, " ")
	//	fmt.Println(words[9], len(words))
		i, _ := strconv.Atoi(words[9])
		data = data + i
		line, isPrefix, err = r.ReadLine()
	}
	return nil, data
}



func main() {
	flag.Parse()
	root := flag.Arg(0)
	err := filepath.Walk(root, visit)
	fmt.Printf("filepath.Walk() returned %v\n", err)
}
