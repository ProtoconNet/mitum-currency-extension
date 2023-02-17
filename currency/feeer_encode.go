package currency

import (
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/encoder"
	"github.com/spikeekips/mitum/util/hint"
)

func (fa *FixedFeeer) unpack(enc encoder.Encoder, ht hint.Hint, rc string, am string, em string) error {
	e := util.StringErrorFunc("failed to unmarshal FixedFeeer")

	fa.BaseHinter = hint.NewBaseHinter(ht)

	switch a, err := base.DecodeAddress(rc, enc); {
	case err != nil:
		return e(err, "")
	default:
		fa.receiver = a
	}

	if big, err := currency.NewBigFromString(am); err != nil {
		return e(err, "")
	} else {
		fa.amount = big
	}

	if exm, err := currency.NewBigFromString(em); err != nil {
		return e(err, "")
	} else {
		fa.exchangeMin = exm
	}

	return nil
}

func (fa *RatioFeeer) unpack(
	enc encoder.Encoder,
	ht hint.Hint,
	rc string,
	ratio float64,
	min, max, em string,
) error {
	e := util.StringErrorFunc("failed to unmarshal RatioFeeer")

	fa.BaseHinter = hint.NewBaseHinter(ht)

	switch a, err := base.DecodeAddress(rc, enc); {
	case err != nil:
		return e(err, "")
	default:
		fa.receiver = a
	}

	fa.ratio = ratio

	if min, err := currency.NewBigFromString(min); err != nil {
		return e(err, "")
	} else {
		fa.min = min
	}

	if max, err := currency.NewBigFromString(max); err != nil {
		return e(err, "")
	} else {
		fa.max = max
	}

	if exm, err := currency.NewBigFromString(em); err != nil {
		return e(err, "")
	} else {
		fa.exchangeMin = exm
	}

	return nil
}
