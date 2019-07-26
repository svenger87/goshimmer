package chosenneighbors

import (
	"sync"

	"github.com/iotaledger/goshimmer/plugins/autopeering/instances/neighborhood"
	"github.com/iotaledger/goshimmer/plugins/autopeering/instances/ownpeer"
	"github.com/iotaledger/goshimmer/plugins/autopeering/types/peerlist"
)

var CANDIDATES peerlist.PeerList
var CANDIDATES_LOCK sync.RWMutex

func configureCandidates() {
	updateNeighborCandidates()

	neighborhood.Events.Update.Attach(updateNeighborCandidates)
}

func updateNeighborCandidates() {
	CANDIDATES_LOCK.Lock()
	neighborhood.LIST_INSTANCE_LOCK.Lock()
	CANDIDATES = neighborhood.LIST_INSTANCE.Sort(DISTANCE(ownpeer.INSTANCE))
	neighborhood.LIST_INSTANCE_LOCK.Unlock()
	CANDIDATES_LOCK.Unlock()
}
