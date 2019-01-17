package device

import (
	"github.com/nalej/system-model/internal/pkg/utils"
	"github.com/onsi/ginkgo"
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
)

/*
docker run --name scylla -p 9042:9042 -d scylladb/scylla
docker exec -it scylla cqlsh

create KEYSPACE nalej WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};
use nalej;

create table IF NOT EXISTS nalej.Devices (organization_id text, device_group_id text, device_id text, register_since int, labels map<text, text>, PRIMARY KEY ( (organization_id, device_group_id), device_id));
create table IF NOT EXISTS nalej.DeviceGroups (organization_id text, device_group_id text, name text, created int, labels map<text, text>, primary KEY (organization_id, device_group_id));
*/

var _ = ginkgo.Describe("Scylla cluster provider", func() {

	if ! utils.RunIntegrationTests() {
		log.Warn().Msg("Integration tests are skipped")
		return
	}
	var scyllaHost = os.Getenv("IT_SCYLLA_HOST")
	if scyllaHost == "" {
		ginkgo.Fail("missing environment variables")
	}
	var nalejKeySpace = os.Getenv("IT_NALEJ_KEYSPACE")
	if scyllaHost == "" {
		ginkgo.Fail("missing environment variables")
	}
	scyllaPort, _ := strconv.Atoi(os.Getenv("IT_SCYLLA_PORT"))
	if scyllaPort <= 0 {
		ginkgo.Fail("missing environment variables")
	}

	// create a provider and connect it
	sp := NewScyllaDeviceProvider(scyllaHost, scyllaPort, nalejKeySpace)

	ginkgo.AfterSuite(func() {
		sp.Disconnect()
	})

	RunTest(sp)
})
