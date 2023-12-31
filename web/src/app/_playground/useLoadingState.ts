import { useState } from 'react';

export const useLoadingState = () => {
  // データ取得の進捗
  const [progress, setProgress] = useState(0);
  const loading = progress != 0 && progress < 1;
  const startLoading = () => setProgress(0.05);
  const stopLoading = () => setProgress(0);
  return {
    progress,
    setProgress,
    loading,
    startLoading,
    stopLoading,
  };
};
