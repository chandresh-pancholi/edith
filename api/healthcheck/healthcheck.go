package healthcheck

import "net/http"

//IsHealthy return 200 OK if application is up and ready to accept request
func IsHealthy(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
