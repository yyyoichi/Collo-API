import { useState } from "react";
import type { NetworkStreamRequest } from "@/api/v3/collo_pb";

// 過去のリクエストを保持する。
export const useReqHistoryState = () => {
	const [histories, setHistories] = useState<NetworkStreamRequest[]>([]);
	return {
		inHistories: (forcusNodeID: number, forcusGroupID: string) => {
			for (const history of histories) {
				if (
					history.config?.forcusGroupId == forcusGroupID &&
					history.forcusNodeId == forcusNodeID
				) {
					return true;
				}
			}
			return false;
		},
		clearHistories: () => setHistories([]),
		addHisotry: (req: NetworkStreamRequest) =>
			setHistories((phv) => [...phv, req]),
	};
};
