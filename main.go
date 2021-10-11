package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gorilla/mux"
)

func containerList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"message" : "error"},`)))
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"message" : "error"},`)))
	}
	c, err := json.Marshal(containers)
	w.Write([]byte(fmt.Sprintf(`{"services" : %+q},`, c)))

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
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"message": "error"},`)))
	}

	io.Copy(w, out)

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

	if val, ok := params["image"]; ok {
		image = val
		if err != nil {
			panic(err)
		}
	}
	out, err := cli.ImagePull(ctx, image, types.ImagePullOptions{RegistryAuth: authStr})
	if err != nil {
		panic(err)
	}

	w.Write([]byte(fmt.Sprintf(`{"message" : "%s pulled.." },`, image)))

	defer out.Close()

}


func stopRunning(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"message" : "error"},`)))
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"message" : "error"}`)))
	}

	for _, container := range containers {
		c, err := json.Marshal(container.ID[:10])
		w.Write([]byte(fmt.Sprintf(`{"Stopping container" : %+q},`, c)))
		if err != nil{
			w.Write([]byte(fmt.Sprintf(`{"message": "error"},`)))
		}
		
		if err := cli.ContainerStop(ctx, container.ID, nil); err != nil {
			w.Write([]byte(fmt.Sprintf(`{"message": "error"},`)))
		}
		w.Write([]byte(fmt.Sprintf(`{"message": "success"},`)))
	}	
}

func listImages(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"message" : "error"},`)))
	}

	images, err := cli.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"message" : "error"},`)))
	}

	for _, image := range images {
		c, err := json.Marshal(image.ID)
		w.Write([]byte(fmt.Sprintf(`{"Image" : +%q},`, c)))
		if err != nil{
			w.Write([]byte(fmt.Sprintf(`{"message" : "error"},`)))
		}
	}

}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"message": "not found"}`))
}

func main() {
	r := mux.NewRouter()
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/container/list", containerList).Methods(http.MethodGet)
	api.HandleFunc("/container/log/{containerID}", loggingContainer).Methods(http.MethodGet)
	api.HandleFunc("/container/stop-all", stopRunning).Methods(http.MethodGet)
	api.HandleFunc("/image/pull/{image}", pullImage).Methods(http.MethodGet)
	api.HandleFunc("/image/list", listImages).Methods(http.MethodGet)

	api.HandleFunc("/", notFound)
	log.Fatal(http.ListenAndServe(":5000", r))
}
