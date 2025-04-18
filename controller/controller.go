package controller

import "github.com/hirano00o/ollama-openai-proxy/provider"

type Router struct {
	prv provider.Provider
}

func NewRouter(p provider.Provider) *Router {
	return &Router{
		prv: p,
	}
}
