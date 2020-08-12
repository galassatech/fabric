/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package gdpr

import (
	"encoding/base64"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric/core/ledger/kvledger/txmgmt/rwsetutil"
	"github.com/hyperledger/fabric/protoutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateBlock(t *testing.T) {
	rawBlock := `CkYIHBIgOBV8evZ1LD96FiZOUTL4LiVUShpscVEPyj8lxqNUg4AaICHjhqvRHgGijEYCaCODiI6FCi7GAHT+5PNYFyEgwlNoEo4fCosfCr8eCssHCnMIAxoLCO6/oPkFEP7h0nsiC3Rlc3RjaGFubmVsKkAyZWYzMjhiMTI3YzNkNzE5MDQ0YzU0YTA4Yzc2MzU3ZWIyYzEyOTI0NGI4NjNkYmUzMDUzMjhiZDUzMDVlYzYxOhMSERIPZGVmYXVsdHBvbGljeWNjEtMGCrYGCgdPcmczTVNQEqoGLS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUNLakNDQWRDZ0F3SUJBZ0lRTlNncVUxczB0M29rYU96UWI3azRwakFLQmdncWhrak9QUVFEQWpCek1Rc3cKQ1FZRFZRUUdFd0pWVXpFVE1CRUdBMVVFQ0JNS1EyRnNhV1p2Y201cFlURVdNQlFHQTFVRUJ4TU5VMkZ1SUVaeQpZVzVqYVhOamJ6RVpNQmNHQTFVRUNoTVFiM0puTXk1bGVHRnRjR3hsTG1OdmJURWNNQm9HQTFVRUF4TVRZMkV1CmIzSm5NeTVsZUdGdGNHeGxMbU52YlRBZUZ3MHlNREE0TURNeE5ESTNNREJhRncwek1EQTRNREV4TkRJM01EQmEKTUd3eEN6QUpCZ05WQkFZVEFsVlRNUk13RVFZRFZRUUlFd3BEWVd4cFptOXlibWxoTVJZd0ZBWURWUVFIRXcxVApZVzRnUm5KaGJtTnBjMk52TVE4d0RRWURWUVFMRXdaamJHbGxiblF4SHpBZEJnTlZCQU1NRmxWelpYSXhRRzl5Clp6TXVaWGhoYlhCc1pTNWpiMjB3V1RBVEJnY3Foa2pPUFFJQkJnZ3Foa2pPUFFNQkJ3TkNBQVNKb0doY2MyKzgKdG0rdXFaODNOVkNHeko5dU5TUmZtYXdEZmtWQTNTZmc5dHB3RGhZem9WbzBHenpUWTBuZ2owMXBRN3pmVUJlMQpMMnN6QWNNYjAzT2JvMDB3U3pBT0JnTlZIUThCQWY4RUJBTUNCNEF3REFZRFZSMFRBUUgvQkFJd0FEQXJCZ05WCkhTTUVKREFpZ0NENnhJcUx3NHFlZGpQSnBSMVFBVXhTQVVGSm4xUGdIZlVieDd3MDBsNFUwakFLQmdncWhrak8KUFFRREFnTklBREJGQWlFQWdQdEU2QjlRanROWlBCZUZhTzFZYWtkdGtEKzJBK2xuMjJ4RVJMRWxIT01DSUVkSApZeExDUnRsd3prMVdSenhZWXUwdFIwUmk2d0ZLRzJoaXR5V29oZ2U4Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0KEhjTCuYjHFm1T4t11hLCi/O7fQwgslfPjz0S7hYK6xYK0wYKtgYKB09yZzNNU1ASqgYtLS0tLUJFR0lOIENFUlRJRklDQVRFLS0tLS0KTUlJQ0tqQ0NBZENnQXdJQkFnSVFOU2dxVTFzMHQzb2thT3pRYjdrNHBqQUtCZ2dxaGtqT1BRUURBakJ6TVFzdwpDUVlEVlFRR0V3SlZVekVUTUJFR0ExVUVDQk1LUTJGc2FXWnZjbTVwWVRFV01CUUdBMVVFQnhNTlUyRnVJRVp5CllXNWphWE5qYnpFWk1CY0dBMVVFQ2hNUWIzSm5NeTVsZUdGdGNHeGxMbU52YlRFY01Cb0dBMVVFQXhNVFkyRXUKYjNKbk15NWxlR0Z0Y0d4bExtTnZiVEFlRncweU1EQTRNRE14TkRJM01EQmFGdzB6TURBNE1ERXhOREkzTURCYQpNR3d4Q3pBSkJnTlZCQVlUQWxWVE1STXdFUVlEVlFRSUV3cERZV3hwWm05eWJtbGhNUll3RkFZRFZRUUhFdzFUCllXNGdSbkpoYm1OcGMyTnZNUTh3RFFZRFZRUUxFd1pqYkdsbGJuUXhIekFkQmdOVkJBTU1GbFZ6WlhJeFFHOXkKWnpNdVpYaGhiWEJzWlM1amIyMHdXVEFUQmdjcWhrak9QUUlCQmdncWhrak9QUU1CQndOQ0FBU0pvR2hjYzIrOAp0bSt1cVo4M05WQ0d6Sjl1TlNSZm1hd0Rma1ZBM1NmZzl0cHdEaFl6b1ZvMEd6elRZMG5najAxcFE3emZVQmUxCkwyc3pBY01iMDNPYm8wMHdTekFPQmdOVkhROEJBZjhFQkFNQ0I0QXdEQVlEVlIwVEFRSC9CQUl3QURBckJnTlYKSFNNRUpEQWlnQ0Q2eElxTHc0cWVkalBKcFIxUUFVeFNBVUZKbjFQZ0hmVWJ4N3cwMGw0VTBqQUtCZ2dxaGtqTwpQUVFEQWdOSUFEQkZBaUVBZ1B0RTZCOVFqdE5aUEJlRmFPMVlha2R0a0QrMkErbG4yMnhFUkxFbEhPTUNJRWRICll4TENSdGx3emsxV1J6eFlZdTB0UjBSaTZ3RktHMmhpdHlXb2hnZTgKLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQoSGNMK5iMcWbVPi3XWEsKL87t9DCCyV8+PPRKSEAotCisKKQgBEhESD2RlZmF1bHRwb2xpY3ljYxoSCgZpbnZva2UKAWEKAWIKAjEwEuAPCtkBCiBCp8Z6vLyITtGvilFIqjkhlaQynIs9QDqJ544QzHgB9xK0AQqUARJACgpfbGlmZWN5Y2xlEjIKMAoqbmFtZXNwYWNlcy9maWVsZHMvZGVmYXVsdHBvbGljeWNjL1NlcXVlbmNlEgIIGRJQCg9kZWZhdWx0cG9saWN5Y2MSPQoWChAA9I+/v2luaXRpYWxpemVkEgIIGgoHCgFhEgIIGgoHCgFiEgIIGhoHCgFhGgI5MBoICgFiGgMyMTAaAwjIASIWEg9kZWZhdWx0cG9saWN5Y2MaAzAuMBL9BgqyBgoHT3JnM01TUBKmBi0tLS0tQkVHSU4gQ0VSVElGSUNBVEUtLS0tLQpNSUlDSnpDQ0FjNmdBd0lCQWdJUWFXMFB6dUxqU0FLUnJQa2o0WGl5SmpBS0JnZ3Foa2pPUFFRREFqQnpNUXN3CkNRWURWUVFHRXdKVlV6RVRNQkVHQTFVRUNCTUtRMkZzYVdadmNtNXBZVEVXTUJRR0ExVUVCeE1OVTJGdUlFWnkKWVc1amFYTmpiekVaTUJjR0ExVUVDaE1RYjNKbk15NWxlR0Z0Y0d4bExtTnZiVEVjTUJvR0ExVUVBeE1UWTJFdQpiM0puTXk1bGVHRnRjR3hsTG1OdmJUQWVGdzB5TURBNE1ETXhOREkzTURCYUZ3MHpNREE0TURFeE5ESTNNREJhCk1Hb3hDekFKQmdOVkJBWVRBbFZUTVJNd0VRWURWUVFJRXdwRFlXeHBabTl5Ym1saE1SWXdGQVlEVlFRSEV3MVQKWVc0Z1JuSmhibU5wYzJOdk1RMHdDd1lEVlFRTEV3UndaV1Z5TVI4d0hRWURWUVFERXhad1pXVnlNQzV2Y21jegpMbVY0WVcxd2JHVXVZMjl0TUZrd0V3WUhLb1pJemowQ0FRWUlLb1pJemowREFRY0RRZ0FFSXFBZ0EzbXVvbFRxCjJJb0FaRHlsSGp3U2dMY1V6Rzkxa1RZME9YRnNIbGV0dTN3M1pWMlgzTUR1Wkd0R2RvUzgwaWp5Ym53NGJJMVcKbFloUC9YdlpmcU5OTUVzd0RnWURWUjBQQVFIL0JBUURBZ2VBTUF3R0ExVWRFd0VCL3dRQ01BQXdLd1lEVlIwagpCQ1F3SW9BZytzU0tpOE9Lbm5ZenlhVWRVQUZNVWdGQlNaOVQ0QjMxRzhlOE5OSmVGTkl3Q2dZSUtvWkl6ajBFCkF3SURSd0F3UkFJZ0dLWld4QlF5TmNMNDFkVEE0bGRLdGgyWTV3NmFsaHM1bjI2THFkb0hSVThDSUhYSHNnUGoKckI1MHpUOWEwTkNhaG9NK0xOVDV5QWdSTGZDNUVMdU8xdjdBCi0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0KEkYwRAIgAJ4FUu5UwvLPN+p0pzruoqA8bUpdSWAvjv2ZvKRc10gCICUhgq6oWXf58Vd5pfso/Gap4fQBHogVfZy4EmbKN5KOEoEHCrYGCgdPcmcxTVNQEqoGLS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUNLVENDQWMrZ0F3SUJBZ0lSQUtjNjNtblRBSzhpZy9TWlBOcFhJYnd3Q2dZSUtvWkl6ajBFQXdJd2N6RUwKTUFrR0ExVUVCaE1DVlZNeEV6QVJCZ05WQkFnVENrTmhiR2xtYjNKdWFXRXhGakFVQmdOVkJBY1REVk5oYmlCRwpjbUZ1WTJselkyOHhHVEFYQmdOVkJBb1RFRzl5WnpFdVpYaGhiWEJzWlM1amIyMHhIREFhQmdOVkJBTVRFMk5oCkxtOXlaekV1WlhoaGJYQnNaUzVqYjIwd0hoY05NakF3T0RBek1UUXlOekF3V2hjTk16QXdPREF4TVRReU56QXcKV2pCcU1Rc3dDUVlEVlFRR0V3SlZVekVUTUJFR0ExVUVDQk1LUTJGc2FXWnZjbTVwWVRFV01CUUdBMVVFQnhNTgpVMkZ1SUVaeVlXNWphWE5qYnpFTk1Bc0dBMVVFQ3hNRWNHVmxjakVmTUIwR0ExVUVBeE1XY0dWbGNqQXViM0puCk1TNWxlR0Z0Y0d4bExtTnZiVEJaTUJNR0J5cUdTTTQ5QWdFR0NDcUdTTTQ5QXdFSEEwSUFCT044STJCWFdyWDAKOXZiVlFMMlgzdjJvL29jM0NXeDJIeXhUcGdUNEVxTkpTcVlhQi9LaE1BUE5QU2I0UEJwLzdoQzZqamdBVTFzRApoYldwRjE2a3B5cWpUVEJMTUE0R0ExVWREd0VCL3dRRUF3SUhnREFNQmdOVkhSTUJBZjhFQWpBQU1Dc0dBMVVkCkl3UWtNQ0tBSU93QVpJbW9XUXFWV1ZDbittWGczRkN2bDQ2RmcxQXRiU0dmb3cvUHQvcFNNQW9HQ0NxR1NNNDkKQkFNQ0EwZ0FNRVVDSVFDYmV5TThsWWFucFdHeFMyY2RWU2tlUUM1Z0h6bFBQK2t5dG51SFAvZWRiQUlnTHNuSgp1YTlsMVdCcjZtQURxek1WS2FaTEI3T1A3WlFwM1FxL29OMkx0cUE9Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0KEkYwRAIgZPOdRDK/06LcJ0WvNOZHRvhah7872xOk1nxlJ1SqxqkCIAy6dDHXFcXCuxW2EjAwxN+l8A2onzOqI7rGcP0lq8P7EkcwRQIhAIKMMSVV+dDY+zv7fGjGSx7IPBC31/UgFKSTgWJqBOv/AiB0nQ81Cp6uHE/9Ft4B//pVAoovcVj3wwATEscFZ80CJxqbBwqMBwoPCgIIExIJCgcKAQEQAhgeEvgGCq0GCpAGCgpPcmRlcmVyTVNQEoEGLS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUNERENDQWJLZ0F3SUJBZ0lRYUlWTnZiZUlkbExtd3V0TUx3NVJtVEFLQmdncWhrak9QUVFEQWpCcE1Rc3cKQ1FZRFZRUUdFd0pWVXpFVE1CRUdBMVVFQ0JNS1EyRnNhV1p2Y201cFlURVdNQlFHQTFVRUJ4TU5VMkZ1SUVaeQpZVzVqYVhOamJ6RVVNQklHQTFVRUNoTUxaWGhoYlhCc1pTNWpiMjB4RnpBVkJnTlZCQU1URG1OaExtVjRZVzF3CmJHVXVZMjl0TUI0WERUSXdNRGd3TXpFME1qY3dNRm9YRFRNd01EZ3dNVEUwTWpjd01Gb3dXREVMTUFrR0ExVUUKQmhNQ1ZWTXhFekFSQmdOVkJBZ1RDa05oYkdsbWIzSnVhV0V4RmpBVUJnTlZCQWNURFZOaGJpQkdjbUZ1WTJsegpZMjh4SERBYUJnTlZCQU1URTI5eVpHVnlaWEl1WlhoaGJYQnNaUzVqYjIwd1dUQVRCZ2NxaGtqT1BRSUJCZ2dxCmhrak9QUU1CQndOQ0FBUmkvc3V1STRhYU5pVVViRDkvQ0paTHUwTFhKME9Lbi9McEwwZk9qdVBBdEhDWGpFUzAKN2Jsc2loTnp4THFweU5HeWF5VDNHcjZ4T3NCdVYxS3dPN0dobzAwd1N6QU9CZ05WSFE4QkFmOEVCQU1DQjRBdwpEQVlEVlIwVEFRSC9CQUl3QURBckJnTlZIU01FSkRBaWdDQkNrQ09UQVdEajFJZWVHOVFqMTlibEE4ZkJuelQ4CndxN0I0ZTRPR0V4SXB6QUtCZ2dxaGtqT1BRUURBZ05JQURCRkFpRUE4TnhCSnIxYllieEhLQk9mR3QwT1lIRSsKYWNFeUt5RVBMTFB4cjg1eVFkWUNJR0lzS3NUK0ZGZHEyWjFrMGlTekExeXgvdnRhR0U2T2Q4Y0l3RzkvQmE5TAotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tChIYboJJI8hvvFogiMVFDzl8DO3kIB1V1w8sEkYwRAIgCXoLmZRTaolTkzjjoRh/osa6s5YezQ2kBI9zuQOJ95sCIE5+McHDHtfmxTgOx+hrtO+PGa8mOtlVGHjkzo1t+ub0CgQKAggTCgAKAAoA`
	blockBytes, err := base64.StdEncoding.DecodeString(rawBlock)
	assert.NoError(t, err)

	block := common.Block{} // common.Block{}

	assert.NoError(t, proto.Unmarshal(blockBytes, &block))
	preImages := extractPreimages(&block)
	//preImages[0] = ([]byte)("123456")
	block.Data.PreimageSpace = preImages
	newBlock, err := validate(&block)
	require.NoError(t, err, "adsfadsfasdf")
	err = checkKvExist(preImages, newBlock)
	assert.NoError(t, err)

	preImages[0] = ([]byte)("123456")
	_, err = validate(&block)
	require.Error(t, err, "Error expected!")
}

