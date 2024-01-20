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
import { useCallback, useMemo, useState } from 'react';
import { Timestamp } from '@bufbuild/protobuf';
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

export type NetworkHandle = {
  start: () => void;
  stream: (process: number) => void;
  end: () => void;
  err: () => void;
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

  // データ取得
  const request = async (req: NetworkStreamRequest, handle: NetworkHandle) => {
    handle.start();
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
        handle.stream(m.process);
        if (m.process < 1) {
          continue;
        }
        console.log(`Get ${m.nodes.length}_Nodes, ${m.edges.length}_Edges. At ${m.meta?.groupId}`);
        // データ追加
        setNetwork((pns) => {
          const key = m.meta?.groupId || 'total';
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
      // map内のnode,edgeをソートしユニークな配列にする。
      setNetwork((pns) => {
        const map: NetworkState = new Map();
        for (const [key, val] of pns.entries()) {
          // sort
          val.nodes.sort(({ rate: arate }, { rate: brate }) => {
            return brate - arate;
          });
          val.edges.sort(({ rate: arate }, { rate: brate }) => {
            return brate - arate;
          });
          const asset: Pick<NetworkStreamResponse, 'nodes' | 'edges' | 'meta'> = {
            nodes: [],
            edges: [],
            meta: val.meta,
          };
          if (val.nodes.length > 0) {
            // uniq
            let reservedNodeID = [val.nodes[0].nodeId];
            asset.nodes.push(val.nodes[0]);
            for (const node of val.nodes) {
              const id = node.nodeId;
              if (reservedNodeID.includes(id)) {
                continue;
              }
              reservedNodeID.push(id);
              asset.nodes.push(node);
            }
          }
          if (val.edges.length > 0) {
            // uniq
            let reservedEdgeID = [val.edges[0].edgeId];
            asset.edges.push(val.edges[0]);
            for (const edge of val.edges) {
              const id = edge.edgeId;
              if (reservedEdgeID.includes(id)) {
                continue;
              }
              reservedEdgeID.push(id);
              asset.edges.push(edge);
            }
          }
          map.set(key, asset);
        }
        return map;
      });
      handle.end();
    } catch (e) {
      console.error(e);
      handle.err();
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
  const newRequest = (param: RequestParamsFromUI, handle: NetworkHandle) => {
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
    config.useNdlCache = true;
    config.createNdlCache = true;
    const req = new NetworkStreamRequest();
    req.config = config;
    req.forcusNodeId = 0;
    return request(req, handle);
  };

  /** ForcusNodeIDとForcusGroupIDを現在のリクエストに追加する */
  const continueRequest = (
    forcusNodeID: RequestParamsFromUI['forcusNodeID'],
    forcusGroupID: RequestParamsFromUI['forcusGroupID'],
    handle: NetworkHandle,
  ) => {
    const req = requestParms.clone();
    req.forcusNodeId = forcusNodeID;
    if (!req.config) {
      req.config = new RequestConfig();
    }
    req.config.forcusGroupId = forcusGroupID;
    return request(req, handle);
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

  // totalの中心性Top3のNodeIDを返す。
  const getTop3NodeIDInTotal = useCallback(() => {
    const totalasset = network.get('total');
    if (!totalasset) return [];
    if (totalasset.nodes.length < 1) return [];
    const ids = [];
    for (let i = 0; i < 3; i++) {
      ids.push(totalasset.nodes[i].nodeId);
    }
    return ids;
  }, [network]);

  // == network.entries()
  const entries = useCallback(
    function* () {
      for (const a of network.entries()) {
        yield a;
      }
    },
    [network],
  );

  // Keyごとにジェネレーターで返す keyIndexはsortされたキーの位置、nodeIndexは渡されたidのうち一致した位置。
  const getWords = useCallback(
    function* (nodeIDs: number[]) {
      let keyIndex = 0;
      for (const [key, { nodes, edges, meta }] of entries()) {
        for (const node of nodes) {
          const nodeIndex = nodeIDs.indexOf(node.nodeId);
          if (nodeIndex === -1) continue;
          yield {
            nodeIndex,
            keyIndex,
            key,
            node,
          };
        }
        keyIndex++;
      }
    },
    [entries],
  );

  const numKeys = useMemo(() => {
    return network.size;
  }, [network]);

  // 日付でソートしたGroupIDsを返す。
  const sortedGroupID = useMemo(() => {
    const keys: string[] = [];
    const times: number[] = [];
    for (const [key, val] of network.entries()) {
      const time = (val.meta?.from?.toDate() || new Date(1970, 0, 1)).getTime();
      if (keys.length === 0) {
        keys.push(key);
        times.push(time);
        continue;
      }
      let added = false;
      for (let i = 0; i < keys.length; i++) {
        if (times[i] <= time) {
          continue;
        }
        keys.splice(i, 0, key);
        times.splice(i, 0, time);
        added = true;
        break;
      }
      if (!added) {
        keys.push(key);
        times.push(time);
      }
    }
    return keys;
  }, [network]);

  const getTotalNodes = useCallback(() => {
    const asset = network.get('total');
    if (!asset || !asset.nodes.length) return [];
    return asset.nodes;
  }, [network]);
  return {
    // network
    entries,
    getNetworkAt,
    sortedGroupID,
    getTop3NodeIDInTotal,
    getWords,
    numKeys,
    getTotalNodes,

    // request
    newRequest,
    continueRequest,
    pickType: requestParms.config?.pickGroupType,

    // histroy
    inRequestHisotries: requestHistories.inHistories,
  };
};
