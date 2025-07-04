/* eslint-disable @typescript-eslint/no-unused-vars */
import { useQuery, type UseQueryOptions, useInfiniteQuery, type UseInfiniteQueryOptions } from "@tanstack/react-query";
import { useClient } from '../useClient';

export default function useZrchainIdentity() {
  const client = useClient();
  const QueryParams = ( options: any) => {
    const key = { type: 'QueryParams',  };    
    return useQuery([key], () => {
      return  client.ZrchainIdentity.query.queryParams().then( res => res.data );
    }, options);
  }
  
  const QueryWorkspaces = (query: any, options: any, perPage: number) => {
    const key = { type: 'QueryWorkspaces', query };    
    return useInfiniteQuery([key], ({pageParam = 1}: { pageParam?: number}) => {
      const {query } = key

      query['pagination.limit']=perPage;
      query['pagination.offset']= (pageParam-1)*perPage;
      query['pagination.count_total']= true;
      return  client.ZrchainIdentity.query.queryWorkspaces(query ?? undefined).then( res => ({...res.data,pageParam}) );
    }, {...options,
      getNextPageParam: (lastPage, allPages) => { if ((lastPage.pagination?.total ?? 0) >((lastPage.pageParam ?? 0) * perPage)) {return lastPage.pageParam+1 } else {return undefined}},
      getPreviousPageParam: (firstPage, allPages) => { if (firstPage.pageParam==1) { return undefined } else { return firstPage.pageParam-1}}
    }
    );
  }
  
  const QueryWorkspaceByAddress = (workspace_addr: string,  options: any) => {
    const key = { type: 'QueryWorkspaceByAddress',  workspace_addr };    
    return useQuery([key], () => {
      const { workspace_addr } = key
      return  client.ZrchainIdentity.query.queryWorkspaceByAddress(workspace_addr).then( res => res.data );
    }, options);
  }
  
  const QueryKeyrings = (query: any, options: any, perPage: number) => {
    const key = { type: 'QueryKeyrings', query };    
    return useInfiniteQuery([key], ({pageParam = 1}: { pageParam?: number}) => {
      const {query } = key

      query['pagination.limit']=perPage;
      query['pagination.offset']= (pageParam-1)*perPage;
      query['pagination.count_total']= true;
      return  client.ZrchainIdentity.query.queryKeyrings(query ?? undefined).then( res => ({...res.data,pageParam}) );
    }, {...options,
      getNextPageParam: (lastPage, allPages) => { if ((lastPage.pagination?.total ?? 0) >((lastPage.pageParam ?? 0) * perPage)) {return lastPage.pageParam+1 } else {return undefined}},
      getPreviousPageParam: (firstPage, allPages) => { if (firstPage.pageParam==1) { return undefined } else { return firstPage.pageParam-1}}
    }
    );
  }
  
  const QueryKeyringByAddress = (keyring_addr: string,  options: any) => {
    const key = { type: 'QueryKeyringByAddress',  keyring_addr };    
    return useQuery([key], () => {
      const { keyring_addr } = key
      return  client.ZrchainIdentity.query.queryKeyringByAddress(keyring_addr).then( res => res.data );
    }, options);
  }
  
  return {QueryParams,QueryWorkspaces,QueryWorkspaceByAddress,QueryKeyrings,QueryKeyringByAddress,
  }
}
