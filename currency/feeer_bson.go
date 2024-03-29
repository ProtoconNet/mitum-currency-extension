package currency

import (
	"go.mongodb.org/mongo-driver/bson"

	bsonenc "github.com/ProtoconNet/mitum-currency/v2/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (fa NilFeeer) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(bsonenc.NewHintedDoc(fa.Hint()))
}

func (fa *NilFeeer) UnmarsahlBSON(b []byte) error {
	e := util.StringErrorFunc("failed to unmarshal bson of NilFeeer")

	var head bsonenc.HintedHead
	if err := bsonenc.Unmarshal(b, &head); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(head.H)
	if err != nil {
		return e(err, "")
	}

	fa.BaseHinter = hint.NewBaseHinter(ht)

	return nil
}

func (fa FixedFeeer) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":               fa.Hint().String(),
			"receiver":            fa.receiver,
			"amount":              fa.amount.String(),
			"exchange_min_amount": fa.exchangeMin.String(),
		},
	)
}

type FixedFeeerBSONUnpacker struct {
	Hint              string `bson:"_hint"`
	Receiver          string `bson:"receiver"`
	Amount            string `bson:"amount"`
	ExchangeMinAmount string `bson:"exchange_min_amount"`
}

func (fa *FixedFeeer) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of FixedFeeer")

	var ufa FixedFeeerBSONUnpacker
	if err := enc.Unmarshal(b, &ufa); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(ufa.Hint)
	if err != nil {
		return e(err, "")
	}

	return fa.unpack(enc, ht, ufa.Receiver, ufa.Amount, ufa.ExchangeMinAmount)
}

func (fa RatioFeeer) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":               fa.Hint().String(),
			"receiver":            fa.receiver,
			"ratio":               fa.ratio,
			"min":                 fa.min.String(),
			"max":                 fa.max.String(),
			"exchange_min_amount": fa.exchangeMin.String(),
		},
	)
}

type RatioFeeerBSONUnpacker struct {
	Hint              string  `bson:"_hint"`
	Receiver          string  `bson:"receiver"`
	Ratio             float64 `bson:"ratio"`
	Min               string  `bson:"min"`
	Max               string  `bson:"max"`
	ExchangeMinAmount string  `bson:"exchange_min_amount"`
}

func (fa *RatioFeeer) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of RatioFeeer")

	var ufa RatioFeeerBSONUnpacker
	if err := enc.Unmarshal(b, &ufa); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(ufa.Hint)
	if err != nil {
		return e(err, "")
	}

	return fa.unpack(enc, ht, ufa.Receiver, ufa.Ratio, ufa.Min, ufa.Max, ufa.ExchangeMinAmount)
}
