import { ColloWebService } from '@/api/v2/collo_connect';
import { ColloWebStreamResponse } from '@/api/v2/collo_pb';
import { createPromiseClient } from '@connectrpc/connect';
import { createConnectTransport } from '@connectrpc/connect-web';
import { useState } from 'react';
import { useRequestState } from './useRequestState';

const transport = createConnectTransport({
  baseUrl: process.env.NEXT_PUBLIC_GRPC_HOST || '',
});

export const useNetworkState = () => {
  // networkデータ
  const [network, setNetwork] = useState<Pick<ColloWebStreamResponse, 'nodes' | 'edges'>>({ nodes: [], edges: [] });
  // データ取得の進捗
  const [progress, setProgress] = useState(0);
  // リクエストフォーム
  const { createRequest, ...requestState } = useRequestState();

  const loading = progress != 0 && progress < 1;
  const startLoading = () => setProgress(0.05);
  const stopLoading = () => setProgress(0);
  // データ取得
  const request = async () => {
    const client = createPromiseClient(ColloWebService, transport);
    const stream = client.colloWebStream(createRequest());
    for await (const m of stream) {
      if (m.needs > m.dones) {
        // データ分析中
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
  };
  return {
    network,
    progress,
    loading,
    startLoading,
    stopLoading,
    request,
    ...requestState,
  };
};
