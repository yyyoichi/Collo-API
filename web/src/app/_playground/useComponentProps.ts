import { PlayGroundComponentProps } from './Component';
import { useNetworkState } from './useNetworkState';

export const useComponentProps = (): PlayGroundComponentProps => {
  const networkState = useNetworkState();
  const fmtDate = (d: Date) => {
    const mm = `0${d.getMonth() + 1}`;
    const dd = `0${d.getDate()}`;
    return `${d.getFullYear()}-${mm.substring(mm.length - 2)}-${dd.substring(dd.length - 2)}`;
  };
  const props: PlayGroundComponentProps = {
    formProps: {
      onSubmit: (event) => {
        event.preventDefault();
        networkState.startLoading();

        const form = new FormData(event.currentTarget);
        const start = form.get('from');
        const end = form.get('until');
        const keyword = form.get('keyword');
        if (!start || !end || !keyword) {
          networkState.stopLoading();
          return;
        }
        const from = new Date(start.toString());
        const until = new Date(end.toString());
        if (from.getTime() > until.getTime()) {
          networkState.stopLoading();
          return;
        }
        networkState.newRequest(from, until, keyword.toString());
      },
    },
    defaultValues: {
      from: fmtDate(networkState.initRequestParams.from),
      until: fmtDate(networkState.initRequestParams.until),
      keyword: networkState.initRequestParams.keyword,
    },
    progressBarProps: {
      progress: networkState.progress,
    },
    loaderProps: networkState,
    loading: networkState.loading,
  };

  return props;
};
