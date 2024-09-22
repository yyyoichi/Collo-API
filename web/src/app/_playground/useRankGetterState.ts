import { useEffect, useState } from "react";
import {
	type NetworkStreamRequest,
	Node,
	NodeRateStreamRequest,
	RequestConfig,
} from "@/api/v3/collo_pb";
import { MintGreenService } from "@/api/v3/collo_connect";
import { createPromiseClient } from "@connectrpc/connect";
import { type NetworkHandle, transport } from "./connect";

// 初期値は3.いつもtop3つはすべて初めのリクエストで返されるため。
const DEFAULT_DONES = 3;
// 現在のリクエストが保持するワード（ランク）数。
let max = -1;
export const useRankGetterState = (params: NetworkStreamRequest) => {
	// 取得ランク数
	const [dones, setDones] = useState(DEFAULT_DONES);
	useEffect(() => {
		max = -1;
		setDones(DEFAULT_DONES);
	}, [params]);

	const hasNext = max == -1 || dones < max;
	const request = async (
		config: RequestConfig | undefined,
		handle: NetworkHandle,
	) => {
		handle.start();
		if (!hasNext || config?.keyword == "") {
			handle.end();
			return;
		}
		const client = createPromiseClient(MintGreenService, transport);
		const req = new NodeRateStreamRequest();
		req.config = config || new RequestConfig();
		req.offset = dones;
		req.limit = 50;
		try {
			const stream = client.nodeRateStream(req);
			console.log(
				`Start request.. Offset:${req.offset},`,
				`Limit:${req.limit}`,
			);
			let returnCount = 0;
			for await (const m of stream) {
				handle.stream(m);
				const key = (max = m.num);
				returnCount = m.count;
			}
			setDones((dones) => (dones += returnCount));
			handle.end();
		} catch (e) {
			return handle.err(e);
		}
	};

	return {
		dones,
		request,
		hasNext,
	};
};
