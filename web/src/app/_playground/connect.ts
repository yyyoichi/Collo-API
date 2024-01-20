import { NetworkStreamResponse, NodeRateStreamResponse } from '@/api/v3/collo_pb';
import { createConnectTransport } from '@connectrpc/connect-web';

export type V3Response = NetworkStreamResponse | NodeRateStreamResponse;
export type NetworkHandle = {
  start: () => void;
  stream: (resp: V3Response) => void;
  end: () => void;
  err: (e: unknown) => Error;
};

export const transport = createConnectTransport({
  baseUrl: process.env.NEXT_PUBLIC_RPC_HOST || '',
});
