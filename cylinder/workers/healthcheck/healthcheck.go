package healthcheck

import (
	"errors"
	"fmt"
	"time"

	"github.com/bandprotocol/chain/v2/cylinder"
	"github.com/bandprotocol/chain/v2/cylinder/client"
	"github.com/bandprotocol/chain/v2/pkg/logger"
	bandtsstypes "github.com/bandprotocol/chain/v2/x/bandtss/types"
)

// HealthCheck is a worker responsible for updating active status to the chain
type HealthCheck struct {
	context *cylinder.Context
	logger  *logger.Logger
	client  *client.Client
}

var _ cylinder.Worker = &HealthCheck{}

// New creates a new instance of the HealthCheck worker.
// It initializes the necessary components and returns the created HealthCheck instance or an error if initialization fails.
func New(ctx *cylinder.Context) (*HealthCheck, error) {
	cli, err := client.New(ctx.Config, ctx.Keyring)
	if err != nil {
		return nil, err
	}

	return &HealthCheck{
		context: ctx,
		logger:  ctx.Logger.With("worker", "HealthCheck"),
		client:  cli,
	}, nil
}

// updateHealthCheck updates last active
func (a *HealthCheck) updateHealthCheck() {
	// Query active information
	status, err := a.client.QueryStatus(a.context.Config.Granter)
	if err != nil {
		a.logger.Error(":cold_sweat: Failed to query status information: %s", err)
		return
	}

	if status.Status == bandtsstypes.MEMBER_STATUS_UNSPECIFIED {
		a.logger.Info(":cold_sweat: User's status not specified; maybe because of no active group")
		return
	} else if status.Status != bandtsstypes.MEMBER_STATUS_ACTIVE {
		a.logger.Error(fmt.Sprintf(":cold_sweat: User's status is not active: %s", status.Status))
		a.context.ErrCh <- errors.New("the status of the address is not active")
		return
	}

	if time.Now().Before(status.LastActive.Add(a.context.Config.ActivePeriod)) {
		return
	}

	// Log
	a.logger.Info(":delivery_truck: Updating last active")

	// Send MsgActive
	a.context.MsgCh <- bandtsstypes.NewMsgHealthCheck(a.context.Config.Granter)
}

// Start starts the healthcheck worker that will check latest healthcheck of validator on BandChain
// and send healthcheck msg if needed every hour.
func (a *HealthCheck) Start() {
	a.logger.Info("start")

	// Update one time when starting worker first time.
	a.updateHealthCheck()

	for range time.Tick(time.Hour * 1) {
		a.updateHealthCheck()
	}
}

// Stop stops the HealthCheck worker.
func (a *HealthCheck) Stop() error {
	a.logger.Info("stop")
	return a.client.Stop()
}
