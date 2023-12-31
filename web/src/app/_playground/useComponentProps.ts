import { PlayGroundComponentProps } from './Component';
import { NetworkGraphLoaderProps, useLoadGraphEffect } from './useLoadGraphEffect';
import { RequestParamsFromUI, useNetworkState } from './useNetworkState';

export const useComponentProps = (): PlayGroundComponentProps => {
  const networkState = useNetworkState();
  const fmtDate = (d: Date) => {
    const mm = `0${d.getMonth() + 1}`;
    const dd = `0${d.getDate()}`;
    return `${d.getFullYear()}-${mm.substring(mm.length - 2)}-${dd.substring(dd.length - 2)}`;
  };
  const clickNode: NetworkGraphLoaderProps['clickNode'] = (payload) => {
    payload.preventSigmaDefault();
    networkState.startLoading();
    let forcusID = 0;
    try {
      forcusID = Number(payload.node);
    } catch (e) {
      console.error(e);
    }
    if (forcusID) {
      networkState.continueRequest(forcusID);
    }
  };
  const updateGraph: NetworkGraphLoaderProps['updateGraph'] = (graph) => {
    if (networkState.progress < 1) return false;
    const asset = networkState.getNetworkAt('all');
    for (const node of asset.nodes) {
      if (graph.hasNode(node.nodeId)) continue;
      graph.addNode(node.nodeId, {
        label: node.word,
        size: node.rate * 10,
        x: Math.random() * 100,
        y: Math.random() * 100,
      });
    }
    for (const edge of asset.edges) {
      if (graph.hasEdge(edge.edgeId)) continue;
      graph.addEdgeWithKey(edge.edgeId, edge.nodeId1, edge.nodeId2, {
        size: 1,
      });
    }
    return true;
  };
  const loaderProps: PlayGroundComponentProps['loaderProps'] = {
    useLoadingGraphEffect: useLoadGraphEffect.bind(this, {
      clickNode,
      updateGraph,
    }),
  };
  const props: PlayGroundComponentProps = {
    formProps: {
      onSubmit: (event) => {
        event.preventDefault();
        networkState.startLoading();

        const form = new FormData(event.currentTarget);
        const start = form.get('from');
        const end = form.get('until');
        const keyword = form.get('keyword');
        if (!start || !end || !keyword) {
          networkState.stopLoading();
          return;
        }
        const from = new Date(start.toString());
        const until = new Date(end.toString());
        if (from.getTime() > until.getTime()) {
          networkState.stopLoading();
          return;
        }
        const checkedPoSpeechTypes: number[] = [];
        for (const checkName of ['noun', 'personName', 'placeName', 'number', 'adjective', 'adjectiveVerb', 'verb']) {
          const value = Number(form.get(checkName) || 0);
          value && checkedPoSpeechTypes.push(value);
        }
        if (!checkedPoSpeechTypes.length) {
          networkState.stopLoading();
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
          mode: Number(form.get('mode')?.toString() || 0),
        };
        networkState.newRequest(params).then((res) => {
          if (res instanceof Error) {
            window.alert(res.message);
          }
        });
      },
    },
    defaultValues: {
      from: fmtDate(networkState.initRequestParams.from),
      until: fmtDate(networkState.initRequestParams.until),
      keyword: networkState.initRequestParams.keyword,
    },
    progressBarProps: {
      progress: networkState.progress,
    },
    loaderProps: loaderProps,
    loading: networkState.loading,
  };

  return props;
};
