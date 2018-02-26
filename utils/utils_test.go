package utils_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/henrytk/aws-service-broker/utils"
)

var _ = Describe("GetMD5Hex", func() {
	It("returns the Hex encoded MD5 hash of the given string", func() {
		md5Hex := GetMD5Hex("ce71b484-d542-40f7-9dd4-5526e38c81ba", 32)
		// Expectation generated with
		// echo -n ce71b484-d542-40f7-9dd4-5526e38c81ba | openssl dgst -md5 -binary | xxd -p
		Expect(md5Hex).To(Equal("3b3501055c9616a1a66fba5be7898f55"))
	})

	It("truncates the result when it's longer than the resuested max size", func() {
		md5Hex := GetMD5Hex("ce71b484-d542-40f7-9dd4-5526e38c81ba", 16)
		Expect(md5Hex).To(Equal("3b3501055c9616a1"))
	})
})
