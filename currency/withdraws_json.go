package currency

import (
	"encoding/json"

	mitumcurrency "github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
)

type TransferFactJSONMarshaler struct {
	base.BaseFactJSONMarshaler
	Sender base.Address    `json:"sender"`
	Items  []WithdrawsItem `json:"items"`
}

func (fact WithdrawsFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(TransferFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Sender:                fact.sender,
		Items:                 fact.items,
	})
}

type WithdrawsFactJSONUnmarshaler struct {
	base.BaseFactJSONUnmarshaler
	Sender string          `json:"sender"`
	Items  json.RawMessage `json:"items"`
}

func (fact *WithdrawsFact) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of WithdrawsFact")

	var uf WithdrawsFactJSONUnmarshaler

	if err := enc.Unmarshal(b, &uf); err != nil {
		return e(err, "")
	}

	fact.BaseFact.SetJSONUnmarshaler(uf.BaseFactJSONUnmarshaler)

	return fact.unpack(enc, uf.Sender, uf.Items)
}

type withdrawsMarshaler struct {
	mitumcurrency.BaseOperationJSONMarshaler
}

func (op Withdraws) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(withdrawsMarshaler{
		BaseOperationJSONMarshaler: op.BaseOperation.JSONMarshaler(),
	})
}

func (op *Withdraws) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of Withdraws")

	var ubo mitumcurrency.BaseOperation
	if err := ubo.DecodeJSON(b, enc); err != nil {
		return e(err, "")
	}

	op.BaseOperation = ubo

	return nil
}
