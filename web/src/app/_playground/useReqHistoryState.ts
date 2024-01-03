import { useState } from 'react';
import { ColloRateWebStreamRequest } from '@/api/v2/collo_pb';

// 過去のリクエストを保持する。
export const useReqHistoryState = () => {
  const [histories, setHistories] = useState<ColloRateWebStreamRequest[]>([]);
  return {
    inHistories: (forcusNodeID: number, forcusGroupID: string) => {
      for (const history of histories) {
        if (history.forcusGroupId == forcusGroupID && history.forcusNodeId == forcusNodeID) {
          return true;
        }
      }
      return false;
    },
    clearHistories: () => setHistories([]),
    addHisotry: (req: ColloRateWebStreamRequest) => setHistories((phv) => [...phv, req]),
  };
};
