/* eslint-disable @typescript-eslint/no-unused-vars */
import { useQuery, type UseQueryOptions, useInfiniteQuery, type UseInfiniteQueryOptions } from "@tanstack/react-query";
import { useClient } from '../useClient';

export default function useZrchainZentp() {
  const client = useClient();
  const QueryParams = ( options: any) => {
    const key = { type: 'QueryParams',  };    
    return useQuery([key], () => {
      return  client.ZrchainZentp.query.queryParams().then( res => res.data );
    }, options);
  }
  
  const QueryMints = (id: string, creator: string, status: string, denom: string, query: any, options: any, perPage: number) => {
    const key = { type: 'QueryMints',  id,  creator,  status,  denom, query };    
    return useInfiniteQuery([key], ({pageParam = 1}: { pageParam?: number}) => {
      const { id,  creator,  status,  denom,query } = key

      query['pagination.limit']=perPage;
      query['pagination.offset']= (pageParam-1)*perPage;
      query['pagination.count_total']= true;
      return  client.ZrchainZentp.query.queryMints(id, creator, status, denom, query ?? undefined).then( res => ({...res.data,pageParam}) );
    }, {...options,
      getNextPageParam: (lastPage, allPages) => { if ((lastPage.pagination?.total ?? 0) >((lastPage.pageParam ?? 0) * perPage)) {return lastPage.pageParam+1 } else {return undefined}},
      getPreviousPageParam: (firstPage, allPages) => { if (firstPage.pageParam==1) { return undefined } else { return firstPage.pageParam-1}}
    }
    );
  }
  
  const QueryBurns = (id: string, denom: string, query: any, options: any, perPage: number) => {
    const key = { type: 'QueryBurns',  id,  denom, query };    
    return useInfiniteQuery([key], ({pageParam = 1}: { pageParam?: number}) => {
      const { id,  denom,query } = key

      query['pagination.limit']=perPage;
      query['pagination.offset']= (pageParam-1)*perPage;
      query['pagination.count_total']= true;
      return  client.ZrchainZentp.query.queryBurns(id, denom, query ?? undefined).then( res => ({...res.data,pageParam}) );
    }, {...options,
      getNextPageParam: (lastPage, allPages) => { if ((lastPage.pagination?.total ?? 0) >((lastPage.pageParam ?? 0) * perPage)) {return lastPage.pageParam+1 } else {return undefined}},
      getPreviousPageParam: (firstPage, allPages) => { if (firstPage.pageParam==1) { return undefined } else { return firstPage.pageParam-1}}
    }
    );
  }
  
  const QueryStats = (address: string, query: any, options: any) => {
    const key = { type: 'QueryStats',  address, query };    
    return useQuery([key], () => {
      const { address,query } = key
      return  client.ZrchainZentp.query.queryStats(address, query ?? undefined).then( res => res.data );
    }, options);
  }
  
  const QueryQuerySolanaROCKSupply = ( options: any) => {
    const key = { type: 'QueryQuerySolanaROCKSupply',  };    
    return useQuery([key], () => {
      return  client.ZrchainZentp.query.queryQuerySolanaRocksupply().then( res => res.data );
    }, options);
  }
  
  return {QueryParams,QueryMints,QueryBurns,QueryStats,QueryQuerySolanaROCKSupply,
  }
}
