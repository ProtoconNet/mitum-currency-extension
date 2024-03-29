package currency

import (
	mitumcurrency "github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (it *BaseCreateContractAccountsItem) unpack(enc encoder.Encoder, ht hint.Hint, bks []byte, bam []byte, sadtype string) error {
	e := util.StringErrorFunc("failed to unmarshal BaseCreateContractAccountsItem")

	it.BaseHinter = hint.NewBaseHinter(ht)

	if hinter, err := enc.Decode(bks); err != nil {
		return e(err, "")
	} else if k, ok := hinter.(mitumcurrency.AccountKeys); !ok {
		return e(util.ErrWrongType.Errorf("expected AccountsKeys, not %T", hinter), "")
	} else {
		it.keys = k
	}

	ham, err := enc.DecodeSlice(bam)
	if err != nil {
		return e(err, "")
	}

	amounts := make([]mitumcurrency.Amount, len(ham))
	for i := range ham {
		j, ok := ham[i].(mitumcurrency.Amount)
		if !ok {
			return e(util.ErrWrongType.Errorf("expected Amount, not %T", ham[i]), "")
		}

		amounts[i] = j
	}

	it.amounts = amounts
	it.addressType = hint.Type(sadtype)

	return nil
}
