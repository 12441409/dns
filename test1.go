package main
import (
         "fmt"
         "time"
         "os"
         "io"
         "bufio"
         "strings"
       )


func byte16PutString(s string) [4]byte {
    var a [4]byte
    if len(s) > 4 {
        copy(a[:], s)
    } else {
        copy(a[4-len(s):], s)
    }
    return a
}

func main () {
     fmt.Println(time.Now())

     maptest1 :=make(map[string][4]byte)
     maptest1["test"]=  [4]byte {220, 181, 38, 150}
     f,err := os.Open(os.Args[1])
     defer f.Close()
    
     if err == nil {
        buf := bufio.NewReader(f)
        for {

              line,err := buf.ReadString('\n')
              if err != nil || io.EOF == nil {

                break
              }
              line1 := strings.Replace(line,"\n","",-1)
              array_line:= strings.Split(line1," ")
              if len(array_line) == 2{
                 fmt.Println(array_line[0],array_line[1])
                 array_line1 := strings.Split(array_line[1],".")
                 if len(array_line1) == 4 {
                    maptest1[array_line[0]]=  byte16PutString(array_line[1])
                 }
              }

        }

     }
     
      for  k ,v := range maptest1 {
       fmt.Println(k,v)
     }
}
