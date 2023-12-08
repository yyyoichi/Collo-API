import { PlayGroundComponentProps } from './Component';
import { useNetworkState } from './useNetworkState';

export const useComponentProps = (): PlayGroundComponentProps => {
  const networkState = useNetworkState();
  const props: PlayGroundComponentProps = {
    formProps: {
      onSubmit: (event) => {
        event.preventDefault();
        networkState.startLoading();

        const form = new FormData(event.currentTarget);
        const start = form.get('start');
        const end = form.get('end');
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
        networkState.resetRequestParams(from, until, keyword.toString());
        networkState.request().catch((e) => {
          console.error(e);
          networkState.stopLoading();
        });
      },
    },
    loaderProps: networkState,
    loading: networkState.loading,
  };

  return props;
};
