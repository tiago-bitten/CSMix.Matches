// Package match is the aggregate this service exists for. It is event sourced: what is stored is what happened, and every read model is derived from that.
//
// The four folders every CSMix slice grows - api, app, domain, infra - appear
// here when there is something to put in them, and not before.
package match
