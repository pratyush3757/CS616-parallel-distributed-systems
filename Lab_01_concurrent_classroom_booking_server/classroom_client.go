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

/*
type request_type int

const (
	Book request_type = iota + 1 // 1
	Cancel
	Get
)

type status_code int

const (
	Invalid status_code = iota - 3 // -3
	Cooldown
	Already_Booked
	OK
)*/

type request struct {
	Timestamp int
	Req_type int
	Room string
	Slot string
}

type class_tuple struct {
	Room string
	Slot string
}

type booked_classes map[class_tuple]int

type response struct {
	Req_type int
	Room int
	Status_code int
	Class_list string
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
	
        // 
        CONNECT := arguments[2]
  //       conn, err := net.Dial("tcp", CONNECT)
  //       if err != nil {
		// fmt.Println("ERROR ::", err)
  //               return
  //       }
        
        // for {
        //         reader := bufio.NewReader(os.Stdin)
        //         fmt.Print(">> ")
        //         text, _ := reader.ReadString('\n')
        //         fmt.Fprintf(c, text+"\n")
        // 
        //         message, _ := bufio.NewReader(c).ReadString('\n')
        //         fmt.Print("->: " + message)
        //         if strings.TrimSpace(string(text)) == "STOP" {
        //                 fmt.Println("TCP client exiting...")
        //                 return
        //         }
        // }

        // encoder := gob.NewEncoder(conn)
	// dec := gob.NewDecoder(conn)
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
		// fmt.Println(rec)
		
		// res := &booked_classes{}
		// err = dec.Decode(res)
		// if err != nil {
			// fmt.Println("ERROR ::", err)
			// return
		// }
                // fmt.Println("res")
		message, _ := bufio.NewReader(conn).ReadString('\n')
                // fmt.Print("->: " + message)
		fmt.Print(message)
		conn.Close()
	}
}
