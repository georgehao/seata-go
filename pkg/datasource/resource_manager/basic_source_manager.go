package resource_manager

import (
	"context"
	"errors"
	"sync"

	"github.com/seata/seata-go/pkg/constant"
	"github.com/seata/seata-go/pkg/protocol/message"
	"github.com/seata/seata-go/pkg/rm"
)

// BasicSourceManager the base resource manager, all the branch type resource need extend it
type BasicSourceManager struct {
	// lock
	lock sync.RWMutex
}

func (dm *BasicSourceManager) BranchCommit(ctx context.Context, req message.BranchCommitRequest) (int, error) {
	return constant.BranchStatusPhaseOneDone, nil
}

func (dm *BasicSourceManager) BranchRollback(ctx context.Context, req message.BranchRollbackRequest) (int, error) {
	return constant.BranchStatusPhaseOneFailed, nil
}

func (dm *BasicSourceManager) BranchRegister(ctx context.Context, clientId string, req message.BranchRegisterRequest) (int64, error) {
	return 0, nil
}

func (dm *BasicSourceManager) BranchReport(ctx context.Context, req message.BranchReportRequest) error {
	return nil
}

func (dm *BasicSourceManager) LockQuery(ctx context.Context, branchType int, resourceId, xid, lockKeys string) (bool, error) {
	return true, nil
}

// RegisterResource register model.Resource to be managed by model.Resource Manager
func (dm *BasicSourceManager) RegisterResource(resource rm.Resource) error {
	err := rm.GetRMRemotingInstance().RegisterResource(resource)
	if err != nil {
		return err
	}
	return nil
}

// UnregisterResource Unregister a model.Resource from the   model.Resource Manager
func (dm *BasicSourceManager) UnregisterResource(resource rm.Resource) error {
	return errors.New("unsupport unregister resource")
}

func (dm *BasicSourceManager) GetManagedResources() *sync.Map {
	return nil
}

func (dm *BasicSourceManager) GetBranchType() int {
	return constant.BranchTypeAT
}
