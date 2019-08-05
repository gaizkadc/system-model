package user

import (
	"fmt"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/nalej/system-model/internal/pkg/utils"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
)
/*
IT_SCYLLA_HOST=127.0.0.1
IT_NALEJ_KEYSPACE=nalej
RUN_INTEGRATION_TEST=true
*/
/*
before past test:

docker run --name scylla -p 9042:9042 -d scylladb/scylla
docker exec -it scylla nodetool status

Prepare the database...

docker exec -it scylla cqlsh
create KEYSPACE nalej WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};
create type IF NOT EXISTS nalej.UserContactInfo (full_name text, address text, phone map<text, text>, alt_email text, company_name text, title text);
create table IF NOT EXISTS nalej.Users (organization_id text, email text, name text, photo_url text, member_since int, contact_info FROZEN<UserContactInfo>, PRIMARY KEY (email));
create table IF NOT EXISTS nalej.AccountUser (account_id text, email text, role_id text, internal boolean, status int, PRIMARY KEY (email, account_id));
create table IF NOT EXISTS nalej.AccountUserInvite (account_id text, email text, role_id text, invited_by text, msg text, expires bigint, PRIMARY KEY (email, account_id));
create table IF NOT EXISTS nalej.ProjectUser (account_id text, project_id text, email text, role_id text, internal boolean, status int, PRIMARY KEY (account_id, project_id, email));
create table IF NOT EXISTS nalej.ProjectUserInvite (account_id text, project_id text, email text, role_id text, invited_by text, msg text, expires bigint, PRIMARY KEY (email, account_id, project_id));

 */


var _ = ginkgo.Describe("Scylla user provider", func(){

	var numUsers = 50

	if ! utils.RunIntegrationTests() {
		log.Warn().Msg("Integration tests are skipped")
		return
	}

	var scyllaHost = os.Getenv("IT_SCYLLA_HOST")
	if scyllaHost == "" {
		ginkgo.Fail("missing environment variables")
	}

	scyllaPort, _  := strconv.Atoi(os.Getenv("IT_SCYLLA_PORT"))

	if scyllaPort <= 0 {
		ginkgo.Fail("missing environment variables")
	}

	var nalejKeySpace = os.Getenv("IT_NALEJ_KEYSPACE")
	if nalejKeySpace == "" {
		ginkgo.Fail("missing environment variables")

	}

	// create a provider and connect it
	sp := NewScyllaUserProvider(scyllaHost, scyllaPort, nalejKeySpace)

	// disconnect
	ginkgo.AfterSuite(func() {
		sp.Disconnect()
	})

	RunTest(sp)

	ginkgo.It("Should be able to add user", func(){

		for i := 0; i < numUsers; i++ {
			email := fmt.Sprintf("email_%d@company.org", i)
			user := &entities.User{OrganizationId: fmt.Sprintf("org_%d", i),
				Email:       email,
				Name:        fmt.Sprintf("name_%d", i),
				MemberSince: int64(i)}
			err := sp.Add(*user)
			gomega.Expect(err).To(gomega.Succeed())
		}

	})

})
