package main

import (
	"fmt"
	"flag"
	"os"
	"io"
	rand "crypto/rand"
	mrand "math/rand"
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
func letterKB(kbs int64,letter int64) ([]byte) {
	rkb := make([]byte, (1024*kbs))
	for index, _ := range rkb {
		rkb[index] = byte(65+int(letter%26))
	}
	return rkb
}
func createDedupFile(filename *string, sizekb int64, block int64) {
        if _, err := os.Stat(*filename); err != nil {
                if os.IsNotExist(err) {
                        fmt.Printf("Creating file %s\n",*filename)
                        f, err := os.Create(*filename)
                        check(err)
                        defer f.Close()
			
			ctr := int64(0)

                        for rb := int64(0);rb<sizekb;rb += block {
                                _, err := f.Write(letterKB(block,ctr))
                                check(err)
				ctr = (ctr+1)%26
                        }
                } else {
                        fmt.Printf("Strange error on %s\n",*filename)
                }
        } else {
                fmt.Printf("File already exists %s\n",*filename)
        }

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

func createmap(banks int64) ([]bool) {
	createmap := make([]bool,banks) 
	for ctr := int64(0);ctr < banks;ctr++ {
		createmap[ctr] = false
	}
	return createmap
}
func randomMoveFile(filename *string, block int) {
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

                blockkb := block*1024
                block64kb := int64(blockkb)

                fsize := fstat.Size()
                moves := int64(0)

		
		//move on every 128MB, should make memory limitations lower and file jumping shorter
		subblocking := int64(1024*1024*128)
                ctr := int64(0)

                for ; ctr < fsize;ctr += subblocking {
			rangeStart := ctr
			rangeEnd := ctr + subblocking
			if (rangeEnd > fsize) {
				rangeEnd = fsize
			}

			subblocks := (rangeEnd-rangeStart)/block64kb
			markmap := createmap(subblocks)
			fmt.Printf("\n%d - %d   --> %d\n",rangeStart,rangeEnd,subblocks)

			reader := make([]byte, blockkb)
			firstblock := make([]byte, blockkb)

			prevblock := mrand.Int63n(subblocks)
			f.ReadAt(firstblock,(rangeStart+(prevblock*block64kb)))
			markmap[prevblock] = true			
			
			for x := int64(1);x < subblocks;x++ {
				newblock := mrand.Int63n(subblocks)
				safety := int64(0)
				for ;markmap[newblock];newblock = (newblock+1)%subblocks {
					if safety > subblocks {
						panic("AAAAAAAAh running around in circles")
					}
					safety++
				}
				markmap[newblock] = true

				blockin := (rangeStart+(newblock*block64kb))
				blockout := (rangeStart+(prevblock*block64kb))
					
				f.ReadAt(reader,blockin)
				f.WriteAt(reader,blockout)

				prevblock = newblock

//				fmt.Printf("( %d > %d ) ----  ",blockin,blockout)
				moves += 1
			} 
			f.WriteAt(firstblock,(rangeStart+(prevblock*block64kb)))
			moves += 1

                        
                }

                fmt.Printf("Did %d moves\n",moves)
        }
}


func moveFile(filename *string, block int) {
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

		blockkb := block*1024
                block64kb := int64(blockkb)

                fsize := fstat.Size()
                moves := int64(0)

		reader := make([]byte, blockkb)
		firstblock := make([]byte, blockkb)
		almostlast := (fsize-block64kb)


		ctr := int64(0)
		f.ReadAt(firstblock,0)
	        for ; ctr < almostlast;ctr += block64kb {
			f.ReadAt(reader,ctr+block64kb)
                        f.WriteAt(reader,ctr)
                        moves += 1
                }
		f.WriteAt(firstblock,ctr)
		moves++

                fmt.Printf("Did %d moves\n",moves)
        }
}


func main() () {
	filename := flag.String("file","/tmp/file","Provide file to work on")
	action := flag.String("action","write","'create' | 'creatededup' | 'write' | 'read' | move allowed")
	sizein := flag.Int64("size",1024,"Size in mb, by default 1024 or 1GB. works with create")
	interval := flag.Int("interval",64,"interval size in kb")
	block := flag.Int("block",64,"Block size in kb to touch, then jump to the next interval") 
	flag.Parse()

	sizekb := int64((*sizein)*1024)

	
	switch *action {
		case "create": {
			createFile(filename,sizekb)	
		}
		case "creatededup": {
			 createDedupFile(filename,sizekb,int64(*block))
		}
		case "read": {
			readFile(filename)
		}
		case "write": {
			writeFile(filename,*interval,*block)
		} 
		case "move": {
			moveFile(filename,*block)
		}
                case "randommove": {
                        randomMoveFile(filename,*block)
                }
 
	}
}
