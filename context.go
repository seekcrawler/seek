package kraken

import "net/url"

type Context struct {
	URL      url.URL
	Params   Params
	fullPath string
	handlers HandlersChain
}

// HandlerFunc defines the handler used by gin middleware as return value.
type HandlerFunc func(*Context) error

type HandlersChain []HandlerFunc
