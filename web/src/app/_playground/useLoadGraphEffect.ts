import { useState, useEffect } from "react";
import {
	useLoadGraph,
	useRegisterEvents,
	useSetSettings,
	useSigma,
} from "@react-sigma/core";
import { useLayoutForceAtlas2 } from "@react-sigma/layout-forceatlas2";
import Graph from "graphology";
import type { Attributes } from "graphology-types";
import type { Edge, Node } from "@/api/v3/collo_pb";

export type NetworkGraphLoaderProps = {
	progress: number;
	asset: {
		nodes: Node[];
		edges: Edge[];
	};
	startLoading: () => void;
	continueRequest: (forcusID: number) => Promise<Error | undefined>;
};
export const useLoadGraphEffect = (props: NetworkGraphLoaderProps) => {
	const { assign } = useLayoutForceAtlas2();
	const registerEvents = useRegisterEvents();
	const loadGraph = useLoadGraph();

	const sigma = useSigma();
	const setSettings = useSetSettings();
	const [hoveredNode, setHoveredNode] = useState<string | null>(null);

	useEffect(() => {
		const graph = new Graph();
		if (props.progress < 1) return;
		const asset = props.asset;
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
			if (
				graph.hasEdge(edge.edgeId) ||
				!graph.hasNode(edge.nodeId1) ||
				!graph.hasNode(edge.nodeId2)
			)
				continue;
			graph.addEdgeWithKey(edge.edgeId, edge.nodeId1, edge.nodeId2, {
				size: 1,
			});
		}
		loadGraph(graph);
		assign();

		registerEvents({
			clickNode: (payload) => {
				payload.preventSigmaDefault();
				props.startLoading();
				let forcusID = 0;
				try {
					forcusID = Number(payload.node);
				} catch (e) {
					console.error(e);
					window.alert();
				}
				if (forcusID == 0) {
					return;
				}
				props.continueRequest(forcusID).then((res) => {
					if (res instanceof Error) {
						window.alert(res.message);
					}
					setHoveredNode(null);
				});
			},
			enterNode: (event) => setHoveredNode(event.node),
			leaveNode: () => setHoveredNode(null),
		});
	}, [loadGraph, registerEvents, assign, props]);

	// hover アクション
	useEffect(() => {
		setSettings({
			nodeReducer: (node, data) => {
				const graph = sigma.getGraph();
				const newData: Attributes = {
					...data,
					highlighted: data.highlighted || false,
				};
				// user doesnot hover any node
				if (!hoveredNode) return newData;

				// hightligth hover node and related node
				if (
					node === hoveredNode ||
					graph.neighbors(hoveredNode).includes(node)
				) {
					newData.highlighted = true;
				} else {
					newData.color = "#E2E2E2";
					newData.highlighted = false;
				}
				return newData;
			},
			edgeReducer: (edge, data) => {
				const graph = sigma.getGraph();
				const newData = { ...data, hidden: false };

				if (hoveredNode && !graph.extremities(edge).includes(hoveredNode)) {
					newData.hidden = true;
				}
				return newData;
			},
		});
	}, [hoveredNode, setSettings, sigma]);

	return null;
};
