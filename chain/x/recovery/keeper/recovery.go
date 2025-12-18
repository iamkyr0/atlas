package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	trainingtypes "github.com/atlas/chain/x/training/types"
	computetypes "github.com/atlas/chain/x/compute/types"
	healthtypes "github.com/atlas/chain/x/health/types"
)

func (k Keeper) RollbackTasksForNode(ctx sdk.Context, nodeID string) error {
	var tasksToRollback []string
	
	k.trainingKeeper.IterateTasks(ctx, func(task trainingtypes.Task) (stop bool) {
		if task.NodeID == nodeID && (task.Status == trainingtypes.TaskStatus_IN_PROGRESS || task.Status == trainingtypes.TaskStatus_ASSIGNED) {
			tasksToRollback = append(tasksToRollback, task.Id)
		}
		return false
	})

	for _, taskID := range tasksToRollback {
		task, found := k.trainingKeeper.GetTask(ctx, taskID)
		if !found {
			continue
		}

		task.Status = trainingtypes.TaskStatus_ROLLBACK
		task.NodeID = ""
		k.trainingKeeper.SetTask(ctx, task)

		task.Status = trainingtypes.TaskStatus_PENDING
		k.trainingKeeper.SetTask(ctx, task)
	}

	return nil
}

func (k Keeper) ReassignTask(ctx sdk.Context, taskID string, newNodeID string) error {
	task, found := k.trainingKeeper.GetTask(ctx, taskID)
	if !found {
		return fmt.Errorf("task not found")
	}

	node, nodeFound := k.computeKeeper.GetNode(ctx, newNodeID)
	if !nodeFound {
		return fmt.Errorf("node not found")
	}

	if node.Status != "online" {
		return fmt.Errorf("node is not online")
	}

	task.NodeID = newNodeID
	task.Status = trainingtypes.TaskStatus_ASSIGNED
	k.trainingKeeper.SetTask(ctx, task)

	return nil
}

func (k Keeper) HandleNodeOffline(ctx sdk.Context, nodeID string) error {
	if err := k.RollbackTasksForNode(ctx, nodeID); err != nil {
		return err
	}

	var tasksToReassign []string
	k.trainingKeeper.IterateTasks(ctx, func(task trainingtypes.Task) (stop bool) {
		if task.Status == trainingtypes.TaskStatus_PENDING && task.NodeID == "" {
			tasksToReassign = append(tasksToReassign, task.Id)
		}
		return false
	})

	var availableNodes []string
	k.computeKeeper.IterateNodes(ctx, func(node computetypes.Node) (stop bool) {
		if node.Status == "online" {
			isHealthy, err := k.healthKeeper.CheckNodeHealth(ctx, node.Id)
			if err == nil && isHealthy {
				availableNodes = append(availableNodes, node.Id)
			}
		}
		return false
	})

	if len(availableNodes) > 0 {
		for i, taskID := range tasksToReassign {
			newNodeID := availableNodes[i%len(availableNodes)]
			if err := k.ReassignTask(ctx, taskID, newNodeID); err != nil {
				continue
			}
		}
	}

	return nil
}

