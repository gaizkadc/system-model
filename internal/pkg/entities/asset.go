/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package entities

import (
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-inventory-go"
	"time"
)


type OperatingSystemClass int
const (
	LINUX   = iota +1
	WINDOWS
	DARWIN
)

var OperatingSystemClassToGRPC = map[OperatingSystemClass]grpc_inventory_go.OperatingSystemClass{
	LINUX:    	grpc_inventory_go.OperatingSystemClass_LINUX,
	WINDOWS:	grpc_inventory_go.OperatingSystemClass_WINDOWS,
	DARWIN: 	grpc_inventory_go.OperatingSystemClass_DARWIN,
}
var OperatingSystemClassFromGRPC = map[grpc_inventory_go.OperatingSystemClass]OperatingSystemClass{
	grpc_inventory_go.OperatingSystemClass_LINUX: 	LINUX,
	grpc_inventory_go.OperatingSystemClass_WINDOWS: WINDOWS,
	grpc_inventory_go.OperatingSystemClass_DARWIN: 	DARWIN,
}
// OperatingSystemInfo contains information about the operating system of an asset. Notice that no
// enums have been used for either the name or the version as to no constraint the elements of the
// inventory even if we do not have an agent running supporting those.
type OperatingSystemInfo struct {
	// Name of the operating system. Expecting full name.
	Name 			string `json:"name,omitempty" cql:"name"`
	// Version installed.
	Version 		string   `json:"version,omitempty" cql:"version"`
	// Class of the operating system - determines the binary format together with architecture
	Class 			OperatingSystemClass `json:"class,omitempty" cql:"op_class"`
	// Architecture of the OS.
	Architecture	string   `json:"architecture,omitempty" cql:"architecture"`
}

func NewOperatingSystemInfoFromGRPC(osInfo *grpc_inventory_go.OperatingSystemInfo) *OperatingSystemInfo{
	if osInfo == nil{
		return nil
	}
	return &OperatingSystemInfo{
		Name:    osInfo.Name,
		Version: osInfo.Version,
		Class: 	 OperatingSystemClassFromGRPC[osInfo.Class],
		Architecture: osInfo.Architecture,
	}
}

func (os * OperatingSystemInfo) ToGRPC() *grpc_inventory_go.OperatingSystemInfo{
	if os == nil{
		return nil
	}
	return &grpc_inventory_go.OperatingSystemInfo{
		Name:           os.Name,
		Version:		os.Version,
		Class:			OperatingSystemClassToGRPC[os.Class],
		Architecture: 	os.Architecture,
	}
}

// HardareInfo contains information about the hardware of an asset.
type HardwareInfo struct {
	// CPUs contains the list of CPU available in a given asset.
	Cpus []*CPUInfo `json:"cpus,omitempty" cql:"cpus"`
	// InstalledRam contains the total RAM available in MB.
	InstalledRam int64 `json:"installed_ram,omitempty" cql:"installed_ram"`
	// NetInterfaces with the list of networking cards.
	NetInterfaces        []*NetworkingHardwareInfo `json:"net_interfaces,omitempty" cql:"net_interfaces"`
}

func NewHardwareInfoFromGRPC(hardwareInfo * grpc_inventory_go.HardwareInfo) * HardwareInfo{
	if hardwareInfo == nil{
		return nil
	}
	cpus := make([]*CPUInfo, 0)
	for _, info := range hardwareInfo.Cpus{
		cpus = append(cpus, NewCPUInfoFromGRPC(info))
	}
	netCards := make([]*NetworkingHardwareInfo, 0)
	for _, net:= range hardwareInfo.NetInterfaces{
		netCards = append(netCards, NewNetworkingHardwareInfoFromGRPC(net))
	}
	return &HardwareInfo{
		Cpus:          cpus,
		InstalledRam:  hardwareInfo.InstalledRam,
		NetInterfaces: netCards,
	}
}

func (hi * HardwareInfo) ToGRPC() *grpc_inventory_go.HardwareInfo{
	if hi == nil{
		return nil
	}
	cpus := make([]*grpc_inventory_go.CPUInfo, 0)
	for _, info := range hi.Cpus{
		cpus = append(cpus, info.ToGRPC())
	}
	netCards := make([]*grpc_inventory_go.NetworkingHardwareInfo, 0)
	for _, net := range hi.NetInterfaces{
		netCards = append(netCards, net.ToGRPC())
	}
	return &grpc_inventory_go.HardwareInfo{
		Cpus:                 cpus,
		InstalledRam:         hi.InstalledRam,
		NetInterfaces:        netCards,
	}
}

// CPUInfo contains information about a CPU.
type CPUInfo struct {
	// Manufacturer of the CPU.
	Manufacturer string `json:"manufacturer,omitempty" cql:"manufacturer"`
	// Model of the CPU.
	Model string `json:"model,omitempty" cql:"model"`
	// Architecture of the CPU.
	Architecture string `json:"architecture,omitempty" cql:"architecture"`
	// NumCores with the number of cores.
	NumCores             int32    `json:"num_cores,omitempty" cql:"num_cores"`
}

