import { FormComps, Label, DateInput, KeywordInput, LoadingButton, StartButton, WrapProps, PoSpeechCheckbox, CheckboxLabel, StopwordsTextarea, AccordionPanel } from "./Forms";

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
            <FormComps.Col>
                <CheckboxLabel htmlFor="mode"><PoSpeechCheckbox id='mode' name="mode" value={1} />{"マルチモード"}</CheckboxLabel>
            </FormComps.Col>
            <AccordionPanel.Head>{"詳細設定"}</AccordionPanel.Head>
            <AccordionPanel.Content>
                <FormComps.Col>
                    <Label htmlFor="">{"出力品詞"}</Label>
                    <div className="flex flex-wrap gap-1 mt-1 p-2 border-b rounded-md w-full">
                        <CheckboxLabel htmlFor="noun"><PoSpeechCheckbox id='noun' name="noun" value={101} defaultChecked />{"普通名詞"}</CheckboxLabel>
                        <CheckboxLabel htmlFor="personName"><PoSpeechCheckbox id='personName' name="personName" value={111} />{"人名"}</CheckboxLabel>
                        <CheckboxLabel htmlFor="placeName"><PoSpeechCheckbox id='placeName' name="placeName" value={121} />{"地名"}</CheckboxLabel>
                        <CheckboxLabel htmlFor="number"><PoSpeechCheckbox id='number' name="number" value={121} />{"数"}</CheckboxLabel>
                        <CheckboxLabel htmlFor="adjective"><PoSpeechCheckbox id='adjective' name="adjective" value={201} />{"形容詞"}</CheckboxLabel>
                        <CheckboxLabel htmlFor="adjectiveVerb"><PoSpeechCheckbox id='adjectiveVerb' name="adjectiveVerb" value={301} />{"形容動詞"}</CheckboxLabel>
                        <CheckboxLabel htmlFor="verb"><PoSpeechCheckbox id='verb' name="verb" value={401} />{"動詞"}</CheckboxLabel>
                    </div>
                </FormComps.Col>
                <FormComps.Col>
                    <Label htmlFor='stopwords'>{"除外ワード"}</Label><StopwordsTextarea id='stopwords' name='stopwords' placeholder={"スペース区切りで複数入力"} />
                </FormComps.Col>
            </AccordionPanel.Content>
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

type NetworkGraphLoaderProps = {
    useLoadingGraphEffect: () => void
}
const NetworkGraphLoader = (props: NetworkGraphLoaderProps) => {
    props.useLoadingGraphEffect();
    return null
}
