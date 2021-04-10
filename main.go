package main

import (
	"context"
	"encoding/base64"
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
	containerID := params["containerID"]
	var err error
	if val, ok := params["containerID"]; ok {
		containerID = val

	}
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	options := types.ContainerLogsOptions{ShowStdout: true, Follow: true}
	out, err := cli.ContainerLogs(ctx, containerID, options)
	//fix it
	if err != nil {
		panic(err)
	}
	dst := os.Stdout

	io.Copy(w, out)

	if err != nil {
		w.Write([]byte(`{"message": "error"}`))
	}
	w.Write([]byte(fmt.Sprintf(`{"data" : %+q, "containerId" : %v}`, dst, containerID)))
}

func pullImage(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	image := params["image"]

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	authConfig := types.AuthConfig{
		Username: "username",
		Password: "password",
	}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		panic(err)
	}
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)

	if val, ok := params["containerID"]; ok {
		image = val
		if err != nil {
			panic(err)
		}
	}
	out, err := cli.ImagePull(ctx, image, types.ImagePullOptions{RegistryAuth: authStr})
	if err != nil {
		panic(err)
	}

	w.Write([]byte(fmt.Sprintf(`{"message" : "%s pulled.." }`, image)))

	defer out.Close()

	//io.Copy(w, out)

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
	api.HandleFunc("/pull-image/{image}", pullImage).Methods(http.MethodGet)
	api.HandleFunc("/log/{containerID}", loggingContainer).Methods(http.MethodGet)
	api.HandleFunc("/", notFound)
	log.Fatal(http.ListenAndServe(":3131", r))
}
