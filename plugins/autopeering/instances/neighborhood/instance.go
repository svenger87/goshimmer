package neighborhood

import (
	"sync"
	"time"

	"github.com/iotaledger/goshimmer/packages/daemon"
	"github.com/iotaledger/goshimmer/packages/node"
	"github.com/iotaledger/goshimmer/packages/timeutil"
	"github.com/iotaledger/goshimmer/plugins/autopeering/instances/knownpeers"
	"github.com/iotaledger/goshimmer/plugins/autopeering/instances/outgoingrequest"
	"github.com/iotaledger/goshimmer/plugins/autopeering/types/peerlist"
	"github.com/iotaledger/goshimmer/plugins/autopeering/types/peerregister"
	"github.com/iotaledger/goshimmer/plugins/autopeering/types/request"
)

var INSTANCE *peerregister.PeerRegister

var LIST_INSTANCE peerlist.PeerList
var LIST_INSTANCE_LOCK sync.RWMutex

// Selects a fixed neighborhood from all known peers - this allows nodes to "stay in the same circles" that share their
// view on the ledger an is a preparation for economic clustering
var NEIGHBORHOOD_SELECTOR = func(this *peerregister.PeerRegister, req *request.Request) *peerregister.PeerRegister {
	filteredPeers := peerregister.New()
	for id, peer := range this.Peers {
		filteredPeers.Peers[id] = peer
	}

	return filteredPeers
}

var lastUpdate = time.Now()

func Configure(plugin *node.Plugin) {
	updateNeighborHood()
}

func Run(plugin *node.Plugin) {
	daemon.BackgroundWorker("Neighborhood Updater", func() {
		timeutil.Ticker(updateNeighborHood, 1*time.Second)
	})
}

func updateNeighborHood() {
	if INSTANCE != nil {
		INSTANCE.Lock()
		defer INSTANCE.Unlock()
	}
	knownpeers.INSTANCE.Lock()
	defer knownpeers.INSTANCE.Unlock()
	if INSTANCE == nil || float64(len(INSTANCE.Peers))*1.2 <= float64(len(knownpeers.INSTANCE.Peers)) || lastUpdate.Before(time.Now().Add(-300*time.Second)) {
		INSTANCE = knownpeers.INSTANCE.Filter(NEIGHBORHOOD_SELECTOR, outgoingrequest.INSTANCE)
		LIST_INSTANCE_LOCK.Lock()
		LIST_INSTANCE = INSTANCE.List()
		LIST_INSTANCE_LOCK.Unlock()

		lastUpdate = time.Now()

		Events.Update.Trigger()
	}
}
