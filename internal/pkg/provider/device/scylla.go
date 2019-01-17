package device

import (
	"github.com/gocql/gocql"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/system-model/internal/pkg/entities/device"
	"github.com/rs/zerolog/log"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
	"sync"
)

// table and field names
const (
	deviceGroupTable = "DeviceGroups"
	deviceTable = "Devices"
	organizationIdField = "organization_id"
	deviceGroupIdField = "device_group_id"
	deviceIdField = "device_id"
	labelsField  = "labels"
	registerSinceField = "register_since"
	rowNotFound = "not found"
)

type ScyllaDeviceProvider struct {
	Address string
	Port int
	Keyspace string
	Session *gocql.Session
	sync.Mutex
}

func NewScyllaDeviceProvider (address string, port int, keyspace string) * ScyllaDeviceProvider {
	provider := ScyllaDeviceProvider{Address: address, Port:  port, Keyspace: keyspace, Session: nil}
	provider.connect()
	return &provider

}

// connect to the database
func (sp *ScyllaDeviceProvider) connect() derrors.Error {

	// connect to the cluster
	conf := gocql.NewCluster(sp.Address)
	conf.Keyspace = sp.Keyspace
	conf.Port = sp.Port

	session, err := conf.CreateSession()
	if err != nil {
		log.Error().Str("provider", "ScyllaDeviceProvider").Str("trace", conversions.ToDerror(err).DebugReport()).Msg("unable to connect")
		return derrors.AsError(err, "cannot connect")
	}

	sp.Session = session

	return nil
}
// disconnect from the database
func (sp *ScyllaDeviceProvider) Disconnect () {

	sp.Lock()
	defer sp.Unlock()

	if sp.Session != nil {
		sp.Session.Close()
		sp.Session = nil
	}
}
func (sp *ScyllaDeviceProvider) checkAndConnect () derrors.Error{

	if sp.Session == nil{
		log.Info().Msg("session no created, trying to reconnect...")
		// try to reconnect
		err := sp.connect()
		if err != nil  {
			return err
		}
	}
	return nil
}

// -------------------------------------------------------------------------------------------------------------------


func (sp * ScyllaDeviceProvider) unsafeExistsGroup (organizationID string, deviceGroupID string) (bool, derrors.Error) {

	var returnedId string

	if err := sp.checkAndConnect(); err != nil{
		return false, err
	}

	stmt, names := qb.Select(deviceGroupTable).Columns(organizationIdField).Where(qb.Eq(organizationIdField)).
		Where(qb.Eq(deviceGroupIdField)).ToCql()

	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		organizationIdField: organizationID,
		deviceGroupIdField:deviceGroupID})

	err := q.GetRelease(&returnedId)
	if err != nil {
		if err.Error() == rowNotFound {
			return false, nil
		}else{
			return false, derrors.AsError(err, "cannot determinate if device group exists")
		}
	}

	return true, nil
}
// AddDeviceGroup adds a new device group
func (sp *ScyllaDeviceProvider) AddDeviceGroup (deviceGroup device.DeviceGroup) derrors.Error {

	sp.Lock()
	defer sp.Unlock()

	// check connection
	if err := sp.checkAndConnect(); err != nil {
		return err
	}

	// check if the group already exists
	exists, err := sp.unsafeExistsGroup(deviceGroup.OrganizationId, deviceGroup.DeviceGroupId)
	if err != nil {
		return err
	}
	if exists {
		return derrors.NewAlreadyExistsError("Add device group").WithParams(deviceGroup.OrganizationId, deviceGroup.DeviceGroupId)
	}
	// add it into database
	stmt, names := qb.Insert(deviceGroupTable).Columns("organization_id",
		"device_group_id", "name", "created", "labels").ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindStruct(deviceGroup)
	cqlErr := q.ExecRelease()

	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot add deviceGroup")
	}

	return nil

}

