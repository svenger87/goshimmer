package ledgerstate

import (
	"github.com/iotaledger/goshimmer/packages/node"
)

// region plugin module setup //////////////////////////////////////////////////////////////////////////////////////////

var PLUGIN = node.NewPlugin("LedgerState", configure, run)

func configure(plugin *node.Plugin) {
	configureConfirmedLedgerDatabase(plugin)
}

func run(plugin *node.Plugin) {
	// this plugin has no background workers
}

// endregion ///////////////////////////////////////////////////////////////////////////////////////////////////////////
