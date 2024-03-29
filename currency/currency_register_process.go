package currency

import (
	"context"
	"sync"

	mitumcurrency "github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/isaac"
	"github.com/ProtoconNet/mitum2/util"
)

var currencyRegisterProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(CurrencyRegisterProcessor)
	},
}

func (CurrencyRegister) Process(
	ctx context.Context, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	// NOTE Process is nil func
	return nil, nil, nil
}

type CurrencyRegisterProcessor struct {
	*base.BaseOperationProcessor
	suffrage  base.Suffrage
	threshold base.Threshold
}

func NewCurrencyRegisterProcessor(threshold base.Threshold) GetNewProcessor {
	return func(height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringErrorFunc("failed to create new CurrencyRegisterProcessor")

		nopp := currencyRegisterProcessorPool.Get()
		opp, ok := nopp.(*CurrencyRegisterProcessor)
		if !ok {
			return nil, e(nil, "expected CurrencyRegisterProcessor, not %T", nopp)
		}

		b, err := base.NewBaseOperationProcessor(
			height, getStateFunc, newPreProcessConstraintFunc, newProcessConstraintFunc)
		if err != nil {
			return nil, e(err, "")
		}

		opp.BaseOperationProcessor = b
		opp.threshold = threshold

		switch i, found, err := getStateFunc(isaac.SuffrageStateKey); {
		case err != nil:
			return nil, e(err, "")
		case !found, i == nil:
			return nil, e(isaac.ErrStopProcessingRetry.Errorf("empty state"), "")
		default:
			sufstv := i.Value().(base.SuffrageNodesStateValue) //nolint:forcetypeassert //...

			suf, err := sufstv.Suffrage()
			if err != nil {
				return nil, e(isaac.ErrStopProcessingRetry.Errorf("failed to get suffrage from state"), "")
			}

			opp.suffrage = suf
		}

		return opp, nil
	}
}

func (opp *CurrencyRegisterProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	e := util.StringErrorFunc("failed to preprocess CurrencyRegister")

	nop, ok := op.(CurrencyRegister)
	if !ok {
		return ctx, nil, e(nil, "expected CurrencyRegister, not %T", op)
	}

	fact, ok := op.Fact().(CurrencyRegisterFact)
	if !ok {
		return ctx, nil, e(nil, "expected CurrencyRegisterFact, not %T", op.Fact())
	}

	if err := base.CheckFactSignsBySuffrage(opp.suffrage, opp.threshold, nop.NodeSigns()); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("not enough signs: %w", err), nil
	}

	item := fact.currency

	if err := checkNotExistsState(StateKeyCurrencyDesign(item.Currency()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("currency design already exists, %q: %w", item.Currency(), err), nil
	}

	if err := checkExistsState(mitumcurrency.StateKeyAccount(item.genesisAccount), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("genesis account not found, %q: %w", item.genesisAccount, err), nil
	}

	if err := checkNotExistsState(StateKeyContractAccount(item.genesisAccount), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("contract account cannot be genesis account of currency, %q: %w", item.genesisAccount, err), nil
	}

	if receiver := item.Policy().Feeer().Receiver(); receiver != nil {
		if err := checkExistsState(mitumcurrency.StateKeyAccount(receiver), getStateFunc); err != nil {
			return ctx, base.NewBaseOperationProcessReasonError("feeer receiver not found, %q: %w", receiver, err), nil
		}

		if err := checkNotExistsState(StateKeyContractAccount(receiver), getStateFunc); err != nil {
			return ctx, base.NewBaseOperationProcessReasonError("contract account cannot be fee receiver, %q: %w", receiver, err), nil
		}
	}

	if err := checkNotExistsState(mitumcurrency.StateKeyBalance(item.genesisAccount, item.Currency()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("account balance already exists, %q: %w", mitumcurrency.StateKeyBalance(item.genesisAccount, item.Currency()), err), nil
	}

	return ctx, nil, nil
}

func (opp *CurrencyRegisterProcessor) Process(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringErrorFunc("failed to process CurrencyRegister")

	fact, ok := op.Fact().(CurrencyRegisterFact)
	if !ok {
		return nil, nil, e(nil, "expected CurrencyRegisterFact, not %T", op.Fact())
	}

	sts := make([]base.StateMergeValue, 4)

	item := fact.currency

	ba := mitumcurrency.NewBalanceStateValue(item.amount)
	sts[0] = mitumcurrency.NewBalanceStateMergeValue(
		mitumcurrency.StateKeyBalance(item.genesisAccount, item.Currency()),
		ba,
	)

	de := NewCurrencyDesignStateValue(item)
	sts[1] = NewCurrencyDesignStateMergeValue(StateKeyCurrencyDesign(item.Currency()), de)

	{
		l, err := createZeroAccount(item.Currency(), getStateFunc)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to create zero account, %q: %w", item.Currency(), err), nil
		}
		sts[2], sts[3] = l[0], l[1]
	}

	return sts, nil, nil
}

func (opp *CurrencyRegisterProcessor) Close() error {
	opp.suffrage = nil
	opp.threshold = 0

	currencyRegisterProcessorPool.Put(opp)

	return nil
}
