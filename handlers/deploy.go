// Copyright (c) OpenFaaS Author(s) 2019. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/foize/go.fifo"
	typesv1 "github.com/openfaas/faas-provider/types"
)

type Status int64
const (
	READY Status = 0
	RUNNING     = 1
	POWEROFF     = 2
)

type Worker struct{
	id int
	ip string
	status Status
}

var job_queue = fifo.NewQueue()

func scheduler(){
//	job := job_queue.Next()
//	the_job := job.(*Payload)
//	log.Info(the_job.Params)
	for job_queue.Len() > 0 {
		job := job_queue.Next()
	 	the_job := job.(*Payload)
		log.Info(the_job.Params)
	}
}

var functions = map[string]*typesv1.FunctionStatus{}

// MakeDeployHandler creates a handler to create new functions in the cluster
func MakeDeployHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		log.Info("deployment request")
		go scheduler()
		defer r.Body.Close()

		body, _ := ioutil.ReadAll(r.Body)

		request := typesv1.FunctionDeployment{}
		if err := json.Unmarshal(body, &request); err != nil {
			log.Errorln("error during unmarshal of create function request. ", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		functions[request.Service] = requestToStatus(request)

		log.Infof("deployment request for function %s", request.Service)

		w.WriteHeader(http.StatusOK)
	}
}
