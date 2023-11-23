export class NodeStrength {
  /**単語識別子とその共起回数の総数を保持する */
  readonly wordsPairCount = new Map<string, number>();
  isAddedNode(wordID: string) {
    return typeof this.wordsPairCount.get(wordID) !== 'undefined';
  }
  add(wordID: string, count: number) {
    this.wordsPairCount.set(wordID, (this.wordsPairCount.get(wordID) || 0) + count);
  }
  get(wordID: string) {
    return this.wordsPairCount.get(wordID) || 0;
  }
  /**もっとも多い合計共起回数を返す */
  getMaxCount() {
    let max = 0;
    for (const count of this.wordsPairCount.values()) {
      if (max < count) {
        max = count;
      }
    }
    return max;
  }
}
