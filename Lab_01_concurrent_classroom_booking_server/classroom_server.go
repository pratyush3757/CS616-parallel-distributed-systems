package main

import (
        "io"
        "fmt"
        "net"
        "os"
        "strings"
	"sync"
	"encoding/gob"
)

type request struct {
	Timestamp int
	Req_type int
	Room string
	Slot string
}
func getReqTypeString(x int) string {
	switch x {
		case 1:
			return "BOOK"
		case 2:
			return "CANCEL"
		case 3:
			return "GET"
		default:
			return ""
	}
}

type class_tuple struct {
	Room string
	Slot string
}

type booked_classes map[class_tuple]int
func (m booked_classes) String() string {
	var sb strings.Builder
	sb.WriteString("{")
	for k, _ := range m {
		sb.WriteString(fmt.Sprintf("('%s', '%s'), ", k.Room, k.Slot))
	}
	str := sb.String()
	if len(str) > 0 {
		str = str[:len(str)-2]
	}
	str = str + "}"
        return str
}

type response struct {
	Req_type int
	Room string
	Slot string
	Status_code int
	Class_list string
}
func (x response) String() string {
	if x.Class_list != "" {
		return fmt.Sprintf("%s,%s,%s,%d,%s\n",
				getReqTypeString(x.Req_type),
				x.Room, x.Slot, x.Status_code, x.Class_list)
	} else {
		return fmt.Sprintf("%s,%s,%s,%d\n",
				getReqTypeString(x.Req_type),
				x.Room, x.Slot, x.Status_code)
	}
}

type status_code int
const (
	Invalid status_code = iota - 3 // -3
	Cooldown
	Already_Booked
	OK
)

var mu sync.RWMutex
var schedule booked_classes = make(booked_classes)

func bookRoom(req request) response {
	mu.Lock() // writers lock
	defer mu.Unlock()
	k := class_tuple{req.Room, req.Slot}
	_, present := schedule[k]
	if present {
		return response{req.Req_type, req.Room, req.Slot, int(Already_Booked), ""}
	} else {
		schedule[k] = req.Timestamp
		return response{req.Req_type, req.Room, req.Slot, int(OK), ""}
	}
}

func cancelRoom(req request) response {
	mu.Lock() // writers lock
	defer mu.Unlock()
	k := class_tuple{req.Room, req.Slot}
	v, present := schedule[k]
	if present {
		if req.Timestamp - v < 20 {
			return response{req.Req_type,
				req.Room, req.Slot, int(Cooldown), ""}
		}
		
		delete(schedule, k)
		return response{req.Req_type, req.Room, req.Slot, int(OK), ""}
	} else {
		return response{req.Req_type, req.Room, req.Slot, int(Invalid), ""}
	}
}


func getRoomStatus(req request) response {
	mu.RLock() // readers lock
	defer mu.RUnlock()
	return response{req.Req_type, req.Room, req.Slot, int(OK), schedule.String()}
}

func isRequestValid(req request) bool {
	if req.Req_type < 1 {
		return false
	}
	
	if req.Req_type != 3 {
		switch req.Room {
			case "1", "2", "3", "4", "5":
				break
			default:
				return false
		}
		
		switch req.Slot {
			case "8:00-9:30", "9:30-11:00", 
				"11:00-12:30", "12:30-14:00",
				"14:00-15:30", "15:30-17:00",
				"17:00-18:30","18:30-20:00":
				break
			default:
				return false
		}
	}
	
	return true
}

func processRequest(req request) response {
	if isRequestValid(req) {
		switch req.Req_type {
			case 1:
				return bookRoom(req)
			case 2:
				return cancelRoom(req)
			case 3:
				return getRoomStatus(req)
			default:
				break
		}
	}
	return response{req.Req_type, req.Room, req.Slot, int(Invalid), ""}
}

func handleConnection(conn net.Conn) {
        defer conn.Close()
	dec := gob.NewDecoder(conn)
	req := &request{}
	err := dec.Decode(req)
	
	if err == io.EOF {
		return
	}
	if err != nil {
		fmt.Println("ERROR ::", err)
		return
	}
	res := processRequest(*req)
	conn.Write([]byte(res.String()))
}

func main() {
        arguments := os.Args
        if len(arguments) == 1 {
                fmt.Println("Please provide a port number!")
                return
        }

        PORT := ":" + arguments[1]
        l, err := net.Listen("tcp4", PORT)
        if err != nil {
                fmt.Println(err)
                return
        }
        defer l.Close()

        for {
                c, err := l.Accept()
                if err != nil {
                        fmt.Println(err)
                        return
                }
                go handleConnection(c)
        }
}
