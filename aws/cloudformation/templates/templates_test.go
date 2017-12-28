package templates_test

import (
	"encoding/json"

	"github.com/henrytk/aws-service-broker/aws/cloudformation/templates"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Templates", func() {
	It("Must contain a MongoDB stack template using valid JSON", func() {
		var data interface{}
		err := json.Unmarshal(templates.MongoDBStack, &data)
		Expect(err).NotTo(HaveOccurred())
	})
})
