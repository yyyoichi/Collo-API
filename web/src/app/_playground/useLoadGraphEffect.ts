import { useState, useEffect } from 'react';
import { useLoadGraph, useRegisterEvents, useSetSettings, useSigma } from '@react-sigma/core';
import { useLayoutForceAtlas2 } from '@react-sigma/layout-forceatlas2';
import { useNetworkState } from './useNetworkState';
import Graph from 'graphology';
import { Attributes } from 'graphology-types';

export const useLoadGraphEffect = (props: ReturnType<typeof useNetworkState>) => {
  const { assign } = useLayoutForceAtlas2();
  const registerEvents = useRegisterEvents();
  const loadGraph = useLoadGraph();

  const sigma = useSigma();
  const setSettings = useSetSettings();
  const [hoveredNode, setHoveredNode] = useState<string | null>(null);

  useEffect(() => {
    if (props.progress < 1) return;
    const graph = new Graph();
    // add nodes
    for (const node of props.network.nodes) {
      if (graph.hasNode(node.nodeId)) continue;
      graph.addNode(node.nodeId, {
        label: node.word,
        size: 3,
        x: Math.random() * 100,
        y: Math.random() * 100,
      });
    }
    for (const edge of props.network.edges) {
      if (graph.hasEdge(edge.edgeId)) continue;
      graph.addEdgeWithKey(edge.edgeId, edge.nodeId1, edge.nodeId2, {
        size: edge.count,
      });
    }
    loadGraph(graph);
    assign();

    registerEvents({
      clickNode: (event) => {
        event.preventSigmaDefault();
        props.startLoading();
        let forcusID = 0;
        try {
          forcusID = Number(event.node);
        } catch (e) {
          console.error(e);
        }
        if (forcusID) {
          props.continueRequest(forcusID);
        }
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
        const newData: Attributes = { ...data, highlighted: data.highlighted || false };
        // user doesnot hover any node
        if (!hoveredNode) return newData;

        // hightligth hover node and related node
        if (node === hoveredNode || graph.neighbors(hoveredNode).includes(node)) {
          newData.highlighted = true;
        } else {
          newData.color = '#E2E2E2';
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