func NewCPUInfoFromGRPC(cpu * grpc_inventory_go.CPUInfo) * CPUInfo{
	if cpu == nil{
		return nil
	}
	return &CPUInfo{
		Manufacturer: cpu.Manufacturer,
		Model:       cpu.Model,
		Architecture: cpu.Architecture,
		NumCores:     cpu.NumCores,
	}
}

func (ci * CPUInfo) ToGRPC() *grpc_inventory_go.CPUInfo{
	if ci == nil{
		return nil
	}
	return &grpc_inventory_go.CPUInfo{
		Manufacturer:         ci.Manufacturer,
		Model:                ci.Model,
		Architecture:         ci.Architecture,
		NumCores:             ci.NumCores,
	}
}

// NetworkingHardwareInfo with the list of network interfaces that are available.
type NetworkingHardwareInfo struct {
	// Type of network interface. For example, ethernet, wifi, infiniband, etc.
	Type string `json:"type,omitempty" cql:"type"`
	// Link capacity in Mbps.
	LinkCapacity         int64    `json:"link_capacity,omitempty" cql:"link_capacity"`
}

func NewNetworkingHardwareInfoFromGRPC(netInfo * grpc_inventory_go.NetworkingHardwareInfo) * NetworkingHardwareInfo{
	if netInfo == nil{
		return nil
	}
	return &NetworkingHardwareInfo{
		Type:         netInfo.Type,
		LinkCapacity: netInfo.LinkCapacity,
	}
}

func (nhi * NetworkingHardwareInfo) ToGRPC() *grpc_inventory_go.NetworkingHardwareInfo{
	if nhi == nil{
		return nil
	}
	return &grpc_inventory_go.NetworkingHardwareInfo{
		Type:                 nhi.Type,
		LinkCapacity:         nhi.LinkCapacity,
	}
}

// StorageHardwareInfo with the storage related information.
type StorageHardwareInfo struct {
	// Type of storage.
	Type string `json:"type,omitempty" cql:"type"`
	// Total capacity in MB.
	TotalCapacity        int64    `json:"total_capacity,omitempty" cql:"total_capacity"`
}

func NewStorageHardwareInfoFromGRPC(storage * grpc_inventory_go.StorageHardwareInfo) * StorageHardwareInfo{
	if storage == nil{
		return nil
	}
	return &StorageHardwareInfo{
		Type:          storage.Type,
		TotalCapacity: storage.TotalCapacity,
	}
}

func (shi * StorageHardwareInfo) ToGRPC() *grpc_inventory_go.StorageHardwareInfo{
	if shi == nil{
		return nil
	}
	return &grpc_inventory_go.StorageHardwareInfo{
		Type:                 shi.Type,
		TotalCapacity:        shi.TotalCapacity,
	}
}

// Enumerate with the type of instances we can deploy in the system.
type AgentOpStatus int32

const (
	AgentOpStatusScheduled AgentOpStatus = iota + 1
	AgentOpStatusSuccess
	AgentOpStatusFail
)

var AgentOpStatusToGRPC = map[AgentOpStatus]grpc_inventory_go.AgentOpStatus{
	AgentOpStatusScheduled : grpc_inventory_go.AgentOpStatus_SCHEDULED,
	AgentOpStatusSuccess : grpc_inventory_go.AgentOpStatus_SUCCESS,
	AgentOpStatusFail : grpc_inventory_go.AgentOpStatus_FAIL,
}

var AgentOpStatusFromGRPC = map[grpc_inventory_go.AgentOpStatus]AgentOpStatus {
	grpc_inventory_go.AgentOpStatus_SCHEDULED:AgentOpStatusScheduled,
	grpc_inventory_go.AgentOpStatus_SUCCESS:AgentOpStatusSuccess,
	grpc_inventory_go.AgentOpStatus_FAIL:AgentOpStatusFail,
}

// AgentOpSummary contains the result of an asset operation
// this is a provisional result!
type AgentOpSummary struct {
	// OperationId with the operation identifier.
	OperationId string `json:"operation_id,omitempty" cql:"operation_id"`
	// Timestamp of the response.
	Timestamp int64 `json:"timestamp,omitempty" cql:"timestamp"`
	// Status indicates if the operation was successfull
	Status AgentOpStatus `json:"status,omitempty" cql:"status"`
	// Info with additional information for an operation.
	Info                 string   `json:"info,omitempty" cql:"info"`
}

func (a * AgentOpSummary) ToGRPC() *grpc_inventory_go.AgentOpSummary {
	if a == nil {
		return nil
	}
	return &grpc_inventory_go.AgentOpSummary{
		OperationId:a.OperationId,
		Timestamp:	a.Timestamp,
		Status: 	AgentOpStatusToGRPC[a.Status],
		Info: 		a.Info,
	}
}
func NewAgentOpSummaryFromGRPC(op *grpc_inventory_go.AgentOpSummary) *AgentOpSummary {
	return &AgentOpSummary{
		OperationId:op.OperationId,
		Timestamp:	op.Timestamp,
		Status: 	AgentOpStatusFromGRPC[op.Status],
		Info: 		op.Info,
	}
}

