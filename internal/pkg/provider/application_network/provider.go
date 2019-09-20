/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package application_network

import (
	"github.com/nalej/derrors"
	"github.com/nalej/system-model/internal/pkg/entities"
)

// Provider for the application networking instances.
type Provider interface {
	// AddConnectionInstance Adds a new connection between applications.
	AddConnectionInstance(connectionInstance entities.ConnectionInstance) derrors.Error
	// ExistsConnectionInstance Checks if the connection instance exists on the system.
	ExistsConnectionInstance(organizationId string, sourceInstanceId string, targetInstanceId string, inboundName string, outboundName string) (bool, derrors.Error)
	// GetConnectionInstance Retrieve the connection instance using organizationId, sourceInstanceId, and targetInstanceId.
	GetConnectionInstance(organizationId string, sourceInstanceId string, targetInstanceId string, inboundName string, outboundName string) (*entities.ConnectionInstance, derrors.Error)
	// ListConnectionInstances Lists all the connection instances.
	ListConnectionInstances(organizationId string) ([]entities.ConnectionInstance, derrors.Error)
	// RemoveConnectionInstance Removes a connection from the system
	RemoveConnectionInstance(organizationId string, sourceInstanceId string, targetInstanceId string, inboundName string, outboundName string) derrors.Error

	// AddConnectionInstanceLink Adds a new connection between applications.
	AddConnectionInstanceLink(connectionInstanceLink entities.ConnectionInstanceLink) derrors.Error
	// ExistsConnectionInstanceLink Checks if the connection instance exists on the system.
	ExistsConnectionInstanceLink(organizationId string, sourceInstanceId string, targetInstanceId string, sourceClusterId string, targetClusterId string, inboundName string, outboundName string) (bool, derrors.Error)
	// GetConnectionInstanceLink Retrieve the connection instance.
	GetConnectionInstanceLink(organizationId string, sourceInstanceId string, targetInstanceId string, sourceClusterId string, targetClusterId string, inboundName string, outboundName string) (*entities.ConnectionInstanceLink, derrors.Error)
	// ListConnectionInstanceLinks Lists all the connection instance links of one connection instance.
	ListConnectionInstanceLinks(organizationId string, sourceInstanceId string, targetInstanceId string, inboundName string, outboundName string) ([]entities.ConnectionInstanceLink, derrors.Error)
	// RemoveConnectionInstanceLinks Removes all connection links from a connection instance.
	RemoveConnectionInstanceLinks(organizationId string, sourceInstanceId string, targetInstanceId string, inboundName string, outboundName string) derrors.Error

	// clear the connections information
	Clear() derrors.Error
}