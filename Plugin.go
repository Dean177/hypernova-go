package hypernova_go

type Plugin interface {
	getViewData(name string, data ReactProps) ReactProps
	prepareRequest(jobs Jobs) Jobs
	shouldSendRequest(jobs Jobs) bool
	willSendRequest(jobs Jobs)
	afterResponse(currentResponse HypernovaResponse, originalResponse HypernovaResponse) HypernovaResponse
	onSuccess(resp HypernovaResponse, jobs Jobs)
	onError(err error, jobs Jobs)
}
