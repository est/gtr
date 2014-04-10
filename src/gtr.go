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

        ttl, err := syscall.GetsockoptInt(fd, syscall.IPPROTO_IP, syscall.IP_TTL)

        fmt.Println("Default ttl =", ttl)

        _, _, _ = err, sendpkt, start

        for c_ttl:=0; c_ttl <= ttl; c_ttl++ {
            syscall.SetsockoptInt(fd, syscall.IPPROTO_IP, syscall.IP_TTL, c_ttl)
            fmt.Println("Now ttl =", c_ttl)
            conn.Write(sendpkt)
            time.Sleep(1000 * time.Millisecond)
        }
        
}