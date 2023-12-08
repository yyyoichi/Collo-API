import { SigmaContainer, ControlsContainer, ZoomControl, FullScreenControl, SearchControl } from "@react-sigma/core";
import { LayoutForceAtlas2Control } from "@react-sigma/layout-forceatlas2";
import { FormComps, Label, DateInput, KeywordInput, LoadingButton, StartButton, WrapProps } from "./Forms";
import { useLoadGraphEffect } from "./useLoadGraphEffect";
import { useNetworkState } from "./useNetworkState";
import '@react-sigma/core/lib/react-sigma.min.css';

export type PlayGroundComponentProps = {
    formProps: Pick<WrapProps, 'onSubmit'>,
    loaderProps: NetworkGraphLoaderProps,
    progressBarProps: ProgressBarProps,
    loading: boolean,
}
export const PlayGroundComponent = (props: PlayGroundComponentProps) => {
    return <>
        <FormComps.Wrap {...props.formProps}>
            <FormComps.Col>
                <Label htmlFor='keyword'>{"開始日"}</Label><DateInput id='start' name='start' defaultValue={"2023-03-01"} />
            </FormComps.Col>
            <FormComps.Col>
                <Label htmlFor='keyword'>{"終了日"}</Label><DateInput id='end' name='end' defaultValue={"2023-03-03"} />
            </FormComps.Col>
            <FormComps.Col>
                <Label htmlFor='keyword'>{"キーワード"}</Label><KeywordInput id='keyword' name='keyword' defaultValue={"アニメ"} />
            </FormComps.Col>
            {props.loading ? <LoadingButton /> : <StartButton />}
        </FormComps.Wrap>
        <ProgressBar {...props.progressBarProps} />
        <SigmaContainer style={{ height: "500px" }}>
            <NetworkGraphLoader {...props.loaderProps} />
            <ControlsContainer position={"bottom-right"}>
                <ZoomControl />
                <FullScreenControl />
                <LayoutForceAtlas2Control settings={{ settings: { slowDown: 10 } }} />
            </ControlsContainer>
            <ControlsContainer position={"top-right"}>
                <SearchControl style={{ width: "200px" }} />
            </ControlsContainer>
        </SigmaContainer>
    </>
}

type ProgressBarProps = {
    progress: number
}
const ProgressBar = ({ progress }: ProgressBarProps) => {
    return (
        <div className="bg-gray-200 h-2 rounded overflow-hidden">
            <div
                className="bg-green-500 h-full transition-transform duration-300"
                style={{ width: `${progress * 100}%` }}
            />
        </div>
    );
};

type NetworkGraphLoaderProps = ReturnType<typeof useNetworkState>
const NetworkGraphLoader = (props: NetworkGraphLoaderProps) => {
    useLoadGraphEffect(props);
    return null
}
