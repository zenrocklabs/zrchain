/* eslint-disable @typescript-eslint/no-unused-vars */
import { useQuery, type UseQueryOptions, useInfiniteQuery, type UseInfiniteQueryOptions } from "@tanstack/react-query";
import { useClient } from '../useClient';

export default function useZrchainTreasury() {
  const client = useClient();
  const QueryParams = ( options: any) => {
    const key = { type: 'QueryParams',  };    
    return useQuery([key], () => {
      return  client.ZrchainTreasury.query.queryParams().then( res => res.data );
    }, options);
  }
  
  const QueryKeyRequests = (keyring_addr: string, status: string, workspace_addr: string, query: any, options: any, perPage: number) => {
    const key = { type: 'QueryKeyRequests',  keyring_addr,  status,  workspace_addr, query };    
    return useInfiniteQuery([key], ({pageParam = 1}: { pageParam?: number}) => {
      const { keyring_addr,  status,  workspace_addr,query } = key

      query['pagination.limit']=perPage;
      query['pagination.offset']= (pageParam-1)*perPage;
      query['pagination.count_total']= true;
      return  client.ZrchainTreasury.query.queryKeyRequests(keyring_addr, status, workspace_addr, query ?? undefined).then( res => ({...res.data,pageParam}) );
    }, {...options,
      getNextPageParam: (lastPage, allPages) => { if ((lastPage.pagination?.total ?? 0) >((lastPage.pageParam ?? 0) * perPage)) {return lastPage.pageParam+1 } else {return undefined}},
      getPreviousPageParam: (firstPage, allPages) => { if (firstPage.pageParam==1) { return undefined } else { return firstPage.pageParam-1}}
    }
    );
  }
  
  const QueryKeyRequestByID = (id: string,  options: any) => {
    const key = { type: 'QueryKeyRequestByID',  id };    
    return useQuery([key], () => {
      const { id } = key
      return  client.ZrchainTreasury.query.queryKeyRequestById(id).then( res => res.data );
    }, options);
  }
  
  const QueryKeys = (workspace_addr: string, wallet_type: string, prefixes: string, query: any, options: any, perPage: number) => {
    const key = { type: 'QueryKeys',  workspace_addr,  wallet_type,  prefixes, query };    
    return useInfiniteQuery([key], ({pageParam = 1}: { pageParam?: number}) => {
      const { workspace_addr,  wallet_type,  prefixes,query } = key

      query['pagination.limit']=perPage;
      query['pagination.offset']= (pageParam-1)*perPage;
      query['pagination.count_total']= true;
      return  client.ZrchainTreasury.query.queryKeys(workspace_addr, wallet_type, prefixes, query ?? undefined).then( res => ({...res.data,pageParam}) );
    }, {...options,
      getNextPageParam: (lastPage, allPages) => { if ((lastPage.pagination?.total ?? 0) >((lastPage.pageParam ?? 0) * perPage)) {return lastPage.pageParam+1 } else {return undefined}},
      getPreviousPageParam: (firstPage, allPages) => { if (firstPage.pageParam==1) { return undefined } else { return firstPage.pageParam-1}}
    }
    );
  }
  
  const QueryKeyByID = (id: string, wallet_type: string, prefixes: string,  options: any) => {
    const key = { type: 'QueryKeyByID',  id,  wallet_type,  prefixes };    
    return useQuery([key], () => {
      const { id,  wallet_type,  prefixes } = key
      return  client.ZrchainTreasury.query.queryKeyById(id, wallet_type, prefixes).then( res => res.data );
    }, options);
  }
  
  const QuerySignatureRequests = (keyring_addr: string, status: string, query: any, options: any, perPage: number) => {
    const key = { type: 'QuerySignatureRequests',  keyring_addr,  status, query };    
    return useInfiniteQuery([key], ({pageParam = 1}: { pageParam?: number}) => {
      const { keyring_addr,  status,query } = key

      query['pagination.limit']=perPage;
      query['pagination.offset']= (pageParam-1)*perPage;
      query['pagination.count_total']= true;
      return  client.ZrchainTreasury.query.querySignatureRequests(keyring_addr, status, query ?? undefined).then( res => ({...res.data,pageParam}) );
    }, {...options,
      getNextPageParam: (lastPage, allPages) => { if ((lastPage.pagination?.total ?? 0) >((lastPage.pageParam ?? 0) * perPage)) {return lastPage.pageParam+1 } else {return undefined}},
      getPreviousPageParam: (firstPage, allPages) => { if (firstPage.pageParam==1) { return undefined } else { return firstPage.pageParam-1}}
    }
    );
  }
  
  const QuerySignatureRequestByID = (id: string,  options: any) => {
    const key = { type: 'QuerySignatureRequestByID',  id };    
    return useQuery([key], () => {
      const { id } = key
      return  client.ZrchainTreasury.query.querySignatureRequestById(id).then( res => res.data );
    }, options);
  }
  
  const QuerySignTransactionRequests = (wallet_type: string, key_id: string, status: string, query: any, options: any, perPage: number) => {
    const key = { type: 'QuerySignTransactionRequests',  wallet_type,  key_id,  status, query };    
    return useInfiniteQuery([key], ({pageParam = 1}: { pageParam?: number}) => {
      const { wallet_type,  key_id,  status,query } = key

      query['pagination.limit']=perPage;
      query['pagination.offset']= (pageParam-1)*perPage;
      query['pagination.count_total']= true;
      return  client.ZrchainTreasury.query.querySignTransactionRequests(wallet_type, key_id, status, query ?? undefined).then( res => ({...res.data,pageParam}) );
    }, {...options,
      getNextPageParam: (lastPage, allPages) => { if ((lastPage.pagination?.total ?? 0) >((lastPage.pageParam ?? 0) * perPage)) {return lastPage.pageParam+1 } else {return undefined}},
      getPreviousPageParam: (firstPage, allPages) => { if (firstPage.pageParam==1) { return undefined } else { return firstPage.pageParam-1}}
    }
    );
  }
  
  const QuerySignTransactionRequestByID = (id: string,  options: any) => {
    const key = { type: 'QuerySignTransactionRequestByID',  id };    
    return useQuery([key], () => {
      const { id } = key
      return  client.ZrchainTreasury.query.querySignTransactionRequestById(id).then( res => res.data );
    }, options);
  }
  
  const QueryZrSignKeys = (address: string, walletType: string, query: any, options: any, perPage: number) => {
    const key = { type: 'QueryZrSignKeys',  address,  walletType, query };    
    return useInfiniteQuery([key], ({pageParam = 1}: { pageParam?: number}) => {
      const { address,  walletType,query } = key

      query['pagination.limit']=perPage;
      query['pagination.offset']= (pageParam-1)*perPage;
      query['pagination.count_total']= true;
      return  client.ZrchainTreasury.query.queryZrSignKeys(address, walletType, query ?? undefined).then( res => ({...res.data,pageParam}) );
    }, {...options,
      getNextPageParam: (lastPage, allPages) => { if ((lastPage.pagination?.total ?? 0) >((lastPage.pageParam ?? 0) * perPage)) {return lastPage.pageParam+1 } else {return undefined}},
      getPreviousPageParam: (firstPage, allPages) => { if (firstPage.pageParam==1) { return undefined } else { return firstPage.pageParam-1}}
    }
    );
  }
  
  const QueryKeyByAddress = (query: any, options: any) => {
    const key = { type: 'QueryKeyByAddress', query };    
    return useQuery([key], () => {
      const {query } = key
      return  client.ZrchainTreasury.query.queryKeyByAddress(query ?? undefined).then( res => res.data );
    }, options);
  }
  
  const QueryZenbtcWallets = (recipient_addr: string, chain_type: string, mint_chain_id: string, return_addr: string, query: any, options: any, perPage: number) => {
    const key = { type: 'QueryZenbtcWallets',  recipient_addr,  chain_type,  mint_chain_id,  return_addr, query };    
    return useInfiniteQuery([key], ({pageParam = 1}: { pageParam?: number}) => {
      const { recipient_addr,  chain_type,  mint_chain_id,  return_addr,query } = key

      query['pagination.limit']=perPage;
      query['pagination.offset']= (pageParam-1)*perPage;
      query['pagination.count_total']= true;
      return  client.ZrchainTreasury.query.queryZenbtcWallets(recipient_addr, chain_type, mint_chain_id, return_addr, query ?? undefined).then( res => ({...res.data,pageParam}) );
    }, {...options,
      getNextPageParam: (lastPage, allPages) => { if ((lastPage.pagination?.total ?? 0) >((lastPage.pageParam ?? 0) * perPage)) {return lastPage.pageParam+1 } else {return undefined}},
      getPreviousPageParam: (firstPage, allPages) => { if (firstPage.pageParam==1) { return undefined } else { return firstPage.pageParam-1}}
    }
    );
  }
  
  return {QueryParams,QueryKeyRequests,QueryKeyRequestByID,QueryKeys,QueryKeyByID,QuerySignatureRequests,QuerySignatureRequestByID,QuerySignTransactionRequests,QuerySignTransactionRequestByID,QueryZrSignKeys,QueryKeyByAddress,QueryZenbtcWallets,
  }
}