func TestGetVanillaBlock(t *testing.T) {
	rawBlock := `CkYIHBIgOBV8evZ1LD96FiZOUTL4LiVUShpscVEPyj8lxqNUg4AaICHjhqvRHgGijEYCaCODiI6FCi7GAHT+5PNYFyEgwlNoEo4fCosfCr8eCssHCnMIAxoLCO6/oPkFEP7h0nsiC3Rlc3RjaGFubmVsKkAyZWYzMjhiMTI3YzNkNzE5MDQ0YzU0YTA4Yzc2MzU3ZWIyYzEyOTI0NGI4NjNkYmUzMDUzMjhiZDUzMDVlYzYxOhMSERIPZGVmYXVsdHBvbGljeWNjEtMGCrYGCgdPcmczTVNQEqoGLS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUNLakNDQWRDZ0F3SUJBZ0lRTlNncVUxczB0M29rYU96UWI3azRwakFLQmdncWhrak9QUVFEQWpCek1Rc3cKQ1FZRFZRUUdFd0pWVXpFVE1CRUdBMVVFQ0JNS1EyRnNhV1p2Y201cFlURVdNQlFHQTFVRUJ4TU5VMkZ1SUVaeQpZVzVqYVhOamJ6RVpNQmNHQTFVRUNoTVFiM0puTXk1bGVHRnRjR3hsTG1OdmJURWNNQm9HQTFVRUF4TVRZMkV1CmIzSm5NeTVsZUdGdGNHeGxMbU52YlRBZUZ3MHlNREE0TURNeE5ESTNNREJhRncwek1EQTRNREV4TkRJM01EQmEKTUd3eEN6QUpCZ05WQkFZVEFsVlRNUk13RVFZRFZRUUlFd3BEWVd4cFptOXlibWxoTVJZd0ZBWURWUVFIRXcxVApZVzRnUm5KaGJtTnBjMk52TVE4d0RRWURWUVFMRXdaamJHbGxiblF4SHpBZEJnTlZCQU1NRmxWelpYSXhRRzl5Clp6TXVaWGhoYlhCc1pTNWpiMjB3V1RBVEJnY3Foa2pPUFFJQkJnZ3Foa2pPUFFNQkJ3TkNBQVNKb0doY2MyKzgKdG0rdXFaODNOVkNHeko5dU5TUmZtYXdEZmtWQTNTZmc5dHB3RGhZem9WbzBHenpUWTBuZ2owMXBRN3pmVUJlMQpMMnN6QWNNYjAzT2JvMDB3U3pBT0JnTlZIUThCQWY4RUJBTUNCNEF3REFZRFZSMFRBUUgvQkFJd0FEQXJCZ05WCkhTTUVKREFpZ0NENnhJcUx3NHFlZGpQSnBSMVFBVXhTQVVGSm4xUGdIZlVieDd3MDBsNFUwakFLQmdncWhrak8KUFFRREFnTklBREJGQWlFQWdQdEU2QjlRanROWlBCZUZhTzFZYWtkdGtEKzJBK2xuMjJ4RVJMRWxIT01DSUVkSApZeExDUnRsd3prMVdSenhZWXUwdFIwUmk2d0ZLRzJoaXR5V29oZ2U4Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0KEhjTCuYjHFm1T4t11hLCi/O7fQwgslfPjz0S7hYK6xYK0wYKtgYKB09yZzNNU1ASqgYtLS0tLUJFR0lOIENFUlRJRklDQVRFLS0tLS0KTUlJQ0tqQ0NBZENnQXdJQkFnSVFOU2dxVTFzMHQzb2thT3pRYjdrNHBqQUtCZ2dxaGtqT1BRUURBakJ6TVFzdwpDUVlEVlFRR0V3SlZVekVUTUJFR0ExVUVDQk1LUTJGc2FXWnZjbTVwWVRFV01CUUdBMVVFQnhNTlUyRnVJRVp5CllXNWphWE5qYnpFWk1CY0dBMVVFQ2hNUWIzSm5NeTVsZUdGdGNHeGxMbU52YlRFY01Cb0dBMVVFQXhNVFkyRXUKYjNKbk15NWxlR0Z0Y0d4bExtTnZiVEFlRncweU1EQTRNRE14TkRJM01EQmFGdzB6TURBNE1ERXhOREkzTURCYQpNR3d4Q3pBSkJnTlZCQVlUQWxWVE1STXdFUVlEVlFRSUV3cERZV3hwWm05eWJtbGhNUll3RkFZRFZRUUhFdzFUCllXNGdSbkpoYm1OcGMyTnZNUTh3RFFZRFZRUUxFd1pqYkdsbGJuUXhIekFkQmdOVkJBTU1GbFZ6WlhJeFFHOXkKWnpNdVpYaGhiWEJzWlM1amIyMHdXVEFUQmdjcWhrak9QUUlCQmdncWhrak9QUU1CQndOQ0FBU0pvR2hjYzIrOAp0bSt1cVo4M05WQ0d6Sjl1TlNSZm1hd0Rma1ZBM1NmZzl0cHdEaFl6b1ZvMEd6elRZMG5najAxcFE3emZVQmUxCkwyc3pBY01iMDNPYm8wMHdTekFPQmdOVkhROEJBZjhFQkFNQ0I0QXdEQVlEVlIwVEFRSC9CQUl3QURBckJnTlYKSFNNRUpEQWlnQ0Q2eElxTHc0cWVkalBKcFIxUUFVeFNBVUZKbjFQZ0hmVWJ4N3cwMGw0VTBqQUtCZ2dxaGtqTwpQUVFEQWdOSUFEQkZBaUVBZ1B0RTZCOVFqdE5aUEJlRmFPMVlha2R0a0QrMkErbG4yMnhFUkxFbEhPTUNJRWRICll4TENSdGx3emsxV1J6eFlZdTB0UjBSaTZ3RktHMmhpdHlXb2hnZTgKLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQoSGNMK5iMcWbVPi3XWEsKL87t9DCCyV8+PPRKSEAotCisKKQgBEhESD2RlZmF1bHRwb2xpY3ljYxoSCgZpbnZva2UKAWEKAWIKAjEwEuAPCtkBCiBCp8Z6vLyITtGvilFIqjkhlaQynIs9QDqJ544QzHgB9xK0AQqUARJACgpfbGlmZWN5Y2xlEjIKMAoqbmFtZXNwYWNlcy9maWVsZHMvZGVmYXVsdHBvbGljeWNjL1NlcXVlbmNlEgIIGRJQCg9kZWZhdWx0cG9saWN5Y2MSPQoWChAA9I+/v2luaXRpYWxpemVkEgIIGgoHCgFhEgIIGgoHCgFiEgIIGhoHCgFhGgI5MBoICgFiGgMyMTAaAwjIASIWEg9kZWZhdWx0cG9saWN5Y2MaAzAuMBL9BgqyBgoHT3JnM01TUBKmBi0tLS0tQkVHSU4gQ0VSVElGSUNBVEUtLS0tLQpNSUlDSnpDQ0FjNmdBd0lCQWdJUWFXMFB6dUxqU0FLUnJQa2o0WGl5SmpBS0JnZ3Foa2pPUFFRREFqQnpNUXN3CkNRWURWUVFHRXdKVlV6RVRNQkVHQTFVRUNCTUtRMkZzYVdadmNtNXBZVEVXTUJRR0ExVUVCeE1OVTJGdUlFWnkKWVc1amFYTmpiekVaTUJjR0ExVUVDaE1RYjNKbk15NWxlR0Z0Y0d4bExtTnZiVEVjTUJvR0ExVUVBeE1UWTJFdQpiM0puTXk1bGVHRnRjR3hsTG1OdmJUQWVGdzB5TURBNE1ETXhOREkzTURCYUZ3MHpNREE0TURFeE5ESTNNREJhCk1Hb3hDekFKQmdOVkJBWVRBbFZUTVJNd0VRWURWUVFJRXdwRFlXeHBabTl5Ym1saE1SWXdGQVlEVlFRSEV3MVQKWVc0Z1JuSmhibU5wYzJOdk1RMHdDd1lEVlFRTEV3UndaV1Z5TVI4d0hRWURWUVFERXhad1pXVnlNQzV2Y21jegpMbVY0WVcxd2JHVXVZMjl0TUZrd0V3WUhLb1pJemowQ0FRWUlLb1pJemowREFRY0RRZ0FFSXFBZ0EzbXVvbFRxCjJJb0FaRHlsSGp3U2dMY1V6Rzkxa1RZME9YRnNIbGV0dTN3M1pWMlgzTUR1Wkd0R2RvUzgwaWp5Ym53NGJJMVcKbFloUC9YdlpmcU5OTUVzd0RnWURWUjBQQVFIL0JBUURBZ2VBTUF3R0ExVWRFd0VCL3dRQ01BQXdLd1lEVlIwagpCQ1F3SW9BZytzU0tpOE9Lbm5ZenlhVWRVQUZNVWdGQlNaOVQ0QjMxRzhlOE5OSmVGTkl3Q2dZSUtvWkl6ajBFCkF3SURSd0F3UkFJZ0dLWld4QlF5TmNMNDFkVEE0bGRLdGgyWTV3NmFsaHM1bjI2THFkb0hSVThDSUhYSHNnUGoKckI1MHpUOWEwTkNhaG9NK0xOVDV5QWdSTGZDNUVMdU8xdjdBCi0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0KEkYwRAIgAJ4FUu5UwvLPN+p0pzruoqA8bUpdSWAvjv2ZvKRc10gCICUhgq6oWXf58Vd5pfso/Gap4fQBHogVfZy4EmbKN5KOEoEHCrYGCgdPcmcxTVNQEqoGLS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUNLVENDQWMrZ0F3SUJBZ0lSQUtjNjNtblRBSzhpZy9TWlBOcFhJYnd3Q2dZSUtvWkl6ajBFQXdJd2N6RUwKTUFrR0ExVUVCaE1DVlZNeEV6QVJCZ05WQkFnVENrTmhiR2xtYjNKdWFXRXhGakFVQmdOVkJBY1REVk5oYmlCRwpjbUZ1WTJselkyOHhHVEFYQmdOVkJBb1RFRzl5WnpFdVpYaGhiWEJzWlM1amIyMHhIREFhQmdOVkJBTVRFMk5oCkxtOXlaekV1WlhoaGJYQnNaUzVqYjIwd0hoY05NakF3T0RBek1UUXlOekF3V2hjTk16QXdPREF4TVRReU56QXcKV2pCcU1Rc3dDUVlEVlFRR0V3SlZVekVUTUJFR0ExVUVDQk1LUTJGc2FXWnZjbTVwWVRFV01CUUdBMVVFQnhNTgpVMkZ1SUVaeVlXNWphWE5qYnpFTk1Bc0dBMVVFQ3hNRWNHVmxjakVmTUIwR0ExVUVBeE1XY0dWbGNqQXViM0puCk1TNWxlR0Z0Y0d4bExtTnZiVEJaTUJNR0J5cUdTTTQ5QWdFR0NDcUdTTTQ5QXdFSEEwSUFCT044STJCWFdyWDAKOXZiVlFMMlgzdjJvL29jM0NXeDJIeXhUcGdUNEVxTkpTcVlhQi9LaE1BUE5QU2I0UEJwLzdoQzZqamdBVTFzRApoYldwRjE2a3B5cWpUVEJMTUE0R0ExVWREd0VCL3dRRUF3SUhnREFNQmdOVkhSTUJBZjhFQWpBQU1Dc0dBMVVkCkl3UWtNQ0tBSU93QVpJbW9XUXFWV1ZDbittWGczRkN2bDQ2RmcxQXRiU0dmb3cvUHQvcFNNQW9HQ0NxR1NNNDkKQkFNQ0EwZ0FNRVVDSVFDYmV5TThsWWFucFdHeFMyY2RWU2tlUUM1Z0h6bFBQK2t5dG51SFAvZWRiQUlnTHNuSgp1YTlsMVdCcjZtQURxek1WS2FaTEI3T1A3WlFwM1FxL29OMkx0cUE9Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0KEkYwRAIgZPOdRDK/06LcJ0WvNOZHRvhah7872xOk1nxlJ1SqxqkCIAy6dDHXFcXCuxW2EjAwxN+l8A2onzOqI7rGcP0lq8P7EkcwRQIhAIKMMSVV+dDY+zv7fGjGSx7IPBC31/UgFKSTgWJqBOv/AiB0nQ81Cp6uHE/9Ft4B//pVAoovcVj3wwATEscFZ80CJxqbBwqMBwoPCgIIExIJCgcKAQEQAhgeEvgGCq0GCpAGCgpPcmRlcmVyTVNQEoEGLS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUNERENDQWJLZ0F3SUJBZ0lRYUlWTnZiZUlkbExtd3V0TUx3NVJtVEFLQmdncWhrak9QUVFEQWpCcE1Rc3cKQ1FZRFZRUUdFd0pWVXpFVE1CRUdBMVVFQ0JNS1EyRnNhV1p2Y201cFlURVdNQlFHQTFVRUJ4TU5VMkZ1SUVaeQpZVzVqYVhOamJ6RVVNQklHQTFVRUNoTUxaWGhoYlhCc1pTNWpiMjB4RnpBVkJnTlZCQU1URG1OaExtVjRZVzF3CmJHVXVZMjl0TUI0WERUSXdNRGd3TXpFME1qY3dNRm9YRFRNd01EZ3dNVEUwTWpjd01Gb3dXREVMTUFrR0ExVUUKQmhNQ1ZWTXhFekFSQmdOVkJBZ1RDa05oYkdsbWIzSnVhV0V4RmpBVUJnTlZCQWNURFZOaGJpQkdjbUZ1WTJsegpZMjh4SERBYUJnTlZCQU1URTI5eVpHVnlaWEl1WlhoaGJYQnNaUzVqYjIwd1dUQVRCZ2NxaGtqT1BRSUJCZ2dxCmhrak9QUU1CQndOQ0FBUmkvc3V1STRhYU5pVVViRDkvQ0paTHUwTFhKME9Lbi9McEwwZk9qdVBBdEhDWGpFUzAKN2Jsc2loTnp4THFweU5HeWF5VDNHcjZ4T3NCdVYxS3dPN0dobzAwd1N6QU9CZ05WSFE4QkFmOEVCQU1DQjRBdwpEQVlEVlIwVEFRSC9CQUl3QURBckJnTlZIU01FSkRBaWdDQkNrQ09UQVdEajFJZWVHOVFqMTlibEE4ZkJuelQ4CndxN0I0ZTRPR0V4SXB6QUtCZ2dxaGtqT1BRUURBZ05JQURCRkFpRUE4TnhCSnIxYllieEhLQk9mR3QwT1lIRSsKYWNFeUt5RVBMTFB4cjg1eVFkWUNJR0lzS3NUK0ZGZHEyWjFrMGlTekExeXgvdnRhR0U2T2Q4Y0l3RzkvQmE5TAotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tChIYboJJI8hvvFogiMVFDzl8DO3kIB1V1w8sEkYwRAIgCXoLmZRTaolTkzjjoRh/osa6s5YezQ2kBI9zuQOJ95sCIE5+McHDHtfmxTgOx+hrtO+PGa8mOtlVGHjkzo1t+ub0CgQKAggTCgAKAAoA`
	blockBytes, err := base64.StdEncoding.DecodeString(rawBlock)
	assert.NoError(t, err)

	block := common.Block{} // common.Block{}

	assert.NoError(t, proto.Unmarshal(blockBytes, &block))
	preImages := extractPreimages(&block)

	block.Data.PreimageSpace = preImages
	clearKVWrites(&block)

	getVanillaBlock(&block)
	err = checkKvExist(preImages, &block)
	assert.NoError(t, err)

}