// Asset represents an element in the network from which we register some type of information. Example of
// assets could be workstations, nodes in a cluster, or other type of hardware.
type Asset struct {
	// OrganizationId with the organization identifier.
	OrganizationId string `json:"organization_id,omitempty"`
	// EdgeControllerId with the EIC identifier
	EdgeControllerId string `json:"edge_controller_id,omitempty"`
	// AssetId with the asset identifier.
	AssetId string `json:"asset_id,omitempty"`
	// AgentId with the agent identifier that is monitoring this asset if any.
	AgentId string `json:"agent_id,omitempty"`
	// Show flag to determine if this asset should be shown on the UI. This flag is internally used
	// for the async uninstall/removal of the asset.
	Show bool `json:"show,omitempty"`
	// Created time
	Created int64 `json:"created,omitempty"`
	// Labels defined by the user.
	Labels map[string]string `json:"labels,omitempty"`
	// OS contains Operating System information.
	Os *OperatingSystemInfo `json:"os,omitempty" cql:"os"`
	// Hardware information.
	Hardware *HardwareInfo `json:"hardware,omitempty" cql:"hardware"`
	// Storage information.
	Storage []StorageHardwareInfo `json:"storage,omitempty" cql:"storage"`
	// EicNetIp contains the current IP address that connects the asset to the EIC.
	EicNetIp             string   `json:"eic_net_ip,omitempty" cql:"eic_net_ip"`
	// AgentOpSummary contains the result of the last operation fr this asset
	LastOpResult *AgentOpSummary `json:"last_op_result,omitempty" cql: "eic_net_ip"`
	// LastAliveTimestamp contains the last alive message received
	LastAliveTimestamp   int64    `json:"last_alive_timestamp,omitempty" cql:"last_alive_timestamp"`
}

func NewAssetFromGRPC(addRequest * grpc_inventory_go.AddAssetRequest) *Asset{

	storage := make ([]StorageHardwareInfo, 0)
	for _, sto := range addRequest.Storage {
		storage = append(storage, * NewStorageHardwareInfoFromGRPC(sto) )
	}

	return &Asset{
		OrganizationId: addRequest.OrganizationId,
		EdgeControllerId: addRequest.EdgeControllerId,
		AssetId:        GenerateUUID(),
		AgentId:        addRequest.AgentId,
		Show:           true,
		Created:        time.Now().Unix(),
		Labels:         addRequest.Labels,
		Os:             NewOperatingSystemInfoFromGRPC(addRequest.Os),
		Hardware:       NewHardwareInfoFromGRPC(addRequest.Hardware),
		Storage:        storage,
	}
}

func (a * Asset) ToGRPC() *grpc_inventory_go.Asset{

	storage := make ([]*grpc_inventory_go.StorageHardwareInfo, 0)
	for _, sto := range a.Storage {
		storage = append(storage,sto.ToGRPC() )
	}

	return &grpc_inventory_go.Asset{
		OrganizationId:       a.OrganizationId,
		EdgeControllerId:     a.EdgeControllerId,
		AssetId:              a.AssetId,
		AgentId:              a.AgentId,
		Show:                 a.Show,
		Created:              a.Created,
		Labels:               a.Labels,
		Os:                   a.Os.ToGRPC(),
		Hardware:             a.Hardware.ToGRPC(),
		Storage:              storage,
		EicNetIp:             a.EicNetIp,
		LastAliveTimestamp:   a.LastAliveTimestamp,
		LastOpResult:         a.LastOpResult.ToGRPC(),
	}
}

func (a * Asset) ApplyUpdate(request * grpc_inventory_go.UpdateAssetRequest){
	if request.AddLabels {
		if a.Labels == nil {
			a.Labels = make(map[string]string, 0)
		}
		for k, v := range request.Labels {
			a.Labels[k] = v
		}
	}
	if request.RemoveLabels {
		for k, _ := range request.Labels {
			delete(a.Labels, k)
		}
	}
	if request.UpdateLastAlive {
		a.LastAliveTimestamp = request.LastAliveTimestamp
	}
	if request.UpdateLastOpSummary {
		a.LastOpResult = NewAgentOpSummaryFromGRPC(request.LastOpSummary)
	}
	if request.UpdateIp {
		a.EicNetIp = request.EicNetIp
	}
}

func ValidAddAssetRequest(addRequest * grpc_inventory_go.AddAssetRequest) derrors.Error{
	if addRequest.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if addRequest.EdgeControllerId == "" {
		return derrors.NewInvalidArgumentError(emptyEdgeControllerId)
	}
	return nil
}

func ValidAssetID(assetID * grpc_inventory_go.AssetId) derrors.Error{
	if assetID.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if assetID.AssetId == "" {
		return derrors.NewInvalidArgumentError(emptyAssetId)
	}

	return nil
}

func ValidUpdateAssetRequest(request *grpc_inventory_go.UpdateAssetRequest) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if request.AssetId == "" {
		return derrors.NewInvalidArgumentError(emptyAssetId)
	}
	return nil
}
