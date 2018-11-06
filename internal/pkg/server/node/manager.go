/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package node

import (
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-infrastructure-go"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/nalej/system-model/internal/pkg/provider/cluster"
	"github.com/nalej/system-model/internal/pkg/provider/node"
	"github.com/nalej/system-model/internal/pkg/provider/organization"
)

// Manager structure with the required providers for node operations.
type Manager struct {
	OrgProvider organization.Provider
	ClusterProvider cluster.Provider
	NodeProvider node.Provider
}

// NewManager creates a Manager using a set of providers.
func NewManager(
	orgProvider organization.Provider,
	clusterProvider cluster.Provider,
	nodeProvider node.Provider) Manager {
	return Manager{orgProvider, clusterProvider, nodeProvider}
}

// AddNode adds a new node to the system.
func (m * Manager) AddNode(addNodeRequest *grpc_infrastructure_go.AddNodeRequest) (*entities.Node, derrors.Error) {
	exists := m.OrgProvider.Exists(addNodeRequest.OrganizationId)
	if !exists{
		return nil, derrors.NewNotFoundError("organizationID").WithParams(addNodeRequest.OrganizationId)
	}
	toAdd := entities.NewNodeFromGRPC(addNodeRequest)
	err := m.NodeProvider.Add(*toAdd)
	if err != nil {
		return nil, err
	}
	err = m.OrgProvider.AddNode(toAdd.OrganizationId, toAdd.NodeId)
	if err != nil {
		return nil, err
	}
	return toAdd, nil
}

func (m * Manager) UpdateNode(updateNodeRequest * grpc_infrastructure_go.UpdateNodeRequest) (*entities.Node, derrors.Error){
	exists := m.OrgProvider.Exists(updateNodeRequest.OrganizationId)
	if !exists{
		return nil, derrors.NewNotFoundError("organizationID").WithParams(updateNodeRequest.OrganizationId)
	}
	exists = m.OrgProvider.NodeExists(updateNodeRequest.OrganizationId, updateNodeRequest.NodeId)
	if !exists{
		return nil, derrors.NewNotFoundError("nodeID").WithParams(updateNodeRequest.NodeId)
	}
	old, err := m.NodeProvider.Get(updateNodeRequest.NodeId)
	if err != nil{
		return nil, err
	}
	old.ApplyUpdate(*updateNodeRequest)
	err = m.NodeProvider.Update(*old)
	if err != nil{
		return nil, err
	}
	return old, nil
}

// AttachNode links a node with a given cluster.
func (m * Manager) AttachNode(attachNodeRequest *grpc_infrastructure_go.AttachNodeRequest) derrors.Error {
	exists := m.OrgProvider.Exists(attachNodeRequest.OrganizationId)
	if !exists{
		return derrors.NewNotFoundError("organizationID").WithParams(attachNodeRequest.OrganizationId)
	}
	exists = m.OrgProvider.ClusterExists(attachNodeRequest.OrganizationId, attachNodeRequest.ClusterId)
	if !exists{
		return derrors.NewNotFoundError("clusterID").WithParams(attachNodeRequest.ClusterId)
	}
	exists = m.OrgProvider.NodeExists(attachNodeRequest.OrganizationId, attachNodeRequest.NodeId)
	if !exists{
		return derrors.NewNotFoundError("nodeID").WithParams(attachNodeRequest.NodeId)
	}
	retrieved, err := m.NodeProvider.Get(attachNodeRequest.NodeId)
	if err != nil {
		return err
	}
	err = m.ClusterProvider.AddNode(attachNodeRequest.ClusterId, attachNodeRequest.NodeId)
	if err != nil {
		return err
	}
	retrieved.ClusterId = attachNodeRequest.ClusterId
	err = m.NodeProvider.Update(*retrieved)
	if err != nil {
		return err
	}
	return nil
}

// ListNodes obtains a list of nodes in a cluster.
func (m * Manager) ListNodes(clusterID *grpc_infrastructure_go.ClusterId) ([] entities.Node, derrors.Error) {
	if !m.OrgProvider.Exists(clusterID.OrganizationId){
		return nil, derrors.NewNotFoundError("organizationID").WithParams(clusterID.OrganizationId)
	}
	nodes, err := m.ClusterProvider.ListNodes(clusterID.ClusterId)
	if err != nil {
		return nil, err
	}
	result := make([] entities.Node, 0)
	for _, nID := range nodes {
		toAdd, err := m.NodeProvider.Get(nID)
		if err != nil {
			return nil, err
		}
		result = append(result, *toAdd)
	}
	return result, nil
}

// RemoveNodes removes a set of nodes from the system.
func (m * Manager) RemoveNodes(removeNodesRequest *grpc_infrastructure_go.RemoveNodesRequest) derrors.Error {
	if ! m.OrgProvider.Exists(removeNodesRequest.OrganizationId){
		return derrors.NewNotFoundError("organizationID").WithParams(removeNodesRequest.OrganizationId)
	}

	for _, nID := range removeNodesRequest.Nodes {
		node, err := m.NodeProvider.Get(nID)
		if err != nil {
			return derrors.NewNotFoundError("nodeID").WithParams(nID)
		}
		if node.ClusterId != "" {
			err := m.ClusterProvider.DeleteNode(node.ClusterId, node.NodeId)
			if err != nil {
				return derrors.NewInternalError("cannot delete node from cluster").CausedBy(err).WithParams(node.ClusterId, node.NodeId)
			}
		}
		err = m.OrgProvider.DeleteNode(node.OrganizationId, node.NodeId)
		if err != nil {
			return derrors.NewInternalError("cannot delete node from organization").CausedBy(err).WithParams(node.OrganizationId, node.NodeId)
		}
		err = m.NodeProvider.Remove(node.NodeId)
		if err != nil {
			return err
		}
	}

	return nil
}