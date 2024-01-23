package grogu

import (
	"sync/atomic"
	"time"

	rpcclient "github.com/cometbft/cometbft/rpc/client"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bandprotocol/chain/v2/grogu/executor"
	"github.com/bandprotocol/chain/v2/pkg/filecache"
	"github.com/bandprotocol/chain/v2/x/oracle/types"
)

type FeeEstimationData struct {
	askCount    int64
	minCount    int64
	callData    []byte
	rawRequests []rawRequest
	clientID    string
}

type ReportMsgWithKey struct {
	msg               *types.MsgReportData
	execVersion       []string
	keyIndex          int64
	feeEstimationData FeeEstimationData
}

type Context struct {
	client           rpcclient.Client
	validator        sdk.ValAddress
	gasPrices        string
	keys             []*keyring.Record
	executor         executor.Executor
	fileCache        filecache.Cache
	broadcastTimeout time.Duration
	maxTry           uint64
	rpcPollInterval  time.Duration
	maxReport        uint64

	pendingSymbols     chan []string
	inProgressSymbols  *InProgressSymbols
	freeKeys           chan int64
	keyRoundRobinIndex int64 // Must use in conjunction with sync/atomic

	home string
}

func (c *Context) nextKeyIndex() int64 {
	keyIndex := atomic.AddInt64(&c.keyRoundRobinIndex, 1) % int64(len(c.keys))
	return keyIndex
}
