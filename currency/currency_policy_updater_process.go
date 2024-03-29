package currency

import (
	"context"
	"sync"

	mitumcurrency "github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/isaac"
	"github.com/ProtoconNet/mitum2/util"
)

var currencyPolicyUpdaterProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(CurrencyPolicyUpdaterProcessor)
	},
}

func (CurrencyPolicy) Process(
	ctx context.Context, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	// NOTE Process is nil func
	return nil, nil, nil
}

type CurrencyPolicyUpdaterProcessor struct {
	*base.BaseOperationProcessor
	suffrage  base.Suffrage
	threshold base.Threshold
}

func NewCurrencyPolicyUpdaterProcessor(threshold base.Threshold) GetNewProcessor {
	return func(height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringErrorFunc("failed to create new CurrencyPolicyUpdaterProcessor")

		nopp := currencyPolicyUpdaterProcessorPool.Get()
		opp, ok := nopp.(*CurrencyPolicyUpdaterProcessor)
		if !ok {
			return nil, e(nil, "expected CurrencyPolicyUpdaterProcessor, not %T", nopp)
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

func (opp *CurrencyPolicyUpdaterProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	e := util.StringErrorFunc("failed to preprocess CurrencyPolicyUpdater")

	nop, ok := op.(CurrencyPolicyUpdater)
	if !ok {
		return ctx, nil, e(nil, "expected CurrencyPolicyUpdater, not %T", op)
	}

	if err := base.CheckFactSignsBySuffrage(opp.suffrage, opp.threshold, nop.NodeSigns()); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("not enough signs: %w", err), nil
	}

	fact, ok := op.Fact().(CurrencyPolicyUpdaterFact)
	if !ok {
		return ctx, nil, e(nil, "expected CurrencyPolicyUpdaterFact, not %T", op.Fact())
	}

	err := checkExistsState(StateKeyCurrencyDesign(fact.currency), getStateFunc)
	if err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("currency not found, %q: %w", fact.currency, err), nil
	}

	if receiver := fact.policy.Feeer().Receiver(); receiver != nil {
		if err := checkExistsState(mitumcurrency.StateKeyAccount(receiver), getStateFunc); err != nil {
			return ctx, base.NewBaseOperationProcessReasonError("feeer receiver not found, %q: %w", fact.policy.Feeer().Receiver(), err), nil
		}
	}

	return ctx, nil, nil
}

func (opp *CurrencyPolicyUpdaterProcessor) Process(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringErrorFunc("failed to process CurrencyPolicyUpdater")

	fact, ok := op.Fact().(CurrencyPolicyUpdaterFact)
	if !ok {
		return nil, nil, e(nil, "expected CurrencyPolicyUpdaterFact, not %T", op.Fact())
	}

	sts := make([]base.StateMergeValue, 1)

	st, err := existsState(StateKeyCurrencyDesign(fact.currency), "key of currency design", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("currency not found, %q: %w", fact.currency, err), nil
	}

	de, err := StateCurrencyDesignValue(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("failed to get currency design value, %q: %w", fact.currency, err), nil
	}

	de.policy = fact.policy

	c := NewCurrencyDesignStateMergeValue(
		st.Key(),
		NewCurrencyDesignStateValue(de),
	)
	sts[0] = c

	return sts, nil, nil
}

func (opp *CurrencyPolicyUpdaterProcessor) Close() error {
	opp.suffrage = nil
	opp.threshold = 0

	currencyPolicyUpdaterProcessorPool.Put(opp)

	return nil
}
