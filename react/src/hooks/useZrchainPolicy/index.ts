/* eslint-disable @typescript-eslint/no-unused-vars */
import { useQuery, type UseQueryOptions, useInfiniteQuery, type UseInfiniteQueryOptions } from "@tanstack/react-query";
import { useClient } from '../useClient';

export default function useZrchainPolicy() {
  const client = useClient();
  const QueryParams = ( options: any) => {
    const key = { type: 'QueryParams',  };    
    return useQuery([key], () => {
      return  client.ZrchainPolicy.query.queryParams().then( res => res.data );
    }, options);
  }
  
  const QueryActions = (query: any, options: any, perPage: number) => {
    const key = { type: 'QueryActions', query };    
    return useInfiniteQuery([key], ({pageParam = 1}: { pageParam?: number}) => {
      const {query } = key

      query['pagination.limit']=perPage;
      query['pagination.offset']= (pageParam-1)*perPage;
      query['pagination.count_total']= true;
      return  client.ZrchainPolicy.query.queryActions(query ?? undefined).then( res => ({...res.data,pageParam}) );
    }, {...options,
      getNextPageParam: (lastPage, allPages) => { if ((lastPage.pagination?.total ?? 0) >((lastPage.pageParam ?? 0) * perPage)) {return lastPage.pageParam+1 } else {return undefined}},
      getPreviousPageParam: (firstPage, allPages) => { if (firstPage.pageParam==1) { return undefined } else { return firstPage.pageParam-1}}
    }
    );
  }
  
  const QueryPolicies = (query: any, options: any, perPage: number) => {
    const key = { type: 'QueryPolicies', query };    
    return useInfiniteQuery([key], ({pageParam = 1}: { pageParam?: number}) => {
      const {query } = key

      query['pagination.limit']=perPage;
      query['pagination.offset']= (pageParam-1)*perPage;
      query['pagination.count_total']= true;
      return  client.ZrchainPolicy.query.queryPolicies(query ?? undefined).then( res => ({...res.data,pageParam}) );
    }, {...options,
      getNextPageParam: (lastPage, allPages) => { if ((lastPage.pagination?.total ?? 0) >((lastPage.pageParam ?? 0) * perPage)) {return lastPage.pageParam+1 } else {return undefined}},
      getPreviousPageParam: (firstPage, allPages) => { if (firstPage.pageParam==1) { return undefined } else { return firstPage.pageParam-1}}
    }
    );
  }
  
  const QueryPolicyById = (id: string,  options: any) => {
    const key = { type: 'QueryPolicyById',  id };    
    return useQuery([key], () => {
      const { id } = key
      return  client.ZrchainPolicy.query.queryPolicyById(id).then( res => res.data );
    }, options);
  }
  
  const QuerySignMethodsByAddress = (address: string, query: any, options: any, perPage: number) => {
    const key = { type: 'QuerySignMethodsByAddress',  address, query };    
    return useInfiniteQuery([key], ({pageParam = 1}: { pageParam?: number}) => {
      const { address,query } = key

      query['pagination.limit']=perPage;
      query['pagination.offset']= (pageParam-1)*perPage;
      query['pagination.count_total']= true;
      return  client.ZrchainPolicy.query.querySignMethodsByAddress(address, query ?? undefined).then( res => ({...res.data,pageParam}) );
    }, {...options,
      getNextPageParam: (lastPage, allPages) => { if ((lastPage.pagination?.total ?? 0) >((lastPage.pageParam ?? 0) * perPage)) {return lastPage.pageParam+1 } else {return undefined}},
      getPreviousPageParam: (firstPage, allPages) => { if (firstPage.pageParam==1) { return undefined } else { return firstPage.pageParam-1}}
    }
    );
  }
  
  const QueryPoliciesByCreator = (creators: string, query: any, options: any, perPage: number) => {
    const key = { type: 'QueryPoliciesByCreator',  creators, query };    
    return useInfiniteQuery([key], ({pageParam = 1}: { pageParam?: number}) => {
      const { creators,query } = key

      query['pagination.limit']=perPage;
      query['pagination.offset']= (pageParam-1)*perPage;
      query['pagination.count_total']= true;
      return  client.ZrchainPolicy.query.queryPoliciesByCreator(creators, query ?? undefined).then( res => ({...res.data,pageParam}) );
    }, {...options,
      getNextPageParam: (lastPage, allPages) => { if ((lastPage.pagination?.total ?? 0) >((lastPage.pageParam ?? 0) * perPage)) {return lastPage.pageParam+1 } else {return undefined}},
      getPreviousPageParam: (firstPage, allPages) => { if (firstPage.pageParam==1) { return undefined } else { return firstPage.pageParam-1}}
    }
    );
  }
  
  const QueryActionDetailsById = (id: string,  options: any) => {
    const key = { type: 'QueryActionDetailsById',  id };    
    return useQuery([key], () => {
      const { id } = key
      return  client.ZrchainPolicy.query.queryActionDetailsById(id).then( res => res.data );
    }, options);
  }
  
  return {QueryParams,QueryActions,QueryPolicies,QueryPolicyById,QuerySignMethodsByAddress,QueryPoliciesByCreator,QueryActionDetailsById,
  }
}
