package hypernova_go

import (
	"net/http"
	"bytes"
	"fmt"
	"encoding/json"
	"log"
	//"io/ioutil"
)

type JobError struct {
	Name    string   `json:"name"`
	Message string   `json:"message"`
	Stack   []string `json:"stack"`
}

type JobResp struct {
	Error      JobError `json:"error"`
	Html       string   `json:"html"`
	Duration   float32  `json:"duration"`
	StatusCode int      `json:"statusCode"`
	Success    bool     `json:"success"`
}

type JobResponse struct {
	Success bool
	Err     string
	Results map[string]JobResp
}

type Plugin interface {
	getViewData(name string, data ReactProps) ReactProps
	prepareRequest(jobs Jobs) Jobs
	shouldSendRequest(jobs Jobs) bool
	willSendRequest(jobs Jobs)
	afterResponse(currentResponse JobResponse, originalResponse JobResponse) JobResponse
	onSuccess(resp JobResponse, jobs Jobs)
	onError(err error, jobs Jobs)
}

func fallbackHtml(viewName string, data ReactProps) string {
	var buffer bytes.Buffer
	jsonString, _ := json.Marshal(data)
	buffer.WriteString(`<div data-hypernova-key="`)
	buffer.WriteString(viewName)
	buffer.WriteString(`"></div><script type="application/json" data-hypernova-key="`)
	buffer.WriteString(viewName)
	buffer.WriteString(`"><!--`)
	buffer.Write(jsonString)
	buffer.WriteString(`}--></script>`)
	return buffer.String()
}

func toHtml(resp JobResponse) string {
	var buffer bytes.Buffer
	for _, jobResponse := range resp.Results {
		buffer.WriteString(jobResponse.Html)
	}
	return buffer.String()
}

type Config struct{}

type Renderer struct {
	Url     string
	Plugins []Plugin
	Config  Config
}

func (r Renderer) Render(jobs Jobs) (html string, err error) {
	jobMap := make(map[string]Job)
	for name, job := range jobs {
		jobMap[name] = job
		for _, plugin := range r.Plugins {
			jobMap[name] = Job{
				Name: name,
				Data: plugin.getViewData(name, jobMap[name].Data),
			}
		}
	}

	for _, plugin := range r.Plugins {
		jobMap = plugin.prepareRequest(jobMap)
	}

	shouldSendRequest := true
	for _, plugin := range r.Plugins {
		shouldSendRequest = shouldSendRequest && plugin.shouldSendRequest(jobMap)
	}

	if !shouldSendRequest {
		results := make(map[string]JobResp)
		for name, jobResp := range jobs {
			results[name] = JobResp{
				Html: fallbackHtml(name, jobResp.Data),
			}
		}
		return toHtml(JobResponse{Results: results, Success: true}), nil
	}

	for _, plugin := range r.Plugins {
		plugin.willSendRequest(jobMap)
	}

	jobMapJsonString, err := json.Marshal(jobMap)
	fmt.Println("request: ", string(jobMapJsonString))
	resp, err := http.Post(
		r.Url,
		"application/json",
		bytes.NewBuffer(jobMapJsonString))
	defer resp.Body.Close()
	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	//byts, err := ioutil.ReadAll(resp.Body)
	//fmt.Println("response BodyRaw:", string(byts))

	originalJobResponse := new(JobResponse)
	err = json.NewDecoder(resp.Body).Decode(originalJobResponse)
	fmt.Println("response Body:", originalJobResponse)

	if err != nil {
		for _, plugin := range r.Plugins {
			plugin.onError(err, jobs)
		}
		log.Fatal(err)
		return "", err
	}

	successfulJobs := make(map[string]Job)
	failedJobs := make(map[string]Job)
	for name, jobResp := range (*originalJobResponse).Results {
		if jobResp.Success {
			successfulJobs[name] = jobs[name]
		} else {
			failedJobs[name] = jobs[name]
		}
	}
	for _, plugin := range r.Plugins {
		plugin.onSuccess(*originalJobResponse, successfulJobs)
	}

	currentJobResponse := *originalJobResponse
	for _, plugin := range r.Plugins {
		currentJobResponse = plugin.afterResponse(currentJobResponse, *originalJobResponse)
	}

	return toHtml(currentJobResponse), nil
}
