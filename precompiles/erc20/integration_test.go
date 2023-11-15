package erc20_test

import (
	"fmt"
	"math/big"

	auth "github.com/evmos/evmos/v15/precompiles/authorization"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/evmos/evmos/v15/precompiles/erc20"
	"github.com/evmos/evmos/v15/precompiles/erc20/testdata"
	"github.com/evmos/evmos/v15/precompiles/testutil"
	"github.com/evmos/evmos/v15/testutil/integration/evmos/factory"
	"github.com/evmos/evmos/v15/testutil/integration/evmos/keyring"
	utiltx "github.com/evmos/evmos/v15/testutil/tx"
	evmtypes "github.com/evmos/evmos/v15/x/evm/types"

	//nolint:revive // dot imports are fine for Ginkgo
	. "github.com/onsi/ginkgo/v2"
	//nolint:revive // dot imports are fine for Ginkgo
	. "github.com/onsi/gomega"
)

var _ = Describe("ERC20 Extension -", func() {
	var (
		contractAddr common.Address
		err          error
		sender       keyring.Key

		execRevertedCheck testutil.LogCheckArgs
		failCheck         testutil.LogCheckArgs
		passCheck         testutil.LogCheckArgs
	)

	BeforeEach(func() {
		s.SetupTest()

		sender = s.keyring.GetKey(0)

		contractAddr, err = s.factory.DeployContract(
			sender.Priv,
			evmtypes.EvmTxArgs{}, // NOTE: passing empty struct to use default values
			factory.ContractDeploymentData{
				Contract:        testdata.ERC20CallerContract,
				ConstructorArgs: []interface{}{s.precompile.Address()},
			},
		)
		Expect(err).ToNot(HaveOccurred(), "failed to deploy contract")

		failCheck = testutil.LogCheckArgs{ABIEvents: s.precompile.Events}
		execRevertedCheck = failCheck.WithErrContains("execution reverted")
		passCheck = failCheck.WithExpPass(true)

		err = s.network.NextBlock()
		Expect(err).ToNot(HaveOccurred(), "failed to advance block")
	})

	When("querying balance", func() {
		DescribeTable("it should return an existing balance", func(callType int) {
			expBalance := big.NewInt(100)

			// Fund account with some tokens
			err := s.network.FundAccount(sender.AccAddr, sdk.Coins{
				{Denom: s.tokenDenom, Amount: sdk.NewIntFromBigInt(expBalance)},
			})
			Expect(err).ToNot(HaveOccurred(), "failed to fund account")

			// Query the balance
			txArgs, balancesArgs := s.getTxAndCallArgs(callType, contractAddr)
			balancesArgs.MethodName = erc20.BalanceOfMethod
			balancesArgs.Args = []interface{}{sender.Addr}

			_, ethRes, err := s.callContractAndCheckLogs(sender.Priv, txArgs, balancesArgs, passCheck)
			Expect(err).ToNot(HaveOccurred(), "unexpected result calling contract")

			var balance *big.Int
			err = s.precompile.UnpackIntoInterface(&balance, erc20.BalanceOfMethod, ethRes.Ret)
			Expect(err).ToNot(HaveOccurred(), "failed to unpack result")
			Expect(balance).To(Equal(expBalance), "expected different balance")
		},
			Entry(" - direct call", directCall),
			Entry(" - through contract", contractCall),
		)

		DescribeTable("it should return zero if balance only exists for other tokens", func(callType int) {
			address := utiltx.GenerateAddress()

			// Fund account with some tokens
			err := s.network.FundAccount(sender.AccAddr, sdk.Coins{sdk.NewInt64Coin(s.network.GetDenom(), 100)})
			Expect(err).ToNot(HaveOccurred(), "failed to fund account")

			// Query the balance
			txArgs, balancesArgs := s.getTxAndCallArgs(callType, contractAddr)
			balancesArgs.MethodName = erc20.BalanceOfMethod
			balancesArgs.Args = []interface{}{address}

			_, ethRes, err := s.callContractAndCheckLogs(sender.Priv, txArgs, balancesArgs, passCheck)
			Expect(err).ToNot(HaveOccurred(), "unexpected result calling contract")

			var balance *big.Int
			err = s.precompile.UnpackIntoInterface(&balance, erc20.BalanceOfMethod, ethRes.Ret)
			Expect(err).ToNot(HaveOccurred(), "failed to unpack result")
			Expect(balance.Int64()).To(BeZero(), "expected zero balance")
		},
			Entry(" - direct call", directCall),
			Entry(" - through contract", contractCall),
		)

		DescribeTable("it should return zero if the account does not exist", func(callType int) {
			address := utiltx.GenerateAddress()

			// Query the balance
			txArgs, balancesArgs := s.getTxAndCallArgs(callType, contractAddr)
			balancesArgs.MethodName = erc20.BalanceOfMethod
			balancesArgs.Args = []interface{}{address}

			_, ethRes, err := s.callContractAndCheckLogs(sender.Priv, txArgs, balancesArgs, passCheck)
			Expect(err).ToNot(HaveOccurred(), "unexpected result calling contract")

			var balance *big.Int
			err = s.precompile.UnpackIntoInterface(&balance, erc20.BalanceOfMethod, ethRes.Ret)
			Expect(err).ToNot(HaveOccurred(), "failed to unpack result")
			Expect(balance.Int64()).To(BeZero(), "expected zero balance")
		},
			Entry(" - direct call", directCall),
			Entry(" - through contract", contractCall),
		)
	})

	When("querying allowance", func() {
		DescribeTable("it should return an existing allowance", func(callType int) {
			grantee := utiltx.GenerateAddress()
			granter := sender
			expAllowance := big.NewInt(100)

			s.setupSendAuthz(grantee.Bytes(), granter.Priv, sdk.Coins{
				{Denom: s.tokenDenom, Amount: sdk.NewIntFromBigInt(expAllowance)},
			})

			txArgs, allowanceArgs := s.getTxAndCallArgs(callType, contractAddr)
			allowanceArgs.MethodName = auth.AllowanceMethod
			allowanceArgs.Args = []interface{}{granter.Addr, grantee}

			_, ethRes, err := s.callContractAndCheckLogs(granter.Priv, txArgs, allowanceArgs, passCheck)
			Expect(err).ToNot(HaveOccurred(), "unexpected result calling contract")

			var allowance *big.Int
			err = s.precompile.UnpackIntoInterface(&allowance, auth.AllowanceMethod, ethRes.Ret)
			Expect(err).ToNot(HaveOccurred(), "failed to unpack result")
			Expect(allowance).To(Equal(expAllowance), "expected different allowance")
		},
			Entry(" - direct call", directCall),
			Entry(" - through contract", contractCall),
		)

		DescribeTable("it should return an error if no allowance exists", func(callType int) {
			grantee := s.keyring.GetAddr(1)
			granter := sender

			balanceGrantee, err := s.grpcHandler.GetBalance(grantee.Bytes(), s.network.GetDenom())
			Expect(err).ToNot(HaveOccurred(), "failed to get balance")
			Expect(balanceGrantee.Balance.Amount.Int64()).ToNot(BeZero(), "expected zero balance")

			txArgs, allowanceArgs := s.getTxAndCallArgs(callType, contractAddr)
			allowanceArgs.MethodName = auth.AllowanceMethod
			allowanceArgs.Args = []interface{}{granter.Addr, grantee}

			noAuthzCheck := failCheck.WithErrContains(
				fmt.Sprintf(auth.ErrAuthzDoesNotExistOrExpired, erc20.SendMsgURL, grantee.String()),
			)
			if callType == contractCall {
				noAuthzCheck = execRevertedCheck
			}

			_, _, err = s.callContractAndCheckLogs(granter.Priv, txArgs, allowanceArgs, noAuthzCheck)
			Expect(err).ToNot(HaveOccurred(), "unexpected result calling contract")
		},
			Entry(" - direct call", directCall),
			Entry(" - through contract", contractCall),
		)

		DescribeTable("it should return zero if an allowance exists for other tokens", func(callType int) {
			grantee := s.keyring.GetAddr(1)
			granter := sender
			amount := big.NewInt(100)

			s.setupSendAuthz(grantee.Bytes(), granter.Priv, sdk.Coins{
				{Denom: s.network.GetDenom(), Amount: sdk.NewIntFromBigInt(amount)},
			})

			txArgs, allowanceArgs := s.getTxAndCallArgs(callType, contractAddr)
			allowanceArgs.MethodName = auth.AllowanceMethod
			allowanceArgs.Args = []interface{}{granter.Addr, grantee}

			_, ethRes, err := s.callContractAndCheckLogs(granter.Priv, txArgs, allowanceArgs, passCheck)
			Expect(err).ToNot(HaveOccurred(), "unexpected result calling contract")

			var allowance *big.Int
			err = s.precompile.UnpackIntoInterface(&allowance, auth.AllowanceMethod, ethRes.Ret)
			Expect(err).ToNot(HaveOccurred(), "failed to unpack result")
			Expect(allowance.Int64()).To(BeZero(), "expected zero allowance")
		},
			Entry(" - direct call", directCall),
			Entry(" - through contract", contractCall),
		)

		DescribeTable("it should return an error if the account does not exist", func(callType int) {
			grantee := utiltx.GenerateAddress()
			granter := sender

			txArgs, allowanceArgs := s.getTxAndCallArgs(callType, contractAddr)
			allowanceArgs.MethodName = auth.AllowanceMethod
			allowanceArgs.Args = []interface{}{granter.Addr, grantee}

			noAuthzCheck := failCheck.WithErrContains(
				fmt.Sprintf(auth.ErrAuthzDoesNotExistOrExpired, erc20.SendMsgURL, grantee.String()),
			)
			if callType == contractCall {
				noAuthzCheck = execRevertedCheck
			}

			_, _, err = s.callContractAndCheckLogs(granter.Priv, txArgs, allowanceArgs, noAuthzCheck)
			Expect(err).ToNot(HaveOccurred(), "unexpected result calling contract")
		},
			Entry(" - direct call", directCall),
			Entry(" - through contract", contractCall),
		)
	})

	When("querying total supply", func() {
		DescribeTable("it should return the total supply", func(callType int) {
			expSupply := big.NewInt(100)

			// Fund account with some tokens
			err := s.network.FundAccount(sender.AccAddr, sdk.Coins{
				{Denom: s.tokenDenom, Amount: sdk.NewIntFromBigInt(expSupply)},
			})
			Expect(err).ToNot(HaveOccurred(), "failed to fund account")

			// Query the balance
			txArgs, supplyArgs := s.getTxAndCallArgs(callType, contractAddr)
			supplyArgs.MethodName = erc20.TotalSupplyMethod

			_, ethRes, err := s.callContractAndCheckLogs(sender.Priv, txArgs, supplyArgs, passCheck)
			Expect(err).ToNot(HaveOccurred(), "unexpected result calling contract")

			var supply *big.Int
			err = s.precompile.UnpackIntoInterface(&supply, erc20.TotalSupplyMethod, ethRes.Ret)
			Expect(err).ToNot(HaveOccurred(), "failed to unpack result")
			Expect(supply).To(Equal(expSupply), "expected different supply")
		},
			Entry(" - direct call", directCall),
			Entry(" - through contract", contractCall),
		)

		DescribeTable("it should return zero if no tokens exist", func(callType int) {
			txArgs, supplyArgs := s.getTxAndCallArgs(callType, contractAddr)
			supplyArgs.MethodName = erc20.TotalSupplyMethod

			_, ethRes, err := s.callContractAndCheckLogs(sender.Priv, txArgs, supplyArgs, passCheck)
			Expect(err).ToNot(HaveOccurred(), "unexpected result calling contract")

			var supply *big.Int
			err = s.precompile.UnpackIntoInterface(&supply, erc20.TotalSupplyMethod, ethRes.Ret)
			Expect(err).ToNot(HaveOccurred(), "failed to unpack result")
			Expect(supply.Int64()).To(BeZero(), "expected zero supply")
		},
			Entry(" - direct call", directCall),
			Entry(" - through contract", contractCall),
		)
	})

	When("transferring tokens", func() {
		DescribeTable("it should transfer tokens to a non-existing address", func(callType int) {
			receiver := utiltx.GenerateAddress()
			fundAmount := big.NewInt(200)
			amount := big.NewInt(100)

			// Fund account with some tokens
			err := s.network.FundAccount(sender.AccAddr, sdk.Coins{
				{Denom: s.tokenDenom, Amount: sdk.NewIntFromBigInt(fundAmount)},
			})
			Expect(err).ToNot(HaveOccurred(), "failed to fund account")

			senderBalancePre, err := s.grpcHandler.GetBalance(sender.AccAddr, s.tokenDenom)
			Expect(err).ToNot(HaveOccurred(), "failed to get balance")
			Expect(senderBalancePre.Balance.Amount.Int64()).To(Equal(fundAmount.Int64()), "expected different balance before transfer")

			receiverBalancePre, err := s.grpcHandler.GetBalance(receiver.Bytes(), s.tokenDenom)
			Expect(err).ToNot(HaveOccurred(), "failed to get balance")
			Expect(receiverBalancePre.Balance.Amount.Int64()).To(BeZero(), "expected zero balance before transfer")

			// Transfer tokens
			txArgs, transferArgs := s.getTxAndCallArgs(callType, contractAddr)
			transferArgs.MethodName = erc20.TransferMethod
			transferArgs.Args = []interface{}{receiver, amount}

			transferCheck := passCheck.WithExpEvents(erc20.EventTypeTransfer)

			res, ethRes, err := s.callContractAndCheckLogs(sender.Priv, txArgs, transferArgs, transferCheck)
			Expect(err).ToNot(HaveOccurred(), "unexpected result calling contract")

			senderBalancePost, err := s.grpcHandler.GetBalance(sender.AccAddr, s.tokenDenom)
			Expect(err).ToNot(HaveOccurred(), "failed to get balance")
			Expect(senderBalancePost.Balance.Amount.Int64()).To(Equal(senderBalancePre.Balance.Amount.Int64()-amount.Int64()), "expected different balance after transfer")

			receiverBalancePost, err := s.grpcHandler.GetBalance(receiver.Bytes(), s.tokenDenom)
			Expect(err).ToNot(HaveOccurred(), "failed to get balance")
			Expect(receiverBalancePost.Balance.Amount.Int64()).To(Equal(amount.Int64()), "expected different balance after transfer")

			// TODO: Check gas
			println("Gas used (res): ", res.GasUsed)
			println("Gas used (ethRes): ", ethRes.GasUsed)
			// Expect(res.GasUsed).To(Equal(uint64(0)), "expected different gas used")
			// Expect(ethRes.GasUsed).To(Equal(1), "expected different gas used")
		},
			Entry(" - direct call", directCall),
			// NOTE: we are not passing the contract call here because that requires an authorization which is
			// a separate test case.
		)

		DescribeTable("it should transfer tokens to an existing address", func(callType int) {
			receiver := s.keyring.GetKey(1)
			fundAmountSender := big.NewInt(300)
			fundAmountReceiver := big.NewInt(500)
			amount := big.NewInt(100)

			// Fund accounts with some tokens
			err = s.network.FundAccount(sender.AccAddr, sdk.Coins{
				{Denom: s.tokenDenom, Amount: sdk.NewIntFromBigInt(fundAmountSender)},
			})
			Expect(err).ToNot(HaveOccurred(), "failed to fund account")
			err = s.network.FundAccount(receiver.AccAddr, sdk.Coins{
				{Denom: s.tokenDenom, Amount: sdk.NewIntFromBigInt(fundAmountReceiver)},
			})
			Expect(err).ToNot(HaveOccurred(), "failed to fund account")

			senderBalancePre, err := s.grpcHandler.GetBalance(sender.AccAddr, s.tokenDenom)
			Expect(err).ToNot(HaveOccurred(), "failed to get balance")
			Expect(senderBalancePre.Balance.Amount.Int64()).To(Equal(fundAmountSender.Int64()), "expected different balance before transfer")

			receiverBalancePre, err := s.grpcHandler.GetBalance(receiver.AccAddr, s.tokenDenom)
			Expect(err).ToNot(HaveOccurred(), "failed to get balance")
			Expect(receiverBalancePre.Balance.Amount.Int64()).To(Equal(fundAmountReceiver.Int64()), "expected different balance before transfer")

			// Transfer tokens
			txArgs, transferArgs := s.getTxAndCallArgs(callType, contractAddr)
			transferArgs.MethodName = erc20.TransferMethod
			transferArgs.Args = []interface{}{receiver.Addr, amount}

			transferCheck := passCheck.WithExpEvents(erc20.EventTypeTransfer)

			_, _, err = s.callContractAndCheckLogs(sender.Priv, txArgs, transferArgs, transferCheck)
			Expect(err).ToNot(HaveOccurred(), "unexpected result calling contract")

			senderBalancePost, err := s.grpcHandler.GetBalance(sender.AccAddr, s.tokenDenom)
			Expect(err).ToNot(HaveOccurred(), "failed to get balance")
			Expect(senderBalancePost.Balance.Amount.Int64()).To(Equal(senderBalancePre.Balance.Amount.Int64()-amount.Int64()), "expected different balance after transfer")

			receiverBalancePost, err := s.grpcHandler.GetBalance(receiver.AccAddr, s.tokenDenom)
			Expect(err).ToNot(HaveOccurred(), "failed to get balance")
			Expect(receiverBalancePost.Balance.Amount.Int64()).To(Equal(receiverBalancePre.Balance.Amount.Int64()+amount.Int64()), "expected different balance after transfer")
		},
			Entry(" - direct call", directCall),
			// NOTE: we are not passing the contract call here because that requires an authorization which is
			// a separate test case.
		)

		// TODO: is this the behavior we want? Makes sense right because the contract is not a wallet?
		DescribeTable("it should return an error trying to call from a smart contract", func(callType int) {
			receiver := s.keyring.GetAddr(1)
			fundAmount := big.NewInt(300)
			amount := big.NewInt(100)

			// Fund account with some tokens
			err = s.network.FundAccount(sender.AccAddr, sdk.Coins{
				{Denom: s.tokenDenom, Amount: sdk.NewIntFromBigInt(fundAmount)},
			})
			Expect(err).ToNot(HaveOccurred(), "failed to fund account")

			// Transfer tokens
			txArgs, transferArgs := s.getTxAndCallArgs(callType, contractAddr)
			transferArgs.MethodName = erc20.TransferMethod
			transferArgs.Args = []interface{}{receiver, amount}

			_, _, err = s.callContractAndCheckLogs(sender.Priv, txArgs, transferArgs, execRevertedCheck)
			Expect(err).ToNot(HaveOccurred(), "unexpected result calling contract")
		},
			// NOTE: we are not passing the direct call here because this test is specific to the contract calls
			Entry(" - through contract", contractCall),
		)

		DescribeTable("it should return an error if the sender does not have enough tokens", func(callType int) {
			receiver := s.keyring.GetAddr(1)
			fundAmount := big.NewInt(100)
			amount := big.NewInt(200)

			// Fund account with some tokens
			err = s.network.FundAccount(sender.AccAddr, sdk.Coins{
				{Denom: s.tokenDenom, Amount: sdk.NewIntFromBigInt(fundAmount)},
			})
			Expect(err).ToNot(HaveOccurred(), "failed to fund account")

			// Transfer tokens
			txArgs, transferArgs := s.getTxAndCallArgs(callType, contractAddr)
			transferArgs.MethodName = erc20.TransferMethod
			transferArgs.Args = []interface{}{receiver, amount}

			insufficientBalanceCheck := failCheck.WithErrContains(
				"spendable balance 100xmpl is smaller than 200xmpl: insufficient funds",
			)

			_, _, err = s.callContractAndCheckLogs(sender.Priv, txArgs, transferArgs, insufficientBalanceCheck)
			Expect(err).ToNot(HaveOccurred(), "unexpected result calling contract")
		},
			Entry(" - direct call", directCall),
			// NOTE: we are not passing the contract call here because this test is for direct calls only
		)
	})

	When("transferring tokens from another account", func() {
		DescribeTable("it should transfer tokens from another account with a sufficient approval set", func(callType int) {
			owner := sender
			spender := s.keyring.GetKey(1)
			receiver := utiltx.GenerateAddress()

			fundCoin := sdk.NewInt64Coin(s.tokenDenom, 300)
			transferCoin := sdk.NewInt64Coin(s.tokenDenom, 100)

			// Fund account with some tokens
			err = s.network.FundAccount(owner.AccAddr, sdk.Coins{fundCoin})

			// Set allowance
			s.setupSendAuthz(spender.AccAddr, owner.Priv, sdk.Coins{transferCoin})

			// Transfer tokens
			txArgs, transferArgs := s.getTxAndCallArgs(callType, contractAddr)
			transferArgs.MethodName = erc20.TransferFromMethod
			transferArgs.Args = []interface{}{owner.Addr, receiver, transferCoin.Amount.BigInt()}

			transferCheck := passCheck.WithExpEvents(erc20.EventTypeTransfer)

			_, _, err = s.callContractAndCheckLogs(spender.Priv, txArgs, transferArgs, transferCheck)
			Expect(err).ToNot(HaveOccurred(), "unexpected result calling contract")

			expBalances := []ExpectedBalance{
				{address: owner.AccAddr, expCoins: sdk.Coins{fundCoin.Sub(transferCoin)}},
				{address: receiver.Bytes(), expCoins: sdk.Coins{transferCoin}},
			}
			s.ExpectBalances(expBalances)

			// Check that the allowance was removed since we authorized only the transferred amount
			s.expectNoSendAuthz(spender.AccAddr, owner.AccAddr)
		},
			Entry(" - direct call", directCall),
			// NOTE: we are not passing the contract call here because this test is for direct calls only
		)

		DescribeTable("it should transfer tokens using a smart contract with a sufficient approval set", func(callType int) {
			owner := sender
			spender := contractAddr // NOTE: in case of a contract call the spender is the contract itself
			receiver := utiltx.GenerateAddress()
			fundCoin := sdk.NewInt64Coin(s.tokenDenom, 300)
			transferCoin := sdk.NewInt64Coin(s.tokenDenom, 100)

			// Fund account with some tokens
			err = s.network.FundAccount(owner.AccAddr, sdk.Coins{fundCoin})

			// Set allowance
			s.setupSendAuthz(spender.Bytes(), owner.Priv, sdk.Coins{transferCoin})

			// Transfer tokens
			txArgs, transferArgs := s.getTxAndCallArgs(callType, contractAddr)
			transferArgs.MethodName = erc20.TransferFromMethod
			transferArgs.Args = []interface{}{owner.Addr, receiver, transferCoin.Amount.BigInt()}

			transferCheck := passCheck.WithExpEvents(erc20.EventTypeTransfer)

			_, _, err = s.callContractAndCheckLogs(owner.Priv, txArgs, transferArgs, transferCheck)
			Expect(err).ToNot(HaveOccurred(), "unexpected result calling contract")

			expBalances := []ExpectedBalance{
				{address: owner.AccAddr, expCoins: sdk.Coins{fundCoin.Sub(transferCoin)}},
				{address: receiver.Bytes(), expCoins: sdk.Coins{transferCoin}},
			}
			s.ExpectBalances(expBalances)

			// Check that the allowance was removed since we authorized only the transferred amount
			s.expectNoSendAuthz(spender.Bytes(), owner.AccAddr)
		},
			// Entry(" - direct call", directCall),
			Entry(" - through contract", contractCall),
		)

		// TODO: This is working right now! We should probably block this.
		DescribeTable("it should return an error trying to send using a smart contract but triggered from another account", func(callType int) {
			msgSender := s.keyring.GetKey(0)
			owner := s.keyring.GetKey(1)
			receiver := utiltx.GenerateAddress()
			spender := contractAddr

			fundCoin := sdk.NewInt64Coin(s.tokenDenom, 300)
			transferCoin := sdk.NewInt64Coin(s.tokenDenom, 100)

			// Fund account with some tokens
			err = s.network.FundAccount(owner.AccAddr, sdk.Coins{fundCoin})
			Expect(err).ToNot(HaveOccurred(), "failed to fund account")

			// Set allowance
			s.setupSendAuthz(spender.Bytes(), owner.Priv, sdk.Coins{transferCoin})

			// Transfer tokens
			txArgs, transferArgs := s.getTxAndCallArgs(callType, contractAddr)
			transferArgs.MethodName = erc20.TransferFromMethod
			transferArgs.Args = []interface{}{owner.Addr, receiver, transferCoin.Amount.BigInt()}

			_, _, err = s.callContractAndCheckLogs(msgSender.Priv, txArgs, transferArgs, execRevertedCheck)
			Expect(err).ToNot(HaveOccurred(), "unexpected result calling contract")
		},
			// NOTE: we are not passing the direct call here because this test is specific to the contract calls
			Entry(" - through contract", contractCall),
		)

		//nolint:dupl // these tests are not duplicates
		DescribeTable("it should return an error when the spender does not have enough allowance", func(callType int) {
			owner := sender
			spender := s.keyring.GetKey(1)
			receiver := utiltx.GenerateAddress()
			fundCoin := sdk.NewInt64Coin(s.tokenDenom, 300)
			authzCoin := sdk.NewInt64Coin(s.tokenDenom, 100)
			transferCoin := sdk.NewInt64Coin(s.tokenDenom, 200)

			// Fund account with some tokens
			err = s.network.FundAccount(owner.AccAddr, sdk.Coins{fundCoin})

			// Set allowance
			s.setupSendAuthz(spender.AccAddr, owner.Priv, sdk.Coins{authzCoin})

			// Transfer tokens
			txArgs, transferArgs := s.getTxAndCallArgs(callType, contractAddr)
			transferArgs.MethodName = erc20.TransferFromMethod
			transferArgs.Args = []interface{}{owner.Addr, receiver, transferCoin.Amount.BigInt()}

			insufficientAllowanceCheck := failCheck.WithErrContains("requested amount is more than spend limit")

			_, _, err = s.callContractAndCheckLogs(spender.Priv, txArgs, transferArgs, insufficientAllowanceCheck)
			Expect(err).ToNot(HaveOccurred(), "unexpected result calling contract")
		},
			Entry(" - direct call", directCall),
			// NOTE: we are not passing the contract call here because this test case only covers direct calls
		)

		//nolint:dupl // these tests are not duplicates
		DescribeTable("it should return an error when using smart contract and the spender does not have enough allowance", func(callType int) {
			from := sender
			spender := contractAddr // NOTE: in case of a contract call the spender is the contract itself
			receiver := utiltx.GenerateAddress()
			fundAmount := big.NewInt(400)
			authzAmount := big.NewInt(200)
			amount := big.NewInt(300)

			// Fund account with some tokens
			err = s.network.FundAccount(from.AccAddr, sdk.Coins{
				{Denom: s.tokenDenom, Amount: sdk.NewIntFromBigInt(fundAmount)},
			})
			Expect(err).ToNot(HaveOccurred(), "failed to fund account")

			// Set allowance
			s.setupSendAuthz(spender.Bytes(), from.Priv, sdk.Coins{
				{Denom: s.tokenDenom, Amount: sdk.NewIntFromBigInt(authzAmount)},
			})

			// Transfer tokens
			txArgs, transferArgs := s.getTxAndCallArgs(callType, contractAddr)
			transferArgs.MethodName = erc20.TransferFromMethod
			transferArgs.Args = []interface{}{from.Addr, receiver, amount}

			insufficientAllowanceCheck := execRevertedCheck

			_, _, err = s.callContractAndCheckLogs(from.Priv, txArgs, transferArgs, insufficientAllowanceCheck)
			Expect(err).ToNot(HaveOccurred(), "unexpected result calling contract")
		},
			// NOTE: we are not passing the direct call here because this test is for contract calls only
			Entry(" - through contract", contractCall),
		)

		DescribeTable("it should return an error if there is no allowance set", func(callType int) {
			from := s.keyring.GetKey(1)
			receiver := utiltx.GenerateAddress()
			fundAmount := big.NewInt(400)
			amount := big.NewInt(300)

			// Fund account with some tokens
			err = s.network.FundAccount(from.AccAddr, sdk.Coins{
				{Denom: s.tokenDenom, Amount: sdk.NewIntFromBigInt(fundAmount)},
			})
			Expect(err).ToNot(HaveOccurred(), "failed to fund account")

			// Transfer tokens
			txArgs, transferArgs := s.getTxAndCallArgs(callType, contractAddr)
			transferArgs.MethodName = erc20.TransferFromMethod
			transferArgs.Args = []interface{}{from.Addr, receiver, amount}

			insufficientAllowanceCheck := failCheck.WithErrContains(
				"authorization not found",
			)

			_, _, err = s.callContractAndCheckLogs(sender.Priv, txArgs, transferArgs, insufficientAllowanceCheck)
			Expect(err).ToNot(HaveOccurred(), "unexpected result calling contract")
		},
			Entry(" - direct call", directCall),
			// NOTE: we are not passing the contract call here because this test case only covers direct calls
		)

		//nolint:dupl // these tests are not duplicates
		DescribeTable("it should return an error if the sender does not have enough tokens", func(callType int) {
			from := s.keyring.GetKey(1)
			receiver := utiltx.GenerateAddress()
			fundAmount := big.NewInt(200)
			authzAmount := big.NewInt(300)
			amount := big.NewInt(300)

			// Fund account with some tokens
			err = s.network.FundAccount(from.AccAddr, sdk.Coins{
				{Denom: s.tokenDenom, Amount: sdk.NewIntFromBigInt(fundAmount)},
			})
			Expect(err).ToNot(HaveOccurred(), "failed to fund account")

			// Set allowance
			s.setupSendAuthz(sender.AccAddr, from.Priv, sdk.Coins{
				{Denom: s.tokenDenom, Amount: sdk.NewIntFromBigInt(authzAmount)},
			})

			// Transfer tokens
			txArgs, transferArgs := s.getTxAndCallArgs(callType, contractAddr)
			transferArgs.MethodName = erc20.TransferFromMethod
			transferArgs.Args = []interface{}{from.Addr, receiver, amount}

			insufficientBalanceCheck := failCheck.WithErrContains(
				"spendable balance 200xmpl is smaller than 300xmpl: insufficient funds",
			)

			_, _, err = s.callContractAndCheckLogs(sender.Priv, txArgs, transferArgs, insufficientBalanceCheck)
			Expect(err).ToNot(HaveOccurred(), "unexpected result calling contract")
		},
			Entry(" - direct call", directCall),
			// NOTE: we are not passing the contract call here because this test case only covers direct calls
		)
	})

	When("approving an allowance", func() {
		DescribeTable("it should approve an allowance", func(callType int) {
			grantee := s.keyring.GetKey(1)
			granter := sender
			amount := big.NewInt(100)

			// Approve allowance
			txArgs, approveArgs := s.getTxAndCallArgs(callType, contractAddr)
			approveArgs.MethodName = auth.ApproveMethod
			approveArgs.Args = []interface{}{grantee.Addr, amount}

			approveCheck := passCheck.WithExpEvents(auth.EventTypeApproval)

			_, _, err = s.callContractAndCheckLogs(granter.Priv, txArgs, approveArgs, approveCheck)
			Expect(err).ToNot(HaveOccurred(), "unexpected result calling contract")

			// Check allowance
			authzs, err := s.grpcHandler.GetAuthorizations(grantee.AccAddr.String(), granter.AccAddr.String())
			Expect(err).ToNot(HaveOccurred(), "failed to get authorizations")
			Expect(authzs).To(HaveLen(1), "expected different number of authorizations")

			sendAuthz, ok := authzs[0].(*banktypes.SendAuthorization)
			Expect(ok).To(BeTrue(), "expected different authorization type")
			Expect(sendAuthz.SpendLimit).To(HaveLen(1), "expected spend limit only in one denomination")
			Expect(sendAuthz.SpendLimit[0].Amount.Int64()).To(Equal(amount.Int64()), "expected different spend limit")
		},
			Entry(" - direct call", directCall),
			// NOTE: we are not passing the contract call here because this test case only covers direct calls
		)

		DescribeTable("it should add a new spend limit to an existing allowance with a different token", func(callType int) {
			grantee := s.keyring.GetKey(1)
			granter := sender
			amount := big.NewInt(100)
			bondCoin := sdk.Coin{Denom: s.network.GetDenom(), Amount: sdk.NewInt(200)}
			tokenCoin := sdk.Coin{Denom: s.tokenDenom, Amount: sdk.NewIntFromBigInt(amount)}

			// Setup a previous authorization
			s.setupSendAuthz(grantee.AccAddr, granter.Priv, sdk.Coins{bondCoin})

			// Approve allowance
			txArgs, approveArgs := s.getTxAndCallArgs(callType, contractAddr)
			approveArgs.MethodName = auth.ApproveMethod
			approveArgs.Args = []interface{}{grantee.Addr, amount}

			approveCheck := passCheck.WithExpEvents(auth.EventTypeApproval)

			_, _, err = s.callContractAndCheckLogs(granter.Priv, txArgs, approveArgs, approveCheck)
			Expect(err).ToNot(HaveOccurred(), "unexpected result calling contract")

			// Check allowance contains both spend limits
			authzs, err := s.grpcHandler.GetAuthorizations(grantee.AccAddr.String(), granter.AccAddr.String())
			Expect(err).ToNot(HaveOccurred(), "failed to get authorizations")
			Expect(authzs).To(HaveLen(1), "expected different number of authorizations")

			sendAuthz, ok := authzs[0].(*banktypes.SendAuthorization)
			Expect(ok).To(BeTrue(), "expected different authorization type")
			Expect(sendAuthz.SpendLimit).To(HaveLen(2), "expected spend limit in two denominations")
			Expect(sendAuthz.SpendLimit).To(Equal(sdk.NewCoins(bondCoin, tokenCoin)), "expected different spend limit")
		},
			Entry(" - direct call", directCall),
			// NOTE: we are not passing the contract call here because this test case only covers direct calls
		)

		DescribeTable("it should set the new spend limit for an existing allowance with the same token", func(callType int) {
			grantee := s.keyring.GetKey(1)
			granter := sender
			amount := big.NewInt(100)
			bondCoin := sdk.Coin{Denom: s.network.GetDenom(), Amount: sdk.NewInt(200)}
			tokenCoin := sdk.Coin{Denom: s.tokenDenom, Amount: sdk.NewIntFromBigInt(amount)}
			doubleTokenCoin := sdk.Coin{Denom: s.tokenDenom, Amount: sdk.NewInt(200)}

			// Setup a previous authorization
			s.setupSendAuthz(grantee.AccAddr, granter.Priv, sdk.NewCoins(bondCoin, doubleTokenCoin))

			// Approve allowance
			txArgs, approveArgs := s.getTxAndCallArgs(callType, contractAddr)
			approveArgs.MethodName = auth.ApproveMethod
			approveArgs.Args = []interface{}{grantee.Addr, amount}

			approveCheck := passCheck.WithExpEvents(auth.EventTypeApproval)

			_, _, err = s.callContractAndCheckLogs(granter.Priv, txArgs, approveArgs, approveCheck)
			Expect(err).ToNot(HaveOccurred(), "unexpected result calling contract")

			// Check allowance contains both spend limits
			authzs, err := s.grpcHandler.GetAuthorizations(grantee.AccAddr.String(), granter.AccAddr.String())
			Expect(err).ToNot(HaveOccurred(), "failed to get authorizations")
			Expect(authzs).To(HaveLen(1), "expected different number of authorizations")

			sendAuthz, ok := authzs[0].(*banktypes.SendAuthorization)
			Expect(ok).To(BeTrue(), "expected different authorization type")
			Expect(sendAuthz.SpendLimit).To(HaveLen(2), "expected spend limit in two denominations")
			Expect(sendAuthz.SpendLimit).To(Equal(sdk.NewCoins(bondCoin, tokenCoin)), "expected different spend limit")
		},
			Entry(" - direct call", directCall),
			// NOTE: we are not passing the contract call here because this test case only covers direct calls
		)

		DescribeTable("it should remove the token from the spend limit of an existing authorization when approving zero", func(callType int) {
			grantee := s.keyring.GetKey(1)
			granter := sender
			amount := big.NewInt(100)
			bondCoin := sdk.Coin{Denom: s.network.GetDenom(), Amount: sdk.NewInt(200)}
			tokenCoin := sdk.Coin{Denom: s.tokenDenom, Amount: sdk.NewIntFromBigInt(amount)}

			// Setup a previous authorization
			s.setupSendAuthz(grantee.AccAddr, granter.Priv, sdk.NewCoins(bondCoin, tokenCoin))

			// Approve allowance
			txArgs, approveArgs := s.getTxAndCallArgs(callType, contractAddr)
			approveArgs.MethodName = auth.ApproveMethod
			approveArgs.Args = []interface{}{grantee.Addr, common.Big0}

			approveCheck := passCheck.WithExpEvents(auth.EventTypeApproval)

			_, _, err = s.callContractAndCheckLogs(granter.Priv, txArgs, approveArgs, approveCheck)
			Expect(err).ToNot(HaveOccurred(), "unexpected result calling contract")

			// Check allowance contains both spend limits
			authzs, err := s.grpcHandler.GetAuthorizations(grantee.AccAddr.String(), granter.AccAddr.String())
			Expect(err).ToNot(HaveOccurred(), "failed to get authorizations")
			Expect(authzs).To(HaveLen(1), "expected different number of authorizations")

			sendAuthz, ok := authzs[0].(*banktypes.SendAuthorization)
			Expect(ok).To(BeTrue(), "expected different authorization type")
			Expect(sendAuthz.SpendLimit).To(HaveLen(1), "expected spend limit in one denomination")
			Expect(sendAuthz.SpendLimit).To(Equal(sdk.NewCoins(bondCoin)), "expected different spend limit")
		},
			Entry(" - direct call", directCall),
			// NOTE: we are not passing the contract call here because this test case only covers direct calls
		)

		DescribeTable("it should delete the authorization when approving zero with no other spend limits", func(callType int) {
			grantee := s.keyring.GetKey(1)
			granter := sender
			amount := big.NewInt(100)
			tokenCoin := sdk.Coin{Denom: s.tokenDenom, Amount: sdk.NewIntFromBigInt(amount)}

			// Setup a previous authorization
			s.setupSendAuthz(grantee.AccAddr, granter.Priv, sdk.NewCoins(tokenCoin))

			// Approve allowance
			txArgs, approveArgs := s.getTxAndCallArgs(callType, contractAddr)
			approveArgs.MethodName = auth.ApproveMethod
			approveArgs.Args = []interface{}{grantee.Addr, common.Big0}

			approveCheck := passCheck.WithExpEvents(auth.EventTypeApproval)

			_, _, err = s.callContractAndCheckLogs(granter.Priv, txArgs, approveArgs, approveCheck)
			Expect(err).ToNot(HaveOccurred(), "unexpected result calling contract")

			// Check allowance contains both spend limits
			authzs, err := s.grpcHandler.GetAuthorizations(grantee.AccAddr.String(), granter.AccAddr.String())
			Expect(err).ToNot(HaveOccurred(), "failed to get authorizations")
			Expect(authzs).To(HaveLen(0), "expected different number of authorizations")
		},
			Entry(" - direct call", directCall),
			// NOTE: we are not passing the contract call here because this test case only covers direct calls
		)

		DescribeTable("it should return an error if approving 0 and no allowance exists", func(callType int) {
			grantee := s.keyring.GetKey(1)
			granter := sender

			// Approve allowance
			txArgs, approveArgs := s.getTxAndCallArgs(callType, contractAddr)
			approveArgs.MethodName = auth.ApproveMethod
			approveArgs.Args = []interface{}{grantee.Addr, common.Big0}

			nonPosCheck := failCheck.WithErrContains("cannot approve non-positive values")

			_, _, err = s.callContractAndCheckLogs(granter.Priv, txArgs, approveArgs, nonPosCheck)
			Expect(err).ToNot(HaveOccurred(), "unexpected result calling contract")

			// Check allowance contains both spend limits
			authzs, err := s.grpcHandler.GetAuthorizations(grantee.AccAddr.String(), granter.AccAddr.String())
			Expect(err).ToNot(HaveOccurred(), "failed to get authorizations")
			Expect(authzs).To(HaveLen(0), "expected different number of authorizations")
		},
			Entry(" - direct call", directCall),
			// NOTE: we are not passing the contract call here because this test case only covers direct calls
		)

		// TODO: this is passing right now?? Should we allow someone to create an authorization for themselves?
		DescribeTable("it should return an error if the grantee is the same as the granter", func(callType int) {
			grantee := sender
			granter := sender
			amount := big.NewInt(100)

			// Approve allowance
			txArgs, approveArgs := s.getTxAndCallArgs(callType, contractAddr)
			approveArgs.MethodName = auth.ApproveMethod
			approveArgs.Args = []interface{}{grantee.Addr, amount}

			_, _, err = s.callContractAndCheckLogs(granter.Priv, txArgs, approveArgs, execRevertedCheck)
			Expect(err).ToNot(HaveOccurred(), "unexpected result calling contract")
		},
			Entry(" - direct call", directCall),
			// NOTE: we are not passing the contract call here because this test case only covers direct calls
		)

		DescribeTable("it should return an error if approving 0 and allowance only exists for other tokens", func(callType int) {
			grantee := s.keyring.GetKey(1)
			granter := sender

			// Setup a previous authorization
			s.setupSendAuthz(grantee.AccAddr, granter.Priv, sdk.NewCoins(sdk.Coin{Denom: s.network.GetDenom(), Amount: sdk.NewInt(200)}))

			// Approve allowance
			txArgs, approveArgs := s.getTxAndCallArgs(callType, contractAddr)
			approveArgs.MethodName = auth.ApproveMethod
			approveArgs.Args = []interface{}{grantee.Addr, common.Big0}

			notFoundCheck := failCheck.WithErrContains(
				fmt.Sprintf("allowance for token %s does not exist", s.tokenDenom),
			)

			_, _, err = s.callContractAndCheckLogs(granter.Priv, txArgs, approveArgs, notFoundCheck)
			Expect(err).ToNot(HaveOccurred(), "unexpected result calling contract")
		},
			Entry(" - direct call", directCall),
			// NOTE: we are not passing the contract call here because this test case only covers direct calls
		)
	})
})