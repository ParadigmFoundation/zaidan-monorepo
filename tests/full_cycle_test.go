package tests_test

import (
	"strconv"

	"github.com/ParadigmFoundation/zaidan-monorepo/tests"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Full cycle", func() {
	var (
		quote *tests.Quote
	)

	Describe("Getting a Quote", func() {
		args := []interface{}{
			"0x871dd7c2b4b25e1aa18728e9d5f2af4c4e431f5c",
			"0x0b1ba0af832d7c05fd64161e0db78e85978e8082",
			nil,
			"1",
			"0xE36Ea790bc9d7AB70C55260C66D52b1eca985f84",
		}
		It("calls dealer_getQuote", func() {
			resp := []interface{}{
				new(tests.Quote),
			}

			err := client.Call(&resp, "dealer_getQuote", args...)
			Expect(err).NotTo(HaveOccurred())
			quote = resp[0].(*tests.Quote)
		})
	})

	Describe("Submit a previously-fetched quote for settlement", func() {
		It("calls dealer_submitFill", func() {
			if quote == nil {
				return
			}

			e, _ := strconv.Atoi(quote.ZeroExTransactionInfo.Transaction.ExpirationTimeSeconds)
			args := []interface{}{
				quote.QuoteId,
				quote.ZeroExTransactionInfo.Transaction.Salt,
				quote.ZeroExTransactionInfo.Order.Signature,
				quote.ZeroExTransactionInfo.Transaction.SignerAddress,
				quote.ZeroExTransactionInfo.Transaction.Data,
				quote.ZeroExTransactionInfo.Transaction.GasPrice,
				e,
			}

			var (
				quoteId     string
				txHash      string
				submittedAt int64
			)

			resp := []interface{}{
				&quoteId, &txHash, &submittedAt,
			}

			err := client.Call(&resp, "dealer_submitFill", args...)
			Expect(err).NotTo(HaveOccurred())
			Expect(quoteId).Should(Equal(quote.QuoteId))
			Expect(submittedAt).Should(BeNumerically(">", quote.ServerTime))
			Expect(submittedAt).Should(BeNumerically("<", quote.Expiration))
		})
	})
})
