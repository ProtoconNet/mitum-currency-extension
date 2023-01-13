package cmds

type SealCommand struct {
	CreateAccount         CreateAccountCommand         `cmd:"" name:"create-account" help:"create new account"`
	KeyUpdater            KeyUpdaterCommand            `cmd:"" name:"key-updater" help:"update account keys"`
	Transfer              TransferCommand              `cmd:"" name:"transfer" help:"transfer amounts to receiver"`
	CreateContractAccount CreateContractAccountCommand `cmd:"" name:"create-contract-account" help:"create new contract account"`
	Withdraw              WithdrawCommand              `cmd:"" name:"withdraw" help:"withdraw amounts from target contract account"`
	// CurrencyRegister      CurrencyRegisterCommand      `cmd:"" name:"currency-register" help:"register new currency"`
	// CurrencyPolicyUpdater CurrencyPolicyUpdaterCommand `cmd:"" name:"currency-policy-updater" help:"update currency policy"`
	// SuffrageInflation     SuffrageInflationCommand     `cmd:"" name:"suffrage-inflation" help:"suffrage inflation operation"` // revive:disable-line:line-length-limit
}

func NewSealCommand() SealCommand {
	return SealCommand{
		CreateAccount:         NewCreateAccountCommand(),
		KeyUpdater:            NewKeyUpdaterCommand(),
		Transfer:              NewTransferCommand(),
		CreateContractAccount: NewCreateContractAccountCommand(),
		Withdraw:              NewWithdrawCommand(),
		// CurrencyRegister:      NewCurrencyRegisterCommand(),
		// CurrencyPolicyUpdater: NewCurrencyPolicyUpdaterCommand(),
		// SuffrageInflation:     NewSuffrageInflationCommand(),
	}
}
