# hypernova-client

> A Go client for Hypernova

## Getting Started

`go get github.com/dean177/hypernova-go`

## Example usage

```go
package main

import (
	"github.com/dean177/hypernova-go"
	"log"
	"net/http"
)

func main() {
	renderer := hypernova_go.Renderer{
		Url: "http://localhost:3030/batch",
		Plugins: []hypernova_go.Plugin{hypernova_go.DevPlugin{}},
	}

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/html")

		html, err := renderer.Render(hypernova_go.Jobs{
			"MyComponent.js": {
				Name: "MyComponent.js",
				Data: map[string]interface{}{ "name": "strange" },
			},
			"Component2.js": {
				Name: "Component2.js",
				Data: map[string]interface{}{ "text": "Hi!" },
			},
		})

		if err != nil {			
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte("500"))
		} else {
			writer.WriteHeader(http.StatusOK)
			writer.Write([]byte(`
				<!DOCTYPE html>
				<html>
				<body>
					<h1>Hi</h1>
			`))
			writer.Write([]byte(html))
			writer.Write([]byte(`		
				</body>
				</html>
			`))
		}
	})

	log.Print("Listening")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

## Plugin API

Hypernova enables you to control and alter requests at different stages of the render lifecycle via a plugin system.

### `getViewData`

```go
type ReactProps interface {}
getViewData(name string, data ReactProps) ReactProps
```

Allows you to alter the data that a "view" will receive.

### `prepareRequest`

```go
type Job struct {
	Name string
	Data ReactProps `json:"data"`
}

type Jobs map[string]Job
prepareRequest(jobs Jobs) Jobs
```

This function is called when preparing the request to Hypernova and receives the current running jobs Object.

### `shouldSendRequest`

```go
shouldSendRequest(jobs Jobs) bool
```
If `false` is returned then the request is canceled and falls back to client rendering.


An `every` type function. If one returns `false` then the request is canceled.

### `willSendRequest`

```go
willSendRequest(jobs: Jobs): void {}
```

An event type function that is called prior to a request being sent.

### `afterResponse`

```go
type JobResponse struct {
	Error      JobError 	
	Html       string   
	Duration   float32  
	StatusCode int      
	Success    bool    
}

type HypernovaResponse struct {
	Success bool                 
	Error   string                 
	Results map[string]JobResponse 
}

afterResponse(currentResponse HypernovaResponse, originalResponse HypernovaResponse) HypernovaResponse
```

A reducer type function which receives the current response and the original response from the
Hypernova service.

### `onSuccess`

```go
onSuccess(resp HypernovaResponse, jobs Jobs)
```

An event type function that is called whenever a request was successful.

### `onSuccess`

```go
onError(err error, jobs Jobs)
```

An event type function that is called whenever any error is encountered.
