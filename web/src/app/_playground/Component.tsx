import { useLoadGraphEffect } from "./useLoadGraphEffect";
import { useNetworkState } from "./useNetworkState";

import { FormComps, Label, DateInput, KeywordInput, LoadingButton, StartButton, WrapProps } from "./Forms";

import { LayoutForceAtlas2Control } from "@react-sigma/layout-forceatlas2";
import { SigmaContainer, ControlsContainer, ZoomControl, FullScreenControl, SearchControl } from "@react-sigma/core";
import '@react-sigma/core/lib/react-sigma.min.css';

export type PlayGroundComponentProps = {
    formProps: Pick<WrapProps, 'onSubmit'>,
    loaderProps: NetworkGraphLoaderProps,
    progressBarProps: ProgressBarProps,
    defaultValues: {
        from: React.ComponentProps<typeof DateInput>["defaultValue"],
        until: React.ComponentProps<typeof DateInput>["defaultValue"],
        keyword: React.ComponentProps<typeof DateInput>["defaultValue"],
    }
    loading: boolean,
}
export const PlayGroundComponent = (props: PlayGroundComponentProps) => {
    return <>
        <ProgressBar {...props.progressBarProps} />
        <FormComps.Wrap {...props.formProps}>
            <FormComps.Col>
                <Label htmlFor='from'>{"開始日"}</Label><DateInput id='from' name='from' defaultValue={props.defaultValues.from} />
            </FormComps.Col>
            <FormComps.Col>
                <Label htmlFor='until'>{"終了日"}</Label><DateInput id='until' name='until' defaultValue={props.defaultValues.until} />
            </FormComps.Col>
            <FormComps.Col>
                <Label htmlFor='keyword'>{"キーワード"}</Label><KeywordInput id='keyword' name='keyword' defaultValue={props.defaultValues.keyword} />
            </FormComps.Col>
            {props.loading ? <LoadingButton /> : <StartButton />}
        </FormComps.Wrap>
        <SigmaContainer style={{ height: "600px" }}>
            <NetworkGraphLoader {...props.loaderProps} />
            <ControlsContainer position={"bottom-right"}>
                <ZoomControl />
                <FullScreenControl />
                <LayoutForceAtlas2Control settings={{ settings: { slowDown: 10 } }} />
            </ControlsContainer>
            <ControlsContainer position={"top-right"}>
                <SearchControl style={{ width: "400px" }} />
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
