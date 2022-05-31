package currency

import (
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
	"github.com/spikeekips/mitum/util/hint"
)

func (fa NilFeeer) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(jsonenc.HintedHead{
		H: fa.Hint(),
	})
}

func (fa *NilFeeer) UnmarsahlJSON(b []byte) error {
	var ht jsonenc.HintedHead
	if err := jsonenc.Unmarshal(b, &ht); err != nil {
		return err
	}

	fa.BaseHinter = hint.NewBaseHinter(ht.H)

	return nil
}

type FixedFeeerJSONPacker struct {
	jsonenc.HintedHead
	RC base.Address `json:"receiver"`
	AM currency.Big `json:"amount"`
	EM currency.Big `json:"exchange-min-amount"`
}

func (fa FixedFeeer) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(FixedFeeerJSONPacker{
		HintedHead: jsonenc.NewHintedHead(fa.Hint()),
		RC:         fa.receiver,
		AM:         fa.amount,
		EM:         fa.exchangeMin,
	})
}

type FixedFeeerJSONUnpacker struct {
	HT hint.Hint           `json:"_hint"`
	RC base.AddressDecoder `json:"receiver"`
	AM currency.Big        `json:"amount"`
	EM currency.Big        `json:"exchange-min-amount"`
}

func (fa *FixedFeeer) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var ufa FixedFeeerJSONUnpacker
	if err := enc.Unmarshal(b, &ufa); err != nil {
		return err
	}

	return fa.unpack(enc, ufa.HT, ufa.RC, ufa.AM, ufa.EM)
}

type RatioFeeerJSONPacker struct {
	jsonenc.HintedHead
	RC base.Address `json:"receiver"`
	RA float64      `json:"ratio"`
	MI currency.Big `json:"min"`
	MA currency.Big `json:"max"`
	EM currency.Big `json:"exchange-min-amount"`
}

func (fa RatioFeeer) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(RatioFeeerJSONPacker{
		HintedHead: jsonenc.NewHintedHead(fa.Hint()),
		RC:         fa.receiver,
		RA:         fa.ratio,
		MI:         fa.min,
		MA:         fa.max,
		EM:         fa.exchangeMin,
	})
}

type RatioFeeerJSONUnpacker struct {
	HT hint.Hint           `json:"_hint"`
	RC base.AddressDecoder `json:"receiver"`
	RA float64             `json:"ratio"`
	MI currency.Big        `json:"min"`
	MA currency.Big        `json:"max"`
	EM currency.Big        `json:"exchange-min-amount"`
}

func (fa *RatioFeeer) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var ufa RatioFeeerJSONUnpacker
	if err := enc.Unmarshal(b, &ufa); err != nil {
		return err
	}

	return fa.unpack(enc, ufa.HT, ufa.RC, ufa.RA, ufa.MI, ufa.MA, ufa.EM)
}
