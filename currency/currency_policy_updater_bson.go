package currency // nolint: dupl

import (
	"go.mongodb.org/mongo-driver/bson"

	mitumcurrency "github.com/ProtoconNet/mitum-currency/v2/currency"
	bsonenc "github.com/ProtoconNet/mitum-currency/v2/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
)

func (fact CurrencyPolicyUpdaterFact) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":    fact.Hint().String(),
			"currency": fact.currency,
			"policy":   fact.policy,
			"hash":     fact.BaseFact.Hash().String(),
			"token":    fact.BaseFact.Token(),
		},
	)
}

type CurrencyPolicyUpdaterFactBSONUnmarshaler struct {
	Hint     string   `bson:"_hint"`
	Currency string   `bson:"currency"`
	Policy   bson.Raw `bson:"policy"`
}

func (fact *CurrencyPolicyUpdaterFact) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of CurrencyPolicyUpdaterFact")

	var ubf mitumcurrency.BaseFactBSONUnmarshaler
	if err := enc.Unmarshal(b, &ubf); err != nil {
		return e(err, "")
	}

	fact.BaseFact.SetHash(valuehash.NewBytesFromString(ubf.Hash))
	fact.BaseFact.SetToken(ubf.Token)

	var uf CurrencyPolicyUpdaterFactBSONUnmarshaler
	if err := bson.Unmarshal(b, &uf); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(uf.Hint)
	if err != nil {
		return e(err, "")
	}
	fact.BaseHinter = hint.NewBaseHinter(ht)

	return fact.unpack(enc, uf.Currency, uf.Policy)
}

func (op CurrencyPolicyUpdater) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint": op.Hint().String(),
			"hash":  op.Hash().String(),
			"fact":  op.Fact(),
			"signs": op.Signs(),
		})
}

func (op *CurrencyPolicyUpdater) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of CurrencyPolicyUpdater")

	var ubo mitumcurrency.BaseNodeOperation
	if err := ubo.DecodeBSON(b, enc); err != nil {
		return e(err, "")
	}

	op.BaseNodeOperation = ubo

	return nil
}
