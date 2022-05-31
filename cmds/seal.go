package cmds

import currencycmds "github.com/spikeekips/mitum-currency/cmds"

type SealCommand struct {
	Send                  SendCommand                           `cmd:"" name:"send" help:"send seal to remote mitum node"`
	CreateAccount         currencycmds.CreateAccountCommand     `cmd:"" name:"create-account" help:"create new account"`
	CreateContractAccount CreateContractAccountCommand          `cmd:"" name:"create-contract-account" help:"create new contract account"`
	Deactivate            DeactivateCommand                     `cmd:"" name:"deactivate" help:"deactivate contract account"`
	Withdraw              WithdrawCommand                       `cmd:"" name:"withdraw" help:"withdraw contract account"`
	Transfer              currencycmds.TransferCommand          `cmd:"" name:"transfer" help:"transfer big"`
	KeyUpdater            currencycmds.KeyUpdaterCommand        `cmd:"" name:"key-updater" help:"update keys"`
	CurrencyRegister      CurrencyRegisterCommand               `cmd:"" name:"currency-register" help:"register new currency"`
	CurrencyPolicyUpdater CurrencyPolicyUpdaterCommand          `cmd:"" name:"currency-policy-updater" help:"update currency policy"`  // revive:disable-line:line-length-limit
	SuffrageInflation     currencycmds.SuffrageInflationCommand `cmd:"" name:"suffrage-inflation" help:"suffrage inflation operation"` // revive:disable-line:line-length-limit
	Sign                  currencycmds.SignSealCommand          `cmd:"" name:"sign" help:"sign seal"`
	SignFact              currencycmds.SignFactCommand          `cmd:"" name:"sign-fact" help:"sign facts of operation seal"`
}

func NewSealCommand() SealCommand {
	return SealCommand{
		Send:                  NewSendCommand(),
		CreateAccount:         currencycmds.NewCreateAccountCommand(),
		CreateContractAccount: NewCreateContractAccountCommand(),
		Deactivate:            NewDeactivateCommand(),
		Withdraw:              NewWithdrawCommand(),
		Transfer:              currencycmds.NewTransferCommand(),
		KeyUpdater:            currencycmds.NewKeyUpdaterCommand(),
		CurrencyRegister:      NewCurrencyRegisterCommand(),
		CurrencyPolicyUpdater: NewCurrencyPolicyUpdaterCommand(),
		SuffrageInflation:     currencycmds.NewSuffrageInflationCommand(),
		Sign:                  currencycmds.NewSignSealCommand(),
		SignFact:              currencycmds.NewSignFactCommand(),
	}
}
