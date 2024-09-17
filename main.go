package main

import (
	"bufio"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os/exec"

	"github.com/r3labs/sse/v2"
)

//go:embed templates
var content embed.FS

type Message struct {
	Data []byte
	Err  error
}

func main() {
	client := sse.New()
	client.CreateStream("messages")
	client.CreateStream("live-data")

	sub, _ := fs.Sub(content, "templates")
	http.Handle("/", http.FileServer(http.FS(sub)))
	//http.HandleFunc("/", renderHome)

	go broadcast(client)

	http.Handle("/events", client)
	http.ListenAndServe(":8200", nil)

}

func renderHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/index.html")
}

func broadcast(client *sse.Server) {
	cmd := exec.Command("vnstat", "--live", "--json")
	stdout, _ := cmd.StdoutPipe()
	cmd.Start()

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		m := scanner.Text()
		log.Printf("vnstat: %s\n", m)

		var msg Message

		var transferInfo LiveEventData
		if err := json.Unmarshal([]byte(m), &transferInfo); err != nil {
			msg.Err = err
		}

		snippet := fmt.Sprintf("<span>Down %s Up %s</span>", transferInfo.Rx.Ratestring, transferInfo.Tx.Ratestring)

		// Broadcasting the data to all clients registered on "/events"
		go client.Publish("messages", &sse.Event{
			Data: []byte(snippet),
		})

		go client.Publish("live-data", &sse.Event{
			Data: scanner.Bytes(),
		})
	}
	cmd.Wait()
}

type LiveEventData struct {
	Index   int `json:"index"`
	Seconds int `json:"seconds"`
	Rx      struct {
		Ratestring       string `json:"ratestring"`
		Bytespersecond   int    `json:"bytespersecond"`
		Packetspersecond int    `json:"packetspersecond"`
		Bytes            int    `json:"bytes"`
		Packets          int    `json:"packets"`
		Totalbytes       int    `json:"totalbytes"`
		Totalpackets     int    `json:"totalpackets"`
	} `json:"rx"`
	Tx struct {
		Ratestring       string `json:"ratestring"`
		Bytespersecond   int    `json:"bytespersecond"`
		Packetspersecond int    `json:"packetspersecond"`
		Bytes            int    `json:"bytes"`
		Packets          int    `json:"packets"`
		Totalbytes       int    `json:"totalbytes"`
		Totalpackets     int    `json:"totalpackets"`
	} `json:"tx"`
}

//func watchLive() {
//	//var out []string
//	var outData []vnstat.LiveData
//
//	cmd := exec.Command("vnstat", "--live", "--json")
//	stdout, _ := cmd.StdoutPipe()
//	cmd.Start()
//
//	scanner := bufio.NewScanner(stdout)
//	for scanner.Scan() {
//		data := scanner.Bytes()
//		var res vnstat.LiveData
//		if err := json.Unmarshal(data, &res); err != nil {
//			log.Printf("error unmarshalling data: %s\n", err)
//		}
//
//		outData = append(outData, res)
//
//		//m := scanner.Text()
//		//out = append(out, m)
//		//fmt.Println(m) // print the output line in console
//	}
//
//	cmd.Wait()
//	// 'out' now contains all the lines outputted by the command
//}
