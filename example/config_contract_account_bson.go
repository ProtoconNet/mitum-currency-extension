package example // nolint: dupl

import (
	"go.mongodb.org/mongo-driver/bson"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
	"github.com/spikeekips/mitum/util/valuehash"
)

func (fact ConfigContractAccountFact) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bsonenc.MergeBSONM(bsonenc.NewHintedDoc(fact.Hint()),
			bson.M{
				"hash":     fact.h,
				"token":    fact.token,
				"sender":   fact.sender,
				"target":   fact.target,
				"currency": fact.currency,
			}))
}

type ConfigContractAccountFactBSONUnpacker struct {
	H  valuehash.Bytes     `bson:"hash"`
	TK []byte              `bson:"token"`
	SD base.AddressDecoder `bson:"sender"`
	TG base.AddressDecoder `bson:"target"`
	CR string              `bson:"currency"`
}

func (fact *ConfigContractAccountFact) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var ufact ConfigContractAccountFactBSONUnpacker
	if err := bson.Unmarshal(b, &ufact); err != nil {
		return err
	}

	return fact.unpack(enc, ufact.H, ufact.TK, ufact.SD, ufact.TG, ufact.CR)
}

func (op *ConfigContractAccount) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var ubo currency.BaseOperation
	if err := ubo.UnpackBSON(b, enc); err != nil {
		return err
	}

	op.BaseOperation = ubo

	return nil
}
