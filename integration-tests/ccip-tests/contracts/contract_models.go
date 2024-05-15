package contracts

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/rand"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"

	"github.com/smartcontractkit/ccip/integration-tests/wrappers"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/arm_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/commit_store"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/commit_store_1_2_0"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_offramp_1_2_0"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_onramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_onramp_1_2_0"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/lock_release_token_pool"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/lock_release_token_pool_1_4_0"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/maybe_revert_message_receiver"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/mock_arm_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/mock_usdc_token_transmitter"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/mock_v3_aggregator_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/price_registry"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/price_registry_1_2_0"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/token_admin_registry"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/token_pool"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/token_pool_1_4_0"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/usdc_token_pool"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/usdc_token_pool_1_4_0"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/burn_mint_erc677"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/erc20"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
)

var (
	FiftyCoins   = new(big.Int).Mul(big.NewInt(1e18), big.NewInt(50))
	HundredCoins = new(big.Int).Mul(big.NewInt(1e18), big.NewInt(100))
)

type ContractVersion string

const (
	Network                               = "Network Name"
	V1_2_0                ContractVersion = "1.2.0"
	V1_4_0                ContractVersion = "1.4.0"
	LatestPoolVersion     ContractVersion = "1.5.0-dev"
	Latest                ContractVersion = "latest"
	PriceRegistryContract                 = "PriceRegistry"
	OffRampContract                       = "OffRamp"
	OnRampContract                        = "OnRamp"
	TokenPoolContract                     = "TokenPool"
	CommitStoreContract                   = "CommitStore"
)

var (
	VersionMap = map[string]ContractVersion{
		PriceRegistryContract: Latest,
		OffRampContract:       Latest,
		OnRampContract:        Latest,
		CommitStoreContract:   Latest,
		TokenPoolContract:     Latest,
	}
	SupportedContracts = map[string]map[ContractVersion]bool{
		PriceRegistryContract: {
			Latest: true,
			V1_2_0: true,
		},
		OffRampContract: {
			Latest: true,
			V1_2_0: true,
		},
		OnRampContract: {
			Latest: true,
			V1_2_0: true,
		},
		CommitStoreContract: {
			Latest: true,
			V1_2_0: true,
		},
		TokenPoolContract: {
			Latest: true,
			V1_4_0: true,
		},
	}
)

type RateLimiterConfig struct {
	IsEnabled bool
	Rate      *big.Int
	Capacity  *big.Int
	Tokens    *big.Int
}

type ARMConfig struct {
	ARMWeightsByParticipants map[string]*big.Int // mapping : ARM participant address => weight
	ThresholdForBlessing     *big.Int
	ThresholdForBadSignal    *big.Int
}

type TokenTransmitter struct {
	client          blockchain.EVMClient
	instance        *mock_usdc_token_transmitter.MockE2EUSDCTransmitter
	ContractAddress common.Address
}

type ERC677Token struct {
	client          blockchain.EVMClient
	instance        *burn_mint_erc677.BurnMintERC677
	ContractAddress common.Address
}

func (token *ERC677Token) GrantMintAndBurn(burnAndMinter common.Address) error {
	opts, err := token.client.TransactionOpts(token.client.GetDefaultWallet())
	if err != nil {
		return fmt.Errorf("failed to get transaction opts: %w", err)
	}
	log.Info().
		Str(Network, token.client.GetNetworkName()).
		Str("BurnAndMinter", burnAndMinter.Hex()).
		Str("Token", token.ContractAddress.Hex()).
		Msg("Granting mint and burn roles")
	tx, err := token.instance.GrantMintAndBurnRoles(opts, burnAndMinter)
	if err != nil {
		return fmt.Errorf("failed to grant mint and burn roles: %w", err)
	}
	return token.client.ProcessTransaction(tx)
}

