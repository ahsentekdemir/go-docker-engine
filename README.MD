[![Go Report Card](https://goreportcard.com/badge/github.com/yikikdeveloper/go-docker-engine)](https://goreportcard.com/report/github.com/yikikdeveloper/go-docker-engine)
[![Go Reference](https://pkg.go.dev/badge/github.com/yikikdeveloper/go-docker-engine.svg)](https://pkg.go.dev/github.com/yikikdeveloper/go-docker-engine)

## About the project

This project is about you can simply manage your docker engine with API calls.

### DOC

* `*/api/v1` is base route.
*  `/container/list` lists all containers which are running in your system.
*  `/container/log/{container-id}` shows logs about your container.
*  `/container/stop-all` stops all containers.
*  `image/pull/{image}` pulling image.
*  `image/list` lists all images in your system.

