package currency

import (
	"encoding/json"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/util"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
	"github.com/spikeekips/mitum/util/hint"
)

type CreateContractAccountsItemJSONMarshaler struct {
	hint.BaseHinter
	Keys    currency.AccountKeys `json:"keys"`
	Amounts []currency.Amount    `json:"amounts"`
}

func (it BaseCreateContractAccountsItem) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(CreateContractAccountsItemJSONMarshaler{
		BaseHinter: it.BaseHinter,
		Keys:       it.keys,
		Amounts:    it.amounts,
	})
}

type CreateContractAccountsItemJSONUnMarshaler struct {
	Hint    hint.Hint       `json:"_hint"`
	Keys    json.RawMessage `json:"keys"`
	Amounts json.RawMessage `json:"amounts"`
}

func (it *BaseCreateContractAccountsItem) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of BaseCreateContractAccountsItem")

	var uit CreateContractAccountsItemJSONUnMarshaler
	if err := enc.Unmarshal(b, &uit); err != nil {
		return e(err, "")
	}

	return it.unpack(enc, uit.Hint, uit.Keys, uit.Amounts)
}