func (token *ERC677Token) GrantMintRole(minter common.Address) error {
	opts, err := token.client.TransactionOpts(token.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	log.Info().
		Str(Network, token.client.GetNetworkName()).
		Str("Minter", minter.Hex()).
		Str("Token", token.ContractAddress.Hex()).
		Msg("Granting mint roles")
	tx, err := token.instance.GrantMintRole(opts, minter)
	if err != nil {
		return fmt.Errorf("failed to grant mint role: %w", err)
	}
	return token.client.ProcessTransaction(tx)
}

func (token *ERC677Token) Mint(to common.Address, amount *big.Int) error {
	opts, err := token.client.TransactionOpts(token.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	log.Info().
		Str(Network, token.client.GetNetworkName()).
		Str("To", to.Hex()).
		Str("Token", token.ContractAddress.Hex()).
		Str("Amount", amount.String()).
		Msg("Minting tokens")
	tx, err := token.instance.Mint(opts, to, amount)
	if err != nil {
		return fmt.Errorf("failed to mint tokens: %w", err)
	}
	return token.client.ProcessTransaction(tx)
}

type ERC20Token struct {
	client          blockchain.EVMClient
	instance        *erc20.ERC20
	ContractAddress common.Address
}

func (token *ERC20Token) Address() string {
	return token.ContractAddress.Hex()
}

func (token *ERC20Token) BalanceOf(ctx context.Context, addr string) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(token.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	balance, err := token.instance.BalanceOf(opts, common.HexToAddress(addr))
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}
	return balance, nil
}

func (token *ERC20Token) Allowance(owner, spender string) (*big.Int, error) {
	allowance, err := token.instance.Allowance(nil, common.HexToAddress(owner), common.HexToAddress(spender))
	if err != nil {
		return nil, err
	}
	return allowance, nil
}

func (token *ERC20Token) Approve(to string, amount *big.Int) error {
	opts, err := token.client.TransactionOpts(token.client.GetDefaultWallet())
	if err != nil {
		return fmt.Errorf("failed to get transaction options: %w", err)
	}
	log.Info().
		Str("From", token.client.GetDefaultWallet().Address()).
		Str("To", to).
		Str("Token", token.Address()).
		Str("Amount", amount.String()).
		Uint64("Nonce", opts.Nonce.Uint64()).
		Str(Network, token.client.GetNetworkConfig().Name).
		Msg("Approving ERC20 Transfer")
	tx, err := token.instance.Approve(opts, common.HexToAddress(to), amount)
	if err != nil {
		return fmt.Errorf("failed to approve ERC20: %w", err)
	}
	return token.client.ProcessTransaction(tx)
}

func (token *ERC20Token) Transfer(to string, amount *big.Int) error {
	opts, err := token.client.TransactionOpts(token.client.GetDefaultWallet())
	if err != nil {
		return fmt.Errorf("failed to get transaction options: %w", err)
	}
	log.Info().
		Str("From", token.client.GetDefaultWallet().Address()).
		Str("To", to).
		Str("Amount", amount.String()).
		Uint64("Nonce", opts.Nonce.Uint64()).
		Str(Network, token.client.GetNetworkConfig().Name).
		Msg("Transferring ERC20")
	tx, err := token.instance.Transfer(opts, common.HexToAddress(to), amount)
	if err != nil {
		return fmt.Errorf("failed to transfer ERC20: %w", err)
	}
	return token.client.ProcessTransaction(tx)
}

type LinkToken struct {
	client     blockchain.EVMClient
	instance   *link_token_interface.LinkToken
	EthAddress common.Address
}

func (l *LinkToken) Address() string {
	return l.EthAddress.Hex()
}

func (l *LinkToken) BalanceOf(ctx context.Context, addr string) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(l.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	balance, err := l.instance.BalanceOf(opts, common.HexToAddress(addr))
	if err != nil {
		return nil, fmt.Errorf("failed to get LINK balance: %w", err)
	}
	return balance, nil
}

func (l *LinkToken) Allowance(owner, spender string) (*big.Int, error) {
	allowance, err := l.instance.Allowance(nil, common.HexToAddress(owner), common.HexToAddress(spender))
	if err != nil {
		return nil, err
	}
	return allowance, nil
}

func (l *LinkToken) Approve(to string, amount *big.Int) error {
	opts, err := l.client.TransactionOpts(l.client.GetDefaultWallet())
	if err != nil {
		return fmt.Errorf("failed to get transaction opts: %w", err)
	}
	log.Info().
		Str("From", l.client.GetDefaultWallet().Address()).
		Str("To", to).
		Str("Token", l.Address()).
		Str("Amount", amount.String()).
		Uint64("Nonce", opts.Nonce.Uint64()).
		Str(Network, l.client.GetNetworkConfig().Name).
		Msg("Approving LINK Transfer")
	tx, err := l.instance.Approve(opts, common.HexToAddress(to), amount)
	if err != nil {
		return fmt.Errorf("failed to approve LINK transfer: %w", err)
	}
	return l.client.ProcessTransaction(tx)
}

func (l *LinkToken) Transfer(to string, amount *big.Int) error {
	opts, err := l.client.TransactionOpts(l.client.GetDefaultWallet())
	if err != nil {
		return fmt.Errorf("failed to get transaction opts: %w", err)
	}
	log.Info().
		Str("From", l.client.GetDefaultWallet().Address()).
		Str("To", to).
		Str("Amount", amount.String()).
		Uint64("Nonce", opts.Nonce.Uint64()).
		Str(Network, l.client.GetNetworkConfig().Name).
		Msg("Transferring LINK")
	tx, err := l.instance.Transfer(opts, common.HexToAddress(to), amount)
	if err != nil {
		return fmt.Errorf("failed to transfer LINK: %w", err)
	}
	return l.client.ProcessTransaction(tx)
}

type LatestPool struct {
	PoolInterface   *token_pool.TokenPool
	LockReleasePool *lock_release_token_pool.LockReleaseTokenPool
	USDCPool        *usdc_token_pool.USDCTokenPool
}

type V1_4_0Pool struct {
	PoolInterface   *token_pool_1_4_0.TokenPool
	LockReleasePool *lock_release_token_pool_1_4_0.LockReleaseTokenPool
	USDCPool        *usdc_token_pool_1_4_0.USDCTokenPool
}

type TokenPoolWrapper struct {
	Latest *LatestPool
	V1_4_0 *V1_4_0Pool
}

func (w TokenPoolWrapper) SetRebalancer(opts *bind.TransactOpts, from common.Address) (*types.Transaction, error) {
	if w.Latest != nil && w.Latest.LockReleasePool != nil {
		return w.Latest.LockReleasePool.SetRebalancer(opts, from)
	}
	if w.V1_4_0 != nil && w.V1_4_0.LockReleasePool != nil {
		return w.V1_4_0.LockReleasePool.SetRebalancer(opts, from)
	}
	return nil, fmt.Errorf("no pool found to set rebalancer")
}

func (w TokenPoolWrapper) SetUSDCDomains(opts *bind.TransactOpts, updates []usdc_token_pool.USDCTokenPoolDomainUpdate) (*types.Transaction, error) {
	if w.Latest != nil && w.Latest.USDCPool != nil {
		return w.Latest.USDCPool.SetDomains(opts, updates)
	}
	if w.V1_4_0 != nil && w.V1_4_0.USDCPool != nil {
		V1_4_0Updates := make([]usdc_token_pool_1_4_0.USDCTokenPoolDomainUpdate, len(updates))
		for i, update := range updates {
			V1_4_0Updates[i] = usdc_token_pool_1_4_0.USDCTokenPoolDomainUpdate{
				AllowedCaller:     update.AllowedCaller,
				DomainIdentifier:  update.DomainIdentifier,
				DestChainSelector: update.DestChainSelector,
				Enabled:           update.Enabled,
			}
		}
		return w.V1_4_0.USDCPool.SetDomains(opts, V1_4_0Updates)
	}
	return nil, fmt.Errorf("no pool found to set USDC domains")
}

func (w TokenPoolWrapper) WithdrawLiquidity(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	if w.Latest != nil && w.Latest.LockReleasePool != nil {
		return w.Latest.LockReleasePool.WithdrawLiquidity(opts, amount)
	}
	if w.V1_4_0 != nil && w.V1_4_0.LockReleasePool != nil {
		return w.V1_4_0.LockReleasePool.WithdrawLiquidity(opts, amount)
	}
	return nil, fmt.Errorf("no pool found to withdraw liquidity")
}

func (w TokenPoolWrapper) ProvideLiquidity(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	if w.Latest != nil && w.Latest.LockReleasePool != nil {
		return w.Latest.LockReleasePool.ProvideLiquidity(opts, amount)
	}
	if w.V1_4_0 != nil && w.V1_4_0.LockReleasePool != nil {
		return w.V1_4_0.LockReleasePool.ProvideLiquidity(opts, amount)
	}
	return nil, fmt.Errorf("no pool found to provide liquidity")
}

func (w TokenPoolWrapper) SetRemotePool(opts *bind.TransactOpts, selector uint64, pool []byte) (*types.Transaction, error) {
	if w.Latest != nil && w.Latest.PoolInterface != nil {
		return w.Latest.PoolInterface.SetRemotePool(opts, selector, pool)
	}
	return nil, fmt.Errorf("no pool found to set remote pool")
}

func (w TokenPoolWrapper) IsSupportedChain(opts *bind.CallOpts, remoteChainSelector uint64) (bool, error) {
	if w.Latest != nil && w.Latest.PoolInterface != nil {
		return w.Latest.PoolInterface.IsSupportedChain(opts, remoteChainSelector)
	}
	if w.V1_4_0 != nil && w.V1_4_0.PoolInterface != nil {
		return w.V1_4_0.PoolInterface.IsSupportedChain(opts, remoteChainSelector)
	}
	return false, fmt.Errorf("no pool found to check if chain is supported")
}

func (w TokenPoolWrapper) ApplyChainUpdates(opts *bind.TransactOpts, update []token_pool.TokenPoolChainUpdate) (*types.Transaction, error) {
	if w.Latest != nil && w.Latest.PoolInterface != nil {
		return w.Latest.PoolInterface.ApplyChainUpdates(opts, update)
	}
	if w.V1_4_0 != nil && w.V1_4_0.PoolInterface != nil {
		V1_4_0Updates := make([]token_pool_1_4_0.TokenPoolChainUpdate, len(update))
		for i, u := range update {
			V1_4_0Updates[i] = token_pool_1_4_0.TokenPoolChainUpdate{
				RemoteChainSelector: u.RemoteChainSelector,
				Allowed:             u.Allowed,
				InboundRateLimiterConfig: token_pool_1_4_0.RateLimiterConfig{
					IsEnabled: u.InboundRateLimiterConfig.IsEnabled,
					Capacity:  u.InboundRateLimiterConfig.Capacity,
					Rate:      u.InboundRateLimiterConfig.Rate,
				},
				OutboundRateLimiterConfig: token_pool_1_4_0.RateLimiterConfig{
					IsEnabled: u.OutboundRateLimiterConfig.IsEnabled,
					Capacity:  u.OutboundRateLimiterConfig.Capacity,
					Rate:      u.OutboundRateLimiterConfig.Rate,
				},
			}
		}
		return w.V1_4_0.PoolInterface.ApplyChainUpdates(opts, V1_4_0Updates)
	}
	return nil, fmt.Errorf("no pool found to apply chain updates")
}

func (w TokenPoolWrapper) SetChainRateLimiterConfig(opts *bind.TransactOpts, selector uint64, out token_pool.RateLimiterConfig, in token_pool.RateLimiterConfig) (*types.Transaction, error) {
	if w.Latest != nil && w.Latest.PoolInterface != nil {
		return w.Latest.PoolInterface.SetChainRateLimiterConfig(opts, selector, out, in)
	}
	if w.V1_4_0 != nil && w.V1_4_0.PoolInterface != nil {
		return w.V1_4_0.PoolInterface.SetChainRateLimiterConfig(opts, selector,
			token_pool_1_4_0.RateLimiterConfig{
				IsEnabled: out.IsEnabled,
				Capacity:  out.Capacity,
				Rate:      out.Rate,
			}, token_pool_1_4_0.RateLimiterConfig{
				IsEnabled: in.IsEnabled,
				Capacity:  in.Capacity,
				Rate:      in.Rate,
			})
	}
	return nil, fmt.Errorf("no pool found to set chain rate limiter config")
}

func (w TokenPoolWrapper) GetCurrentOutboundRateLimiterState(opts *bind.CallOpts, selector uint64) (*RateLimiterConfig, error) {
	if w.Latest != nil && w.Latest.PoolInterface != nil {
		rl, err := w.Latest.PoolInterface.GetCurrentOutboundRateLimiterState(opts, selector)
		if err != nil {
			return nil, err
		}
		return &RateLimiterConfig{
			IsEnabled: rl.IsEnabled,
			Capacity:  rl.Capacity,
			Rate:      rl.Rate,
			Tokens:    rl.Tokens,
		}, nil
	}
	if w.V1_4_0 != nil && w.V1_4_0.PoolInterface != nil {
		rl, err := w.V1_4_0.PoolInterface.GetCurrentOutboundRateLimiterState(opts, selector)
		if err != nil {
			return nil, err
		}
		return &RateLimiterConfig{
			IsEnabled: rl.IsEnabled,
			Capacity:  rl.Capacity,
			Rate:      rl.Rate,
			Tokens:    rl.Tokens,
		}, nil
	}
	return nil, fmt.Errorf("no pool found to get current outbound rate limiter state")
}

func (w TokenPoolWrapper) GetCurrentInboundRateLimiterState(opts *bind.CallOpts, selector uint64) (*RateLimiterConfig, error) {
	if w.Latest != nil && w.Latest.PoolInterface != nil {
		rl, err := w.Latest.PoolInterface.GetCurrentInboundRateLimiterState(opts, selector)
		if err != nil {
			return nil, err
		}
		return &RateLimiterConfig{
			IsEnabled: rl.IsEnabled,
			Capacity:  rl.Capacity,
			Rate:      rl.Rate,
			Tokens:    rl.Tokens,
		}, nil
	}
	if w.V1_4_0 != nil && w.V1_4_0.PoolInterface != nil {
		rl, err := w.V1_4_0.PoolInterface.GetCurrentInboundRateLimiterState(opts, selector)
		if err != nil {
			return nil, err
		}
		return &RateLimiterConfig{
			IsEnabled: rl.IsEnabled,
			Capacity:  rl.Capacity,
			Rate:      rl.Rate,
			Tokens:    rl.Tokens,
		}, nil
	}
	return nil, fmt.Errorf("no pool found to get current outbound rate limiter state")
}

// TokenPool represents a TokenPool address
type TokenPool struct {
	client     blockchain.EVMClient
	Instance   *TokenPoolWrapper
	EthAddress common.Address
}

func (pool *TokenPool) Address() string {
	return pool.EthAddress.Hex()
}

func (pool *TokenPool) IsUSDC() bool {
	if pool.Instance.Latest != nil && pool.Instance.Latest.USDCPool != nil {
		return true
	}
	if pool.Instance.V1_4_0 != nil && pool.Instance.V1_4_0.USDCPool != nil {
		return true
	}
	return false
}

func (pool *TokenPool) IsLockRelease() bool {
	if pool.Instance.Latest != nil && pool.Instance.Latest.LockReleasePool != nil {
		return true
	}
	if pool.Instance.V1_4_0 != nil && pool.Instance.V1_4_0.LockReleasePool != nil {
		return true
	}
	return false
}

func (pool *TokenPool) SyncUSDCDomain(destTokenTransmitter *TokenTransmitter, destPoolAddr common.Address, destChainSelector uint64) error {
	if !pool.IsUSDC() {
		return fmt.Errorf("pool is not a USDC pool, cannot sync domain")
	}

	var allowedCallerBytes [32]byte
	copy(allowedCallerBytes[12:], destPoolAddr.Bytes())
	destTokenTransmitterIns, err := mock_usdc_token_transmitter.NewMockE2EUSDCTransmitter(destTokenTransmitter.ContractAddress, destTokenTransmitter.client.Backend())
	if err != nil {
		return fmt.Errorf("failed to create mock USDC token transmitter: %w", err)
	}
	domain, err := destTokenTransmitterIns.LocalDomain(nil)
	if err != nil {
		return fmt.Errorf("failed to get local domain: %w", err)
	}
	opts, err := pool.client.TransactionOpts(pool.client.GetDefaultWallet())
	if err != nil {
		return fmt.Errorf("failed to get transaction opts: %w", err)
	}
	log.Info().
		Str("Token Pool", pool.Address()).
		Str(Network, pool.client.GetNetworkName()).
		Uint32("Domain", domain).
		Str("Allowed Caller", destPoolAddr.Hex()).
		Str("Dest Chain Selector", fmt.Sprintf("%d", destChainSelector)).
		Msg("Syncing USDC Domain")
	tx, err := pool.Instance.SetUSDCDomains(opts, []usdc_token_pool.USDCTokenPoolDomainUpdate{
		{
			AllowedCaller:     allowedCallerBytes,
			DomainIdentifier:  domain,
			DestChainSelector: destChainSelector,
			Enabled:           true,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to set domain: %w", err)
	}
	return pool.client.ProcessTransaction(tx)
}

func (pool *TokenPool) RemoveLiquidity(amount *big.Int) error {
	if !pool.IsLockRelease() {
		return fmt.Errorf("pool is not a lock release pool, cannot remove liquidity")
	}
	opts, err := pool.client.TransactionOpts(pool.client.GetDefaultWallet())
	if err != nil {
		return fmt.Errorf("failed to get transaction opts: %w", err)
	}
	log.Info().
		Str("Token Pool", pool.Address()).
		Str("Amount", amount.String()).
		Msg("Initiating removing funds from pool")
	tx, err := pool.Instance.WithdrawLiquidity(opts, amount)
	if err != nil {
		return fmt.Errorf("failed to withdraw liquidity: %w", err)
	}
	log.Info().
		Str("Token Pool", pool.Address()).
		Str("Amount", amount.String()).
		Str(Network, pool.client.GetNetworkConfig().Name).
		Msg("Liquidity removed")
	return pool.client.ProcessTransaction(tx)
}

type tokenApproveFn func(string, *big.Int) error

func (pool *TokenPool) AddLiquidity(approveFn tokenApproveFn, tokenAddr string, amount *big.Int) error {
	if !pool.IsLockRelease() {
		return fmt.Errorf("pool is not a lock release pool, cannot add liquidity")
	}
	log.Info().
		Str("Link Token", tokenAddr).
		Str("Token Pool", pool.Address()).
		Msg("Initiating transferring of token to token pool")
	err := approveFn(pool.Address(), amount)
	if err != nil {
		return fmt.Errorf("failed to approve token transfer: %w", err)
	}
	err = pool.client.WaitForEvents()
	if err != nil {
		return fmt.Errorf("failed to wait for events: %w", err)
	}
	opts, err := pool.client.TransactionOpts(pool.client.GetDefaultWallet())
	if err != nil {
		return fmt.Errorf("failed to get transaction opts: %w", err)
	}
	_, err = pool.Instance.SetRebalancer(opts, opts.From)
	if err != nil {
		return fmt.Errorf("failed to set rebalancer: %w", err)
	}
	opts, err = pool.client.TransactionOpts(pool.client.GetDefaultWallet())
	if err != nil {
		return fmt.Errorf("failed to get transaction opts: %w", err)
	}
	log.Info().
		Str("Token Pool", pool.Address()).
		Msg("Initiating adding Tokens in pool")
	tx, err := pool.Instance.ProvideLiquidity(opts, amount)
	if err != nil {
		return fmt.Errorf("failed to provide liquidity: %w", err)
	}
	log.Info().
		Str("Token Pool", pool.Address()).
		Str("Link Token", tokenAddr).
		Str(Network, pool.client.GetNetworkConfig().Name).
		Msg("Liquidity added")
	return pool.client.ProcessTransaction(tx)
}

func (pool *TokenPool) SetRemotePool(remoteChainSelector uint64, remotePool common.Address) error {
	// if pool is of version 1.4.0, no need to set remote pool
	if pool.Instance.V1_4_0 != nil {
		return nil
	}
	log.Info().
		Str("Token Pool", pool.Address()).
		Str("Dest Pool", remotePool.Hex()).
		Uint64("RemoteChain", remoteChainSelector).
		Msg("Setting remote pool")
	abiEncodedRemotePool, err := abihelpers.EncodeAddress(remotePool)
	if err != nil {
		return fmt.Errorf("error getting abiEncodedRemotePool %s : %w", remotePool.Hex(), err)
	}
	opts, err := pool.client.TransactionOpts(pool.client.GetDefaultWallet())
	if err != nil {
		return fmt.Errorf("failed to get transaction opts: %w", err)
	}
	tx, err := pool.Instance.SetRemotePool(opts, remoteChainSelector, abiEncodedRemotePool)

	if err != nil {
		return fmt.Errorf("failed to set remote pool %s on token pool: %w", remotePool.Hex(), err)
	}
	log.Info().
		Str("Token Pool", pool.Address()).
		Str("Dest Pool", remotePool.Hex()).
		Uint64("RemoteChain", remoteChainSelector).
		Msg("Done setting remote pool")
	return pool.client.ProcessTransaction(tx)
}

func (pool *TokenPool) SetRemoteChainOnPool(remoteChainSelectors []uint64) error {
	log.Info().
		Str("Token Pool", pool.Address()).
		Msg("Setting remote chain on pool")
	var selectorsToUpdate []token_pool.TokenPoolChainUpdate
	for _, remoteChainSelector := range remoteChainSelectors {
		isSupported, err := pool.Instance.IsSupportedChain(nil, remoteChainSelector)
		if err != nil {
			return fmt.Errorf("failed to get if chain is supported: %w", err)
		}
		// Check if remote chain is already supported , if yes continue
		if isSupported {
			log.Info().
				Str("Token Pool", pool.Address()).
				Str(Network, pool.client.GetNetworkName()).
				Uint64("Remote Chain Selector", remoteChainSelector).
				Msg("Remote chain is already supported")
			continue
		}
		selectorsToUpdate = append(selectorsToUpdate, token_pool.TokenPoolChainUpdate{
			RemoteChainSelector: remoteChainSelector,
			Allowed:             true,
			InboundRateLimiterConfig: token_pool.RateLimiterConfig{
				IsEnabled: true,
				Capacity:  new(big.Int).Mul(big.NewInt(1e18), big.NewInt(1e9)),
				Rate:      new(big.Int).Mul(big.NewInt(1e18), big.NewInt(1e5)),
			},
			OutboundRateLimiterConfig: token_pool.RateLimiterConfig{
				IsEnabled: true,
				Capacity:  new(big.Int).Mul(big.NewInt(1e18), big.NewInt(1e9)),
				Rate:      new(big.Int).Mul(big.NewInt(1e18), big.NewInt(1e5)),
			},
		})
	}
	// if none to update return
	if len(selectorsToUpdate) == 0 {
		return nil
	}
	// If remote chain is not supported , add it
	opts, err := pool.client.TransactionOpts(pool.client.GetDefaultWallet())
	if err != nil {
		return fmt.Errorf("failed to get transaction opts: %w", err)
	}
	tx, err := pool.Instance.ApplyChainUpdates(opts, selectorsToUpdate)

	if err != nil {
		return fmt.Errorf("failed to set chain updates on token pool: %w", err)
	}

	log.Info().
		Str("Token Pool", pool.Address()).
		Uints64("Chain selectors", remoteChainSelectors).
		Str(Network, pool.client.GetNetworkConfig().Name).
		Msg("Remote chains set on token pool")
	return pool.client.ProcessTransaction(tx)
}

func (pool *TokenPool) SetRemoteChainRateLimits(remoteChainSelector uint64, rl token_pool.RateLimiterConfig) error {
	opts, err := pool.client.TransactionOpts(pool.client.GetDefaultWallet())
	if err != nil {
		return fmt.Errorf("error getting transaction opts: %w", err)
	}
	log.Info().
		Str("Token Pool", pool.Address()).
		Str("Remote chain selector", strconv.FormatUint(remoteChainSelector, 10)).
		Interface("RateLimiterConfig", rl).
		Msg("Setting Rate Limit on token pool")
	tx, err := pool.Instance.SetChainRateLimiterConfig(opts, remoteChainSelector, rl, rl)

	if err != nil {
		return fmt.Errorf("error setting rate limit token pool: %w", err)
	}

	log.Info().
		Str("Token Pool", pool.Address()).
		Str("Remote chain selector", strconv.FormatUint(remoteChainSelector, 10)).
		Interface("RateLimiterConfig", rl).
		Msg("Rate Limit on token pool is set")
	return pool.client.ProcessTransaction(tx)
}

type ARM struct {
	client     blockchain.EVMClient
	Instance   *arm_contract.ARMContract
	EthAddress common.Address
}

func (arm *ARM) Address() string {
	return arm.EthAddress.Hex()
}

type MockARM struct {
	client     blockchain.EVMClient
	Instance   *mock_arm_contract.MockARMContract
	EthAddress common.Address
}

func (arm *MockARM) SetClient(client blockchain.EVMClient) {
	arm.client = client
}
func (arm *MockARM) Address() string {
	return arm.EthAddress.Hex()
}

type CommitStoreReportAccepted struct {
	Min        uint64
	Max        uint64
	MerkleRoot [32]byte
	Raw        types.Log
}

type CommitStoreWrapper struct {
	Latest *commit_store.CommitStore
	V1_2_0 *commit_store_1_2_0.CommitStore
}

func (w CommitStoreWrapper) SetOCR2Config(opts *bind.TransactOpts,
	signers []common.Address,
	transmitters []common.Address,
	f uint8,
	onchainConfig []byte,
	offchainConfigVersion uint64,
	offchainConfig []byte,
) (*types.Transaction, error) {
	if w.Latest != nil {
		return w.Latest.SetOCR2Config(opts, signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
	}
	if w.V1_2_0 != nil {
		return w.V1_2_0.SetOCR2Config(opts, signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
	}
	return nil, fmt.Errorf("no instance found to set OCR2 config")
}

func (w CommitStoreWrapper) GetExpectedNextSequenceNumber(opts *bind.CallOpts) (uint64, error) {
	if w.Latest != nil {
		return w.Latest.GetExpectedNextSequenceNumber(opts)
	}
	if w.V1_2_0 != nil {
		return w.V1_2_0.GetExpectedNextSequenceNumber(opts)
	}
	return 0, fmt.Errorf("no instance found to get expected next sequence number")
}

type CommitStore struct {
	client     blockchain.EVMClient
	Instance   *CommitStoreWrapper
	EthAddress common.Address
}

func (b *CommitStore) Address() string {
	return b.EthAddress.Hex()
}

// SetOCR2Config sets the offchain reporting protocol configuration
func (b *CommitStore) SetOCR2Config(
	signers []common.Address,
	transmitters []common.Address,
	f uint8,
	onchainConfig []byte,
	offchainConfigVersion uint64,
	offchainConfig []byte,
) error {
	log.Info().Str("Contract Address", b.Address()).Msg("Configuring OCR config for CommitStore Contract")
	// Set Config
	opts, err := b.client.TransactionOpts(b.client.GetDefaultWallet())
	if err != nil {
		return fmt.Errorf("error getting transaction opts: %w", err)
	}

	log.Info().
		Interface("signerAddresses", signers).
		Interface("transmitterAddresses", transmitters).
		Str(Network, b.client.GetNetworkConfig().Name).
		Msg("Configuring CommitStore")
	tx, err := b.Instance.SetOCR2Config(
		opts,
		signers,
		transmitters,
		f,
		onchainConfig,
		offchainConfigVersion,
		offchainConfig,
	)

	if err != nil {
		return fmt.Errorf("error setting OCR2 config: %w", err)
	}
	return b.client.ProcessTransaction(tx)
}

// WatchReportAccepted watches for report accepted events
// There is no need to differentiate between the two versions of the contract as the event signature is the same
// we can cast the contract to the latest version
func (b *CommitStore) WatchReportAccepted(opts *bind.WatchOpts, acceptedEvent chan *commit_store.CommitStoreReportAccepted) (event.Subscription, error) {
	if b.Instance.Latest != nil {
		return b.Instance.Latest.WatchReportAccepted(opts, acceptedEvent)
	}
	if b.Instance.V1_2_0 != nil {
		newCommitStore, err := commit_store.NewCommitStore(b.EthAddress, wrappers.MustNewWrappedContractBackend(b.client, nil))
		if err != nil {
			return nil, fmt.Errorf("failed to create new CommitStore contract: %w", err)
		}
		return newCommitStore.WatchReportAccepted(opts, acceptedEvent)
	}
	log.Fatal().Msg("No instance found to watch for report accepted")
	return nil, fmt.Errorf("no instance found to watch for report accepted")
}

type ReceiverDapp struct {
	client     blockchain.EVMClient
	instance   *maybe_revert_message_receiver.MaybeRevertMessageReceiver
	EthAddress common.Address
}

func (rDapp *ReceiverDapp) Address() string {
	return rDapp.EthAddress.Hex()
}

func (rDapp *ReceiverDapp) ToggleRevert(revert bool) error {
	opts, err := rDapp.client.TransactionOpts(rDapp.client.GetDefaultWallet())
	if err != nil {
		return fmt.Errorf("error getting transaction opts: %w", err)
	}
	tx, err := rDapp.instance.SetRevert(opts, revert)
	if err != nil {
		return fmt.Errorf("error setting revert: %w", err)
	}
	log.Info().
		Bool("revert", revert).
		Str("tx", tx.Hash().String()).
		Str("ReceiverDapp", rDapp.Address()).
		Str(Network, rDapp.client.GetNetworkConfig().Name).
		Msg("ReceiverDapp revert set")
	return rDapp.client.ProcessTransaction(tx)
}

type InternalTimestampedPackedUint224 struct {
	Value     *big.Int
	Timestamp uint32
}

type PriceRegistryUsdPerUnitGasUpdated struct {
	DestChain uint64
	Value     *big.Int
	Timestamp *big.Int
	Raw       types.Log
}

type PriceRegistryWrappers struct {
	Latest *price_registry.PriceRegistry
	V1_2_0 *price_registry_1_2_0.PriceRegistry
}

func (p *PriceRegistryWrappers) GetTokenPrice(opts *bind.CallOpts, token common.Address) (*big.Int, error) {
	if p.Latest != nil {
		price, err := p.Latest.GetTokenPrice(opts, token)
		if err != nil {
			return nil, err
		}
		return price.Value, nil
	}
	if p.V1_2_0 != nil {
		p, err := p.V1_2_0.GetTokenPrice(opts, token)
		if err != nil {
			return nil, err
		}
		return p.Value, nil
	}
	return nil, fmt.Errorf("no instance found to get token price")
}

func (p *PriceRegistryWrappers) AddPriceUpdater(opts *bind.TransactOpts, addr common.Address) (*types.Transaction, error) {
	if p.Latest != nil {
		return p.Latest.ApplyPriceUpdatersUpdates(opts, []common.Address{addr}, []common.Address{})
	}
	if p.V1_2_0 != nil {
		return p.V1_2_0.ApplyPriceUpdatersUpdates(opts, []common.Address{addr}, []common.Address{})
	}
	return nil, fmt.Errorf("no instance found to add price updater")
}

func (p *PriceRegistryWrappers) AddFeeToken(opts *bind.TransactOpts, addr common.Address) (*types.Transaction, error) {
	if p.Latest != nil {
		return p.Latest.ApplyFeeTokensUpdates(opts, []common.Address{addr}, []common.Address{})
	}
	if p.V1_2_0 != nil {
		return p.V1_2_0.ApplyFeeTokensUpdates(opts, []common.Address{addr}, []common.Address{})
	}
	return nil, fmt.Errorf("no instance found to add fee token")
}

func (p *PriceRegistryWrappers) GetDestinationChainGasPrice(opts *bind.CallOpts, chainselector uint64) (InternalTimestampedPackedUint224, error) {
	if p.Latest != nil {
		price, err := p.Latest.GetDestinationChainGasPrice(opts, chainselector)
		if err != nil {
			return InternalTimestampedPackedUint224{}, err
		}
		return InternalTimestampedPackedUint224{
			Value:     price.Value,
			Timestamp: price.Timestamp,
		}, nil
	}
	if p.V1_2_0 != nil {
		price, err := p.V1_2_0.GetDestinationChainGasPrice(opts, chainselector)
		if err != nil {
			return InternalTimestampedPackedUint224{}, err
		}
		return InternalTimestampedPackedUint224{
			Value:     price.Value,
			Timestamp: price.Timestamp,
		}, nil
	}
	return InternalTimestampedPackedUint224{}, fmt.Errorf("no instance found to add fee token")
}

type InternalGasPriceUpdate struct {
	DestChainSelector uint64
	UsdPerUnitGas     *big.Int
}

type InternalTokenPriceUpdate struct {
	SourceToken common.Address
	UsdPerToken *big.Int
}

type PriceRegistry struct {
	client     blockchain.EVMClient
	Instance   *PriceRegistryWrappers
	EthAddress common.Address
}

func (c *PriceRegistry) Address() string {
	return c.EthAddress.Hex()
}

func (c *PriceRegistry) AddPriceUpdater(addr common.Address) error {
	opts, err := c.client.TransactionOpts(c.client.GetDefaultWallet())
	if err != nil {
		return fmt.Errorf("error getting transaction opts: %w", err)
	}
	tx, err := c.Instance.AddPriceUpdater(opts, addr)
	if err != nil {
		return fmt.Errorf("error adding price updater: %w", err)
	}
	log.Info().
		Str("updaters", addr.Hex()).
		Str(Network, c.client.GetNetworkConfig().Name).
		Msg("PriceRegistry updater added")
	return c.client.ProcessTransaction(tx)
}

func (c *PriceRegistry) AddFeeToken(addr common.Address) error {
	opts, err := c.client.TransactionOpts(c.client.GetDefaultWallet())
	if err != nil {
		return fmt.Errorf("error getting transaction opts: %w", err)
	}
	tx, err := c.Instance.AddFeeToken(opts, addr)
	if err != nil {
		return fmt.Errorf("error adding fee token: %w", err)
	}
	log.Info().
		Str("feeTokens", addr.Hex()).
		Str(Network, c.client.GetNetworkConfig().Name).
		Msg("PriceRegistry feeToken set")
	return c.client.ProcessTransaction(tx)
}

func (c *PriceRegistry) UpdatePrices(tokenUpdates []InternalTokenPriceUpdate, gasUpdates []InternalGasPriceUpdate) error {
	opts, err := c.client.TransactionOpts(c.client.GetDefaultWallet())
	if err != nil {
		return fmt.Errorf("error getting transaction opts: %w", err)
	}
	var tx *types.Transaction
	if c.Instance.Latest != nil {
		var tokenUpdatesLatest []price_registry.InternalTokenPriceUpdate
		var gasUpdatesLatest []price_registry.InternalGasPriceUpdate
		for _, update := range tokenUpdates {
			tokenUpdatesLatest = append(tokenUpdatesLatest, price_registry.InternalTokenPriceUpdate{
				SourceToken: update.SourceToken,
				UsdPerToken: update.UsdPerToken,
			})
		}
		for _, update := range gasUpdates {
			gasUpdatesLatest = append(gasUpdatesLatest, price_registry.InternalGasPriceUpdate{
				DestChainSelector: update.DestChainSelector,
				UsdPerUnitGas:     update.UsdPerUnitGas,
			})
		}
		tx, err = c.Instance.Latest.UpdatePrices(opts, price_registry.InternalPriceUpdates{
			TokenPriceUpdates: tokenUpdatesLatest,
			GasPriceUpdates:   gasUpdatesLatest,
		})
		if err != nil {
			return fmt.Errorf("error updating prices: %w", err)
		}
	}
	if c.Instance.V1_2_0 != nil {
		var tokenUpdates_1_2_0 []price_registry_1_2_0.InternalTokenPriceUpdate
		var gasUpdates_1_2_0 []price_registry_1_2_0.InternalGasPriceUpdate
		for _, update := range tokenUpdates {
			tokenUpdates_1_2_0 = append(tokenUpdates_1_2_0, price_registry_1_2_0.InternalTokenPriceUpdate{
				SourceToken: update.SourceToken,
				UsdPerToken: update.UsdPerToken,
			})
		}
		for _, update := range gasUpdates {
			gasUpdates_1_2_0 = append(gasUpdates_1_2_0, price_registry_1_2_0.InternalGasPriceUpdate{
				DestChainSelector: update.DestChainSelector,
				UsdPerUnitGas:     update.UsdPerUnitGas,
			})
		}
		tx, err = c.Instance.V1_2_0.UpdatePrices(opts, price_registry_1_2_0.InternalPriceUpdates{
			TokenPriceUpdates: tokenUpdates_1_2_0,
			GasPriceUpdates:   gasUpdates_1_2_0,
		})
		if err != nil {
			return fmt.Errorf("error updating prices: %w", err)
		}
	}
	if tx == nil {
		return fmt.Errorf("no instance found to update prices")
	}
	log.Info().
		Str(Network, c.client.GetNetworkConfig().Name).
		Interface("tokenUpdates", tokenUpdates).
		Interface("gasUpdates", gasUpdates).
		Msg("Prices updated")
	return c.client.ProcessTransaction(tx)
}

func (c *PriceRegistry) WatchUsdPerUnitGasUpdated(opts *bind.WatchOpts, latest chan *price_registry.PriceRegistryUsdPerUnitGasUpdated, destChain []uint64) (event.Subscription, error) {
	if c.Instance.Latest != nil {
		return c.Instance.Latest.WatchUsdPerUnitGasUpdated(opts, latest, destChain)
	}
	if c.Instance.V1_2_0 != nil {
		newP, err := price_registry.NewPriceRegistry(c.Instance.V1_2_0.Address(), wrappers.MustNewWrappedContractBackend(c.client, nil))
		if err != nil {
			return nil, fmt.Errorf("failed to create new PriceRegistry contract: %w", err)
		}
		return newP.WatchUsdPerUnitGasUpdated(opts, latest, destChain)
	}
	log.Fatal().Msg("No instance found to watch for price updates")
	return nil, fmt.Errorf("no instance found to watch for price updates")
}

type TokenAdminRegistry struct {
	client     blockchain.EVMClient
	Instance   *token_admin_registry.TokenAdminRegistry
	EthAddress common.Address
}

func (r *TokenAdminRegistry) Address() string {
	return r.EthAddress.Hex()
}

func (r *TokenAdminRegistry) SetAdminAndRegisterPool(tokenAddr, poolAddr common.Address) error {
	opts, err := r.client.TransactionOpts(r.client.GetDefaultWallet())
	if err != nil {
		return fmt.Errorf("error getting transaction opts: %w", err)
	}
	tx, err := r.Instance.RegisterAdministratorPermissioned(opts, tokenAddr, opts.From)
	if err != nil {
		return fmt.Errorf("error setting admin for token %s : %w", tokenAddr.Hex(), err)
	}
	err = r.client.ProcessTransaction(tx)
	if err != nil {
		return fmt.Errorf("error processing tx for setting admin on token %w", err)
	}
	log.Info().
		Str("Admin", opts.From.Hex()).
		Str("Token", tokenAddr.Hex()).
		Str("TokenAdminRegistry", r.Address()).
		Msg("Admin is set for token on TokenAdminRegistry")
	err = r.client.WaitForEvents()
	if err != nil {
		return fmt.Errorf("error waiting for tx for setting admin on pool %w", err)
	}
	opts, err = r.client.TransactionOpts(r.client.GetDefaultWallet())
	if err != nil {
		return fmt.Errorf("error getting transaction opts: %w", err)
	}
	tx, err = r.Instance.SetPool(opts, tokenAddr, poolAddr)
	if err != nil {
		return fmt.Errorf("error setting token %s and pool %s : %w", tokenAddr.Hex(), poolAddr.Hex(), err)
	}
	log.Info().
		Str("token", tokenAddr.Hex()).
		Str("Pool", poolAddr.Hex()).
		Str("TokenAdminRegistry", r.Address()).
		Msg("token and pool are set on TokenAdminRegistry")
	err = r.client.ProcessTransaction(tx)
	if err != nil {
		return fmt.Errorf("error processing tx for setting token %s and pool %s : %w", tokenAddr.Hex(), poolAddr.Hex(), err)
	}
	return nil
}

type Router struct {
	client     blockchain.EVMClient
	Instance   *router.Router
	EthAddress common.Address
}

func (r *Router) Address() string {
	return r.EthAddress.Hex()
}

func (r *Router) SetOnRamp(chainSelector uint64, onRamp common.Address) error {
	opts, err := r.client.TransactionOpts(r.client.GetDefaultWallet())
	if err != nil {
		return fmt.Errorf("error getting transaction opts: %w", err)
	}
	log.Info().
		Str("Router", r.Address()).
		Str("OnRamp", onRamp.Hex()).
		Str(Network, r.client.GetNetworkName()).
		Str("ChainSelector", strconv.FormatUint(chainSelector, 10)).
		Msg("Setting on ramp for r")

	tx, err := r.Instance.ApplyRampUpdates(opts, []router.RouterOnRamp{{DestChainSelector: chainSelector, OnRamp: onRamp}}, nil, nil)
	if err != nil {
		return fmt.Errorf("error applying ramp updates: %w", err)
	}
	log.Info().
		Str("onRamp", onRamp.Hex()).
		Str("Network Name", r.client.GetNetworkConfig().Name).
		Msg("Router is configured")
	return r.client.ProcessTransaction(tx)
}

func (r *Router) CCIPSend(destChainSelector uint64, msg router.ClientEVM2AnyMessage, valueForNative *big.Int) (*types.Transaction, error) {
	opts, err := r.client.TransactionOpts(r.client.GetDefaultWallet())
	if err != nil {
		return nil, fmt.Errorf("error getting transaction opts: %w", err)
	}
	if valueForNative != nil {
		opts.Value = valueForNative
	}

	log.Info().
		Str(Network, r.client.GetNetworkName()).
		Str("Router", r.Address()).
		Interface("TokensAndAmounts", msg.TokenAmounts).
		Str("FeeToken", msg.FeeToken.Hex()).
		Str("ExtraArgs", fmt.Sprintf("0x%x", msg.ExtraArgs[:])).
		Str("Receiver", fmt.Sprintf("0x%x", msg.Receiver[:])).
		Msg("Sending msg")
	return r.Instance.CcipSend(opts, destChainSelector, msg)
}

func (r *Router) CCIPSendAndProcessTx(destChainSelector uint64, msg router.ClientEVM2AnyMessage, valueForNative *big.Int) (*types.Transaction, error) {
	tx, err := r.CCIPSend(destChainSelector, msg, valueForNative)
	if err != nil {
		return nil, fmt.Errorf("failed to send msg: %w", err)
	}
	log.Info().
		Str("router", r.Address()).
		Str("txHash", tx.Hash().Hex()).
		Str(Network, r.client.GetNetworkConfig().Name).
		Str("chain selector", strconv.FormatUint(destChainSelector, 10)).
		Msg("msg is sent")
	return tx, r.client.ProcessTransaction(tx)
}

func (r *Router) AddOffRamp(offRamp common.Address, sourceChainId uint64) (*types.Transaction, error) {
	opts, err := r.client.TransactionOpts(r.client.GetDefaultWallet())
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction opts: %w", err)
	}
	tx, err := r.Instance.ApplyRampUpdates(opts, nil, nil, []router.RouterOffRamp{{SourceChainSelector: sourceChainId, OffRamp: offRamp}})
	if err != nil {
		return nil, fmt.Errorf("failed to add offRamp: %w", err)
	}
	log.Info().
		Str("offRamp", offRamp.Hex()).
		Str(Network, r.client.GetNetworkConfig().Name).
		Msg("offRamp is added to Router")
	return tx, r.client.ProcessTransaction(tx)
}

func (r *Router) SetWrappedNative(wNative common.Address) (*types.Transaction, error) {
	opts, err := r.client.TransactionOpts(r.client.GetDefaultWallet())
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction opts: %w", err)
	}
	tx, err := r.Instance.SetWrappedNative(opts, wNative)
	if err != nil {
		return nil, fmt.Errorf("failed to set wrapped native: %w", err)
	}
	log.Info().
		Str("wrapped native", wNative.Hex()).
		Str("router", r.Address()).
		Str(Network, r.client.GetNetworkConfig().Name).
		Msg("wrapped native is added for Router")
	return tx, r.client.ProcessTransaction(tx)
}

func (r *Router) GetFee(destChainSelector uint64, message router.ClientEVM2AnyMessage) (*big.Int, error) {
	return r.Instance.GetFee(nil, destChainSelector, message)
}

type SendReqEventData struct {
	MessageId      [32]byte
	SequenceNumber uint64
	Raw            types.Log
}

type OnRampWrapper struct {
	Latest *evm_2_evm_onramp.EVM2EVMOnRamp
	V1_2_0 *evm_2_evm_onramp_1_2_0.EVM2EVMOnRamp
}

func (w OnRampWrapper) SetNops(opts *bind.TransactOpts, owner common.Address) (*types.Transaction, error) {
	if w.Latest != nil {
		return w.Latest.SetNops(opts, []evm_2_evm_onramp.EVM2EVMOnRampNopAndWeight{
			{
				Nop:    owner,
				Weight: 1,
			},
		})
	}
	if w.V1_2_0 != nil {
		return w.V1_2_0.SetNops(opts, []evm_2_evm_onramp_1_2_0.EVM2EVMOnRampNopAndWeight{
			{
				Nop:    owner,
				Weight: 1,
			},
		})
	}
	return nil, fmt.Errorf("no instance found to set nops")
}

func (w OnRampWrapper) SetTokenTransferFeeConfig(
	opts *bind.TransactOpts,
	config []evm_2_evm_onramp.EVM2EVMOnRampTokenTransferFeeConfigArgs,
	addresses []common.Address,
) (*types.Transaction, error) {
	if w.Latest != nil {
		return w.Latest.SetTokenTransferFeeConfig(opts, config, addresses)
	}
	if w.V1_2_0 != nil {
		var configV12 []evm_2_evm_onramp_1_2_0.EVM2EVMOnRampTokenTransferFeeConfigArgs
		for _, c := range config {
			configV12 = append(configV12, evm_2_evm_onramp_1_2_0.EVM2EVMOnRampTokenTransferFeeConfigArgs{
				Token:             c.Token,
				MinFeeUSDCents:    c.MinFeeUSDCents,
				MaxFeeUSDCents:    c.MaxFeeUSDCents,
				DeciBps:           c.DeciBps,
				DestGasOverhead:   c.DestGasOverhead,
				DestBytesOverhead: c.DestBytesOverhead,
			})
		}
		return w.V1_2_0.SetTokenTransferFeeConfig(opts, configV12)
	}
	return nil, fmt.Errorf("no instance found to set token transfer fee config")
}

func (w OnRampWrapper) PayNops(opts *bind.TransactOpts) (*types.Transaction, error) {
	if w.Latest != nil {
		return w.Latest.PayNops(opts)
	}
	if w.V1_2_0 != nil {
		return w.V1_2_0.PayNops(opts)
	}
	return nil, fmt.Errorf("no instance found to pay nops")
}

func (w OnRampWrapper) WithdrawNonLinkFees(opts *bind.TransactOpts, native common.Address, owner common.Address) (*types.Transaction, error) {
	if w.Latest != nil {
		return w.Latest.WithdrawNonLinkFees(opts, native, owner)
	}
	if w.V1_2_0 != nil {
		return w.V1_2_0.WithdrawNonLinkFees(opts, native, owner)
	}
	return nil, fmt.Errorf("no instance found to withdraw non link fees")
}

func (w OnRampWrapper) SetRateLimiterConfig(opts *bind.TransactOpts, config evm_2_evm_onramp.RateLimiterConfig) (*types.Transaction, error) {
	if w.Latest != nil {
		return w.Latest.SetRateLimiterConfig(opts, config)
	}
	if w.V1_2_0 != nil {
		return w.V1_2_0.SetRateLimiterConfig(opts, evm_2_evm_onramp_1_2_0.RateLimiterConfig{
			IsEnabled: config.IsEnabled,
			Capacity:  config.Capacity,
			Rate:      config.Rate,
		})
	}
	return nil, fmt.Errorf("no instance found to set rate limiter config")
}

func (w OnRampWrapper) ParseCCIPSendRequested(l types.Log) (uint64, error) {
	if w.Latest != nil {
		sendReq, err := w.Latest.ParseCCIPSendRequested(l)
		if err != nil {
			return 0, err
		}
		return sendReq.Message.SequenceNumber, nil
	}
	if w.V1_2_0 != nil {
		sendReq, err := w.V1_2_0.ParseCCIPSendRequested(l)
		if err != nil {
			return 0, err
		}
		return sendReq.Message.SequenceNumber, nil
	}
	return 0, fmt.Errorf("no instance found to parse CCIPSendRequested")
}

func (w OnRampWrapper) GetDynamicConfig(opts *bind.CallOpts) (uint32, error) {
	if w.Latest != nil {
		cfg, err := w.Latest.GetDynamicConfig(opts)
		if err != nil {
			return 0, err
		}
		return cfg.MaxDataBytes, nil
	}
	if w.V1_2_0 != nil {
		cfg, err := w.V1_2_0.GetDynamicConfig(opts)
		if err != nil {
			return 0, err
		}
		return cfg.MaxDataBytes, nil
	}
	return 0, fmt.Errorf("no instance found to get dynamic config")
}

func (w OnRampWrapper) ApplyPoolUpdates(opts *bind.TransactOpts, tokens []common.Address, pools []common.Address) (*types.Transaction, error) {
	if w.Latest != nil {
		return nil, fmt.Errorf("latest version does not support ApplyPoolUpdates")
	}
	if w.V1_2_0 != nil {
		var poolUpdates []evm_2_evm_onramp_1_2_0.InternalPoolUpdate
		if len(tokens) != len(pools) {
			return nil, fmt.Errorf("tokens and pools length mismatch")
		}
		for i, token := range tokens {
			poolUpdates = append(poolUpdates, evm_2_evm_onramp_1_2_0.InternalPoolUpdate{
				Token: token,
				Pool:  pools[i],
			})
		}
		return w.V1_2_0.ApplyPoolUpdates(opts, []evm_2_evm_onramp_1_2_0.InternalPoolUpdate{}, poolUpdates)
	}
	return nil, fmt.Errorf("no instance found to apply pool updates")
}

func (w OnRampWrapper) CurrentRateLimiterState(opts *bind.CallOpts) (*RateLimiterConfig, error) {
	if w.Latest != nil {
		rlConfig, err := w.Latest.CurrentRateLimiterState(opts)
		if err != nil {
			return nil, err
		}
		return &RateLimiterConfig{
			IsEnabled: rlConfig.IsEnabled,
			Rate:      rlConfig.Rate,
			Capacity:  rlConfig.Capacity,
			Tokens:    rlConfig.Tokens,
		}, err
	}
	if w.V1_2_0 != nil {
		rlConfig, err := w.V1_2_0.CurrentRateLimiterState(opts)
		if err != nil {
			return nil, err
		}
		return &RateLimiterConfig{
			IsEnabled: rlConfig.IsEnabled,
			Rate:      rlConfig.Rate,
			Capacity:  rlConfig.Capacity,
			Tokens:    rlConfig.Tokens,
		}, err
	}
	return nil, fmt.Errorf("no instance found to get current rate limiter state")
}

type OnRamp struct {
	client     blockchain.EVMClient
	Instance   *OnRampWrapper
	EthAddress common.Address
}

// WatchCCIPSendRequested returns a subscription to watch for CCIPSendRequested events
// there is no difference in the event between the two versions
// so we can use the latest version to watch for events
func (onRamp *OnRamp) WatchCCIPSendRequested(opts *bind.WatchOpts, sendReqEvent chan *evm_2_evm_onramp.EVM2EVMOnRampCCIPSendRequested) (event.Subscription, error) {
	if onRamp.Instance.Latest != nil {
		return onRamp.Instance.Latest.WatchCCIPSendRequested(opts, sendReqEvent)
	}
	// cast the contract to the latest version so that we can watch for events with latest wrapper
	if onRamp.Instance.V1_2_0 != nil {
		newRamp, err := evm_2_evm_onramp.NewEVM2EVMOnRamp(onRamp.EthAddress, wrappers.MustNewWrappedContractBackend(onRamp.client, nil))
		if err != nil {
			return nil, fmt.Errorf("failed to cast to latest version: %w", err)
		}
		return newRamp.WatchCCIPSendRequested(opts, sendReqEvent)
	}
	// should never reach here
	log.Fatal().Msg("no instance found to watch for CCIPSendRequested")
	return nil, fmt.Errorf("no instance found to watch for CCIPSendRequested")
}

func (onRamp *OnRamp) Address() string {
	return onRamp.EthAddress.Hex()
}

func (onRamp *OnRamp) SetNops() error {
	opts, err := onRamp.client.TransactionOpts(onRamp.client.GetDefaultWallet())
	if err != nil {
		return fmt.Errorf("failed to get transaction opts: %w", err)
	}
	owner := common.HexToAddress(onRamp.client.GetDefaultWallet().Address())
	// set the payee to the default wallet
	tx, err := onRamp.Instance.SetNops(opts, owner)
	if err != nil {
		return fmt.Errorf("failed to set nops: %w", err)
	}
	return onRamp.client.ProcessTransaction(tx)
}

func (onRamp *OnRamp) SetTokenTransferFeeConfig(tokenTransferFeeConfig []evm_2_evm_onramp.EVM2EVMOnRampTokenTransferFeeConfigArgs) error {
	opts, err := onRamp.client.TransactionOpts(onRamp.client.GetDefaultWallet())
	if err != nil {
		return fmt.Errorf("failed to get transaction opts: %w", err)
	}
	tx, err := onRamp.Instance.SetTokenTransferFeeConfig(opts, tokenTransferFeeConfig, []common.Address{})
	if err != nil {
		return fmt.Errorf("failed to set token transfer fee config: %w", err)
	}
	log.Info().
		Interface("tokenTransferFeeConfig", tokenTransferFeeConfig).
		Str("onRamp", onRamp.Address()).
		Str(Network, onRamp.client.GetNetworkConfig().Name).
		Msg("TokenTransferFeeConfig set in OnRamp")
	return onRamp.client.ProcessTransaction(tx)
}

func (onRamp *OnRamp) PayNops() error {
	opts, err := onRamp.client.TransactionOpts(onRamp.client.GetDefaultWallet())
	if err != nil {
		return fmt.Errorf("failed to get transaction opts: %w", err)
	}
	tx, err := onRamp.Instance.PayNops(opts)
	if err != nil {
		return fmt.Errorf("failed to pay nops: %w", err)
	}
	return onRamp.client.ProcessTransaction(tx)
}

func (onRamp *OnRamp) WithdrawNonLinkFees(wrappedNative common.Address) error {
	opts, err := onRamp.client.TransactionOpts(onRamp.client.GetDefaultWallet())
	if err != nil {
		return fmt.Errorf("failed to get transaction opts: %w", err)
	}
	owner := common.HexToAddress(onRamp.client.GetDefaultWallet().Address())
	tx, err := onRamp.Instance.WithdrawNonLinkFees(opts, wrappedNative, owner)
	if err != nil {
		return fmt.Errorf("failed to withdraw non link fees: %w", err)
	}
	return onRamp.client.ProcessTransaction(tx)
}

func (onRamp *OnRamp) SetRateLimit(rlConfig evm_2_evm_onramp.RateLimiterConfig) error {
	opts, err := onRamp.client.TransactionOpts(onRamp.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := onRamp.Instance.SetRateLimiterConfig(opts, rlConfig)
	if err != nil {
		return fmt.Errorf("failed to set rate limit: %w", err)
	}
	log.Info().
		Bool("Enabled", rlConfig.IsEnabled).
		Str("capacity", rlConfig.Capacity.String()).
		Str("rate", rlConfig.Rate.String()).
		Str("onRamp", onRamp.Address()).
		Str(Network, onRamp.client.GetNetworkConfig().Name).
		Msg("Setting Rate limit in OnRamp")
	return onRamp.client.ProcessTransaction(tx)
}

func (onRamp *OnRamp) ApplyPoolUpdates(tokens []common.Address, pools []common.Address) error {
	// if the latest version is used, no need to apply pool updates
	if onRamp.Instance.Latest != nil {
		return nil
	}
	opts, err := onRamp.client.TransactionOpts(onRamp.client.GetDefaultWallet())
	if err != nil {
		return fmt.Errorf("failed to get transaction opts: %w", err)
	}
	tx, err := onRamp.Instance.ApplyPoolUpdates(opts, tokens, pools)
	if err != nil {
		return fmt.Errorf("failed to apply pool updates: %w", err)
	}
	log.Info().
		Interface("tokens", tokens).
		Interface("pools", pools).
		Str("onRamp", onRamp.Address()).
		Str(Network, onRamp.client.GetNetworkConfig().Name).
		Msg("poolUpdates set in OnRamp")
	return onRamp.client.ProcessTransaction(tx)
}

type OffRampWrapper struct {
	Latest *evm_2_evm_offramp.EVM2EVMOffRamp
	V1_2_0 *evm_2_evm_offramp_1_2_0.EVM2EVMOffRamp
}

func (offRamp *OffRampWrapper) CurrentRateLimiterState(opts *bind.CallOpts) (RateLimiterConfig, error) {
	if offRamp.Latest != nil {
		rlConfig, err := offRamp.Latest.CurrentRateLimiterState(opts)
		if err != nil {
			return RateLimiterConfig{}, err
		}
		return RateLimiterConfig{
			IsEnabled: rlConfig.IsEnabled,
			Capacity:  rlConfig.Capacity,
			Rate:      rlConfig.Rate,
		}, nil
	}
	if offRamp.V1_2_0 != nil {
		rlConfig, err := offRamp.V1_2_0.CurrentRateLimiterState(opts)
		if err != nil {
			return RateLimiterConfig{}, err
		}
		return RateLimiterConfig{
			IsEnabled: rlConfig.IsEnabled,
			Capacity:  rlConfig.Capacity,
			Rate:      rlConfig.Rate,
		}, nil
	}
	return RateLimiterConfig{}, fmt.Errorf("no instance found to get rate limiter state")
}

type EVM2EVMOffRampExecutionStateChanged struct {
	SequenceNumber uint64
	MessageId      [32]byte
	State          uint8
	ReturnData     []byte
	Raw            types.Log
}

type OffRamp struct {
	client     blockchain.EVMClient
	Instance   *OffRampWrapper
	EthAddress common.Address
}

func (offRamp *OffRamp) Address() string {
	return offRamp.EthAddress.Hex()
}

// WatchExecutionStateChanged returns a subscription to watch for ExecutionStateChanged events
// there is no difference in the event between the two versions
// so we can use the latest version to watch for events
func (offRamp *OffRamp) WatchExecutionStateChanged(opts *bind.WatchOpts, execEvent chan *evm_2_evm_offramp.EVM2EVMOffRampExecutionStateChanged, sequenceNumber []uint64, messageId [][32]byte) (event.Subscription, error) {
	if offRamp.Instance.Latest != nil {
		return offRamp.Instance.Latest.WatchExecutionStateChanged(opts, execEvent, sequenceNumber, messageId)
	}
	if offRamp.Instance.V1_2_0 != nil {
		newOffRamp, err := evm_2_evm_offramp.NewEVM2EVMOffRamp(offRamp.EthAddress, wrappers.MustNewWrappedContractBackend(offRamp.client, nil))
		if err != nil {
			return nil, fmt.Errorf("failed to cast to latest version of OffRamp from v1_2_0: %w", err)
		}
		return newOffRamp.WatchExecutionStateChanged(opts, execEvent, sequenceNumber, messageId)
	}
	log.Fatal().Msg("no instance found to watch for ExecutionStateChanged")
	return nil, fmt.Errorf("no instance found to watch for ExecutionStateChanged")
}

// SetOCR2Config sets the offchain reporting protocol configuration
func (offRamp *OffRamp) SetOCR2Config(
	signers []common.Address,
	transmitters []common.Address,
	f uint8,
	onchainConfig []byte,
	offchainConfigVersion uint64,
	offchainConfig []byte,
) error {
	log.Info().Str("Contract Address", offRamp.Address()).Msg("Configuring OffRamp Contract")
	// Set Config
	opts, err := offRamp.client.TransactionOpts(offRamp.client.GetDefaultWallet())
	if err != nil {
		return fmt.Errorf("failed to get transaction options: %w", err)
	}
	log.Info().
		Interface("signerAddresses", signers).
		Interface("transmitterAddresses", transmitters).
		Str(Network, offRamp.client.GetNetworkConfig().Name).
		Msg("Configuring OffRamp")
	if offRamp.Instance.Latest != nil {
		tx, err := offRamp.Instance.Latest.SetOCR2Config(
			opts,
			signers,
			transmitters,
			f,
			onchainConfig,
			offchainConfigVersion,
			offchainConfig,
		)
		if err != nil {
			return fmt.Errorf("failed to set OCR2 config: %w", err)
		}
		return offRamp.client.ProcessTransaction(tx)
	}
	if offRamp.Instance.V1_2_0 != nil {
		tx, err := offRamp.Instance.V1_2_0.SetOCR2Config(
			opts,
			signers,
			transmitters,
			f,
			onchainConfig,
			offchainConfigVersion,
			offchainConfig,
		)
		if err != nil {
			return fmt.Errorf("failed to set OCR2 config: %w", err)
		}
		return offRamp.client.ProcessTransaction(tx)
	}
	return fmt.Errorf("no instance found to set OCR2 config")
}

func (offRamp *OffRamp) UpdateRateLimitTokens(sourceTokens, destTokens []common.Address) error {
	if offRamp.Instance.V1_2_0 != nil {
		return nil
	}
	opts, err := offRamp.client.TransactionOpts(offRamp.client.GetDefaultWallet())
	if err != nil {
		return fmt.Errorf("failed to get transaction opts: %w", err)
	}
	if offRamp.Instance.Latest != nil {
		rateLimitTokens := make([]evm_2_evm_offramp.EVM2EVMOffRampRateLimitToken, len(sourceTokens))
		for i, sourceToken := range sourceTokens {
			rateLimitTokens[i] = evm_2_evm_offramp.EVM2EVMOffRampRateLimitToken{
				SourceToken: sourceToken,
				DestToken:   destTokens[i],
			}
		}

		tx, err := offRamp.Instance.Latest.UpdateRateLimitTokens(opts, []evm_2_evm_offramp.EVM2EVMOffRampRateLimitToken{}, rateLimitTokens)
		if err != nil {
			return fmt.Errorf("failed to apply rate limit tokens updates: %w", err)
		}
		log.Info().
			Interface("rateLimitToken adds", rateLimitTokens).
			Str("offRamp", offRamp.Address()).
			Str(Network, offRamp.client.GetNetworkConfig().Name).
			Msg("rateLimitTokens set in OffRamp")
		return offRamp.client.ProcessTransaction(tx)
	}
	return fmt.Errorf("no instance found to update rate limit tokens")
}

func (offRamp *OffRamp) SyncTokensAndPools(sourceTokens, pools []common.Address) error {
	if offRamp.Instance.Latest != nil {
		return nil
	}
	opts, err := offRamp.client.TransactionOpts(offRamp.client.GetDefaultWallet())
	if err != nil {
		return fmt.Errorf("failed to get transaction opts: %w", err)
	}
	if offRamp.Instance.V1_2_0 != nil {
		var tokenUpdates []evm_2_evm_offramp_1_2_0.InternalPoolUpdate
		for i, srcToken := range sourceTokens {
			tokenUpdates = append(tokenUpdates, evm_2_evm_offramp_1_2_0.InternalPoolUpdate{
				Token: srcToken,
				Pool:  pools[i],
			})
		}
		tx, err := offRamp.Instance.V1_2_0.ApplyPoolUpdates(opts, []evm_2_evm_offramp_1_2_0.InternalPoolUpdate{}, tokenUpdates)
		if err != nil {
			return fmt.Errorf("failed to apply pool updates: %w", err)
		}
		log.Info().
			Interface("tokenUpdates", tokenUpdates).
			Str("offRamp", offRamp.Address()).
			Str(Network, offRamp.client.GetNetworkConfig().Name).
			Msg("tokenUpdates set in OffRamp")
		return offRamp.client.ProcessTransaction(tx)
	}
	return fmt.Errorf("no instance found to sync tokens and pools")
}

type MockAggregator struct {
	client          blockchain.EVMClient
	Instance        *mock_v3_aggregator_contract.MockV3Aggregator
	ContractAddress common.Address
}

func (a *MockAggregator) ChainID() uint64 {
	return a.client.GetChainID().Uint64()
}

func (a *MockAggregator) UpdateRoundData(answer *big.Int) error {
	opts, err := a.client.TransactionOpts(a.client.GetDefaultWallet())
	if err != nil {
		return fmt.Errorf("unable to get transaction opts: %w", err)
	}
	log.Info().
		Str("Contract Address", a.ContractAddress.Hex()).
		Str("Network Name", a.client.GetNetworkConfig().Name).
		Msg("Updating Round Data")
	// we get the round from latest round data
	// if there is any error in fetching the round , we set the round with a random number
	// otherwise increase the latest round by 1 and set the value for the next round
	round, err := a.Instance.LatestRound(nil)
	if err != nil {
		rand.Seed(uint64(time.Now().UnixNano()))
		round = big.NewInt(int64(rand.Uint64()))
	}
	round = new(big.Int).Add(round, big.NewInt(1))
	tx, err := a.Instance.UpdateRoundData(opts, round, answer, big.NewInt(time.Now().UTC().UnixNano()), big.NewInt(time.Now().UTC().UnixNano()))
	if err != nil {
		return fmt.Errorf("unable to update round data: %w", err)
	}
	log.Info().
		Str("Contract Address", a.ContractAddress.Hex()).
		Str("Network Name", a.client.GetNetworkConfig().Name).
		Str("Round", round.String()).
		Str("Answer", answer.String()).
		Msg("Updated Round Data")
	_, err = bind.WaitMined(context.Background(), a.client.DeployBackend(), tx)
	if err != nil {
		return fmt.Errorf("error waiting for tx %s to be mined", tx.Hash().Hex())
	}

	return a.client.MarkTxAsSentOnL2(tx)
}
