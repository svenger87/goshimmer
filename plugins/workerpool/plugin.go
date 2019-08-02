package workerpool

import (
	"runtime"

	"github.com/gammazero/workerpool"
	"github.com/iotaledger/goshimmer/packages/node"
)

var PLUGIN = node.NewPlugin("workerpool", node.Enabled, configure, run)
var WP *workerpool.WorkerPool

func configure(plugin *node.Plugin) {
	WP = workerpool.New(runtime.NumCPU())
}

func run(plugin *node.Plugin) {
}