// ExistsDeviceGroup checks if a group exists on the system.
func (sp *ScyllaDeviceProvider) ExistsDeviceGroup(organizationID string, deviceGroupID string) (bool, derrors.Error) {
	if err := sp.checkAndConnect(); err != nil{
		return false, err
	}

	var returnedId string

	stmt, names := qb.Select(deviceGroupTable).Columns(organizationIdField).Where(qb.Eq(organizationIdField)).
		Where(qb.Eq(deviceGroupIdField)).ToCql()

	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		organizationIdField: organizationID,
		deviceGroupIdField:deviceGroupID})

	err := q.GetRelease(&returnedId)
	if err != nil {
		if err.Error() == rowNotFound {
			return false, nil
		}else{
			return false, derrors.AsError(err, "cannot determinate if device group exists")
		}
	}

	return true, nil

}
// GetDeviceGroup returns a device Group.
func (sp *ScyllaDeviceProvider) GetDeviceGroup(organizationID string, deviceGroupID string) (* device.DeviceGroup, derrors.Error) {

	sp.Lock()
	defer sp.Unlock()

	if err := sp.checkAndConnect(); err != nil{
		return nil, err
	}

	var deviceGroup device.DeviceGroup

	stmt, names := qb.Select(deviceGroupTable).Where(qb.Eq(organizationIdField)).
		Where(qb.Eq(deviceGroupIdField)).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		organizationIdField: organizationID,
		deviceGroupIdField:deviceGroupID})

	err := q.GetRelease(&deviceGroup)
	if err != nil {
		if err.Error() == rowNotFound {
			return nil, derrors.NewNotFoundError("device group").WithParams(organizationID, deviceGroupID)
		} else {
			return nil, derrors.AsError(err, "cannot get device group")
		}
	}

	return &deviceGroup, nil

}
// ListDeviceGroups returns a list of device groups in a organization.
func (sp *ScyllaDeviceProvider) ListDeviceGroups(organizationID string) ([]device.DeviceGroup, derrors.Error) {

	sp.Lock()
	defer sp.Unlock()

	if err := sp.checkAndConnect(); err != nil{
		return nil, err
	}

	stmt, names := qb.Select(deviceGroupTable).Columns("organization_id", "device_group_id",
		"created", "labels", "name").Where(qb.Eq("organization_id")).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"organization_id": organizationID,
	})

	groups := make([]device.DeviceGroup, 0)
	cqlErr := gocqlx.Select(&groups, q.Query)

	if cqlErr != nil {
		return nil, derrors.AsError(cqlErr, "cannot list device groups of an organization")
	}

	return groups, nil

}
// Remove a device group
func (sp *ScyllaDeviceProvider) RemoveDeviceGroup(organizationID string, deviceGroupID string) derrors.Error {
	sp.Lock()
	defer sp.Unlock()

	if err := sp.checkAndConnect(); err != nil{
		return err
	}

	// check if the group exists
	exists, err := sp.unsafeExistsGroup(organizationID, deviceGroupID)
	if err != nil {
		return err
	}
	if !exists {
		return derrors.NewNotFoundError("device group").WithParams(organizationID, deviceGroupID)
	}

	stmt, _ := qb.Delete(deviceGroupTable).Where(qb.Eq(organizationIdField)).Where(qb.Eq(deviceGroupIdField)).ToCql()
	cqlErr := sp.Session.Query(stmt, organizationID, deviceGroupID).Exec()

	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot delete device group")
	}

	return nil

}

// -------------------------------------------------------------------------------------------------------------------

func (sp * ScyllaDeviceProvider) unsafeExistsDevice (organizationID string, deviceGroupID string, deviceID string) (bool, derrors.Error) {

	var returnedId string

	if err := sp.checkAndConnect(); err != nil{
		return false, err
	}

	stmt, names := qb.Select(deviceTable).Columns(organizationIdField).
		Where(qb.Eq(organizationIdField)).
		Where(qb.Eq(deviceGroupIdField)).
		Where(qb.Eq(deviceIdField)).ToCql()

	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		organizationIdField: organizationID,
		deviceGroupIdField:deviceGroupID,
		deviceIdField:deviceID})

	err := q.GetRelease(&returnedId)
	if err != nil {
		if err.Error() == rowNotFound {
			return false, nil
		}else{
			return false, derrors.AsError(err, "cannot determinate if device exists")
		}
	}

	return true, nil
}

