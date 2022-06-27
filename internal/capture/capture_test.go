package capture_test

import (
	. "github.com/jimmykarily/open-ocr-reader/internal/capture"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TakePhoto", func() {
	It("does not return an error", func() {
		_, err := TakePhoto()
		Expect(err).To(BeNil())
	})
})
