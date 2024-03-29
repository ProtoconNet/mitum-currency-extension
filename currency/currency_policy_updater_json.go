package currency

import (
	"encoding/json"

	mitumcurrency "github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
)

type CurrencyPolicyUpdaterFactJSONMarshaler struct {
	base.BaseFactJSONMarshaler
	Currency mitumcurrency.CurrencyID `json:"currency"`
	Policy   CurrencyPolicy           `json:"policy"`
}

func (fact CurrencyPolicyUpdaterFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(CurrencyPolicyUpdaterFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Currency:              fact.currency,
		Policy:                fact.policy,
	})
}

type CurrencyPolicyUpdaterFactJSONUnMarshaler struct {
	base.BaseFactJSONUnmarshaler
	Currency string          `json:"currency"`
	Policy   json.RawMessage `json:"policy"`
}

func (fact *CurrencyPolicyUpdaterFact) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of CurrencyPolicyUpdaterFact")

	var uf CurrencyPolicyUpdaterFactJSONUnMarshaler
	if err := enc.Unmarshal(b, &uf); err != nil {
		return e(err, "")
	}

	fact.BaseFact.SetJSONUnmarshaler(uf.BaseFactJSONUnmarshaler)

	return fact.unpack(enc, uf.Currency, uf.Policy)
}

type currencyPolicyUpdaterMarshaler struct {
	mitumcurrency.BaseOperationJSONMarshaler
}

func (op CurrencyPolicyUpdater) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(currencyPolicyUpdaterMarshaler{
		BaseOperationJSONMarshaler: op.BaseOperation.JSONMarshaler(),
	})
}

func (op *CurrencyPolicyUpdater) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of CurrencyPolicyUpdater")

	var ubo mitumcurrency.BaseNodeOperation
	if err := ubo.DecodeJSON(b, enc); err != nil {
		return e(err, "")
	}

	op.BaseNodeOperation = ubo

	return nil
}
