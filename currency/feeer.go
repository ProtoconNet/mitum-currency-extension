package currency

import (
	"bytes"
	"encoding/binary"

	mitumcurrency "github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

const (
	FeeerNil   = "nil"
	FeeerFixed = "fixed"
	FeeerRatio = "ratio"
)

var (
	NilFeeerHint   = hint.MustNewHint("mitum-currency-nil-feeer-v0.0.1")
	FixedFeeerHint = hint.MustNewHint("mitum-currency-fixed-feeer-v0.0.1")
	RatioFeeerHint = hint.MustNewHint("mitum-currency-ratio-feeer-v0.0.1")
)

var UnlimitedMaxFeeAmount = mitumcurrency.NewBig(-1)

type Feeer interface {
	util.IsValider
	hint.Hinter
	Type() string
	Bytes() []byte
	Receiver() base.Address
	Min() mitumcurrency.Big
	ExchangeMin() mitumcurrency.Big
	Fee(mitumcurrency.Big) (mitumcurrency.Big, error)
}

type NilFeeer struct {
	hint.BaseHinter
}

func NewNilFeeer() NilFeeer {
	return NilFeeer{BaseHinter: hint.NewBaseHinter(NilFeeerHint)}
}

func (NilFeeer) Type() string {
	return FeeerNil
}

func (NilFeeer) Bytes() []byte {
	return nil
}

func (NilFeeer) Receiver() base.Address {
	return nil
}

func (NilFeeer) Min() mitumcurrency.Big {
	return mitumcurrency.ZeroBig
}

func (NilFeeer) ExchangeMin() mitumcurrency.Big {
	return mitumcurrency.ZeroBig
}

func (NilFeeer) Fee(mitumcurrency.Big) (mitumcurrency.Big, error) {
	return mitumcurrency.ZeroBig, nil
}

func (fa NilFeeer) IsValid([]byte) error {
	return fa.BaseHinter.IsValid(nil)
}

type FixedFeeer struct {
	hint.BaseHinter
	receiver    base.Address
	amount      mitumcurrency.Big
	exchangeMin mitumcurrency.Big
}

func NewFixedFeeer(receiver base.Address, amount mitumcurrency.Big, exchangeMin mitumcurrency.Big) FixedFeeer {
	return FixedFeeer{
		BaseHinter:  hint.NewBaseHinter(FixedFeeerHint),
		receiver:    receiver,
		amount:      amount,
		exchangeMin: exchangeMin,
	}
}

func (FixedFeeer) Type() string {
	return FeeerFixed
}

func (fa FixedFeeer) Bytes() []byte {
	return util.ConcatBytesSlice(fa.receiver.Bytes(), fa.amount.Bytes(), fa.exchangeMin.Bytes())
}

func (fa FixedFeeer) Receiver() base.Address {
	return fa.receiver
}

func (fa FixedFeeer) Min() mitumcurrency.Big {
	return fa.amount
}

func (fa FixedFeeer) ExchangeMin() mitumcurrency.Big {
	return fa.exchangeMin
}

func (fa FixedFeeer) Fee(mitumcurrency.Big) (mitumcurrency.Big, error) {
	if fa.isZero() {
		return mitumcurrency.ZeroBig, nil
	}

	return fa.amount, nil
}

func (fa FixedFeeer) IsValid([]byte) error {
	if err := fa.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := util.CheckIsValiders(nil, false, fa.receiver); err != nil {
		return util.ErrInvalid.Errorf("invalid receiver for fixed feeer: %w", err)
	}

	if !fa.amount.OverNil() {
		return util.ErrInvalid.Errorf("fixed feeer amount under zero")
	}

	return nil
}

func (fa FixedFeeer) isZero() bool {
	return fa.amount.IsZero()
}

type RatioFeeer struct {
	hint.BaseHinter
	receiver    base.Address
	ratio       float64 // 0 >=, or <= 1.0
	min         mitumcurrency.Big
	max         mitumcurrency.Big
	exchangeMin mitumcurrency.Big
}

func NewRatioFeeer(receiver base.Address, ratio float64, min, max, exchangeMin mitumcurrency.Big) RatioFeeer {
	return RatioFeeer{
		BaseHinter:  hint.NewBaseHinter(RatioFeeerHint),
		receiver:    receiver,
		ratio:       ratio,
		min:         min,
		max:         max,
		exchangeMin: exchangeMin,
	}
}

func (RatioFeeer) Type() string {
	return FeeerRatio
}

func (fa RatioFeeer) Bytes() []byte {
	var rb bytes.Buffer
	_ = binary.Write(&rb, binary.BigEndian, fa.ratio)

	return util.ConcatBytesSlice(fa.receiver.Bytes(), rb.Bytes(), fa.min.Bytes(), fa.max.Bytes(), fa.exchangeMin.Bytes())
}

func (fa RatioFeeer) Receiver() base.Address {
	return fa.receiver
}

func (fa RatioFeeer) Min() mitumcurrency.Big {
	return fa.min
}

func (fa RatioFeeer) ExchangeMin() mitumcurrency.Big {
	return fa.exchangeMin
}

func (fa RatioFeeer) Fee(a mitumcurrency.Big) (mitumcurrency.Big, error) {
	if fa.isZero() {
		return mitumcurrency.ZeroBig, nil
	} else if a.IsZero() {
		return fa.min, nil
	}

	if fa.isOne() {
		return a, nil
	} else if f := a.MulFloat64(fa.ratio); f.Compare(fa.min) < 0 {
		return fa.min, nil
	} else {
		if !fa.isUnlimited() && f.Compare(fa.max) > 0 {
			return fa.max, nil
		}
		return f, nil
	}
}

func (fa RatioFeeer) IsValid([]byte) error {
	if err := fa.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := util.CheckIsValiders(nil, false, fa.receiver); err != nil {
		return util.ErrInvalid.Errorf("invalid receiver for ratio feeer: %w", err)
	}

	if fa.ratio < 0 || fa.ratio > 1 {
		return util.ErrInvalid.Errorf("invalid ratio, %v; it should be 0 >=, <= 1", fa.ratio)
	}

	if !fa.min.OverNil() {
		return util.ErrInvalid.Errorf("ratio feeer min amount under zero")
	} else if !fa.max.Equal(UnlimitedMaxFeeAmount) {
		if !fa.max.OverNil() {
			return util.ErrInvalid.Errorf("ratio feeer max amount under zero")
		} else if fa.min.Compare(fa.max) > 0 {
			return util.ErrInvalid.Errorf("ratio feeer min should over max")
		}
	}

	return nil
}

func (fa RatioFeeer) isUnlimited() bool {
	return fa.max.Equal(UnlimitedMaxFeeAmount)
}

func (fa RatioFeeer) isZero() bool {
	return fa.ratio == 0
}

func (fa RatioFeeer) isOne() bool {
	return fa.ratio == 1
}

func NewFeeToken(feeer Feeer, height base.Height) []byte {
	return util.ConcatBytesSlice(feeer.Bytes(), height.Bytes())
}
