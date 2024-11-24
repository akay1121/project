// Package biz takes the responsibility of assembling business logic. You should place
// all the domain objects in the biz package, which should include the following structures:
//
//   - Entity represents a distinct concept or thing in the domain, and its identity remains
//     constant over time, even if its attributes change.
//   - Event is something significant that happens within the domain that needs to be recorded and possibly acted upon.
//   - Service contains business logic that does not naturally fit within a particular entity or value object.
//     Services encapsulate operations that involve multiple entities or aggregates.
//   - Repository is responsible for retrieving and storing aggregates.
package biz

import "github.com/google/wire"

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(
	NewProjectManager,
)
