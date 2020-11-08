package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"log"
	"net"
	"os"
)

type Item struct {
	ID     	 		uint   		`json:"id"`
//	IpAdrClient   	string 		`json:"ipClient"`
	IpInfo			IpInfo		`json:"ipInfo"`
}
type IpInfo struct {
	Ip				string		`json:"ip"`
	Title			string		`json:"info"`
}

func main() {
//Проверка входхого аргумента
	if len(os.Args) != 2 {
		log.Println("Usage: ", os.Args[0], "host:port")
		os.Exit(1)
	}

//Создание структуры для передачи
	item := Item{
		ID: 1,
		IpInfo: IpInfo{
			Ip: "192.168.3.12/24",
			Title: "use",
		},

	}

//Преобразование структуры в JSON
	b, err := json.Marshal(item)
	checkError(err)

//Compress JSON
	compressedData, compressedDataErr := gZipData(b)
	if compressedDataErr != nil {
		log.Fatal(compressedDataErr)
	}

//	log.Println("compressed data:", compressedData)
//	log.Println("compressed data len:", len(compressedData))

//Создние Коннекта
	service := os.Args[1]

	conn, err := net.Dial("udp", service)
	checkError(err)

//Отправка по UDP
	_, err = conn.Write(compressedData)
	checkError(err)


	os.Exit(0)
}

//========================================================//
//HandlerErr
func checkError(err error) {
	if err != nil {
		log.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}

//Compress Data
func gZipData(data []byte) (compressedData []byte, err error) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)

	_, err = gz.Write(data)
	if err != nil {
		return
	}

	if err = gz.Flush(); err != nil {
		return
	}

	if err = gz.Close(); err != nil {
		return
	}

	compressedData = b.Bytes()

	return
}


