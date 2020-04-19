package main

// gotunnel is a HTTP PROXY compatible proxy that forwards connections via
// SSH to a remote host

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/elazarl/goproxy"
	"golang.org/x/crypto/ssh"
)

var (
	USER        = flag.String("user", "root", "ssh username")
	HOST        = flag.String("host", "", "ssh server hostname")
	PORT        = flag.Int("port", 22, "ssh server port")
	PROXY_ADDR  = flag.String("proxy_addr", "0:8888", "local http proxy address")
	LOCAL_ADDR  = flag.String("local_addr", "127.0.0.1:18083", "local provider listening address")
	REMOTE_ADDR = flag.String("remote_addr", "0:18083", "remote provider listening address")
	PASS        = flag.String("pass", "", "ssh password")
	KEY         = flag.String("key", os.Getenv("HOME")+"/.ssh/id_rsa", "ssh key file path")
)

func init() { flag.Parse() }

type Dialer interface {
	DialTCP(net string, laddr, raddr *net.TCPAddr) (net.Conn, error)
}

func main() {
	log.SetOutput(os.Stdout)
	if *HOST == "" {
		log.Fatalf("must provide remote ssh host!")
	}
	var auths []ssh.AuthMethod
	if *PASS != "" {
		auths = append(auths, ssh.Password(*PASS))
	}
	if *KEY != "" {
		key, err := ioutil.ReadFile(*KEY)
		if err != nil {
			log.Fatalf("unable to read private key: %v", err)
		}
		// Create the Signer for this private key.
		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			log.Fatalf("unable to parse private key: %v", err)
		}
		auths = append(auths, ssh.PublicKeys(signer))
	}

	config := &ssh.ClientConfig{
		User: *USER,
		Auth: auths,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	addr := fmt.Sprintf("%s:%d", *HOST, *PORT)
	conn, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		log.Fatalf("unable to connect to [%s]: %v", addr, err)
	}
	defer conn.Close()
	go ListenRemote(conn)
	prxy := goproxy.NewProxyHttpServer()
	prxy.Tr = &http.Transport{Dial: conn.Dial}
	log.Printf("listening for local HTTP PROXY connections on [%s]\n", *PROXY_ADDR)
	log.Println(http.ListenAndServe(*PROXY_ADDR, prxy))
	log.Println("shutting down")
}

func ListenRemote(client *ssh.Client) {
	// Request the remote side to open port 8080 on all interfaces.
	log.Printf("listening for remote tcp connections on [%s]\n", *REMOTE_ADDR)
	l, err := client.Listen("tcp", *REMOTE_ADDR)
	if err != nil {
		log.Fatal("unable to register tcp forward: ", err)
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal("accept remote failed!err:=%v", err)
			return
		}
		go serveTcp(conn)
	}
}

func serveTcp(conn net.Conn) {
	defer conn.Close()
	localConn, err := net.Dial("tcp", *LOCAL_ADDR)
	if err != nil {
		log.Printf("dial local addr %v failed!err:=%v", *LOCAL_ADDR, err)
		return
	}
	defer localConn.Close()
	ch := make(chan struct{}, 0)
	go func() {
		io.Copy(conn, localConn)
		close(ch)
	}()
	io.Copy(localConn, conn)
	<-ch
}