func checkKvExist(preimages [][]byte, block *common.Block) error {
	m := map[string]struct{}{}
	for _, pi := range preimages {
		temp := (string)(pi)
		m[temp] = struct{}{}
	}

	for _, envBytes := range block.Data.Data {
		env, err := protoutil.GetEnvelopeFromBlock(envBytes)
		if err != nil {
			//logger.Warning("Invalid envelope:", err)
			return err
		}

		payload, err := protoutil.UnmarshalPayload(env.Payload) //protoutil.GetBytesPayload()
		if err != nil {
			//logger.Warning("Invalid payload:", err)
			return err
		}

		tx, err := protoutil.UnmarshalTransaction(payload.Data)
		if err != nil {
			return err
		}

		_, respPayload, err := protoutil.GetPayloads(tx.Actions[0])
		if err != nil {
			return err
		}

		txRWSet := &rwsetutil.TxRwSet{}

		if err = txRWSet.FromProtoBytes(respPayload.Results); err != nil {
			return err
		}
		for _, nsRWSet := range txRWSet.NsRwSets {
			//nsRWSet.KvRwSet
			for _, kvWrite := range nsRWSet.KvRwSet.Writes {
				temp := (string)(kvWrite.Value)
				if memberOf(temp, m) == false {
					return ErrVal
				}
			}

		}
	}
	return nil

}
