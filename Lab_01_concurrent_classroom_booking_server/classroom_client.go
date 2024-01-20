package main

import (
        "bufio"
	"encoding/csv"
	"encoding/gob"
        "fmt"
        "net"
        "os"
	"io"
	"strconv"
)

type request struct {
	Timestamp int
	Req_type int
	Room string
	Slot string
}

func getReqType(x string) int {
	switch x {
		case "BOOK":
			return 1
		case "CANCEL":
			return 2
		case "GET":
			return 3
		default:
			return -1
	}
}

func main() {
        arguments := os.Args
        if len(arguments) < 3 {
                fmt.Println("Usage: infile host:port")
                return
        }
        
	f, err := os.Open(arguments[1])
	if err != nil {
		fmt.Println("ERROR ::", err)
		return
	}
	// close the file at the end of the program
	defer f.Close()
	
	csvReader := csv.NewReader(f)
	_, err = csvReader.Read()
	if err != nil {
		fmt.Println("ERROR ::", err)
		return
	}
	fmt.Println("Type,Room,Timeslot,Status")
	
        CONNECT := arguments[2]

	for {
		conn, err := net.Dial("tcp", CONNECT)
		if err != nil {
			fmt.Println("ERROR ::", err)
			return
		}
		encoder := gob.NewEncoder(conn)

		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("An error encountered ::", err)
			return
		}
		
		timestamp, err := strconv.Atoi(rec[0])
		if err != nil {
			fmt.Println("An error encountered ::", err)
			return
		}
		req_type := getReqType(rec[1])
		
		req := &request{Timestamp: timestamp,
			Req_type: req_type,
			Room: rec[2], Slot: rec[3]}
		err = encoder.Encode(req)
		if err != nil {
			fmt.Println("ERROR ::", err)
			return
		}

		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print(message)
		conn.Close()
	}
}
