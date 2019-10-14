package node

import (
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

func RunTest(provider Provider) {

	labels := make(map[string]string)
	labels["label1"] = "label1"

	ginkgo.BeforeEach(func() {
		provider.Clear()
	})

	// Add
	ginkgo.It("Should be able to add role", func() {

		node := &entities.Node{OrganizationId: "org", ClusterId: "cluster_id", NodeId: "node",
			Ip: "0.0.0.0", Labels: labels, Status: entities.InfraStatusRunning, State: 0}

		err := provider.Add(*node)
		gomega.Expect(err).To(gomega.Succeed())

	})

	// Update
	ginkgo.It("Should be able to update role", func() {

		// add a role
		node := &entities.Node{OrganizationId: "org", ClusterId: "cluster_id", NodeId: "node",
			Ip: "0.0.0.0", Labels: labels, Status: entities.InfraStatusRunning, State: 0}

		err := provider.Add(*node)
		gomega.Expect(err).To(gomega.Succeed())

		// uodate it
		labels["label2"] = "label2"
		node.OrganizationId = "org_MOD"
		node.Labels = labels

		err = provider.Update(*node)
		gomega.Expect(err).To(gomega.Succeed())

	})
	ginkgo.It("Should not be able to update role", func() {

		node := &entities.Node{OrganizationId: "org", ClusterId: "clusterMODD", NodeId: "node",
			Ip: "127.0.0.1", Labels: labels, Status: entities.InfraStatusInstalling, State: 1}

		err := provider.Update(*node)
		gomega.Expect(err).NotTo(gomega.Succeed())

	})

	// Exists
	ginkgo.It("Should be able to find role", func() {

		// add a role
		node := &entities.Node{OrganizationId: "org", ClusterId: "cluster_id", NodeId: "node",
			Ip: "0.0.0.0", Labels: labels, Status: entities.InfraStatusRunning, State: 0}

		err := provider.Add(*node)
		gomega.Expect(err).To(gomega.Succeed())

		// ask if it exists
		exits, err := provider.Exists(node.NodeId)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(exits).To(gomega.BeTrue())

	})

	ginkgo.It("Should not be able to find role", func() {

		exits, err := provider.Exists("node")
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(exits).NotTo(gomega.BeTrue())

	})

	// Get
	ginkgo.It("Should be able to get the role", func() {

		// add a role
		node := &entities.Node{OrganizationId: "org", ClusterId: "cluster_id", NodeId: "node",
			Ip: "0.0.0.0", Labels: labels, Status: entities.InfraStatusRunning, State: 0}

		err := provider.Add(*node)
		gomega.Expect(err).To(gomega.Succeed())

		// ask for it
		node, err = provider.Get(node.NodeId)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(node).NotTo(gomega.BeNil())
	})
	ginkgo.It("Should not be able to get the role", func() {

		node, err := provider.Get("node")
		gomega.Expect(err).NotTo(gomega.Succeed())
		gomega.Expect(node).To(gomega.BeNil())
	})

	// Remove
	ginkgo.It("Should be able to remove the role", func() {

		// add a role
		node := &entities.Node{OrganizationId: "org", ClusterId: "cluster_id", NodeId: "node",
			Ip: "0.0.0.0", Labels: labels, Status: entities.InfraStatusRunning, State: 0}

		err := provider.Add(*node)
		gomega.Expect(err).To(gomega.Succeed())

		// remove it
		err = provider.Remove(node.NodeId)
		gomega.Expect(err).To(gomega.Succeed())
	})
	ginkgo.It("Should not be able to remove the role", func() {

		err := provider.Remove("node")
		gomega.Expect(err).NotTo(gomega.Succeed())
	})
}
