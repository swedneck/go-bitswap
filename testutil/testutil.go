package testutil

import (
	"math/rand"

	bsmsg "github.com/ipfs/go-bitswap/message"
	bssd "github.com/ipfs/go-bitswap/sessiondata"
	"github.com/ipfs/go-bitswap/wantlist"
	blocks "github.com/ipfs/go-block-format"
	cid "github.com/ipfs/go-cid"
	blocksutil "github.com/ipfs/go-ipfs-blocksutil"
	peer "github.com/libp2p/go-libp2p-core/peer"
)

var blockGenerator = blocksutil.NewBlockGenerator()
var prioritySeq int

// GenerateBlocksOfSize generates a series of blocks of the given byte size
func GenerateBlocksOfSize(n int, size int64) []blocks.Block {
	generatedBlocks := make([]blocks.Block, 0, n)
	for i := 0; i < n; i++ {
		// rand.Read never errors
		buf := make([]byte, size)
		rand.Read(buf)
		b := blocks.NewBlock(buf)
		generatedBlocks = append(generatedBlocks, b)

	}
	return generatedBlocks
}

// GenerateCids produces n content identifiers.
func GenerateCids(n int) []cid.Cid {
	cids := make([]cid.Cid, 0, n)
	for i := 0; i < n; i++ {
		c := blockGenerator.Next().Cid()
		cids = append(cids, c)
	}
	return cids
}

// GenerateWantlist makes a populated wantlist.
func GenerateWantlist(n int, ses uint64) *wantlist.SessionTrackedWantlist {
	wl := wantlist.NewSessionTrackedWantlist()
	for i := 0; i < n; i++ {
		prioritySeq++
		entry := wantlist.NewRefEntry(blockGenerator.Next().Cid(), prioritySeq)
		wl.AddEntry(entry, ses)
	}
	return wl
}

// GenerateMessageEntries makes fake bitswap message entries.
func GenerateMessageEntries(n int, isCancel bool) []bsmsg.Entry {
	bsmsgs := make([]bsmsg.Entry, 0, n)
	for i := 0; i < n; i++ {
		prioritySeq++
		msg := bsmsg.Entry{
			Entry:  wantlist.NewRefEntry(blockGenerator.Next().Cid(), prioritySeq),
			Cancel: isCancel,
		}
		bsmsgs = append(bsmsgs, msg)
	}
	return bsmsgs
}

var peerSeq int

// GeneratePeers creates n peer ids.
func GeneratePeers(n int) []peer.ID {
	peerIds := make([]peer.ID, 0, n)
	for i := 0; i < n; i++ {
		peerSeq++
		p := peer.ID(peerSeq)
		peerIds = append(peerIds, p)
	}
	return peerIds
}

// GenerateOptimizedPeers creates n peer ids,
// with optimization fall off up to optCount, curveFunc to scale it
func GenerateOptimizedPeers(n int, optCount int, curveFunc func(float64) float64) []bssd.OptimizedPeer {
	peers := GeneratePeers(n)
	optimizedPeers := make([]bssd.OptimizedPeer, 0, n)
	for i, peer := range peers {
		var optimizationRating float64
		if i <= optCount {
			optimizationRating = 1.0 - float64(i)/float64(optCount)
		} else {
			optimizationRating = 0.0
		}
		optimizationRating = curveFunc(optimizationRating)
		optimizedPeers = append(optimizedPeers, bssd.OptimizedPeer{Peer: peer, OptimizationRating: optimizationRating})
	}
	return optimizedPeers
}

var nextSession uint64

// GenerateSessionID make a unit session identifier.
func GenerateSessionID() uint64 {
	nextSession++
	return uint64(nextSession)
}

// ContainsPeer returns true if a peer is found n a list of peers.
func ContainsPeer(peers []peer.ID, p peer.ID) bool {
	for _, n := range peers {
		if p == n {
			return true
		}
	}
	return false
}

// IndexOf returns the index of a given cid in an array of blocks
func IndexOf(blks []blocks.Block, c cid.Cid) int {
	for i, n := range blks {
		if n.Cid() == c {
			return i
		}
	}
	return -1
}

// ContainsBlock returns true if a block is found n a list of blocks
func ContainsBlock(blks []blocks.Block, block blocks.Block) bool {
	return IndexOf(blks, block.Cid()) != -1
}
