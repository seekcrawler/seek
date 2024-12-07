package kraken

import (
	"math"
	"net/http"
)

const abortIndex int8 = math.MaxInt8 >> 1

func NewRouter(defaultHandler HandlerFunc) *Router {
	engine := &Router{
		RouterGroup: RouterGroup{
			Handlers: nil,
			basePath: "/",
			root:     true,
		},
		trees:              make(methodTrees, 0, 9),
		UseRawPath:         false,
		UnescapePathValues: true,
		RemoveExtraSlash:   false,
		defaultHandler:     defaultHandler,
	}
	engine.RouterGroup.engine = engine
	return engine
}

type Router struct {
	RouterGroup
	trees              methodTrees
	maxParams          uint16
	maxSections        uint16
	UseRawPath         bool
	UnescapePathValues bool
	RemoveExtraSlash   bool
	defaultHandler     HandlerFunc
}

func (engine *Router) addRoute(path string, handlers HandlersChain) {
	assert1(path[0] == '/', "path must begin with '/'")
	assert1(len(handlers) > 0, "there must be at least one handler")

	root := engine.trees.get(http.MethodGet)
	if root == nil {
		root = new(node)
		root.fullPath = "/"
		engine.trees = append(engine.trees, methodTree{method: http.MethodGet, root: root})
	}
	root.addRoute(path, handlers)

	if paramsCount := countParams(path); paramsCount > engine.maxParams {
		engine.maxParams = paramsCount
	}

	if sectionsCount := countSections(path); sectionsCount > engine.maxSections {
		engine.maxSections = sectionsCount
	}
}

func (engine *Router) prepareContext(c *Context, extractor *Extractor) (err error) {
	httpMethod := http.MethodGet
	rPath := c.URL.Path
	unescape := false
	if engine.UseRawPath && len(c.URL.RawPath) > 0 {
		rPath = c.URL.RawPath
		unescape = engine.UnescapePathValues
	}
	if engine.RemoveExtraSlash {
		rPath = cleanPath(rPath)
	}

	c.reset()
	c.Extractor = extractor

	params := make(Params, 0, engine.maxParams)
	skip := make([]skippedNode, 0, engine.maxSections)

	// Find root of the tree for the given HTTP method
	t := engine.trees
	for i, tl := 0, len(t); i < tl; i++ {
		if t[i].method != httpMethod {
			continue
		}
		root := t[i].root
		// Find route in tree
		value := root.getValue(rPath, &params, &skip, unescape)
		if value.params != nil {
			c.Params = *value.params
		}
		if value.handlers != nil {
			c.handlers = value.handlers
			return
		}
		break
	}
	err = handlerNotFoundErr
	return
}

func (engine *Router) handle(c *Context) (err error) {
	err = engine.prepareContext(c, nil)
	if err != nil {
		return
	}
	c.Next()
	return
}
