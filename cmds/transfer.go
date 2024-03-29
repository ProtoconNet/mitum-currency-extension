package cmds

import (
	"context"

	"github.com/pkg/errors"

	currency "github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
)

type TransferCommand struct {
	baseCommand
	OperationFlags
	Sender   AddressFlag          `arg:"" name:"sender" help:"sender address" required:"true"`
	Receiver AddressFlag          `arg:"" name:"receiver" help:"receiver address" required:"true"`
	Amounts  []CurrencyAmountFlag `arg:"" name:"currency-amount" help:"amount (ex: \"<currency>,<amount>\")"`
	sender   base.Address
	receiver base.Address
}

func NewTransferCommand() TransferCommand {
	cmd := NewbaseCommand()
	return TransferCommand{
		baseCommand: *cmd,
	}
}

func (cmd *TransferCommand) Run(pctx context.Context) error {
	if _, err := cmd.prepare(pctx); err != nil {
		return err
	}

	encs = cmd.encs
	enc = cmd.enc

	if err := cmd.parseFlags(); err != nil {
		return err
	}

	op, err := cmd.createOperation()
	if err != nil {
		return err
	}

	PrettyPrint(cmd.Out, op)

	return nil
}

func (cmd *TransferCommand) parseFlags() error {
	if err := cmd.OperationFlags.IsValid(nil); err != nil {
		return err
	}

	if len(cmd.Amounts) < 1 {
		return errors.Errorf("empty currency-amount, must be given at least one")
	}

	if sender, err := cmd.Sender.Encode(enc); err != nil {
		return errors.Wrapf(err, "invalid sender format, %q", cmd.Sender.String())
	} else if receiver, err := cmd.Receiver.Encode(enc); err != nil {
		return errors.Wrapf(err, "invalid receiver format, %q", cmd.Receiver.String())
	} else {
		cmd.sender = sender
		cmd.receiver = receiver
	}

	return nil
}

func (cmd *TransferCommand) createOperation() (base.Operation, error) { // nolint:dupl
	var items []currency.TransfersItem

	ams := make([]currency.Amount, len(cmd.Amounts))
	for i := range cmd.Amounts {
		a := cmd.Amounts[i]
		am := currency.NewAmount(a.Big, a.CID)
		if err := am.IsValid(nil); err != nil {
			return nil, err
		}

		ams[i] = am
	}

	item := currency.NewTransfersItemMultiAmounts(cmd.receiver, ams)
	if err := item.IsValid(nil); err != nil {
		return nil, err
	}
	items = append(items, item)

	fact := currency.NewTransfersFact([]byte(cmd.Token), cmd.sender, items)

	op, err := currency.NewTransfers(fact)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create transfers operation")
	}
	err = op.HashSign(cmd.Privatekey, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create transfers operation")
	}

	return op, nil
}
