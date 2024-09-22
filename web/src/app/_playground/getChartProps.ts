import { RequestConfig_PickGroupType } from "@/api/v3/collo_pb";
import type { CenteralityChartProps } from "./Chart";
import type { useNetworkState } from "./useNetworkState";

export const getChratProps = (
	networkState: ReturnType<typeof useNetworkState>,
	selectedNodeIDs: number[],
) => {
	const fmtDate = (d: Date) => {
		const mm = `0${d.getMonth() + 1}`;
		const dd = `0${d.getDate()}`;
		return `${d.getFullYear()}-${mm.substring(mm.length - 2)}-${dd.substring(dd.length - 2)}`;
	};
	const numKeys = networkState.numKeys;
	const isMonth =
		networkState.pickType === RequestConfig_PickGroupType["MONTH"];
	const getMonthNames = () => {
		const names = [];
		for (const groupID of networkState.sortedGroupID) {
			if (groupID === "total") {
				names.push("すべての期間");
			} else {
				names.push(groupID);
			}
		}
		return names;
	};
	const getIssueNames = () => {
		const names = [];
		for (const groupID of networkState.sortedGroupID) {
			if (groupID === "total") {
				names.push("すべての期間");
				continue;
			}
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
			categories: isMonth ? getMonthNames() : getIssueNames(), // 横軸名（key.length）
			title: isMonth ? "月" : "会議ごと",
		},
	};
	for (const { nodeIndex, keyIndex, node } of networkState.getWords(
		selectedNodeIDs,
	)) {
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
