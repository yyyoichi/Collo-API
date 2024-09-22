import type { PlayGroundComponentProps } from "./Component";

type LoaderProps = PlayGroundComponentProps["networkProps"]["loaderProps"];

type GetNetworkLoaderProps = (id: string) => LoaderProps;

// subnetworkのGroupIDを変更するとloaderも更新されてしまうための対策。
// Memoはトップレベルでのみ呼ばれるのでGroupIDsループ内では使用できないので代用。
// 同じオブジェクトを参照するためのグローバル変数。
const store = new Map<string, LoaderProps>();

export const getLoaderProps = (id: string, getter: GetNetworkLoaderProps) => {
	const got = store.get(id);
	if (got) return got;
	const props = getter(id);
	store.set(id, props);
	return store.get(id)!;
};
export const clearLoaderPropsMemo = () => {
	store.clear();
};
