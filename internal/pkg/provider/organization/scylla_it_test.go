package organization

import (
	"github.com/nalej/system-model/internal/pkg/utils"
	"github.com/onsi/ginkgo"
	"github.com/rs/zerolog/log"
	"os"
)

/*
docker run --name scylla -p 9042:9042 -d scylladb/scylla
docker exec -it scylla nodetool status

docker exec -it scylla cqlsh

create KEYSPACE nalej WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};
use nalej;
create table nalej.organizations (id text, name text, created bigint, PRIMARY KEY (id));
create table nalej.Organization_Clusters (organization_id text, cluster_id text, PRIMARY KEY (organization_id, cluster_id));
create table nalej.Organization_Nodes (organization_id text, node_id text, PRIMARY KEY (organization_id, node_id));
create table nalej.Organization_AppDescriptors (organization_id text, app_descriptor_id text, PRIMARY KEY (organization_id, app_descriptor_id));
create table nalej.Organization_AppInstances (organization_id text, app_instance_id text, PRIMARY KEY (organization_id, app_instance_id));
create table nalej.Organization_Users (organization_id text, email text, PRIMARY KEY (organization_id, email));
create table nalej.Organization_Roles (organization_id text, role_id text, PRIMARY KEY (organization_id, role_id));
*/



var _ = ginkgo.Describe("Scylla organization provider", func() {

	if ! utils.RunIntegrationTests() {
		log.Warn().Msg("Integration tests are skipped")
		return
	}

	var scyllaHost= os.Getenv("IT_SCYLLA_HOST")
	if scyllaHost == "" {
		ginkgo.Fail("missing environment variables")
	}
	var nalejKeySpace= os.Getenv("IT_NALEJ_KEYSPACE")
	if nalejKeySpace == "" {
		ginkgo.Fail("missing environment variables")
	}

	// create a provider and connect it
	sp := NewScyllaOrganizationProvider(scyllaHost, nalejKeySpace)
	err := sp.Connect()
	if err != nil {
		ginkgo.Fail("unable to connect")
	}

	ginkgo.AfterSuite(func() {
		sp.Disconnect()
	})

	RunTest(sp)

})