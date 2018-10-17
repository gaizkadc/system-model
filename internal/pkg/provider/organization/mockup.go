/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package organization

import (
	"github.com/nalej/derrors"
	"github.com/nalej/system-model/internal/pkg/entities"
	"sync"
)

type MockupOrganizationProvider struct {
	sync.Mutex
	// organizations contains the organization indexed per organization identifier.
	organizations map[string] entities.Organization
	// clusters attached to an organization.
	clusters map[string][]string
	// nodes attached to an organization
	nodes map[string][]string
	// Descriptors contains the application descriptors ids per organization.
	descriptors map[string][]string
	// Instances contains the application instances ids per organization.
	instances map[string][]string
}


func NewMockupOrganizationProvider() * MockupOrganizationProvider {
	return &MockupOrganizationProvider{
		organizations:make(map[string]entities.Organization, 0),
		clusters:make(map[string][]string, 0),
		nodes:make(map[string][]string, 0),
		descriptors:make(map[string][]string, 0),
		instances:make(map[string][]string, 0),
	}
}

// Clear cleans the contents of the mockup.
func (m * MockupOrganizationProvider) Clear() {
	m.Lock()
	m.organizations = make(map[string] entities.Organization, 0)
	m.clusters = make(map[string] []string, 0)
	m.nodes = make(map[string] []string, 0)
	m.descriptors = make(map[string] []string, 0)
	m.instances = make(map[string] []string, 0)
	m.Unlock()
}

func (m *MockupOrganizationProvider) unsafeExists(organizationID string) bool {
	_, exists := m.organizations[organizationID]
	return exists
}

func (m *MockupOrganizationProvider) unsafeExistsAppDesc(organizationID string, descriptorID string) bool {
	descriptors, ok := m.descriptors[organizationID]
	if ok {
		for _, descriptor := range descriptors {
			if descriptor == descriptorID {
				return true
			}
		}
		return false
	}
	return false
}

func (m *MockupOrganizationProvider) unsafeExistsAppInst(organizationID string, instanceID string) bool {
	instances, ok := m.instances[organizationID]
	if ok {
		for _, inst := range instances {
			if inst == instanceID {
				return true
			}
		}
		return false
	}
	return false
}

func (m *MockupOrganizationProvider) unsafeExistsCluster(organizationID string, clusterID string) bool {
	clusterList, ok := m.clusters[organizationID]
	if ok {
		for _, cID := range clusterList {
			if cID == clusterID {
				return true
			}
		}
		return false
	}
	return false
}

func (m *MockupOrganizationProvider) unsafeExistsNode(organizationID string, nodeID string) bool {
	nodeList, ok := m.nodes[organizationID]
	if ok {
		for _, nID := range  nodeList {
			if nID == nodeID {
				return true
			}
		}
		return false
	}
	return false
}

// Add a new organization to the system.
func (m *MockupOrganizationProvider) Add(org entities.Organization) derrors.Error {
	m.Lock()
	defer m.Unlock()
	if !m.unsafeExists(org.ID){
		m.organizations[org.ID] = org
		return nil
	}
	return derrors.NewAlreadyExistsError(org.ID)
}

// Check if an organization exists on the system.
func (m *MockupOrganizationProvider) Exists(organizationID string) bool {
	m.Lock()
	defer m.Unlock()
	return m.unsafeExists(organizationID)
}

// Get an organization.
func (m *MockupOrganizationProvider) Get(organizationID string) (*entities.Organization, derrors.Error) {
	m.Lock()
	defer m.Unlock()
	org, exists := m.organizations[organizationID]
	if exists {
		return &org, nil
	}
	return nil, derrors.NewNotFoundError(organizationID)
}

