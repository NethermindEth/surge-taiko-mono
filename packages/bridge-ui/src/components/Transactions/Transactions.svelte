<script lang="ts">
  import { onMount } from 'svelte';
  import { t } from 'svelte-i18n';
  import type { Address } from 'viem';

  import { activeBridge } from '$components/Bridge/state';
  import { destNetwork } from '$components/Bridge/state';
  import { BridgeTypes } from '$components/Bridge/types';
  import { Button } from '$components/Button';
  import { Card } from '$components/Card';
  import { ChainSelectorDirection, ChainSelectorType } from '$components/ChainSelectors';
  import ChainSelector from '$components/ChainSelectors/ChainSelector.svelte';
  import { DesktopOrLarger } from '$components/DesktopOrLarger';
  import { Icon } from '$components/Icon';
  import RotatingIcon from '$components/Icon/RotatingIcon.svelte';
  import { warningToast } from '$components/NotificationToast';
  import OnAccount from '$components/OnAccount/OnAccount.svelte';
  import { Paginator } from '$components/Paginator';
  import { Spinner } from '$components/Spinner';
  import StatusDot from '$components/StatusDot/StatusDot.svelte';
  import { transactionConfig } from '$config';
  import { type BridgeTransaction, fetchTransactions, MessageStatus } from '$libs/bridge';
  import { chainIdToChain } from '$libs/chain';
  import { getAlternateNetwork } from '$libs/network';
  import { bridgeTxService } from '$libs/storage';
  import { TokenType } from '$libs/token';
  import { isDesktop, isTablet } from '$libs/util/responsiveCheck';
  import { account } from '$stores';
  import type { Account } from '$stores/account';

  import { StatusFilterDialog, StatusFilterDropdown } from './Filter';
  import { FungibleTransactionRow, NftTransactionRow } from './Rows/';
  import { StatusInfoDialog } from './Status';
  import { connectedSourceChain } from '$stores/network';

  let transactions: BridgeTransaction[] = [];

  let currentPage = 1;

  let isBlurred = false;
  const transitionTime = transactionConfig.blurTransitionTime;

  let totalItems = 0;
  let pageSize = transactionConfig.pageSizeDesktop;

  let loadingTxs = false;

  let isDesktopOrLarger: boolean;

  let selectedStatus: MessageStatus | null = null; // null indicates no filter is applied

  let menuOpen = false;

  const toggleMenu = () => {
    menuOpen = !menuOpen;
  };

  const handlePageChange = (detail: number) => {
    isBlurred = true;
    setTimeout(() => {
      currentPage = detail;
      isBlurred = false;
    }, transitionTime);
  };

  const getTransactionsToShow = (page: number, pageSize: number, bridgeTx: BridgeTransaction[]) => {
    const start = (page - 1) * pageSize;
    const end = start + pageSize;
    return bridgeTx.slice(start, end);
  };

  const onAccountChange = async (newAccount: Account, oldAccount?: Account) => {
    // We want to make sure that we are connected and only
    // fetch if the account has changed
    if (newAccount?.isConnected && newAccount.address && newAccount.address !== oldAccount?.address) {
      await updateTransactions(newAccount.address);
    }
  };

  const refresh = async () => {
    if ($account?.address) {
      await updateTransactions($account.address);
    }
  };

  const handleTransactionRemoved = () => {
    refresh();
  };

  const updateTransactions = async (address: Address) => {
    if (loadingTxs) return;
    loadingTxs = true;
    // Wait for connectedSourceChain to be set if it's null
    if (!$connectedSourceChain) {
      await new Promise<void>(resolve => {
      const unsubscribe = connectedSourceChain.subscribe((value: any) => {
        if (value) {
        unsubscribe();
        resolve();
        }
      });
      });
    }
    
    const { mergedTransactions, outdatedLocalTransactions, error } = await fetchTransactions(
      address, 
      $connectedSourceChain?.id
    );
    transactions = mergedTransactions;

    if (outdatedLocalTransactions.length > 0) {
      await bridgeTxService.removeTransactions(address, outdatedLocalTransactions);
    }
    if (error) {
      warningToast({ title: $t('transactions.errors.relayer_offline') });
    }
    loadingTxs = false;
  };

  let previousAccount: Account | null = null;
  // refresh only if previous account is different from current account
  $: if (($account && previousAccount && $account.address !== previousAccount.address) || !previousAccount) {
    refresh();
    previousAccount = $account;
  }

  $: statusFilteredTransactions =
    selectedStatus !== null ? transactions.filter((tx) => tx.msgStatus === selectedStatus) : transactions;

  $: tokenAndStatusFilteredTransactions = statusFilteredTransactions.filter((tx) =>
    displayTokenTypesBasedOnType.includes(tx.tokenType),
  );

  $: transactionsToShow = getTransactionsToShow(currentPage, pageSize, tokenAndStatusFilteredTransactions);
  

  $: fungibleView = $activeBridge === BridgeTypes.FUNGIBLE;
  $: nftView = $activeBridge === BridgeTypes.NFT;

  $: fungibleTokens = [TokenType.ERC20, TokenType.ETH];
  $: nftTokens = [TokenType.ERC721, TokenType.ERC1155];
  $: allTokens = [...fungibleTokens, ...nftTokens];

  $: displayTokenTypesBasedOnType = fungibleView ? fungibleTokens : nftView ? nftTokens : allTokens;

  $: filteredTransactions = transactions.filter((tx) => displayTokenTypesBasedOnType.includes(tx.tokenType));

  $: pageSize = isDesktopOrLarger ? transactionConfig.pageSizeDesktop : transactionConfig.pageSizeMobile;

  $: totalItems = filteredTransactions.length;

  // Some shortcuts to make the code more readable
  $: isConnected = $account?.isConnected;
  $: hasTxs = filteredTransactions.length > 0;

  // Controls what we render on the page
  $: renderLoading = loadingTxs && isConnected;
  $: renderTransactions = !renderLoading && isConnected && hasTxs;
  $: renderNoTransactions = !renderLoading && transactionsToShow.length === 0;

  onMount(() => {
    const alternateChainID = getAlternateNetwork();
    if (!$destNetwork && alternateChainID) {
      // if only two chains are available, set the destination chain to the other one
      $destNetwork = chainIdToChain(alternateChainID);
    }
  });

