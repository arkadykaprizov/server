package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
)

type Item struct {
	ID     	 		uint   		`json:"id"`
	IpAdrClient   	string 		`json:"ipClient"`
	IpInfo			IpInfo		`json:"ipInfo"`
}
type IpInfo struct {
	Ip				string		`json:"ip"`
	Title			string		`json:"info"`
}
type udpSrv struct {
	ipPort			string
	udpAddr			*net.UDPAddr
	listener		*net.UDPConn
	err				error
}

func main() {
//Создаем коннект
	var udpSrv udpSrv
	var item  Item

	udpSrv.ipPort = "127.0.0.1:6000"

	udpSrv.udpAddr, udpSrv.err = net.ResolveUDPAddr("udp4", udpSrv.ipPort)
	if udpSrv.err  != nil {
		log.Fatal(udpSrv.err )
	}
//Слушаем UDP
	udpSrv.listener, udpSrv.err = net.ListenUDP("udp4", udpSrv.udpAddr)
	if udpSrv.err != nil {
		log.Fatal(udpSrv.err)
	}

	fmt.Println("UDP server up and listening on port 6000")

	defer udpSrv.listener.Close()
//Обработка коннекта
	for {
		// wait for UDP client to connect
		handleUDPConnection(udpSrv.listener, &item)
	}
}
//====================================================//
//Decompress data
func gUnzipData(data []byte) (resData []byte, err error) {
	b := bytes.NewBuffer(data)

	var r io.Reader
	r, err = gzip.NewReader(b)
	if err != nil {
		return
	}

	var resB bytes.Buffer
	_, err = resB.ReadFrom(r)
	if err != nil {
		return
	}

	resData = resB.Bytes()

	return
}

//HandlerUdpConnection
func handleUDPConnection(conn *net.UDPConn, item *Item) {

	buffer := make([]byte, 1024)
//Запись данных в буфер
	n, addr, err := conn.ReadFromUDP(buffer)
	if err != nil {
		log.Fatal(err)
	}

//	fmt.Println("UDP client : ", addr)
//	fmt.Println("Received from UDP client :  ", buffer[:n])

	item.IpAdrClient = addr.String()

//Decompress
	uncompressedData, uncompressedDataErr := gUnzipData(buffer[:n])
	if uncompressedDataErr != nil {
		log.Fatal(uncompressedDataErr)
	}
//	os.Stdout.Write(uncompressedData) //Вывести в консоль данные в формате JSON после Decompress

//Unmarshal JSON и запись в структуру
	_ = json.Unmarshal(uncompressedData, item)
	log.Println(*item)

}


