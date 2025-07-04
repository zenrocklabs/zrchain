/* eslint-disable @typescript-eslint/no-unused-vars */
import { useQuery, type UseQueryOptions, useInfiniteQuery, type UseInfiniteQueryOptions } from "@tanstack/react-query";
import { useClient } from '../useClient';

export default function useZrchainZenbtc() {
  const client = useClient();
  const QueryQueryParams = ( options: any) => {
    const key = { type: 'QueryQueryParams',  };    
    return useQuery([key], () => {
      return  client.ZrchainZenbtc.query.queryQueryParams().then( res => res.data );
    }, options);
  }
  
  const QueryGetLockTransactions = ( options: any) => {
    const key = { type: 'QueryGetLockTransactions',  };    
    return useQuery([key], () => {
      return  client.ZrchainZenbtc.query.queryGetLockTransactions().then( res => res.data );
    }, options);
  }
  
  const QueryGetRedemptions = (query: any, options: any) => {
    const key = { type: 'QueryGetRedemptions', query };    
    return useQuery([key], () => {
      const {query } = key
      return  client.ZrchainZenbtc.query.queryGetRedemptions(query ?? undefined).then( res => res.data );
    }, options);
  }
  
  const QueryQueryPendingMintTransactions = (query: any, options: any) => {
    const key = { type: 'QueryQueryPendingMintTransactions', query };    
    return useQuery([key], () => {
      const {query } = key
      return  client.ZrchainZenbtc.query.queryQueryPendingMintTransactions(query ?? undefined).then( res => res.data );
    }, options);
  }
  
  const QueryQueryPendingMintTransaction = (tx_hash: string,  options: any) => {
    const key = { type: 'QueryQueryPendingMintTransaction',  tx_hash };    
    return useQuery([key], () => {
      const { tx_hash } = key
      return  client.ZrchainZenbtc.query.queryQueryPendingMintTransaction(tx_hash).then( res => res.data );
    }, options);
  }
  
  const QueryQuerySupply = ( options: any) => {
    const key = { type: 'QueryQuerySupply',  };    
    return useQuery([key], () => {
      return  client.ZrchainZenbtc.query.queryQuerySupply().then( res => res.data );
    }, options);
  }
  
  const QueryQueryBurnEvents = (query: any, options: any) => {
    const key = { type: 'QueryQueryBurnEvents', query };    
    return useQuery([key], () => {
      const {query } = key
      return  client.ZrchainZenbtc.query.queryQueryBurnEvents(query ?? undefined).then( res => res.data );
    }, options);
  }
  
  return {QueryQueryParams,QueryGetLockTransactions,QueryGetRedemptions,QueryQueryPendingMintTransactions,QueryQueryPendingMintTransaction,QueryQuerySupply,QueryQueryBurnEvents,
  }
}
