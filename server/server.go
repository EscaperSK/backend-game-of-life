package server

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var n, m = 100, 100
var grid = []bool{}

func Run() {
	for i := 0; i < n*m; i++ {
		grid = append(grid, false)
	}

	upgrader := websocket.Upgrader{}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		setup()

		http.ServeFile(w, r, "public/page.html")
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Upgrade failed:", err)
			return
		}
		defer conn.Close()

		for {
			data := struct {
				N int
				M int

				Grid []bool
			}{
				N: n,
				M: m,

				Grid: grid,
			}

			err = conn.WriteJSON(data)
			if err != nil {
				log.Println("Write faile:", err)
				break
			}

			step()

			time.Sleep(time.Millisecond * 0)
		}
	})

	http.ListenAndServe("127.0.0.1:8080", nil)
}

func setup() {
	for i := 0; i < n*m; i++ {
		grid[i] = false
	}

	grid[n*m/2+m/2] = true
	grid[n*m/2+m/2+n] = true
	grid[n*m/2+m/2+n*2] = true

	grid[n*m/2+m/2+n-1] = true
	grid[n*m/2+m/2+n*2+1] = true
}

func step() {
	newGrid := []bool{}

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			index := j + i*m

			neighbours := []bool{
				grid[(j-1+m)%m+((i-1+m)%n)*m],
				grid[(j-1+m)%m+(i%n)*m],
				grid[(j-1+m)%m+((i+1)%n)*m],

				grid[(j)%m+((i-1+m)%n)*m],
				grid[(j)%m+((i+1)%n)*m],

				grid[(j+1)%m+((i-1+m)%n)*m],
				grid[(j+1)%m+(i%n)*m],
				grid[(j+1)%m+((i+1)%n)*m],
			}

			count := 0
			for _, v := range neighbours {
				if v {
					count++
				}
			}

			value := grid[index]
			if value && (count < 2 || count > 3) {
				value = false
			} else if count == 3 {
				value = true
			}

			newGrid = append(newGrid, value)
		}
	}

	grid = newGrid
}
