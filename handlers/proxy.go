package handlers

import (
	"encoding/json"
	// b64 "encoding/base64"
	// "fmt"
	"io/ioutil"
	"net/http"
	// "os"
	"time"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/btittelbach/go-bbhw"
	"bytes"
	"strconv"

)

type WorkerResponse struct {
	Fid         string        `json:"fid"`
	TimeElapsed time.Duration `json:"time_elapsed"`
	Result      string        `json:"result"`
	Error       string        `json:"error,omitempty"`
}

type response struct {
	Function     string
	ResponseBody string
	HostName     string
}

func gpio_turn_on(pin_num uint) error {
	pin, err := bbhw.NewSysfsGPIO(pin_num, bbhw.OUT)
	err = pin.SetState(true)
	time.Sleep(500 * time.Millisecond)
	err = pin.SetState(false)
	time.Sleep(500 * time.Millisecond)
	err = pin.SetState(true)
	return err
}
type Payload struct{
	Fid string `json:"fid"`
	Src string `json:"src"`
	Params string `json:"params,omitempty"`
	Lang string `json:"lang"`
	Worker string `json:"worker"`
}

type FuncCall struct{
	Params string `json:"params,omitempty"`
	Lang string `json:"lang,omitempty"`
	Worker string `json:"worker"`
}

type Status int64

const (
	READY    Status = 0
	RUNNING         = 1
	POWEROFF        = 2
)

type Worker struct {
	id     int
	ip     string
	status Status
	// last_ready time.Time
}
//List of workers
var allWorkers = []Worker{
	Worker{0, "192.168.1.20", READY},
	Worker{1, "192.168.1.21", READY},
	Worker{2, "192.168.1.22", READY},
	// Worker{3, "192.168.1.23", READY, time.Now()},
	// Worker{4, "192.168.1.24", READY, time.Now()},
	// Worker{5, "192.168.1.25", READY, time.Now()},
	// Worker{6, "192.168.1.26", READY, time.Now()},
	// Worker{7, "192.168.1.27", READY, time.Now()},
	// Worker{8, "192.168.1.28", READY, time.Now()},
	// Worker{9, "192.168.1.29", READY, time.Now()},
}

func find_worker() string {
	for {
		for i := range allWorkers {
			if(allWorkers[i].status == READY){
				log.Info("Chose worker: " + strconv.Itoa(allWorkers[i].id))
				allWorkers[i].status = RUNNING
				log.Info(allWorkers[0].status, allWorkers[1].status, allWorkers[2].status, allWorkers[3].status, allWorkers[4].status, allWorkers[5].status, allWorkers[6].status, allWorkers[7].status, allWorkers[8].status, allWorkers[9].status)
				return allWorkers[i].ip
			}
		}
  	}
}
// func shutdown() {
// 	for {
// 		for i := range allWorkers {
// 			if(allWorkers[i].status == READY && time.Since(allWorkers[i].last_ready) > 15 * time.Second){
// 				log.Info("worker " + strconv.Itoa(i) + " is shutting down")
// 				// resp, err := client.Post(url, "application/json",
// 				// 	bytes.NewBuffer(packet))

// 				// if err != nil ||  marshal_err != nil {
// 				// 	// log.Fatal(err)
// 				// 	log.Info("HIT AN ERROR HERE: ", err)
// 				// 	return
// 				// }
// 				// resp_body, _ := ioutil.ReadAll(resp.Body)
// 				// log.Info(string(resp_body))
// 			}
// 		}
//   	}
// }

// MakeProxy creates a proxy for HTTP web requests which can be routed to a function.
func MakeProxy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// r.Close = true
		vars := mux.Vars(r)
		name := vars["name"]
		log.Info("proxy request: " + name)

		v, okay := functions[name]
		if !okay {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("{ \"status\" : \"Not found\"}"))
			log.Errorf("%s not found", name)
			return
		}

		// Working GPIO pins
		worker_list := map[int]uint{
			1: 48, // works
			2: 67, // works
			3: 68, // works
		}

		gpio_turn_on(worker_list[3])

		v.InvocationCount = v.InvocationCount + 1

		defer r.Body.Close()
		body, _ := ioutil.ReadAll(r.Body)
		log.Info(string(body))
		var payload Payload
		var func_call FuncCall
		json.Unmarshal([]byte(body), &func_call)
		payload.Fid = name
		payload.Src = v.Image
		payload.Params = func_call.Params
		payload.Worker = func_call.Worker
		payload.Lang = func_call.Lang
		packet, marshal_err := json.Marshal(payload)
		client := http.Client{
			Timeout: 30 * time.Second,
		}
		
		url := "http://" + find_worker() + ":8080/run"
		resp, err := client.Post(url, "application/json",
			bytes.NewBuffer(packet))

		if err != nil ||  marshal_err != nil {
			// log.Fatal(err)
			log.Info("HIT AN ERROR HERE: ", err)
			return
		}
		resp_body, _ := ioutil.ReadAll(resp.Body)
		log.Info(string(resp_body))


		// hostName, _ := os.Hostname()
		// d := &response{
		// 	Function:     name,
		// 	ResponseBody: string(resp_body) ,
		// 	HostName:     hostName,
		// }
		var worker_resp WorkerResponse
		json.Unmarshal([]byte(resp_body), &worker_resp)



		responseBody, err := json.Marshal(worker_resp)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			log.Errorf("error invoking %s. %v", name, err)
			return
		}

		w.Write(responseBody)

		log.Info("proxy request completed: ", name)
	}
}
