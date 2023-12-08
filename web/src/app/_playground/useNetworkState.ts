import { ColloWebService } from '@/api/v2/collo_connect';
import { ColloWebStreamRequest, ColloWebStreamResponse } from '@/api/v2/collo_pb';
import { createPromiseClient } from '@connectrpc/connect';
import { createConnectTransport } from '@connectrpc/connect-web';
import { useState } from 'react';
import { Timestamp } from '@bufbuild/protobuf';

type RequestParams = ReturnType<typeof getInitRequestParams>;

const transport = createConnectTransport({
  baseUrl: process.env.NEXT_PUBLIC_GRPC_HOST || '',
});

export const useNetworkState = () => {
  // networkデータ
  const [network, setNetwork] = useState<Pick<ColloWebStreamResponse, 'nodes' | 'edges'>>({ nodes: [], edges: [] });
  // データ取得の進捗
  const [progress, setProgress] = useState(0);

  const initRequestParams = getInitRequestParams();

  const loading = progress != 0 && progress < 1;
  const startLoading = () => setProgress(0.05);
  const stopLoading = () => setProgress(0);

  // データ取得
  const request = async (req: ColloWebStreamRequest) => {
    const client = createPromiseClient(ColloWebService, transport);
    const stream = client.colloWebStream(req);
    try {
      for await (const m of stream) {
        if (m.needs > m.dones) {
          // データ分析中
          console.log(m.dones / m.needs);
          setProgress(m.dones / m.needs);
          continue;
        }
        console.log('get: ', m.nodes.length);
        // データ追加
        setNetwork((pn) => ({
          nodes: pn.nodes.concat(m.nodes),
          edges: pn.edges.concat(m.edges),
        }));
        // 完了
        setProgress(1);
      }
    } catch (e) {
      console.error(e);
      setProgress(0);
    }
  };
  /** 引数のパラメータにリセットする */
  const newRequest = (
    from: RequestParams['from'],
    until: RequestParams['until'],
    keyword: RequestParams['keyword'],
  ) => {
    setNetwork({ nodes: [], edges: [] });
    const req = new ColloWebStreamRequest();
    req.from = Timestamp.fromDate(from);
    req.until = Timestamp.fromDate(until);
    req.keyword = keyword;
    req.forcusNodeId = 0;
    request(req);
  };

  /** ForcusNodeIDを現在のリクエストに追加する */
  const continueRequest = (forcusNodeID: RequestParams['forcusNodeID']) => {
    const req = new ColloWebStreamRequest();
    req.forcusNodeId = forcusNodeID;
    request(req);
  };

  return {
    network,
    progress,
    loading,
    startLoading,
    stopLoading,
    newRequest,
    continueRequest,
    initRequestParams,
  };
};

function getInitRequestParams() {
  return {
    from: new Date(2023, 2, 1),
    until: new Date(2023, 3, 30),
    keyword: 'アニメ',
    forcusNodeID: 0,
  };
}
