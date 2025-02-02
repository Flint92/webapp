package handler

import (
	"github.com/flint92/webapp/context"
	"log"
	"net/http"
	"sort"
	"strings"
)

type nodeType int

const (
	Root nodeType = iota
	Any
	Param
	Static
)

const ANY = "*"

type matchFunc func(string, *context.Context) bool

type BasedOnTreeHandler struct {
	getRoot     *node
	postRoot    *node
	putRoot     *node
	patchRoot   *node
	deleteRoot  *node
	optionsRoot *node
	headRoot    *node
	connectRoot *node
	traceRoot   *node
}

type node struct {
	pattern  string
	children []*node
	h        HandlerFunc
	m        matchFunc
	typ      nodeType
}

func newRootNode(path string) *node {
	return &node{
		pattern:  path,
		children: make([]*node, 0, 2),
		m: func(p string, ctx *context.Context) bool {
			panic("should never call this")
		},
		typ: Root,
	}
}

func newStaticNode(path string) *node {
	return &node{
		pattern:  path,
		children: make([]*node, 0, 2),
		m: func(p string, ctx *context.Context) bool {
			return path == p && path != ANY
		},
		typ: Static,
	}
}

func newAnyNode(path string) *node {
	return &node{
		pattern: path,
		m: func(p string, ctx *context.Context) bool {
			return true
		},
		typ: Any,
	}
}

func newParamNode(path string) *node {
	paramName := path[1:]
	return &node{
		pattern: path,
		m: func(p string, ctx *context.Context) bool {
			if ctx != nil {
				ctx.PathParams[paramName] = p
			}
			return p != ANY
		},
		typ: Param,
	}
}

func newNode(path string) *node {
	if path == ANY {
		return newAnyNode(path)
	}
	if strings.HasPrefix(path, ":") {
		return newParamNode(path)
	}
	return newStaticNode(path)
}

func (h *BasedOnTreeHandler) Route(method string, pattern string, handler HandlerFunc) {
	paths := strings.Split(strings.Trim(pattern, "/"), "/")

	curr := getRoot(h, method)

	for _, path := range paths {
		matchChild, ok := findMatchChild(curr, path, nil)
		if ok {
			curr = matchChild
		} else {
			newChild := newNode(path)
			curr.children = append(curr.children, newChild)
			curr = newChild
		}
	}

	curr.h = handler
}

func (h *BasedOnTreeHandler) ServeHTTP(ctx *context.Context) {
	handler, found := findRouter(getRoot(h, ctx.R.Method), ctx.R.URL.Path, ctx)
	if !found {
		err := ctx.NotFoundJson()
		if err != nil {
			log.Printf("error writing json: %v", err)
		}
	} else {
		handler(ctx)
	}
}

func getRoot(h *BasedOnTreeHandler, method string) *node {
	switch method {
	case http.MethodGet:
		return h.getRoot
	case http.MethodPost:
		return h.postRoot
	case http.MethodPut:
		return h.putRoot
	case http.MethodPatch:
		return h.patchRoot
	case http.MethodDelete:
		return h.deleteRoot
	case http.MethodOptions:
		return h.optionsRoot
	case http.MethodHead:
		return h.headRoot
	case http.MethodConnect:
		return h.connectRoot
	case http.MethodTrace:
		return h.traceRoot
	default:
		return h.getRoot
	}
}

func findRouter(parent *node, pattern string, ctx *context.Context) (HandlerFunc, bool) {
	paths := strings.Split(strings.Trim(pattern, "/"), "/")

	curr := parent

	for _, path := range paths {
		matchChild, ok := findMatchChild(curr, path, ctx)
		if !ok {
			return nil, false
		}
		curr = matchChild
	}

	if curr.h == nil {
		return nil, false
	}

	return curr.h, true
}

func findMatchChild(parent *node, path string, ctx *context.Context) (*node, bool) {
	candidates := make([]*node, 0, len(parent.children))
	for _, child := range parent.children {
		if child.m(path, ctx) {
			candidates = append(candidates, child)
		}
	}

	if len(candidates) == 0 {
		return nil, false
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].typ < candidates[j].typ
	})

	return candidates[len(candidates)-1], true
}
