package hypernova_go

import (
	"log"
	//"bytes"
	"encoding/json"
	"bytes"
)

type DevPlugin struct{}

func renderStack(stack []string) string {
	var resBuffer bytes.Buffer
	for _, msg := range stack {
		resBuffer.WriteString("<li>" + msg + "</li>")
	}
	return resBuffer.String()
}

func renderError(name string, err JobError) string {
	return `
    <div style="background-color: #ff5a5f; color: #fff; padding: 12px;">
      <p style="margin: 0">
        <strong>Development Warning!</strong>
        The <code>` + name + `</code> component failed to render with Hypernova.
				Message: ` + err.Message + `
				Stack:
      </p>
      <ul style="padding: 0 20px">
        ` + renderStack(err.Stack) + `
      </ul>
    </div>
  `
}

func (_ DevPlugin) getViewData(name string, data ReactProps) ReactProps {
	return data
}

func (_ DevPlugin) prepareRequest(request Jobs) Jobs {
	for key, job := range request {
		dataStr, _ := json.Marshal(job.Data)
		log.Print("Preparing: " + key + " with data: " + string(dataStr))
	}
	return request
}

func (_ DevPlugin) shouldSendRequest(jobs Jobs) bool {
	return true
}

func (_ DevPlugin) willSendRequest(request Jobs) {
	for key, job := range request {
		dataStr, _ := json.Marshal(job.Data)
		log.Print("Requesting: " + key + " with data: " + string(dataStr))
	}
}

func (_ DevPlugin) afterResponse(response HypernovaResponse, originalResponse HypernovaResponse) HypernovaResponse {
	results := make(map[string]JobResponse)
	for name, jobResp := range response.Results {
		if !jobResp.Success {
			results[name] = JobResponse{Html: renderError(name, jobResp.Error)}
		} else {
			results[name] = jobResp
		}
	}
	return HypernovaResponse{
		Success:true,
		Results:results,
	}
}

func (_ DevPlugin) onError(err error, jobs Jobs) {
	log.Fatal(err)
}

func (_ DevPlugin) onSuccess(resp HypernovaResponse, jobs Jobs) {}