// AddCluster adds a new cluster ID to the organization.
func (m *MockupOrganizationProvider) AddCluster(organizationID string, clusterID string) derrors.Error {
	m.Lock()
	defer m.Unlock()
	if m.unsafeExists(organizationID) {
		if !m.unsafeExistsCluster(organizationID, clusterID) {
			clusterList, _ := m.clusters[organizationID]
			m.clusters[organizationID] = append(clusterList, clusterID)
			return nil
		}
		return derrors.NewAlreadyExistsError("cluster").WithParams(organizationID, clusterID)
	}
	return derrors.NewNotFoundError("organization").WithParams(organizationID)
}

// ClusterExists checks if a cluster is linked to an organization.
func (m *MockupOrganizationProvider) ClusterExists(organizationID string, clusterID string) bool {
	m.Lock()
	defer m.Unlock()
	return m.unsafeExistsCluster(organizationID, clusterID)
}

// ListClusters returns a list of clusters in an organization.
func (m *MockupOrganizationProvider) ListClusters(organizationID string) ([]string, derrors.Error) {
	m.Lock()
	defer m.Unlock()

	if !m.unsafeExists(organizationID) {
		return nil, derrors.NewNotFoundError("organization").WithParams(organizationID)
	}

	clusterList, ok := m.clusters[organizationID]
	if ok {
		return clusterList, nil
	}
	return make([]string, 0), nil
}

// DeleteCluster removes a cluster from an organization.
func (m *MockupOrganizationProvider) DeleteCluster(organizationID string, clusterID string) derrors.Error {
	m.Lock()
	defer m.Unlock()
	if m.unsafeExistsCluster(organizationID, clusterID) {
		previous := m.clusters[organizationID]
		newList := make([] string, 0, len(previous)-1)
		for _, id := range previous {
			if id != clusterID {
				newList = append(newList, id)
			}
		}
		m.clusters[organizationID] = newList
		return nil
	}
	return derrors.NewNotFoundError("cluster").WithParams(organizationID, clusterID)
}

// AddNode adds a new node ID to the organization.
func (m *MockupOrganizationProvider) AddNode(organizationID string, nodeID string) derrors.Error {
	m.Lock()
	defer m.Unlock()
	if m.unsafeExists(organizationID) {
		if !m.unsafeExistsNode(organizationID, nodeID) {
			nodeList, _ := m.nodes[organizationID]
			m.nodes[organizationID] = append(nodeList, nodeID)
			return nil
		}
		return derrors.NewAlreadyExistsError("node").WithParams(organizationID, nodeID)
	}
	return derrors.NewNotFoundError("organization").WithParams(organizationID)
}

// NodeExists checks if a node is linked to an organization.
func (m *MockupOrganizationProvider) NodeExists(organizationID string, nodeID string) bool {
	m.Lock()
	defer m.Unlock()
	return m.unsafeExistsNode(organizationID, nodeID)
}

// ListNodes returns a list of nodes in an organization.
func (m *MockupOrganizationProvider) ListNodes(organizationID string) ([]string, derrors.Error) {
	m.Lock()
	defer m.Unlock()

	if !m.unsafeExists(organizationID) {
		return nil, derrors.NewNotFoundError("organization").WithParams(organizationID)
	}

	nodeList, ok := m.nodes[organizationID]
	if ok {
		return nodeList, nil
	}
	return make([]string, 0), nil
}

// DeleteNode removes a node from an organization.
func (m *MockupOrganizationProvider) DeleteNode(organizationID string, nodeID string) derrors.Error {
	m.Lock()
	defer m.Unlock()
	if m.unsafeExistsNode(organizationID, nodeID) {
		previous := m.nodes[organizationID]
		newList := make([] string, 0, len(previous)-1)
		for _, id := range previous {
			if id != nodeID {
				newList = append(newList, id)
			}
		}
		m.nodes[organizationID] = newList
		return nil
	}
	return derrors.NewNotFoundError("node").WithParams(organizationID, nodeID)
}

