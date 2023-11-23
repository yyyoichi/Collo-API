import { useState } from 'react';
import { createPromiseClient } from '@connectrpc/connect';
import { createConnectTransport } from '@connectrpc/connect-web';
import { ColloService } from '@/collo/v1/collo_connect';
import { ColloStreamRequest } from '@/collo/v1/collo_pb';
import { Timestamp } from '@bufbuild/protobuf';

const transport = createConnectTransport({
  baseUrl: process.env.NEXT_PUBLIC_GRPC_HOST || '',
});

export const useColloState = () => {
  const [collo, setCollo] = useState(new ColloNetwork());

  const requestCollo = async (from: Date = new Date(), until: Date = new Date(), keyword: string = '言葉') => {
    setCollo(new ColloNetwork());

    const client = createPromiseClient(ColloService, transport);
    const req = new ColloStreamRequest();
    const f = Timestamp.fromDate(from);
    const u = Timestamp.fromDate(until);
    req.from = f;
    req.until = u;
    req.keyword = keyword;
    const stream = client.colloStream(req);

    let store = new ColloNetwork();
    for await (const m of stream) {
      const copyWord = new Map(store.wordsMap);
      const copyPair = new Map(store.pairsList);

      for (const [key, word] of Object.entries(m.words)) {
        copyWord.set(key, word);
      }
      for (const wordIDs of m.pairs) {
        if (!wordIDs) continue;
        const count = copyPair.get(wordIDs) || 0;
        copyPair.set(wordIDs, count + 1);
      }

      store = new ColloNetwork(copyWord, copyPair);
    }
    setCollo(store);
  };

  return {
    collo,
    requestCollo,
    ready: collo.wordsMap.size === 0,
  };
};

class ColloNetwork {
  /**
   * 共起ペアを表す
   * @param wordsMap 識別子をキーに持つ単語のマップ
   * @param pairsList 識別子のペア(ColloState.joinWordsID)をキーに持つ出現回数のマップ
   */
  constructor(
    readonly wordsMap = new Map<string, string>(),
    readonly pairsList = new Map<string, number>(),
  ) {}
  /**カウントの多い順でペアリストを返す */
  getSortedPairs() {
    const pairs: {
      id1: string;
      id2: string;
      count: number;
    }[] = [];
    for (const [ids, count] of this.pairsList.entries()) {
      const [id1, id2] = this.splitWordsID(ids);
      const pair = { id1, id2, count };
      pairs.push(pair);
    }
    pairs.sort((a, b) => b.count - a.count);
    return pairs;
  }
  /**識別子から単語を取得する */
  getWordByID(id: string) {
    return this.wordsMap.get(id) || '';
  }
  private splitWordsID(ids: string): [string, string] {
    const [id1, id2] = ids.split(',');
    return [id1, id2];
  }
  static joinWordsID(id1: string, id2: string): string {
    return `${id1}:${id2}`;
  }
}
