Gpackage main 
import (
        "fmt"
        "time"
        "sync"
        "os"
        "bufio"
        "io"
        "strings"
        "net"
        "golang.org/x/net/dns/dnsmessage"
        "math/big"
       )


type  Storage struct {
      maptest1   map[string][4]byte
      mutex      sync.Mutex
}

func InitStorage() *Storage {
     return &Storage {
           mutex:sync.Mutex{},
           maptest1:map[string][4]byte{},
     } 
}



func (s *Storage )InitData(fname string) {
  f,err := os.Open(fname)
  fmt.Println(err,"OOO")
     defer f.Close()

     if err == nil {
        buf := bufio.NewReader(f)
        for {

              line,err := buf.ReadString('\n')
              if err != nil || io.EOF == nil {

                break
              }
              line1 := strings.Replace(line,"\n","",-1)
              array_line:= strings.Split(line1,"|")
              if len(array_line) == 3{
                 array_line1 := strings.Split(array_line[2],".")
                 fmt.Println(array_line[0],array_line[1],array_line[2],len(array_line1))
                 if len(array_line1) == 4 {
                    s.maptest1[array_line[0]+array_line[1]]= inet_ntoa(InetAtoN(array_line[2]))
                  
                 }
              }

        }

     }
 
 

}

// ServerDNS serve
func (s *Storage )ServerDNS(addr *net.UDPAddr, conn *net.UDPConn, msg dnsmessage.Message) {
        // query info
        s.maptest1["172.19.220.208www.baidu.com."]= [4]byte{202,99,160,68}
        if len(msg.Questions) < 1 {
                return
        }
        question := msg.Questions[0]
        var (
                queryTypeStr = question.Type.String()
                queryNameStr = question.Name.String()
                queryType    = question.Type
                queryName, _ = dnsmessage.NewName(queryNameStr)
        )
        fmt.Printf("[%s] queryName: [%s]\n", queryTypeStr, queryNameStr)
        // 获取客户端的IP
        str1 := fmt.Sprintf("%s", addr) 
        arrary_str1 := strings.Split(str1,":")
        addrIp :=  arrary_str1[0]
        fmt.Println("ADD",addrIp,queryNameStr)
        

       k := s.maptest1[addrIp+queryNameStr]
       fmt.Println("++++++++++++",k)
       for k1,v1 := range s.maptest1 {
           fmt.Println("==============",k1,v1)
        }
        // find record
        var resource dnsmessage.Resource
        switch queryType {
        case dnsmessage.TypeA:
                if rst, ok := s.maptest1[addrIp+queryNameStr]; ok {
                        resource = NewAResource(queryName, rst)
                } else {
                        fmt.Printf("not fount A record queryName: [%s] \n", queryNameStr)
                        Response(addr, conn, msg)
                        return
                }
        case dnsmessage.TypePTR:
        default:
                fmt.Printf("not support dns queryType: [%s] \n", queryTypeStr)
                return
        }

        // send response
        msg.Response = true
        msg.Answers = append(msg.Answers, resource)
        Response(addr, conn, msg)
}

// Response return
func Response(addr *net.UDPAddr, conn *net.UDPConn, msg dnsmessage.Message) {
        packed, err := msg.Pack()
        if err != nil {
                fmt.Println(err)
                return
        }
        if _, err := conn.WriteToUDP(packed, addr); err != nil {
                fmt.Println(err)
        }
}

// NewAResource A record
func NewAResource(query dnsmessage.Name, a [4]byte) dnsmessage.Resource {
        return dnsmessage.Resource{
                Header: dnsmessage.ResourceHeader{
                        Name:  query,
                        Class: dnsmessage.ClassINET,
                        TTL:   600,
                },
                Body: &dnsmessage.AResource{
                        A: a,
                },
        }
}



func InetAtoN(ip string) int64 {
    ret := big.NewInt(0)
    ret.SetBytes(net.ParseIP(ip).To4())
    return ret.Int64()
}
func inet_ntoa(ipnr int64) [4]byte {
    var bytes [4]byte
    bytes[3] = byte(ipnr & 0xFF)
    bytes[2] = byte((ipnr >> 8) & 0xFF)
    bytes[1] = byte((ipnr >> 16) & 0xFF)
    bytes[0] = byte((ipnr >> 24) & 0xFF)
    fmt.Println("++++++++",bytes)
    return  bytes
}



func main () {
     fmt.Println(time.Now(),"开始读数据")
     Storage := InitStorage()
     Storage.InitData(os.Args[1])  // 初始化数据到内存
    conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: 53})
        if err != nil {
                panic(err)
        }
        defer conn.Close()
        fmt.Println("Listing ...")
        for {
                buf := make([]byte, 512)
                _, addr, _ := conn.ReadFromUDP(buf)

                var msg dnsmessage.Message
                if err := msg.Unpack(buf); err != nil {
                        fmt.Println(err)
                        continue
                }
                go Storage.ServerDNS(addr, conn, msg)
        }
    




}
