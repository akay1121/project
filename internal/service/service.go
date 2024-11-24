// Package service should be in charge of the conversion from data transfer objects to domain objects.
// It provides the functionality of manipulating domain objects in the biz package
// and handling the requests sent to the servers (either GRPC or HTTP).
package service

import "github.com/google/wire"

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(
	NewProjectService,
)
