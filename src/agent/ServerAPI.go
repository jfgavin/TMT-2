package agent

type ServerAPI interface {
	GetEliminationCount() int
	RequestSacrifice(*TMTAgent)
}
