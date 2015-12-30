package main

import(
    "fmt"
    "crypto/sha256"
    "crypto/aes"
	"crypto/cipher"
    "os"
    "net"
    "io"
    "sync"
)

type processedConfig struct{
    Server *net.TCPAddr
    Local *net.TCPAddr
    Key []byte
}

var conf processedConfig

func main(){
    fmt.Println("SecureTunnel Server (Alpha).")
    //config part
    if(len(os.Args)!=2){
        fmt.Println("[FATAL] param error.")
        return
    }
    err := readConfig(os.Args[1])
    if(err!=nil){
        fmt.Println("[FATAL] read config fail.")
    }
    hash := sha256.New()
    hash.Write([]byte(cfg.Key))
    conf.Key = hash.Sum(nil)
    conf.Local, err = net.ResolveTCPAddr("tcp", cfg.Local)
	if(err!=nil){
        fmt.Println("[FATAL] local_address error.")
        return
    }
    conf.Server, err = net.ResolveTCPAddr("tcp", cfg.Server)
	if(err!=nil){
        fmt.Println("[FATAL] server_address error.")
        return
    }
    //tcp part
	listen, err := net.ListenTCP("tcp", conf.Server)
	if(err!=nil){
        fmt.Println("[ERROR] listen fail.")
    }
	fmt.Println("[INFO] Start server...")
	for {
		conn, err := listen.Accept()
		if(err!=nil){
            fmt.Println("[ERROR] accept fail.")
        }
		go tcpHandle(conn)
	}
}

func tcpHandle(conn net.Conn) {
    //cipher part
	block, err := aes.NewCipher(conf.Key)
	if err != nil {
		fmt.Println("[ERROR] aes cipher init fail.")
        conn.Close()
        return
	}
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(conn, iv); err != nil {
		fmt.Println("[ERROR] read iv fail.")
        conn.Close()
        return
	}
    //tcp part
    connR, err := net.DialTCP("tcp", nil, conf.Local)
	if(err!=nil){
        fmt.Println("[ERROR] connect server fail.")
        conn.Close()
        return
    }
    encodeStream := cipher.NewCFBEncrypter(block, iv)
	decodeStream := cipher.NewCFBDecrypter(block, iv)
	reader := &cipher.StreamReader{S: encodeStream, R: conn}
	writer := &cipher.StreamWriter{S: decodeStream, W: conn}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		io.Copy(connR, reader)
		connR.Close()
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		io.Copy(writer, connR)
		writer.Close()
		wg.Done()
	}()
	wg.Wait()
}

