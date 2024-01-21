import { useCallback, useEffect, useMemo } from 'react';
import { PlayGroundComponentProps } from './Component';
import { useLoadGraphEffect } from './useLoadGraphEffect';
import { RequestParamsFromUI, useNetworkState } from './useNetworkState';
import { useSubNetworkState } from './useSubNetworkState';
import { clearLoaderPropsMemo, getLoaderProps } from './useSubLoaderPropsMemo';
import { getChratProps } from './getChartProps';
import { useLoadingState } from './useLoadingState';
import { useSelectNodeState } from './useSelectNodeState';
import { NetworkHandle } from './connect';
import { NetworkStreamResponse, NodeRateStreamResponse } from '@/api/v3/collo_pb';
import { ConnectError } from '@connectrpc/connect';
import { useRankGetterState } from './useRankGetterState';

export const useComponentProps = (): PlayGroundComponentProps => {
  const { getNetworkAt, ...networkState } = useNetworkState();
  const { progress, loading, startLoading, stopLoading, ...stream } = useLoadingState(); // データ取得の進捗
  const subnetworkState = useSubNetworkState();
  const selectedNodes = useSelectNodeState(networkState.requestParms);
  const rankGetterState = useRankGetterState(networkState.requestParms);
  const fmtDate = (d: Date) => {
    const mm = `0${d.getMonth() + 1}`;
    const dd = `0${d.getDate()}`;
    return `${d.getFullYear()}-${mm.substring(mm.length - 2)}-${dd.substring(dd.length - 2)}`;
  };
  // APIからのレスポンス受付時起動するfunctions
  const nwHandle: NetworkHandle = {
    start: function (): void {},
    stream: function (resp): void {
      stream.setProcess(resp.process);
      if (resp.process < 1) return;
      if (resp instanceof NetworkStreamResponse) {
        if (!resp.nodes.length && !resp.edges.length) return;
        console.log(`Get ${resp.nodes.length}_Nodes, ${resp.edges.length}_Edges. At ${resp.meta?.groupId}`);
      } else {
        if (!resp.nodes.length) return;
        console.log(`Get ${resp.nodes.length}_Nodes. At ${resp.meta?.groupId}`);
      }
      networkState.addNetwork(resp);
    },
    end: function (): void {
      networkState.sortNetworkState();
      stream.endStreaming();
    },
    err: function (e) {
      stopLoading();
      console.error(e);
      if (e instanceof ConnectError) {
        return Error(e.rawMessage);
      }
      if (e instanceof Error) {
        return e;
      }
      return Error('予期せぬエラーが発生しました。');
    },
  };

  const getNetwrokLoaderProps: (id: string) => PlayGroundComponentProps['networkProps']['loaderProps'] = useCallback(
    (id: string) => {
      return {
        useLoadingGraphEffect: useLoadGraphEffect.bind(this, {
          asset: getNetworkAt(id),
          progress: progress,
          continueRequest: (forcusNodeID: number) => {
            return networkState.continueRequest(forcusNodeID, 'total', nwHandle);
          },
          startLoading: startLoading,
        }),
      };
    },
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [getNetworkAt],
  );
  useEffect(() => {
    clearLoaderPropsMemo();
  }, [getNetwrokLoaderProps]);

  const networkProps: PlayGroundComponentProps['networkProps'] = useMemo(() => {
    return {
      loaderProps: getNetwrokLoaderProps('total'),
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [getNetwrokLoaderProps]);

  const groupOptions: PlayGroundComponentProps['subNetworksProps'][number]['selectProps']['groupOptionProps'] = [];
  for (const [groupID, { meta }] of networkState.entries()) {
    if (!meta || meta.metas.length < 1) {
      continue;
    }
    if (meta.groupId === 'total') {
      groupOptions.push({
        value: groupID,
        children: `すべての期間`,
      });
      continue;
    }
    const name = meta.metas[0].name;
    const date = meta.metas[0].at?.toDate();
    groupOptions.push({
      value: groupID,
      children: `【${groupID}】${name} ${date ? fmtDate(date) : ''} ${meta.metas.length > 1 ? 'ほか' : ''}`,
    });
  }
  groupOptions.sort((a, b) => {
    const aa = a.children.toString();
    const bb = b.children.toString();
    if (aa > bb) {
      return 1;
    } else if (aa < bb) {
      return -1;
    }
    return 0;
  });

  const subNetworksProps: PlayGroundComponentProps['subNetworksProps'] = subnetworkState.groupIDs.map((id, at) => {
    const metas = getNetworkAt(id)?.meta?.metas || [];
    const metaMap = new Map<string, null>();
    const metaProps: PlayGroundComponentProps['subNetworksProps'][number]['contentsProps']['metaProps'] = [];
    for (const meta of metas) {
      if (typeof metaMap.get(meta.key) != 'undefined') {
        continue;
      }
      metaMap.set(meta.key, null);
      const date = meta.at?.toDate();
      metaProps.push({
        href: `https://kokkai.ndl.go.jp/#/detail?minId=${meta.key}&current=1`,
        children: `${date ? fmtDate(date) : ''} ${meta.name}`,
      });
    }
    console.log('subnetwork id', id, '.');
    const props: PlayGroundComponentProps['subNetworksProps'][number] = {
      contentsProps: {
        metaProps,
        loading: loading,
        top3Button: {
          disabled: !id || networkState.inRequestHisotries(0, id),
          onClick: () => {
            networkState.continueRequest(0, id, nwHandle).then((res) => {
              if (res instanceof Error) {
                window.alert(res.message);
              }
            });
          },
        },
      },
      deleteButtonProps: {
        onClick: () => {
          subnetworkState.removeAt(at);
        },
      },
      loaderProps: getLoaderProps(id, getNetwrokLoaderProps),
      selectProps: {
        groupSelectProps: {
          onChange: (e) => {
            const groupID = e.currentTarget.value;
            subnetworkState.changeID(at, groupID);
          },
        },
        groupOptionProps: groupOptions,
      },
    };
    return props;
  });

  const selectedNodeProps: PlayGroundComponentProps['selectedNodeProps'] = useMemo(
    () =>
      networkState.getTotalNodes(rankGetterState.dones).map((node) => {
        return {
          labelProps: {
            children: node.word,
          },
          checkboxProps: {
            checked: selectedNodes.isSelected(node.nodeId),
            onChange: (e) => {
              if (e.target.checked) {
                selectedNodes.add(node.nodeId);
              } else {
                selectedNodes.remove(node.nodeId);
              }
            },
          },
        };
      }),
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [networkState.getTotalNodes, selectedNodes, rankGetterState.dones],
  );

  const updateNodeRankProps: PlayGroundComponentProps['updateNodeRankProps'] = useMemo(
    () => ({
      buttonProps: {
        // すべて取り切ったあるいは、ローディング中は使用不可
        disabled: !rankGetterState.hasNext || loading,
        onClick: () => {
          if (!rankGetterState.hasNext || loading) return;
          startLoading();
          rankGetterState.request(networkState.requestParms.config, nwHandle).then((resp) => {
            if (resp instanceof Error) {
              window.alert(resp.message);
            }
          });
        },
      },
      spinner: {
        animate: loading,
      },
    }),

    // eslint-disable-next-line react-hooks/exhaustive-deps
    [loading, rankGetterState, networkState.requestParms],
  );

  const chartProps = useMemo(
    () => getChratProps({ getNetworkAt, ...networkState }, selectedNodes.ids),
    [networkState, getNetworkAt, selectedNodes],
  );
  const props: PlayGroundComponentProps = {
    formProps: {
      onSubmit: (event) => {
        event.preventDefault();
        startLoading();

        const form = new FormData(event.currentTarget);
        const start = form.get('from');
        const end = form.get('until');
        const keyword = form.get('keyword');
        if (!start || !end || !keyword) {
          stopLoading();
          return;
        }
        const from = new Date(start.toString());
        const until = new Date(end.toString());
        if (from.getTime() > until.getTime()) {
          stopLoading();
          return;
        }
        const checkedPoSpeechTypes: number[] = [];
        for (const checkName of ['noun', 'personName', 'placeName', 'number', 'adjective', 'adjectiveVerb', 'verb']) {
          const value = Number(form.get(checkName) || 0);
          value && checkedPoSpeechTypes.push(value);
        }
        if (!checkedPoSpeechTypes.length) {
          stopLoading();
          return;
        }
        const stopwords = form.get('stopwords')?.toString().trim().split(/\s+/) || [];
        const params: RequestParamsFromUI = {
          from,
          until,
          keyword: keyword.toString(),
          forcusNodeID: 0,
          poSpeechType: checkedPoSpeechTypes,
          stopwords,
          forcusGroupID: '',
          apiType: form.get('api')?.toString() === '2' ? 2 : 1,
          pickGroupType: form.get('pick')?.toString() === '2' ? 2 : 1,
        };
        // subnetwork reset onClick "submit" botton
        subnetworkState.clearSubnetwork();
        networkState.newRequest(params, nwHandle).then((res) => {
          if (res instanceof Error) {
            window.alert(res.message);
          }
        });
      },
    },
    progressBarProps: {
      progress: progress,
    },
    networkProps: networkProps,
    subNetworksProps: subNetworksProps,
    loading: loading,
    selectedNodeProps: selectedNodeProps,
    updateNodeRankProps: updateNodeRankProps,
    appendNetworkButtonProps: {
      onClick: () => {
        subnetworkState.appendGroupID('');
      },
    },
    chartProps: chartProps,
  };
  return props;
};
