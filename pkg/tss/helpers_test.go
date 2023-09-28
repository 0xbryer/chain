package tss_test

import (
	"github.com/bandprotocol/chain/v2/pkg/tss"
)

func (suite *TSSTestSuite) TestConcatBytes() {
	res := tss.ConcatBytes([]byte("abc"), []byte("de"), []byte("fghi"))
	suite.Require().Equal([]byte("abcdefghi"), res)
}

func (suite *TSSTestSuite) TestGenerateKeyPairs() {
	kps, err := tss.GenerateKeyPairs(3)
	suite.Require().NoError(err)
	suite.Require().Equal(3, len(kps))

	for _, kp := range kps {
		pubKey := kp.PrivKey.Point()
		suite.Require().Equal(kp.PubKey, pubKey)
	}
}

func (suite *TSSTestSuite) TestGenerateKeyPair() {
	kp, err := tss.GenerateKeyPair()
	suite.Require().NoError(err)

	pubKey := kp.PrivKey.Point()
	suite.Require().Equal(kp.PubKey, pubKey)
}