// AddDevice adds a new device group
func (sp *ScyllaDeviceProvider) AddDevice (device device.Device) derrors.Error {

	sp.Lock()
	defer sp.Unlock()

	// check connection
	if err := sp.checkAndConnect(); err != nil {
		return err
	}

	// check if the group already exists
	exists, err := sp.unsafeExistsDevice(device.OrganizationId, device.DeviceGroupId, device.DeviceId)
	if err != nil {
		return err
	}
	if exists {
		return derrors.NewAlreadyExistsError("Add device ").WithParams(device.OrganizationId, device.DeviceGroupId, device.DeviceId)
	}
	// add it into database
	stmt, names := qb.Insert(deviceTable).Columns(organizationIdField, deviceGroupIdField, deviceIdField,
		labelsField, registerSinceField).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindStruct(device)
	cqlErr := q.ExecRelease()

	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot add device")
	}

	return nil

}
// ExistsDevice checks if a device exists on the system.
func (sp *ScyllaDeviceProvider) ExistsDevice(organizationID string, deviceGroupID string, deviceID string) (bool, derrors.Error) {

	sp.Lock()
	defer sp.Unlock()

	var returnedId string

	if err := sp.checkAndConnect(); err != nil{
		return false, err
	}

	stmt, names := qb.Select(deviceTable).Columns(organizationIdField).
		Where(qb.Eq(organizationIdField)).
		Where(qb.Eq(deviceGroupIdField)).
		Where(qb.Eq(deviceIdField)).ToCql()

	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		organizationIdField: organizationID,
		deviceGroupIdField:deviceGroupID,
		deviceIdField:deviceID})

	err := q.GetRelease(&returnedId)
	if err != nil {
		if err.Error() == rowNotFound {
			return false, nil
		}else{
			return false, derrors.AsError(err, "cannot determinate if device exists")
		}
	}

	return true, nil
}
// GetDevice returns a device .
func (sp *ScyllaDeviceProvider) GetDevice(organizationID string, deviceGroupID string, deviceID string) (* device.Device, derrors.Error) {

	sp.Lock()
	defer sp.Unlock()

	if err := sp.checkAndConnect(); err != nil{
		return nil, err
	}

	var device device.Device

	stmt, names := qb.Select(deviceTable).Where(qb.Eq(organizationIdField)).
		Where(qb.Eq(deviceGroupIdField)).Where(qb.Eq(deviceIdField)).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		organizationIdField: organizationID,
		deviceGroupIdField:deviceGroupID,
		deviceIdField:deviceID,
	})

	err := q.GetRelease(&device)
	if err != nil {
		if err.Error() == rowNotFound {
			return nil, derrors.NewNotFoundError("device").WithParams(organizationID, deviceGroupID, deviceID)
		} else {
			return nil, derrors.AsError(err, "cannot get device")
		}
	}

	return &device, nil
}
// ListDevice returns a list of device in a group.
func (sp *ScyllaDeviceProvider) ListDevice(organizationID string, deviceGroupID string) ([]device.Device, derrors.Error) {

	sp.Lock()
	defer sp.Unlock()

	if err := sp.checkAndConnect(); err != nil{
		return nil, err
	}

	stmt, names := qb.Select(deviceTable).Columns(organizationIdField, deviceGroupIdField, deviceIdField,
		labelsField, registerSinceField).Where(qb.Eq(organizationIdField)).
		Where(qb.Eq(deviceGroupIdField)).ToCql()

	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		organizationIdField: organizationID,
		deviceGroupIdField: deviceGroupID,
	})

	devices := make([]device.Device, 0)
	cqlErr := gocqlx.Select(&devices, q.Query)

	if cqlErr != nil {
		return nil, derrors.AsError(cqlErr, "cannot list devices of a group")
	}
	if len(devices) == 0 {
		return nil, derrors.NewNotFoundError("unable to recover the devices list").WithParams(organizationID, deviceGroupID)
	}

	return devices, nil
}
// Remove a device
func (sp *ScyllaDeviceProvider) RemoveDevice(organizationID string, deviceGroupID string, deviceID string) derrors.Error {

	sp.Lock()
	defer sp.Unlock()

	if err := sp.checkAndConnect(); err != nil{
		return err
	}

	// check if the device exists
	exists, err := sp.unsafeExistsDevice(organizationID, deviceGroupID, deviceID)
	if err != nil {
		return err
	}
	if !exists {
		return derrors.NewNotFoundError("device").WithParams(organizationID, deviceGroupID, deviceID)
	}

	stmt, _ := qb.Delete(deviceTable).
		Where(qb.Eq(organizationIdField)).
		Where(qb.Eq(deviceGroupIdField)).
		Where(qb.Eq(deviceIdField)).ToCql()
	cqlErr := sp.Session.Query(stmt, organizationID, deviceGroupID, deviceID).Exec()

	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot delete device group")
	}

	return nil
}
