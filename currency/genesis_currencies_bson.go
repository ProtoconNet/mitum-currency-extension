package currency

import (
	mitumcurrency "github.com/ProtoconNet/mitum-currency/v2/currency"
	bsonenc "github.com/ProtoconNet/mitum-currency/v2/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
	"go.mongodb.org/mongo-driver/bson"
)

func (fact GenesisCurrenciesFact) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":            fact.Hint().String(),
			"genesis_node_key": fact.genesisNodeKey.String(),
			"keys":             fact.keys,
			"currencies":       fact.cs,
			"hash":             fact.BaseFact.Hash().String(),
			"token":            fact.BaseFact.Token(),
		},
	)
}

type GenesisCurrenciesFactBSONUnMarshaler struct {
	Hint           string   `bson:"_hint"`
	GenesisNodeKey string   `bson:"genesis_node_key"`
	Keys           bson.Raw `bson:"keys"`
	Currencies     bson.Raw `bson:"currencies"`
}

func (fact *GenesisCurrenciesFact) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of GenesisCurrenciesFact")

	var ubf mitumcurrency.BaseFactBSONUnmarshaler
	if err := enc.Unmarshal(b, &ubf); err != nil {
		return e(err, "")
	}

	fact.BaseFact.SetHash(valuehash.NewBytesFromString(ubf.Hash))
	fact.BaseFact.SetToken(ubf.Token)

	var uf GenesisCurrenciesFactBSONUnMarshaler
	if err := bson.Unmarshal(b, &uf); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(uf.Hint)
	if err != nil {
		return e(err, "")
	}
	fact.BaseHinter = hint.NewBaseHinter(ht)

	return fact.unpack(enc, uf.GenesisNodeKey, uf.Keys, uf.Currencies)
}

func (op GenesisCurrencies) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(op.BaseOperation)
}

func (op *GenesisCurrencies) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of GenesisCurrencies")

	var ubo mitumcurrency.BaseOperation
	if err := ubo.DecodeBSON(b, enc); err != nil {
		return e(err, "")
	}

	op.BaseOperation = ubo

	return nil
}
