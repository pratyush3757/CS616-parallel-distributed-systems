package main

import (
	"encoding/gob"
	"fmt"
	"io"
	"math"
	"net"
	"os"
	"time"
	// "strings"
	"sync"
)

type request struct {
    Player_ID		int
    Timestamp		int
    Gender		string
    Game_preference	string
}

func (x court_t) String(court_no int) string {
    if x.game_type == "S" {
        return fmt.Sprintf("%d,%d,%d,%d,%d",
            x.startTime, x.endTime, court_no,
            x.players[0].Player_ID, x.players[1].Player_ID)
    } else {
        return fmt.Sprintf("%d,%d,%d,%d,%d,%d,%d",
            x.startTime, x.endTime, court_no, 
            x.players[0].Player_ID, x.players[1].Player_ID,
            x.players[2].Player_ID, x.players[3].Player_ID)
    }
}

type player_t request

type queue_t []player_t


func (x queue_t) topOrDummy() player_t {
    if x.isEmpty() {
        return player_t{Player_ID: math.MaxInt, Timestamp: math.MaxInt, Gender: "F", Game_preference: "B"}
    }
    return x[0]
}

func (x queue_t) discard_top() queue_t {
    return x[1:]
}

func (x queue_t) isEmpty() bool {
    return len(x) == 0
}

func (x queue_t) pop() (player_t, queue_t) {
    return x.topOrDummy(), x.discard_top()
}

type court_t struct {
    players [4]player_t
    startTime int
    endTime int
    game_type string
}

func getGameLength(game_type string, gender string) int {
    if (game_type == "S") {
        if (gender == "M") {
            return 10
        } else {
            return 5
        }
    } else {
        if (gender == "M") {
            return 15
        } else {
            return 10
        }
    }
}

func getEndTime(start_time int, game_type string, gender string, mixed bool) int {
    if (mixed) {
        return start_time + getGameLength(game_type, "M")
    } else {
        return start_time + getGameLength(game_type, gender)
    }
}

func max3(a int, b int, c int) int {
    if (a > c) {
        a,c = c,a
    }
    if (a > b) {
        a,b = b,a
    }
    if (b > c) {
        b,c = c,b
    }
    return c
}

var mu sync.Mutex
var GlobalTime = 0
var singles queue_t = make(queue_t, 0)
var doubles queue_t = make(queue_t, 0)
var both queue_t = make(queue_t, 0)
var courts [4]court_t;

func getNextTimeJump() int {
    nextTime := courts[0].endTime
    for _,court := range courts {
        if (court.endTime < nextTime) {
            nextTime = court.endTime
        }
    }
    return nextTime
}

func getEmptyCourt() int {
    for i, court := range courts {
        if(court.endTime <= GlobalTime) {
            return i
        }
    }
    return -1
}

func isGameFeasible(game_type string) bool {
    if (game_type == "S") {
        return len(singles) + len(both) >= 2
    } else {
        return len(doubles) + len(both) >= 4
    }
}

func choosePlayer(game_type string) player_t {
    var x player_t;
    if (game_type == "S") {
        single_time := singles.topOrDummy().Timestamp
        both_time := both.topOrDummy().Timestamp
        if (single_time <= both_time) {
            x, singles = singles.pop()
        } else {
            x, both = both.pop()
        }
    } else {
        double_time := doubles.topOrDummy().Timestamp
        both_time := both.topOrDummy().Timestamp
        if (double_time <= both_time) {
            x, doubles = doubles.pop()
        } else {
            x, both = both.pop()
        }
    }
    return x
}

func getFirstPlayerDetails() player_t {
    a := singles.topOrDummy()
    b := doubles.topOrDummy()
    c := both.topOrDummy()

    if (a.Timestamp > c.Timestamp) {
        a,c = c,a
    }
    if (a.Timestamp > b.Timestamp) {
        a,b = b,a
    }
    // would have been needed for a sort3
    // if (b.Timestamp > c.Timestamp) {
    //     b,c = c,b
    // }
    return a
}

