import { RequestConfig_PickGroupType } from '@/api/v3/collo_pb';
import { CenteralityChartProps } from './Chart';
import { useNetworkState } from './useNetworkState';

export const getChratProps = (networkState: ReturnType<typeof useNetworkState>) => {
  const fmtDate = (d: Date) => {
    const mm = `0${d.getMonth() + 1}`;
    const dd = `0${d.getDate()}`;
    return `${d.getFullYear()}-${mm.substring(mm.length - 2)}-${dd.substring(dd.length - 2)}`;
  };
  // chartProps: {
  //     series: [
  //       {
  //         data: [0.5, 0.3, 0.6, 0.7, 0.4],
  //         name: 'Word',
  //       },
  //     ],
  //     xaxis: {
  //       categories: networkState.sortedGroupID().splice(1),
  //       title: '月',
  //     },
  // },
  const numKeys = networkState.numKeys;
  const isMonth = networkState.pickType === RequestConfig_PickGroupType['MONTH'];
  const getIssueNames = () => {
    const groupIDs = networkState.sortedGroupID();
    const names = [];
    for (const groupID of groupIDs) {
      const { meta } = networkState.getNetworkAt(groupID);
      if (!meta) {
        names.push(groupID);
        continue;
      }
      const date = fmtDate(meta.from?.toDate() || new Date());

      if (meta.metas.length === 0) {
        names.push(date);
        continue;
      }
      const firstIssue = meta.metas[0].name;
      names.push(`${date}${firstIssue}`);
    }
    return names;
  };

  const centeralityChartProps: CenteralityChartProps = {
    series: [], // {data: keyごとの中心性[], name: 単語}[]
    xaxis: {
      categories: isMonth ? networkState.sortedGroupID() : getIssueNames(), // 横軸名（key.length）
      title: isMonth ? '月' : '会議ごと',
    },
  };
  const top3NodeIDs = networkState.getTop3NodeIDInTotal();
  for (const { nodeIndex, keyIndex, node } of networkState.getWords(top3NodeIDs)) {
    if (!centeralityChartProps.series[nodeIndex]) {
      centeralityChartProps.series[nodeIndex] = {
        data: new Array<number | null>(numKeys).fill(null),
        name: node.word,
      };
    }
    centeralityChartProps.series[nodeIndex].data[keyIndex] = node.rate;
  }

  return centeralityChartProps;
};
