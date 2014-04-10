// 121.254.177.105
// tcpdump -n icmp and icmp[icmptype] != icmp-echo or icmp[icmptype] != icmp-echoreply

// CAP_NET_RAW

// https://github.com/atomaths/gtug8/blob/master/ping/ping.go

package main

import (
    // "bytes"
    // "flag"
    "fmt"
    // "log"
    "net"
    // "os"
    "time"
    "syscall"
    "reflect"
    "errors"
)

const (
    ICMP_ECHO_REQUEST = 8
    ICMP_ECHO_REPLY   = 0
)


// copy paste from http://code.google.com/p/go/source/browse/ipv4/helper_unix.go?repo=net#38
var errInvalidConnType = errors.New("invalid conn type")
func sysfd(c net.Conn) (int, error) {
        cv := reflect.ValueOf(c)
        switch ce := cv.Elem(); ce.Kind() {
        case reflect.Struct:
                netfd := ce.FieldByName("conn").FieldByName("fd")
                switch fe := netfd.Elem(); fe.Kind() {
                case reflect.Struct:
                        fd := fe.FieldByName("sysfd")
                        return int(fd.Int()), nil
                }
        }
        return 0, errInvalidConnType
}

func main() {

    sendpkt := []byte("HEAD / HTTP/1.1")

    start := time.Now().Nanosecond()

    raddr, err := net.ResolveTCPAddr("tcp4", "www.google.com:80")
    conn, err := net.DialTCP("tcp", nil, raddr)
    fd, err := sysfd(conn)
    // connSockAddr, err := syscall.Getsockname(fd)
    // connIPAddr := net.IPAddr{ connSockAddr.Addr}
    laddr, err := net.ResolveTCPAddr("tcp4", conn.LocalAddr().String())
    connIPAddr := net.IPAddr{laddr.IP, laddr.Zone}


    ttl, err := syscall.GetsockoptInt(fd, syscall.IPPROTO_IP, syscall.IP_TTL)

    fmt.Println("Local addr is", connIPAddr, "Default ttl =", ttl)

    _, _, _ = err, sendpkt, start



    go func(){
        ipConn, err := net.ListenIP("ip4:icmp", &connIPAddr)
        resp := make([]byte, 1024)
        for {
            n, faddr, err := ipConn.ReadFrom(resp)
            fmt.Println("Answer from", faddr, "ttl =",)


            _, _ = err, n
        }
        _ = err
    }()


    for cTtl:=0; cTtl <= ttl; cTtl++ {
        syscall.SetsockoptInt(fd, syscall.IPPROTO_IP, syscall.IP_TTL, cTtl)
        // fmt.Println("Now ttl =", cTtl)
        conn.Write(sendpkt)

        time.Sleep(1000 * time.Millisecond)
    }

}