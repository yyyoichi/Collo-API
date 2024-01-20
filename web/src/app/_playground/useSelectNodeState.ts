import { useState } from 'react';

export const useSelectNodeState = () => {
  const [selecedNodeIDs, setState] = useState<Map<number, true>>(new Map());

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
