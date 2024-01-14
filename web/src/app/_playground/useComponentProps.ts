import { useCallback, useEffect, useMemo } from 'react';
import { PlayGroundComponentProps } from './Component';
import { useLoadGraphEffect } from './useLoadGraphEffect';
import { RequestParamsFromUI, useNetworkState } from './useNetworkState';
import { useSubNetworkState } from './useSubNetworkState';
import { clearLoaderPropsMemo, getLoaderProps } from './useSubLoaderPropsMemo';

export const useComponentProps = (): PlayGroundComponentProps => {
  const { getNetworkAt, ...networkState } = useNetworkState();
  const subnetworkState = useSubNetworkState();
  const fmtDate = (d: Date) => {
    const mm = `0${d.getMonth() + 1}`;
    const dd = `0${d.getDate()}`;
    return `${d.getFullYear()}-${mm.substring(mm.length - 2)}-${dd.substring(dd.length - 2)}`;
  };

  const getNetwrokLoaderProps: (id: string) => PlayGroundComponentProps['networkProps']['loaderProps'] = useCallback(
    (id: string) => {
      return {
        useLoadingGraphEffect: useLoadGraphEffect.bind(this, {
          asset: getNetworkAt(id),
          progress: networkState.progress,
          continueRequest: (forcusNodeID: number) => {
            return networkState.continueRequest(forcusNodeID, 'total');
          },
          startLoading: networkState.startLoading,
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
        loading: networkState.loading,
        top3Button: {
          disabled: !id || networkState.inRequestHisotries(0, id),
          onClick: () => {
            networkState.continueRequest(0, id).then((res) => {
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
          forcusGroupID: '',
          apiType: form.get('api')?.toString() === '2' ? 2 : 1,
          pickGroupType: form.get('pick')?.toString() === '2' ? 2 : 1,
        };
        // subnetwork reset onClick "submit" botton
        subnetworkState.clearSubnetwork();
        networkState.newRequest(params).then((res) => {
          if (res instanceof Error) {
            window.alert(res.message);
          }
        });
      },
    },
    progressBarProps: {
      progress: networkState.progress,
    },
    networkProps: networkProps,
    subNetworksProps: subNetworksProps,
    loading: networkState.loading,
    appendNetworkButtonProps: {
      onClick: () => {
        subnetworkState.appendGroupID('');
      },
    },
    chartProps: {
      series: [
        {
          data: [0.5, 0.3, 0.6, 0.7, 0.4],
          name: 'Word',
        },
      ],
      xaxis: {
        categories: ['2023-01', '2023-02', '2023-03', '2023-04', '2023-05'],
        title: '月',
      },
    },
  };

  return props;
};
