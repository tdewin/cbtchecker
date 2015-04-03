package main

import (
	"fmt"
	"flag"
	"os"
	"io"
	"crypto/rand"
)
 	

func check(e error) {
    if e != nil {
        panic(e)
    }
}
func randomKB(kbs int64) ([]byte) {
	rkb := make([]byte, (1024*kbs))
	rand.Read(rkb)
	return rkb 
}
func createFile(filename *string,sizekb int64) {
	if _, err := os.Stat(*filename); err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("Creating file %s\n",*filename)
			f, err := os.Create(*filename)
			check(err)
			defer f.Close()

			for rb := int64(0);rb<sizekb;rb++ {
				_, err := f.Write(randomKB(1))
				check(err)
			}
                } else {
			fmt.Printf("Strange error on %s\n",*filename)
		}
	} else {
		fmt.Printf("File already exists %s\n",*filename)
  	}
}
func readFile(filename *string) {
	if _, err := os.Stat(*filename); err != nil {
                if os.IsNotExist(err) {		
			fmt.Printf("File does not exist %s\n",*filename)
		} else {
			fmt.Printf("Strange error %s\n",*filename)
		}
	} else {
		f, err := os.Open(*filename)
		check(err)
		defer f.Close()

		reader := make([]byte, 1024)
		var bcounter = int64(0)
		for bytesread,err := f.Read(reader);err != io.EOF;bytesread,err = f.Read(reader) {
			bcounter += int64(bytesread)
		} 
		fmt.Printf("Read %d bytes\n",bcounter)
	}
}
func writeFile(filename *string,interval int, block int) {
	fstat, err := os.Stat(*filename)
	if err != nil {
                if os.IsNotExist(err) {		
			fmt.Printf("File does not exist %s\n",*filename)
		} else {
			fmt.Printf("Strange error %s\n",*filename)
		}
	} else {
		f, err := os.OpenFile(*filename,os.O_RDWR,os.FileMode(0666))
		check(err)
		defer f.Close()
		
		intervalbytes := int64(interval)*1024
		block64 := int64(block)

		fsize := fstat.Size()
		updates := int64(0)

		for ctr := int64(0); ctr < fsize;ctr += intervalbytes {
			f.WriteAt(randomKB(block64),ctr)
			updates += 1
		}
		fmt.Printf("Did %d\n",updates)
	}
}

func main() () {
	tmpfile := os.ExpandEnv("${TEMP}\\file")
	filename := flag.String("file",tmpfile,"Provide file to work on")
	action := flag.String("action","write","provide 'create' to make the file, or 'write' or 'read'  to do action")
	sizein := flag.Int64("size",1024,"Size in mb, by default 1024 or 1GB. works with create")
	interval := flag.Int("interval",64,"interval size in kb")
	block := flag.Int("block",64,"Block size in kb to touch, then jump to the next interval") 
	flag.Parse()

	sizekb := int64((*sizein)*1024)

	fmt.Printf("Working on : %s\n",*filename)
	
	switch *action {
		case "create": {
			createFile(filename,sizekb)	
		}
		case "read": {
			readFile(filename)
		}
		case "write": {
			writeFile(filename,*interval,*block)
		} 
	}
}