// AddDescriptor adds a new descriptor ID to a given organization.
func (m *MockupOrganizationProvider) AddDescriptor(organizationID string, appDescriptorID string) derrors.Error {
	m.Lock()
	defer m.Unlock()
	if m.unsafeExists(organizationID) {
		if !m.unsafeExistsAppDesc(organizationID, appDescriptorID) {
			descriptors, _ := m.descriptors[organizationID]
			m.descriptors[organizationID] = append(descriptors, appDescriptorID)
			return nil
		}
		return derrors.NewAlreadyExistsError("descriptor").WithParams(organizationID, appDescriptorID)
	}
	return derrors.NewNotFoundError("organization").WithParams(organizationID)
}

// DescriptorExists checks if an application descriptor exists on the system.
func (m *MockupOrganizationProvider) DescriptorExists(organizationID string, appDescriptorID string) bool {
	m.Lock()
	defer m.Unlock()
	return m.unsafeExistsAppDesc(organizationID, appDescriptorID)
}

// ListDescriptors returns the identifiers of the application descriptors associated with an organization.
func (m *MockupOrganizationProvider) ListDescriptors(organizationID string) ([]string, derrors.Error) {
	m.Lock()
	defer m.Unlock()

	if !m.unsafeExists(organizationID) {
		return nil, derrors.NewNotFoundError("organization").WithParams(organizationID)
	}

	descriptors, ok := m.descriptors[organizationID]
	if ok {
		return descriptors, nil
	}
	return make([]string, 0), nil
}

// DeleteDescriptor removes a descriptor from an organization
func (m *MockupOrganizationProvider) DeleteDescriptor(organizationID string, appDescriptorID string) derrors.Error {
	m.Lock()
	defer m.Unlock()
	if m.unsafeExistsAppDesc(organizationID, appDescriptorID) {
		previous := m.descriptors[organizationID]
		newList := make([] string, 0, len(previous)-1)
		for _, id := range previous {
			if id != appDescriptorID {
				newList = append(newList, id)
			}
		}
		m.descriptors[organizationID] = newList
		return nil
	}
	return derrors.NewNotFoundError("descriptor").WithParams(organizationID, appDescriptorID)
}

// AddInstance adds a new application instance ID to a given organization.
func (m *MockupOrganizationProvider) AddInstance(organizationID string, appInstanceID string) derrors.Error {
	m.Lock()
	defer m.Unlock()
	if m.unsafeExists(organizationID) {
		if !m.unsafeExistsAppInst(organizationID, appInstanceID) {
			instances, _ := m.instances[organizationID]
			m.instances[organizationID] = append(instances, appInstanceID)
			return nil
		}
		return derrors.NewAlreadyExistsError("instance").WithParams(organizationID, appInstanceID)
	}
	return derrors.NewNotFoundError("organization").WithParams(organizationID)
}

// InstanceExists checks if an application instance exists on the system.
func (m *MockupOrganizationProvider) InstanceExists(organizationID string, appInstanceID string) bool {
	m.Lock()
	defer m.Unlock()
	return m.unsafeExistsAppInst(organizationID, appInstanceID)
}

// ListInstances returns a the identifiers associate with a given organization.
func (m *MockupOrganizationProvider) ListInstances(organizationID string) ([]string, derrors.Error) {
	m.Lock()
	defer m.Unlock()

	if !m.unsafeExists(organizationID) {
		return nil, derrors.NewNotFoundError("organization").WithParams(organizationID)
	}

	instances, ok := m.instances[organizationID]
	if ok {
		return instances, nil
	}
	return make([]string, 0), nil
}

// DeleteInstance removes an instance from an organization
func (m *MockupOrganizationProvider) DeleteInstance(organizationID string, appInstanceID string) derrors.Error {
	m.Lock()
	defer m.Unlock()
	if m.unsafeExistsAppDesc(organizationID, appInstanceID) {
		previous := m.instances[organizationID]
		newList := make([] string, 0, len(previous)-1)
		for _, id := range previous {
			if id != appInstanceID {
				newList = append(newList, id)
			}
		}
		m.instances[organizationID] = newList
		return nil
	}
	return derrors.NewNotFoundError("instance").WithParams(organizationID, appInstanceID)
}