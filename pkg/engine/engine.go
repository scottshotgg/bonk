package engine

import (
	"net"
)

type (
	Engine interface {
		// we could also have a set of 'extractor' regexs that run to extract values
		// TODO: idk if we need this function
		// ValidLog(line []bool) bool
		Run(line []byte) (net.IP, bool, error)

		// TODO: should this be here?
		// I think we should make this an interface that others can implement
		// We COULD make this a grpc interface that just opens a stream that you implement; aka an agent
		// Then the agent can distill the *understanding* of the log (i.e, 'get ip from json', etc)

		// To make this as generic as possible we need to be able to pass data back from the Run() func
		// Should we make this a map and then it is up to the action to figure it out?
		// Action()
	}
)
