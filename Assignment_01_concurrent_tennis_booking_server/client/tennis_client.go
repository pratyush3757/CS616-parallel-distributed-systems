package main

import (
    // "bufio"
    "encoding/csv"
    "encoding/gob"
    "fmt"
    "io"
    "net"
    "os"
    "strconv"
)

type request struct {
    Player_ID		int
    Timestamp		int
    Gender		string
    Game_preference	string
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
            fmt.Println("ERROR ::", err)
            return
        }

        Player_ID, err := strconv.Atoi(rec[0])
        if err != nil {
            fmt.Println("ERROR ::", err)
            return
        }
        Timestamp, err := strconv.Atoi(rec[1])
        if err != nil {
            fmt.Println("ERROR ::", err)
            return
        }

        req := &request{Player_ID: Player_ID,
            Timestamp: Timestamp,
            Gender: rec[2], Game_preference: rec[3]}
        err = encoder.Encode(req)
        if err != nil {
            fmt.Println("ERROR ::", err)
            return
        }

        // message, _ := bufio.NewReader(conn).ReadString('\n')
        // fmt.Println(message)
        conn.Close()
    }
}
