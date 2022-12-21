package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

type SSEChannel struct {
	Clients []chan string
	Notifier chan string
}

type PCInfo struct {
	DateTime string `json:"dateTime"`
	Os string `json:"os"`
	MemLoad int `json:"memload"`
	CpuLoad int `json:"cpuload"`
}

var sseChannel SSEChannel

func main(){
	fmt.Println("SSE-GO")

	sseChannel = SSEChannel{
		Clients: make([]chan string, 0),
		Notifier: make(chan string),
	}

	done := make(chan interface{})
	defer close(done)

	go RunForever()

	go broadcaster(done)

	d := make(chan interface{})
	defer close(d)
	defer fmt.Println("Closing channel.")

	http.HandleFunc("/sse", func(w http.ResponseWriter, r *http.Request){
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Connection does not support streaming", http.StatusBadRequest)
			return
		}

		sseChan := make(chan string)
		sseChannel.Clients = append(sseChannel.Clients, sseChan)

		// this goroute reads and streams into the channel
		d := make(chan interface{})
		defer close(d)
		defer fmt.Println("Closing channel.")

		for {
			select {
			case <-d:
				close(sseChan)
				return
			case data := <-sseChan:
				fmt.Printf("data: %v \n\n", data)
				fmt.Fprintf(w, "data: %v \n\n", data)
				flusher.Flush()
			}
		}
	})

	http.HandleFunc("/log", logHTTPRequest)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./public"))))

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	fmt.Println("Listening to " + port)
	http.ListenAndServe(":" + port, nil)

}

func logHTTPRequest(w http.ResponseWriter, r *http.Request){
	buf := new(strings.Builder)
	if _, err := io.Copy(buf, r.Body); err != nil {
		fmt.Printf("Error: %v", err)
	}
	method := r.Method

	logMsg := fmt.Sprintf("Method %v, Body: %v", method, buf.String())
	fmt.Println(logMsg)
	sseChannel.Notifier <- logMsg
}

func broadcaster(done <-chan interface{}){
	fmt.Println("Broadcaster Started.")
	for {
		select {
		case <-done:
			return
		case data := <-sseChannel.Notifier:
			for _, channel := range sseChannel.Clients {
				channel <- data
			}
		}
	}
}

func showStatistic(ch chan struct{}){
	time.Sleep(5 * time.Second)
	t := time.Now()

	ch <- struct{}{}
	
	pcInfo := PCInfo{
		Os: runtime.GOARCH, 
		CpuLoad: getCpuUsage(),
		MemLoad: getMemoryUsage(),
		DateTime: t.Format(time.UnixDate),
	}

	data,  err := json.Marshal(pcInfo)
	
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%s\n", string(data))
	sseChannel.Notifier <- string(data)

}

func RunForever(){
	wait := make(chan struct{})
	for {
		go showStatistic(wait)
		<-wait
	}
}

func getCpuUsage()int{
	percent, err := cpu.Percent(time.Second, false)
	if err != nil{
		log.Fatal(err)
	}
	return int(math.Ceil(percent[0]))
}

func getMemoryUsage()int{
	memory, err := mem.VirtualMemory()
	if err != nil{
		log.Fatal(err)
	}
	return int(math.Ceil(memory.UsedPercent))
}
