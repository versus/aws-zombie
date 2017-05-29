package main

import (
	"log"
	"net/http"

	"encoding/json"

	"fmt"

	"github.com/googollee/go-socket.io"
	"github.com/ivch/aws-zombie/socket/sns"
	"github.com/ivch/aws-zombie/socket/sqs"
	"github.com/rs/cors"
)

var sockets map[string][]socketio.Socket = make(map[string][]socketio.Socket)

type User struct {
	PushID string `json:"push_id"`
}

func sendMessage(sr *sqs.SendRequest) {
	go func() {
		sent := false
		for k, socket := range sockets {
			if k == sr.Phone {
				for _, s := range socket {
					bts, _ := json.Marshal(sr)
					s.Emit("message", string(bts))
					sent = true
				}
			}
		}

		if !sent {
			rsp, err := http.Get("https://api.eu.zombiegram.tk/user/" + sr.Phone)
			if err == nil {
				var u User
				if err := json.NewDecoder(rsp.Body).Decode(&u); err == nil {
					sns.SendPush(sns.New(), u.PushID)
				}
			}
		}
	}()
}

func main() {

	//snsConn := sns.New()
	//sns.SendPush(snsConn)
	//
	//return

	conn := sqs.New()

	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	server.On("connection", func(so socketio.Socket) {
		log.Println("on connection")

		so.Emit("foo", "bar")

		so.On("phone", func(msg string) {
			if _, ok := sockets[msg]; ok {
				sockets[msg] = append(sockets[msg], so)
				return
			}

			sockets[msg] = []socketio.Socket{so}
			fmt.Println(msg)
		})

		so.On("disconnection", func() {
			for j, vArr := range sockets {
				for i, v := range vArr {
					if v == so {
						sockets[j] = append(vArr[:i], vArr[i+1:]...)
					}
				}
			}

			log.Println("on disconnect")
		})
	})

	server.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})

	http.Handle("/socket.io/", c.Handler(server))
	http.HandleFunc("/status", func(rsp http.ResponseWriter, rq *http.Request) {
		rsp.WriteHeader(200)
	})
	http.HandleFunc("/pm", func(rsp http.ResponseWriter, rq *http.Request) {
		var sr sqs.SendRequest
		if err := json.NewDecoder(rq.Body).Decode(&sr); err != nil {
			rsp.WriteHeader(http.StatusBadRequest)
			rsp.Write([]byte(err.Error()))
			return
		}

		rsp.WriteHeader(200)
		go func() {
			sendMessage(&sr)
		}()
	})

	sqs.Listen(conn, func(sr *sqs.SendRequest) {
		sendMessage(sr)
	})

	log.Println("Serving at localhost:8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
