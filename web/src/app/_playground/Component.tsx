"use client";
import React from "react";
import { FormComps, Label, DateInput, KeywordInput, LoadingButton, StartButton, WrapProps, PoSpeechCheckbox, CheckboxLabel, StopwordsTextarea, AccordionPanel } from "./Forms";
import { MainNetworkGraphProps, NetworkGraph, SubNetworkGraph, SubNetworkGraphProps } from "./NetworkGraph";

export type PlayGroundComponentProps = {
    formProps: Pick<WrapProps, 'onSubmit'>,
    networkProps: MainNetworkGraphProps,
    isMultiMode: boolean,
    progressBarProps: ProgressBarProps,
    defaultValues: {
        from: React.ComponentProps<typeof DateInput>["defaultValue"],
        until: React.ComponentProps<typeof DateInput>["defaultValue"],
        keyword: React.ComponentProps<typeof DateInput>["defaultValue"],
    }
    loading: boolean,
    subNetworksProps: Array<SubNetworkGraphProps>,
    appendNetworkButtonProps: AppendNetworkButtonProps,
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
                <CheckboxLabel htmlFor="mode"><PoSpeechCheckbox id='mode' name="mode" value={1} defaultChecked={props.isMultiMode} />{"マルチモード"}</CheckboxLabel>
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
        <NetworkGraph {...props.networkProps} />
        {
            props.subNetworksProps.map((subProps, i) => {
                return <SubNetworkGraph key={i} {...subProps} />
            })
        }
        {
            props.isMultiMode && <AppendNetworkButton {...props.appendNetworkButtonProps} />
        }
    </>
}

type AppendNetworkButtonProps = NonNullablePick<React.ComponentProps<"div">, "onClick">
const AppendNetworkButton = (props: AppendNetworkButtonProps) => {
    return <div className="w-full p-4 mt-10">
        <div
            {...props}
            className={`flex items-center justify-center 
                border-gray-600 border-4 border-dashed 
                rounded-md text-lg font-bold text-gray-600 w-full h-60
                cursor-pointer
                hover:bg-blue-50 transition`
            }>
            {"+ Add Network Graph"}
        </div>
    </div>
}

type ProgressBarProps = {
    progress: number
}
const ProgressBar = ({ progress }: ProgressBarProps) => {
    return (
        <div className="bg-gray-200 h-3 rounded-sm sticky top-0 z-[110]">
            <div
                className="bg-green-500 h-full transition-transform duration-300"
                style={{ width: `${progress * 100}%` }}
            />
        </div>
    );
};

