import { NetworkStreamRequest } from '@/api/v3/collo_pb';
import { use, useEffect, useState } from 'react';

export const useSelectNodeState = (params: NetworkStreamRequest) => {
  const [selecedNodeIDs, setState] = useState<Map<number, true>>(new Map());
  useEffect(() => {
    setState(new Map());
  }, [params]);

  return {
    ids: Array.from(selecedNodeIDs.keys()),
    add: (id: number) => setState((pv) => new Map<number, true>(pv.set(id, true))),
    remove: (id: number) =>
      setState((pv) => {
        pv.delete(id);
        return new Map(pv);
      }),
    isSelected: (id: number) => selecedNodeIDs.get(id) || false,
  };
};
