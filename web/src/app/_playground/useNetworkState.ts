import { MintGreenService } from '@/api/v3/collo_connect';
import {
  NetworkStreamRequest,
  NetworkStreamResponse,
  NodeRateStreamRequest,
  NodeRateStreamResponse,
  RequestConfig,
  RequestConfig_NdlApiType,
  RequestConfig_PickGroupType,
} from '@/api/v3/collo_pb';
import { ConnectError, createPromiseClient } from '@connectrpc/connect';
import { createConnectTransport } from '@connectrpc/connect-web';
import { useCallback, useState } from 'react';
import { Timestamp } from '@bufbuild/protobuf';
import { useLoadingState } from './useLoadingState';
import { useReqHistoryState } from './useReqHistoryState';

export type RequestParamsFromUI = {
  from: Date;
  until: Date;
  keyword: string;
  forcusNodeID: number;
  forcusGroupID: string;
  poSpeechType: number[];
  stopwords: string[];
  apiType: RequestConfig_NdlApiType;
  pickGroupType: RequestConfig_PickGroupType;
};

export type NetworkState = Map<string, Pick<NetworkStreamResponse, 'nodes' | 'edges' | 'meta'>>;

const transport = createConnectTransport({
  baseUrl: process.env.NEXT_PUBLIC_RPC_HOST || '',
});

export const useNetworkState = () => {
  // networkデータ
  const [network, setNetwork] = useState<NetworkState>(new Map());
  const [requestParms, setRequestParams] = useState<NetworkStreamRequest>(new NetworkStreamRequest());
  const requestHistories = useReqHistoryState();
  // データ取得の進捗
  const { progress, setProcess, endStreaming, loading, startLoading, stopLoading } = useLoadingState();

  // データ取得
  const request = async (req: NetworkStreamRequest) => {
    setRequestParams(req);
    requestHistories.addHisotry(req);
    const client = createPromiseClient(MintGreenService, transport);
    try {
      const stream = client.networkStream(req);
      console.log(
        `Start request.. Keyword:${req.config?.keyword},`,
        `From:${req.config?.from?.toJsonString()},`,
        `Until:${req.config?.until?.toJsonString()},`,
        `ForcusNodeID:${req.forcusNodeId},`,
        `PartOfSpeechTypes:${req.config?.partOfSpeechTypes},`,
        `Stopwords:${req.config?.stopwords},`,
        `Api:${req.config?.ndlApiType}, Pick:${req.config?.pickGroupType},`,
      );
      for await (const m of stream) {
        setProcess(m.process);
        if (m.process < 1) {
          continue;
        }
        console.log(`Get ${m.nodes.length}_Nodes, ${m.edges.length}_Edges. At ${m.meta?.groupId}`);
        // データ追加
        setNetwork((pns) => {
          const key = m.meta?.groupId || 'all';
          const pn = pns.get(key) || {
            nodes: [],
            edges: [],
            meta: m.meta,
          };
          pn.nodes = pn.nodes.concat(m.nodes);
          pn.edges = pn.edges.concat(m.edges);
          return new Map(pns.set(key, pn));
        });
      }
      endStreaming();
    } catch (e) {
      console.error(e);
      stopLoading();
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
    setNetwork(new Map()); // 取得結果リセット
    requestHistories.clearHistories(); // 取得履歴リセット
    const config = new RequestConfig();
    config.from = Timestamp.fromDate(param.from);
    config.until = Timestamp.fromDate(param.until);
    config.keyword = param.keyword;
    config.partOfSpeechTypes = param.poSpeechType;
    config.stopwords = param.stopwords;
    config.ndlApiType = param.apiType;
    config.pickGroupType = param.pickGroupType;
    const req = new NetworkStreamRequest();
    req.config = config;
    req.forcusNodeId = 0;
    return request(req);
  };

  /** ForcusNodeIDとForcusGroupIDを現在のリクエストに追加する */
  const continueRequest = (
    forcusNodeID: RequestParamsFromUI['forcusNodeID'],
    forcusGroupID: RequestParamsFromUI['forcusGroupID'],
  ) => {
    const req = requestParms.clone();
    req.forcusNodeId = forcusNodeID;
    if (!req.config) {
      req.config = new RequestConfig();
    }
    req.config.forcusGroupId = forcusGroupID;
    return request(req);
  };
  const getNetworkAt = useCallback(
    (groupID: string) => {
      const asset = network.get(groupID);
      return {
        nodes: asset?.nodes || [],
        edges: asset?.edges || [],
        meta: asset?.meta,
      };
    },
    [network],
  );

  const entries = useCallback(
    function* () {
      for (const a of network.entries()) {
        yield a;
      }
    },
    [network],
  );
  return {
    entries,
    getNetworkAt,
    progress,
    loading,
    startLoading,
    stopLoading,
    newRequest,
    continueRequest,
    inRequestHisotries: requestHistories.inHistories,
  };
};
