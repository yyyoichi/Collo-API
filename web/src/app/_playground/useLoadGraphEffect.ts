import { useState, useEffect } from 'react';
import { useLoadGraph, useRegisterEvents, useSetSettings, useSigma } from '@react-sigma/core';
import { useLayoutForceAtlas2 } from '@react-sigma/layout-forceatlas2';
import Graph from 'graphology';
import { Attributes } from 'graphology-types';
import { SigmaNodeEventPayload } from 'sigma/sigma';

export type NetworkGraphLoaderProps = {
  clickNode: (payload: SigmaNodeEventPayload) => void;
  updateGraph: (graph: Graph) => boolean;
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
    if (!props.updateGraph(graph)) {
      return;
    }
    loadGraph(graph);
    assign();

    registerEvents({
      clickNode: props.clickNode,
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
