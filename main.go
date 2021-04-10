package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gorilla/mux"
)

func runningDockerServices(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		w.Write([]byte(`{"message" : "error"}`))
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		w.Write([]byte(`{"message" : "error"}`))
	}
	c, err := json.Marshal(containers)
	w.Write([]byte(fmt.Sprintf(`{"services" : %+q}`, c)))

}

func loggingContainer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	containerID := " "
	var err error
	if val, ok := params["containerID"]; ok {
		containerID = val

	}
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	options := types.ContainerLogsOptions{ShowStdout: true}
	out, err := cli.ContainerLogs(ctx, "00c7d0cf6591", options)
	//fix it
	io.Copy(os.Stdout, out)

	c, err := json.Marshal(out)

	if err != nil {
		w.Write([]byte(`{"message": "error"}`))
	}
	w.Write([]byte(fmt.Sprintf(`{"data" : %+q, "" : %v}`, c, containerID)))
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"message": "not found"}`))
}

func main() {
	r := mux.NewRouter()
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/docker-services", runningDockerServices).Methods(http.MethodGet)
	api.HandleFunc("/log/{containerID}", loggingContainer).Methods(http.MethodGet)
	api.HandleFunc("/", notFound)
	log.Fatal(http.ListenAndServe(":3131", r))
}
