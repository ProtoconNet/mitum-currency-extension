package cmds

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/pkg/errors"
	currencycmds "github.com/spikeekips/mitum-currency/cmds"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/operation"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/isvalid"
)

type CurrencyFixedFeeerFlags struct {
	Receiver          AddressFlag          `name:"receiver" help:"fee receiver account address"`
	Amount            currencycmds.BigFlag `name:"amount" help:"fee amount"`
	ExchangeMinAmount currencycmds.BigFlag `name:"exchange-min-amount" help:"exchange fee amount"`
	feeer             extensioncurrency.Feeer
}

func (fl *CurrencyFixedFeeerFlags) IsValid([]byte) error {
	if len(fl.Receiver.String()) < 1 {
		return nil
	}

	var receiver base.Address
	if a, err := fl.Receiver.Encode(jenc); err != nil {
		return isvalid.InvalidError.Errorf("invalid receiver format, %q: %w", fl.Receiver.String(), err)
	} else if err := a.IsValid(nil); err != nil {
		return isvalid.InvalidError.Errorf("invalid receiver address, %q: %w", fl.Receiver.String(), err)
	} else {
		receiver = a
	}

	fl.feeer = extensioncurrency.NewFixedFeeer(receiver, fl.Amount.Big, fl.ExchangeMinAmount.Big)
	return fl.feeer.IsValid(nil)
}

type CurrencyRatioFeeerFlags struct {
	Receiver          AddressFlag          `name:"receiver" help:"fee receiver account address"`
	Ratio             float64              `name:"ratio" help:"fee ratio, multifly by operation amount"`
	Min               currencycmds.BigFlag `name:"min" help:"minimum fee"`
	Max               currencycmds.BigFlag `name:"max" help:"maximum fee"`
	ExchangeMinAmount currencycmds.BigFlag `name:"exchange-min-amount" help:"exchange fee amount"`
	feeer             extensioncurrency.Feeer
}

func (fl *CurrencyRatioFeeerFlags) IsValid([]byte) error {
	if len(fl.Receiver.String()) < 1 {
		return nil
	}

	var receiver base.Address
	if a, err := fl.Receiver.Encode(jenc); err != nil {
		return isvalid.InvalidError.Errorf("invalid receiver format, %q: %w", fl.Receiver.String(), err)
	} else if err := a.IsValid(nil); err != nil {
		return isvalid.InvalidError.Errorf("invalid receiver address, %q: %w", fl.Receiver.String(), err)
	} else {
		receiver = a
	}

	fl.feeer = extensioncurrency.NewRatioFeeer(receiver, fl.Ratio, fl.Min.Big, fl.Max.Big, fl.ExchangeMinAmount.Big)
	return fl.feeer.IsValid(nil)
}

type CurrencyPolicyFlags struct {
	NewAccountMinBalance currencycmds.BigFlag `name:"new-account-min-balance" help:"minimum balance for new account"` // nolint lll
}

func (*CurrencyPolicyFlags) IsValid([]byte) error {
	return nil
}

type CurrencyDesignFlags struct {
	Currency                currencycmds.CurrencyIDFlag `arg:"" name:"currency-id" help:"currency id" required:"true"`
	GenesisAmount           currencycmds.BigFlag        `arg:"" name:"genesis-amount" help:"genesis amount" required:"true"`
	GenesisAccount          AddressFlag                 `arg:"" name:"genesis-account" help:"genesis-account address for genesis balance" required:"true"` // nolint lll
	CurrencyPolicyFlags     `prefix:"policy-" help:"currency policy" required:"true"`
	FeeerString             string `name:"feeer" help:"feeer type, {nil, fixed, ratio}" required:"true"`
	CurrencyFixedFeeerFlags `prefix:"feeer-fixed-" help:"fixed feeer"`
	CurrencyRatioFeeerFlags `prefix:"feeer-ratio-" help:"ratio feeer"`
	CurrencyDesign          extensioncurrency.CurrencyDesign
}

