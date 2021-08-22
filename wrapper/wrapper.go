package wrapper

import "net/http"

// https://github.com/justinas/alice/blob/master/chain.go

type CallFunc func(*http.Response, *http.Request) error

// CallWrapper A constructor for a piece of middleware.
// Some middleware use this constructor out of the box,
// so in most cases you can just pass somepackage.New
type CallWrapper func(CallFunc) CallFunc

// Chain acts as a list of http.Handler middlewares.
// Chain is effectively immutable:
// once created, it will always hold
// the same set of middlewares in the same order.
type Chain struct {
	wrappers []CallWrapper
}

// New creates a new chain,
// memorizing the given list of middleware middlewares.
// New serves no other function,
// middlewares are only called upon a call to Then().
func New(wrappers ...CallWrapper) Chain {
	return Chain{append(([]CallWrapper)(nil), wrappers...)}
}

// Then chains the middleware and returns the final http.Handler.
//     New(m1, m2, m3).Then(h)
// is equivalent to:
//     m1(m2(m3(h)))
// When the request comes in, it will be passed to m1, then m2, then m3
// and finally, the given handler
// (assuming every middleware calls the following one).
//
// A chain can be safely reused by calling Then() several times.
//     stdStack := alice.New(ratelimitHandler, csrfHandler)
//     indexPipe = stdStack.Then(indexHandler)
//     authPipe = stdStack.Then(authHandler)
// Note that middlewares are called on every call to Then()
// and thus several instances of the same middleware will be created
// when a chain is reused in this way.
// For proper middleware, this should cause no problems.
//
// Then() treats nil as http.DefaultServeMux.
func (c Chain) Then(h CallFunc) CallFunc {
	if h == nil {
		//h = http.DefaultServeMux
	}

	for i := range c.wrappers {
		h = c.wrappers[len(c.wrappers)-1-i](h)
	}

	return h
}

// Append extends a chain, adding the specified middlewares
// as the last ones in the request flow.
//
// Append returns a new chain, leaving the original one untouched.
//
//     stdChain := alice.New(m1, m2)
//     extChain := stdChain.Append(m3, m4)
//     // requests in stdChain go m1 -> m2
//     // requests in extChain go m1 -> m2 -> m3 -> m4
func (c Chain) Append(wrappers ...CallWrapper) Chain {
	newCons := make([]CallWrapper, 0, len(c.wrappers)+len(wrappers))
	newCons = append(newCons, c.wrappers...)
	newCons = append(newCons, wrappers...)

	return Chain{newCons}
}

// Extend extends a chain by adding the specified chain
// as the last one in the request flow.
//
// Extend returns a new chain, leaving the original one untouched.
//
//     stdChain := alice.New(m1, m2)
//     ext1Chain := alice.New(m3, m4)
//     ext2Chain := stdChain.Extend(ext1Chain)
//     // requests in stdChain go  m1 -> m2
//     // requests in ext1Chain go m3 -> m4
//     // requests in ext2Chain go m1 -> m2 -> m3 -> m4
//
// Another example:
//  aHtmlAfterNosurf := alice.New(m2)
// 	aHtml := alice.New(m1, func(h http.Handler) http.Handler {
// 		csrf := nosurf.New(h)
// 		csrf.SetFailureHandler(aHtmlAfterNosurf.ThenFunc(csrfFail))
// 		return csrf
// 	}).Extend(aHtmlAfterNosurf)
//		// requests to aHtml hitting nosurfs success handler go m1 -> nosurf -> m2 -> target-handler
//		// requests to aHtml hitting nosurfs failure handler go m1 -> nosurf -> m2 -> csrfFail
func (c Chain) Extend(chain Chain) Chain {
	return c.Append(chain.wrappers...)
}
