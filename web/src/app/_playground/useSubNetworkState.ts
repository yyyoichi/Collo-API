import { useState } from 'react';

export const useSubNetworkState = () => {
  const [groupIDs, setGroupIDs] = useState<string[]>([]);
  const removeAt = (at: number) => {
    setGroupIDs((pg) => {
      pg.splice(at, 1);
      return [...pg];
    });
  };
  const changeID = (at: number, id: string) => {
    setGroupIDs((pg) => {
      pg[at] = id;
      return [...pg];
    });
  };
  const appendGroupID = (groupID: string) => {
    setGroupIDs((pg) => {
      return [...pg, groupID];
    });
  };
  const generateGroupID = function* () {
    yield* groupIDs.map((id, at) => ({ id, at }));
  };
  return {
    groupIDs,
    changeID,
    removeAt,
    appendGroupID,
    generateGroupID,
  };
};