func getOppositeGameType(x string) string {
    if (x == "S") {
        return "D"
    } else {
        return "S"
    }
}

func startGame() {
    mu.Lock()

    empty_court := getEmptyCourt()
    // Range check not needed as empty court uses ranged for loop
    // if (empty_court == -1 || empty_court > 3) { 
    if (empty_court == -1) {
        GlobalTime = getNextTimeJump()
        mu.Unlock()
        return
    }

    singleFeasible := isGameFeasible("S")
    doubleFeasible := isGameFeasible("D")
    if !(singleFeasible || doubleFeasible) {
        mu.Unlock()
        return
    }

    var game_type string;
    a := getFirstPlayerDetails()
    switch a.Game_preference {
    case "b":
        if singleFeasible {
            game_type = "S"
        } else {
            game_type = "D"
        }
        break
    case "B":
        if doubleFeasible {
            game_type = "D"
        } else {
            game_type = "S"
        }
        break
    default:
        if isGameFeasible(a.Game_preference) {
            game_type = a.Game_preference
        } else {
            game_type = getOppositeGameType(a.Game_preference)
        }
    }

    var mixed bool = false;
    var startTime int;
    // game_type will always be feasible to play at this point (due to the early returns)
    if (game_type == "S") {
        player_1 := choosePlayer("S")
        player_2 := choosePlayer("S")
        courts[empty_court].players[0], courts[empty_court].players[1] = player_1, player_2
        courts[empty_court].players[2], courts[empty_court].players[3] = player_t{}, player_t{}
        if player_1.Gender != player_2.Gender {
            mixed = true
        }
        startTime = max3(player_1.Timestamp, player_2.Timestamp, GlobalTime)
    } else {
        player_1 := choosePlayer("D")
        player_2 := choosePlayer("D")
        player_3 := choosePlayer("D")
        player_4 := choosePlayer("D")
        courts[empty_court].players[0], courts[empty_court].players[1] = player_1, player_2
        courts[empty_court].players[2], courts[empty_court].players[3] = player_3, player_4
        for _, v := range courts[empty_court].players {
            if (v.Gender != player_1.Gender) {
                mixed = true
            }
        }
        startTime = max3(player_1.Timestamp, player_2.Timestamp, GlobalTime)
        startTime = max3(player_3.Timestamp, player_4.Timestamp, startTime)
    }
    courts[empty_court].startTime = startTime
    courts[empty_court].endTime = getEndTime(startTime, game_type, courts[empty_court].players[0].Gender, mixed)
    courts[empty_court].game_type = game_type

    // to start court no from 1
    fmt.Println(courts[empty_court].String(empty_court + 1))
    // to fill other courts too
    mu.Unlock()
    startGame()
}

func isRequestValid(req request) bool {
    if (req.Player_ID < 1) {
        return false
    }

    if !(req.Gender == "M" || req.Gender == "F") {
            return false
    }
    gp := req.Game_preference
    if !(gp == "S" || gp == "D" || gp == "B" || gp == "b") {
        return false
    }

    return true
}

func processRequest(req request) {
    if !(isRequestValid(req)) {
        return
    }
    mu.Lock()
    defer mu.Unlock()
    switch req.Game_preference {
    case "S":
        singles = append(singles, player_t(req))
        break
    case "D":
        doubles = append(doubles, player_t(req))
        break
    default:
        both = append(both, player_t(req))
        break
    }
    return
}

func cronJob() {
    for {
        time.Sleep(time.Microsecond * 50)
        startGame()
    }
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
    processRequest(*req)
    // conn.Write([]byte("Done"))
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

    fmt.Println("Game-start-time,Game-end-time,Court-Number,Player-ids")
    go cronJob()
    for {
        c, err := l.Accept()
        if err != nil {
            fmt.Println(err)
            return
        }
        go handleConnection(c)
    }
}
