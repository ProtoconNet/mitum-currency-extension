package example // nolint: dupl

import (
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
	"github.com/spikeekips/mitum/util/valuehash"
)

type ConfigContractAccountFactJSONPacker struct {
	jsonenc.HintedHead
	H  valuehash.Hash      `json:"hash"`
	TK []byte              `json:"token"`
	SD base.Address        `json:"sender"`
	TG base.Address        `json:"target"`
	CR currency.CurrencyID `json:"currency"`
}

func (fact ConfigContractAccountFact) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(ConfigContractAccountFactJSONPacker{
		HintedHead: jsonenc.NewHintedHead(fact.Hint()),
		H:          fact.h,
		TK:         fact.token,
		SD:         fact.sender,
		TG:         fact.target,
		CR:         fact.currency,
	})
}

type ConfigContractAccountFactJSONUnpacker struct {
	H  valuehash.Bytes     `json:"hash"`
	TK []byte              `json:"token"`
	SD base.AddressDecoder `json:"sender"`
	TG base.AddressDecoder `json:"target"`
	CR string              `json:"currency"`
}

func (fact *ConfigContractAccountFact) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var ufact ConfigContractAccountFactJSONUnpacker
	if err := enc.Unmarshal(b, &ufact); err != nil {
		return err
	}

	return fact.unpack(enc, ufact.H, ufact.TK, ufact.SD, ufact.TG, ufact.CR)
}

func (op *ConfigContractAccount) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var ubo currency.BaseOperation
	if err := ubo.UnpackJSON(b, enc); err != nil {
		return err
	}

	op.BaseOperation = ubo

	return nil
}
