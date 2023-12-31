import { ColloRateWebService } from '@/api/v2/collo_connect';
import { ColloRateWebStreamRequest, ColloRateWebStreamResponse } from '@/api/v2/collo_pb';
import { ConnectError, createPromiseClient } from '@connectrpc/connect';
import { createConnectTransport } from '@connectrpc/connect-web';
import { useState } from 'react';
import { Timestamp } from '@bufbuild/protobuf';

export type RequestParamsFromUI = {
  from: Date;
  until: Date;
  keyword: string;
  forcusNodeID: number;
  poSpeechType: number[];
  stopwords: string[];
};

const transport = createConnectTransport({
  baseUrl: process.env.NEXT_PUBLIC_RPC_HOST || '',
});

export const useNetworkState = () => {
  // networkデータ
  const [network, setNetwork] = useState<Pick<ColloRateWebStreamResponse, 'nodes' | 'edges'>>({ nodes: [], edges: [] });
  const [requestParms, setRequestParams] = useState<ColloRateWebStreamRequest>(new ColloRateWebStreamRequest());
  // データ取得の進捗
  const [progress, setProgress] = useState(0);

  const initRequestParams = getInitRequestParams();

  const loading = progress != 0 && progress < 1;
  const startLoading = () => setProgress(0.05);
  const stopLoading = () => setProgress(0);

  // データ取得
  const request = async (req: ColloRateWebStreamRequest) => {
    setRequestParams(req);
    const client = createPromiseClient(ColloRateWebService, transport);
    const stream = client.colloRateWebStream(req);
    console.log(
      `Start request.. Keyword:${req.keyword},`,
      `From:${req.from?.toJsonString()},`,
      `Until:${req.until?.toJsonString()},`,
      `ForcusNodeID${req.forcusNodeId}`,
      `PartOfSpeechTypes${req.partOfSpeechTypes}`,
      `Stopwords${req.stopwords}`,
    );
    try {
      for await (const m of stream) {
        if (m.needs > m.dones) {
          // データ分析中
          console.log(m.dones / m.needs);
          if (m.dones > 0) {
            // 進捗があったときのみ更新
            setProgress(m.dones / m.needs);
          }
          continue;
        }
        console.log(`Get ${m.nodes.length}_Nodes, ${m.edges.length}_Edges.`);
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
      if (e instanceof ConnectError) {
        return Error(e.rawMessage);
      }
      if (e instanceof Error) {
        return e;
      }
      return Error('予期せぬエラーが発生しました。');
    }
  };
  /** 引数のパラメータにリセットする */
  const newRequest = (param: RequestParamsFromUI) => {
    setNetwork({ nodes: [], edges: [] });
    const req = new ColloRateWebStreamRequest();
    req.from = Timestamp.fromDate(param.from);
    req.until = Timestamp.fromDate(param.until);
    req.keyword = param.keyword;
    req.forcusNodeId = 0;
    req.partOfSpeechTypes = param.poSpeechType;
    req.stopwords = param.stopwords;
    return request(req);
  };

  /** ForcusNodeIDを現在のリクエストに追加する */
  const continueRequest = (forcusNodeID: RequestParamsFromUI['forcusNodeID']) => {
    const req = requestParms.clone();
    req.forcusNodeId = forcusNodeID;
    return request(req);
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
    from: new Date(2023, 3, 20),
    until: new Date(2023, 3, 30),
    keyword: 'アニメ',
    forcusNodeID: 0,
  };
}
