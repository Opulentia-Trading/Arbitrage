package ethHandler

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/Opulentia-Trading/Arbitrage/contracts/erc20"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const txMineWaitTimeout = 5 * time.Minute

type ERC20Handler struct {
	*EthHandler
	Token    *Token
	Contract *erc20.Erc20
}

func NewERC20Handler(ethHandler *EthHandler, token *Token) (*ERC20Handler, error) {
	if token.Type != ERC20 {
		return nil, fmt.Errorf("invalid token type: %v", token.Type)
	}

	contract, err := erc20.NewErc20(token.AddressForGeth(), ethHandler.Client)
	if err != nil {
		panic(err)
	}

	return &ERC20Handler{
		EthHandler: ethHandler,
		Token:      token,
		Contract:   contract,
	}, nil
}

func (e *ERC20Handler) TotalSupply() (*big.Int, error) {
	return e.Contract.TotalSupply(&bind.CallOpts{})
}

func (e *ERC20Handler) BalanceOf(account common.Address) (*big.Int, error) {
	return e.Contract.BalanceOf(&bind.CallOpts{}, account)
}

func (e *ERC20Handler) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return e.Contract.Allowance(&bind.CallOpts{}, owner, spender)
}

func (e *ERC20Handler) validateApproveTx(wallet *Wallet, tx *types.Transaction, spender common.Address, amount *big.Int) error {
	txReceipt, err := e.WaitTxMined(tx, wallet.Address, txMineWaitTimeout)
	if err != nil {
		panic(err)
	}

	tokenAddress := e.Token.AddressForGeth()

	// keccak256 hash of event signature "Approval(address,address,uint256)"
	approvalSigHashHex := "0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925"
	approvalSigHash := common.HexToHash(approvalSigHashHex)

	// Use the txReceipt bloom filter to check if the approval log is not present
	tokenAddressInLogs := types.BloomLookup(txReceipt.Bloom, tokenAddress)
	approvalSigHashInLogs := types.BloomLookup(txReceipt.Bloom, approvalSigHash)
	if !tokenAddressInLogs || !approvalSigHashInLogs {
		panic("cannot find approval log")
	}

	// Search logs for the approval event
	for _, log := range txReceipt.Logs {
		if log.Address != tokenAddress {
			continue
		}

		if log.Topics[0] != approvalSigHash {
			continue
		}

		approvalInfo, err := e.Contract.ParseApproval(*log)
		if err != nil {
			panic(err)
		}

		valueMatch := approvalInfo.Value.Cmp(amount) == 0
		if approvalInfo.Owner == wallet.Address && approvalInfo.Spender == spender && valueMatch {
			fmt.Println("\n[Approval event log]")
			fmt.Printf("owner: %v\n", approvalInfo.Owner)
			fmt.Printf("spender: %v\n", approvalInfo.Spender)
			fmt.Printf("amount: 0x%v\n", approvalInfo.Value.Text(16))
			return nil
		}
	}

	panic("cannot find approval log")
}

// Changing the allowance directly may allow an attacker to use both the old and new allowance.
// To mitigate this, we have to set the allowance to 0 and then set the desired amount afterwards.
// https://github.com/ethereum/EIPs/issues/20#issuecomment-263524729
func (e *ERC20Handler) unsafeApprove(wallet *Wallet, spender common.Address, amount *big.Int) error {
	nonce, err := e.Client.PendingNonceAt(context.Background(), wallet.Address)
	if err != nil {
		panic(err)
	}

	chainId := big.NewInt(int64(e.Network.ChainId))
	auth, err := bind.NewKeyedTransactorWithChainID(wallet.PrivateKey, chainId)
	if err != nil {
		panic(err)
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = nil
	auth.NoSend = false

	// TODO: Get gas estimates from the gasEstimator module
	auth.GasPrice = nil
	auth.GasFeeCap = nil
	auth.GasTipCap = nil
	auth.GasLimit = uint64(0) // in units

	tx, err := e.Contract.Approve(auth, spender, amount)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\n[[ %v approve tx ]]\n", e.Token.Symbol)
	fmt.Printf("tx hash: %s\n", tx.Hash())
	fmt.Printf("gas priority fee: %v\n", tx.GasTipCap())
	fmt.Printf("gas max fee: %v\n", tx.GasFeeCap())
	fmt.Printf("gas limit: %v\n", tx.Gas())
	if auth.NoSend {
		fmt.Println("Note: transaction not sent on blockchain")
		return nil
	}

	err = e.validateApproveTx(wallet, tx, spender, amount)
	if err != nil {
		panic(err)
	}

	return nil
}

func (e *ERC20Handler) Approve(wallet *Wallet, spender common.Address, amount *big.Int, approveOnlyZero bool) error {
	curAllowance, err := e.Allowance(wallet.Address, spender)
	if err != nil {
		panic(err)
	}

	if curAllowance.Cmp(common.Big0) > 0 && approveOnlyZero {
		return nil
	}

	if curAllowance.Cmp(common.Big0) > 0 {
		err := e.unsafeApprove(wallet, spender, common.Big0)
		if err != nil {
			panic(err)
		}
	}

	if amount.Cmp(common.Big0) > 0 {
		err := e.unsafeApprove(wallet, spender, amount)
		if err != nil {
			panic(err)
		}
	}

	return nil
}

func (e *ERC20Handler) MaxApprove(wallet *Wallet, spender common.Address, approveOnlyZero bool) error {
	// maxAmount = (2^256) - 1
	maxAmount := new(big.Int).Lsh(common.Big1, 256)
	maxAmount.Sub(maxAmount, common.Big1)
	return e.Approve(wallet, spender, maxAmount, approveOnlyZero)
}
