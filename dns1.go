package main 
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
        "github.com/gin-gonic/gin"
       )


type  Storage struct {
      addressBookOfA   map[string][4]byte
      addressBookOfPTR  map[string]string 
      mutex      sync.Mutex
}

func InitStorage() *Storage {
     return &Storage {
           mutex:sync.Mutex{},
           addressBookOfA:map[string][4]byte{},
           addressBookOfPTR:map[string]string{},
     } 
}



func (s *Storage )InitData(fname string) {
  f,err := os.Open(fname)
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
                    s.addressBookOfA[array_line[0]+array_line[1]]= inet_ntoa(InetAtoN(array_line[2]))
                    array_line2 :=  strings.Split(array_line[2],".")
                    if len(array_line2)  == 4 {
                  fmt.Println(array_line2[3]+"."+array_line2[2]+"."+array_line2[1]+"."+array_line2[0]+".in-addr.arpa.",array_line[1])
            s.addressBookOfPTR[array_line2[3]+"."+array_line2[2]+"."+array_line2[1]+"."+array_line2[0]+".in-addr.arpa."]=array_line[1]
                    }
                  
                 }
              }

        }

     }
 
 

}

// ServerDNS serve
func (s *Storage )ServerDNS(addr *net.UDPAddr, conn *net.UDPConn, msg dnsmessage.Message) {
        // query info
        s.addressBookOfA["172.19.220.208www.baidu.com."]= [4]byte{202,99,160,68}
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
        fmt.Println(time.Now(),"收到解析IP的请求 源地址",addrIp, "需要解析的字符串", queryNameStr)
/*        
       for k1,v1 := range s.addressBookOfA {
           fmt.Println("==============",k1,v1)
        }
       for k2,v2 := range s.addressBookOfPTR {
           fmt.Println("==============",k2,v2)
        }

*/


        // find record
        var resource dnsmessage.Resource
        switch queryType {
        case dnsmessage.TypeA:
                if rst, ok := s.addressBookOfA[addrIp+queryNameStr]; ok {
                        resource = NewAResource(queryName, rst)
                        fmt.Println("1111111111",resource,queryName, rst)
                } else {
                        
                        // 开始请求远程地址解析
                         rip := GetIp(queryNameStr) 
                         fmt.Println(time.Now(),"@@@@@@@@@@@@@@远端解析地址如下",queryNameStr,rip)
                         array_line := strings.Split(rip,".")
                         if len(array_line)  == 4 {
                              s.addressBookOfPTR[array_line[3]+"."+array_line[2]+"."+array_line[1]+"."+array_line[0]+".in-addr.arpa."]=queryNameStr
                              resource = NewAResource(queryName, inet_ntoa(InetAtoN(rip)) )
                              fmt.Println("222222222",resource,queryName, inet_ntoa(InetAtoN(rip)) )
                          }




/*
                        fmt.Printf("not fount A record queryName: [%s] \n", queryNameStr)
                        Response(addr, conn, msg)
                        return
*/
                }
        case dnsmessage.TypePTR:
        if rst, ok := s.addressBookOfPTR[queryName.String()]; ok {
			resource = NewPTRResource(queryName, rst)
		} else {
			fmt.Printf("not fount PTR record queryName: [%s] \n", queryNameStr)
			Response(addr, conn, msg)
			return
		}


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
    return  bytes
}

func NewPTRResource(query dnsmessage.Name, ptr string) dnsmessage.Resource {
	name, _ := dnsmessage.NewName(ptr)
	return dnsmessage.Resource{
		Header: dnsmessage.ResourceHeader{
			Name:  query,
			Class: dnsmessage.ClassINET,
		},
		Body: &dnsmessage.PTRResource{
			PTR: name,
		},
	}
}


func GetIp(url string)(ips string) {

       var  ip string
        ns, err := net.LookupHost(url)
        if err != nil {
                fmt.Fprintf(os.Stderr, "Err: %s", err.Error())
                return "0"
        }

        for _, n := range ns {
        //      fmt.Fprintf(os.Stdout, "--%s\n", n)
       //         fmt.Println(url,n)
                array_line := strings.Split(n,".")
                if  len(array_line) == 4 {
                 ip = n
                } 

        }
        return strings.Replace(ip,"\n","",-1)
}




func (s *Storage) ShowA ()(iplist  string) {
    var  str1  string
    for  k,v := range s.addressBookOfA  {
       fmt.Println("++++++++++++++++",k,"=====" ,string(v[:]))
       str1 = k+string(v[:])+"\n" +str1
     }
     return  str1
}
 

func (s *Storage)  StartDns () {
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
                go s.ServerDNS(addr, conn, msg)
        }
    
 }
func main () {
         Storage := InitStorage()
       fmt.Println(time.Now(),"开始读数据")
       Storage.InitData(os.Args[1]) 
    go    Storage.StartDns()

     router := gin.Default()
     router.GET("/showa" ,func(c *gin.Context){
        
       
     c.String(200,Storage.ShowA())
     })    
     router.Run(":80")

}