</script>

<div class="flex flex-col justify-center w-full">
  <Card title={$t('transactions.title')} text={$t('transactions.description')}>
    <div class="space-y-[35px]">
      {#if $isDesktop}
        <div class="my-[30px] f-between-center max-h-[36px] gap-2">
          <ChainSelector
            type={ChainSelectorType.SMALL}
            direction={ChainSelectorDirection.SOURCE}
            label={$t('chain_selector.currently_on')}
            switchWallet />
          <div class="flex gap-2">
            <Button
              type="neutral"
              shape="circle"
              class="bg-neutral rounded-full !min-w-[36px] !min-h-[36px] !max-w-[36px] !max-h-[36px] border-none"
              on:click={async () => await refresh()}>
              <RotatingIcon loading={loadingTxs} type="refresh" size={16} />
            </Button>
            <StatusFilterDropdown bind:selectedStatus />
          </div>
        </div>
      {:else}
        <div class="f-row justify-between my-[30px]">
          <div class="f-row items-center gap-[10px]">
            <StatusDot type="success" simple={false} />
            <ChainSelector type={ChainSelectorType.SMALL} direction={ChainSelectorDirection.SOURCE} switchWallet />
          </div>
          <div class="f-row items-center gap-[5px]">
            {#if $account && $account?.address}
              <button
                class="grid place-items-center bg-neutral min-w-[36px] max-w-[36px] min-h-[36px] max-h-[36px] rounded-full"
                on:click|stopPropagation={toggleMenu}>
                <Icon type="settings" fillClass="fill-primary-icon" size={18} class="self-center" />
              </button>
              <Button
                type="neutral"
                shape="circle"
                class="bg-neutral rounded-full !min-w-[36px] !min-h-[36px] !max-w-[36px] !max-h-[36px] border-none"
                on:click={async () => await refresh()}>
                <RotatingIcon loading={loadingTxs} type="refresh" size={16} />
              </Button>
            {/if}
          </div>
        </div>
      {/if}

      <div
        class="flex flex-col"
        style={`min-height: calc(${transactionsToShow.length} * ${isDesktopOrLarger ? '80px' : '66px'});`}>
        <div class="h-sep !mb-0 display-inline" />

        <div class="text-primary-content flex text-primary-content w-full my-[5px] md:my-[0px] px-[14px] py-[10px]">
          {#if $activeBridge === BridgeTypes.FUNGIBLE}
            {#if $isDesktop}
              <div class="w-1/6 py-2 text-secondary-content">{$t('transactions.header.from')}</div>
              <div class="w-1/6 py-2 text-secondary-content">{$t('transactions.header.to')}</div>
              <div class="w-1/6 py-2 text-secondary-content">{$t('transactions.header.amount')}</div>
              <div class="w-1/6 py-2 text-secondary-content flex flex-row">
                {$t('transactions.header.status')}
                <StatusInfoDialog />
              </div>
              <div class="w-1/6 py-2 text-secondary-content">{$t('transactions.header.date')}</div>
              <div class="w-1/6 py-2 text-secondary-content"></div>
            {:else if $isTablet}
              <div class="w-1/4 py-2 text-secondary-content">{$t('transactions.header.from')}</div>
              <div class="w-1/4 py-2 text-secondary-content">{$t('transactions.header.to')}</div>
              <div class="w-1/4 py-2 text-secondary-content">{$t('transactions.header.amount')}</div>
              <div class="w-1/4 py-2 text-secondary-content flex flex-row">
                {$t('transactions.header.status')}
                <StatusInfoDialog />
              </div>
            {:else}
              <div class="w-1/2 text-center text-secondary-content">{$t('transactions.header.amount')}</div>
              <div class="w-1/2 pr-[14px] f-row items-center justify-end text-secondary-content">
                {$t('transactions.header.status')}
                <StatusInfoDialog />
              </div>
            {/if}
          {:else if $activeBridge === BridgeTypes.NFT}
            {#if $isDesktop}
              <div class="w-1/6 py-2 text-secondary-content">{$t('transactions.header.nft')}</div>
              <div class="w-1/6 py-2 text-secondary-content">{$t('transactions.header.from')}</div>
              <div class="w-1/6 py-2 text-secondary-content">{$t('transactions.header.to')}</div>
              <div class="w-1/6 py-2 text-secondary-content flex flex-row">
                {$t('transactions.header.status')}
                <StatusInfoDialog />
              </div>
              <div class="w-1/6 py-2 text-secondary-content">{$t('transactions.header.date')}</div>
              <div class="w-1/6 py-2 text-secondary-content"></div>
            {:else if $isTablet}
              <div class="w-1/4 py-2 text-secondary-content">{$t('transactions.header.nft')}</div>
              <div class="w-1/4 py-2 text-secondary-content">{$t('transactions.header.from')}</div>
              <div class="w-1/4 py-2 text-secondary-content">{$t('transactions.header.to')}</div>

              <div class="w-1/4 py-2 text-secondary-content flex flex-row">
                {$t('transactions.header.status')}
                <StatusInfoDialog />
              </div>
            {:else}
              <div class="w-1/3 text-left pl-[11px] text-secondary-content">
                {$t('transactions.header.details')}
              </div>
              <div class="w-1/3 text-center text-secondary-content">{$t('transactions.header.nft')}</div>
              <div class="w-1/3 pr-[14px] f-row items-center justify-end text-secondary-content">
                {$t('transactions.header.status')}
                <StatusInfoDialog />
              </div>
            {/if}
          {/if}
        </div>
        <div class="h-sep !my-0" />

        {#if renderLoading}
          <div class="flex items-center justify-center text-primary-content h-[80px]">
            <Spinner /> <span class="pl-3">{$t('common.loading')}...</span>
          </div>
        {/if}

        {#if renderTransactions}
          <div
            class="flex flex-col items-center"
            style={isBlurred ? `filter: blur(5px); transition: filter ${transitionTime / 1000}s ease-in-out` : ''}>
            {#each transactionsToShow as bridgeTx (bridgeTx.srcTxHash)}
              {@const status = bridgeTx.msgStatus}
              {@const isFungible = bridgeTx.tokenType === TokenType.ERC20 || bridgeTx.tokenType === TokenType.ETH}
              {#if isFungible}
                <FungibleTransactionRow bind:bridgeTx {handleTransactionRemoved} bridgeTxStatus={status} />
              {:else}
                <NftTransactionRow bind:bridgeTx {handleTransactionRemoved} bridgeTxStatus={status} />
              {/if}
              <div class="h-sep !my-0 display-inline" />
            {/each}
          </div>
        {/if}

        {#if renderNoTransactions}
          <div class="flex items-center justify-center text-primary-content h-[80px]">
            <span class="pl-3">{$t('transactions.no_transactions')}</span>
          </div>
        {/if}
      </div>
    </div>
  </Card>

  <div class="flex justify-center lg:justify-end pb-5">
    <Paginator {pageSize} {totalItems} on:pageChange={({ detail }) => handlePageChange(detail)} />
  </div>

  <StatusFilterDialog bind:selectedStatus bind:menuOpen />
</div>

<OnAccount change={onAccountChange} />

<DesktopOrLarger bind:is={isDesktopOrLarger} />