func (fl *CurrencyDesignFlags) IsValid([]byte) error {
	if err := fl.CurrencyPolicyFlags.IsValid(nil); err != nil {
		return err
	} else if err := fl.CurrencyFixedFeeerFlags.IsValid(nil); err != nil {
		return err
	} else if err := fl.CurrencyRatioFeeerFlags.IsValid(nil); err != nil {
		return err
	}

	var feeer extensioncurrency.Feeer
	switch t := fl.FeeerString; t {
	case extensioncurrency.FeeerNil, "":
		feeer = extensioncurrency.NewNilFeeer()
	case extensioncurrency.FeeerFixed:
		feeer = fl.CurrencyFixedFeeerFlags.feeer
	case extensioncurrency.FeeerRatio:
		feeer = fl.CurrencyRatioFeeerFlags.feeer
	default:
		return isvalid.InvalidError.Errorf("unknown feeer type, %q", t)
	}

	if feeer == nil {
		return isvalid.InvalidError.Errorf("empty feeer flags")
	} else if err := feeer.IsValid(nil); err != nil {
		return err
	}

	po := extensioncurrency.NewCurrencyPolicy(fl.CurrencyPolicyFlags.NewAccountMinBalance.Big, feeer)
	if err := po.IsValid(nil); err != nil {
		return err
	}

	var genesisAccount base.Address
	if a, err := fl.GenesisAccount.Encode(jenc); err != nil {
		return isvalid.InvalidError.Errorf("invalid genesis-account format, %q: %w", fl.GenesisAccount.String(), err)
	} else if err := a.IsValid(nil); err != nil {
		return isvalid.InvalidError.Errorf("invalid genesis-account address, %q: %w", fl.GenesisAccount.String(), err)
	} else {
		genesisAccount = a
	}

	am := currency.NewAmount(fl.GenesisAmount.Big, fl.Currency.CID)
	if err := am.IsValid(nil); err != nil {
		return err
	}

	fl.CurrencyDesign = extensioncurrency.NewCurrencyDesign(am, genesisAccount, po)
	return fl.CurrencyDesign.IsValid(nil)
}

type CurrencyRegisterCommand struct {
	*BaseCommand
	OperationFlags
	CurrencyDesignFlags
}

func NewCurrencyRegisterCommand() CurrencyRegisterCommand {
	return CurrencyRegisterCommand{
		BaseCommand: NewBaseCommand("currency-register-operation"),
	}
}

func (cmd *CurrencyRegisterCommand) Run(version util.Version) error { // nolint:dupl
	if err := cmd.Initialize(cmd, version); err != nil {
		return errors.Wrap(err, "failed to initialize command")
	}

	if err := cmd.parseFlags(); err != nil {
		return err
	}

	var op operation.Operation
	if i, err := cmd.createOperation(); err != nil {
		return errors.Wrap(err, "failed to create currency-register operation")
	} else if err := i.IsValid([]byte(cmd.OperationFlags.NetworkID)); err != nil {
		return errors.Wrap(err, "invalid currency-register operation")
	} else {
		cmd.Log().Debug().Interface("operation", i).Msg("operation loaded")

		op = i
	}

	i, err := operation.NewBaseSeal(
		cmd.OperationFlags.Privatekey,
		[]operation.Operation{op},
		[]byte(cmd.OperationFlags.NetworkID),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create operation.Seal")
	}
	cmd.Log().Debug().Interface("seal", i).Msg("seal loaded")

	PrettyPrint(cmd.Out, cmd.Pretty, i)

	return nil
}

func (cmd *CurrencyRegisterCommand) parseFlags() error {
	if err := cmd.OperationFlags.IsValid(nil); err != nil {
		return err
	} else if err := cmd.CurrencyDesignFlags.IsValid(nil); err != nil {
		return err
	}

	cmd.Log().Debug().Interface("currency-design", cmd.CurrencyDesignFlags.CurrencyDesign).Msg("currency design loaded")

	return nil
}

func (cmd *CurrencyRegisterCommand) createOperation() (extensioncurrency.CurrencyRegister, error) {
	fact := extensioncurrency.NewCurrencyRegisterFact([]byte(cmd.Token), cmd.CurrencyDesign)

	var fs []base.FactSign
	sig, err := base.NewFactSignature(
		cmd.OperationFlags.Privatekey,
		fact,
		[]byte(cmd.OperationFlags.NetworkID),
	)
	if err != nil {
		return extensioncurrency.CurrencyRegister{}, err
	}
	fs = append(fs, base.NewBaseFactSign(cmd.OperationFlags.Privatekey.Publickey(), sig))

	return extensioncurrency.NewCurrencyRegister(fact, fs, cmd.OperationFlags.Memo)
}