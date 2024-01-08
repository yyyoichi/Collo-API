import { useState } from 'react';

export const useLoadingState = () => {
  // データ取得の進捗
  const [progress, setProgress] = useState(0);
  const loading = progress != 0 && progress < 1;
  const startLoading = () => setProgress(0.1);
  const setProcess = (p: number) => setProgress(p * 0.8 + 0.1);
  const endStreaming = () => setProgress(1);
  const stopLoading = () => setProgress(0);
  return {
    progress,
    loading,
    startLoading,
    setProcess,
    endStreaming,
    stopLoading,
  };
};
