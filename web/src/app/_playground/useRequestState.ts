import { ColloWebStreamRequest } from '@/api/v2/collo_pb';
import { Timestamp } from '@bufbuild/protobuf';
import { useState } from 'react';

type RequestParams = ReturnType<typeof getInitRequestParams>;

export const useRequestState = () => {
  const [params, setParams] = useState<RequestParams>(getInitRequestParams());

  /**RPCの為のリクエストパラメータを返す */
  const createRequest = () => {
    const req = new ColloWebStreamRequest();
    const u = Timestamp.fromDate(params.until);
    req.from = Timestamp.fromDate(params.from);
    req.until = Timestamp.fromDate(params.until);
    req.keyword = params.keyword;
    req.forcusNodeId = params.forcusNodeID;
    return req;
  };

  /** 引数のパラメータにリセットする */
  const resetRequestParams = (
    f: RequestParams['from'],
    u: RequestParams['until'],
    keyword: RequestParams['keyword'],
  ) => {
    setParams({
      from: f,
      until: u,
      keyword: keyword,
      forcusNodeID: 0,
    });
  };

  /** ForcusNodeIDを現在のリクエストに追加する */
  const setForcusNodeID = (id: RequestParams['forcusNodeID']) => {
    setParams((op) => ({
      ...op,
      forcusNodeID: id,
    }));
  };

  return {
    createRequest,
    resetRequestParams,
    setForcusNodeID,
  };
};

function getInitRequestParams() {
  const f = new Date();
  f.setMonth(f.getMonth() - 6);
  const u = new Date();
  u.setMonth(f.getMonth() - 3);
  return {
    from: f,
    until: u,
    keyword: '自動車',
    forcusNodeID: 0,
  };
}
