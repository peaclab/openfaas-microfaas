package handlers

import (
	"net/http"
	"net"
	log "github.com/sirupsen/logrus"
)

// MakeHealthHandler returns 200/OK when healthy
func MakeHealthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		keys, ok := r.URL.Query()["poweroff"]

		ip, _, err:=net.SplitHostPort(r.RemoteAddr)
		if(err != nil){
			log.Info("error during health check");
		}
		userIP := net.ParseIP(ip)
		defer r.Body.Close()
		value := "0"
		if ok {
			value = keys[0]
		}
	
		if(string(value) == "1"){
			log.Info("poweroff notification")
			log.Info("poweroff from " + (userIP.String()))
			for i := range allWorkers{
				if (allWorkers[i].ip == userIP.String()){
					allWorkers[i].status = POWEROFF
					log.Info(allWorkers[i].ip+" is POWEROFF")
					log.Info(allWorkers[0].status, allWorkers[1].status, allWorkers[2].status, allWorkers[3].status, allWorkers[4].status, allWorkers[5].status, allWorkers[6].status, allWorkers[7].status, allWorkers[8].status, allWorkers[9].status)


				}
			}
		} else {
			log.Info("health check request")
			log.Info("healthy from " + (userIP.String()))
			for i := range allWorkers{
				if (allWorkers[i].ip == userIP.String()){
					allWorkers[i].status = READY
					log.Info(allWorkers[i].ip+" is READY")
					log.Info(allWorkers[0].status, allWorkers[1].status, allWorkers[2].status, allWorkers[3].status, allWorkers[4].status, allWorkers[5].status, allWorkers[6].status, allWorkers[7].status, allWorkers[8].status, allWorkers[9].status)



				}
			}
		}
		w.WriteHeader(http.StatusOK)
	}
}
