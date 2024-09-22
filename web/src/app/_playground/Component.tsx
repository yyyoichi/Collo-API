"use client";
import type React from "react";
import { memo, useState } from "react";
import {
	FormComps,
	Label,
	DateInput,
	KeywordInput,
	LoadingButton,
	StartButton,
	type WrapProps,
	PoSpeechCheckbox,
	CheckboxLabel,
	StopwordsTextarea,
	AccordionPanel,
	Radio,
	ChooseBox,
} from "./Forms";
import {
	type MainNetworkGraphProps,
	NetworkGraph,
	SubNetworkGraph,
	type SubNetworkGraphProps,
} from "./NetworkGraph";
import dynamic from "next/dynamic";
import type { CenteralityChartProps } from "./Chart";
export const CenteralityChart = dynamic(() => import("./Chart"), {
	ssr: false,
});

export type PlayGroundComponentProps = {
	formProps: Pick<WrapProps, "onSubmit">;
	networkProps: MainNetworkGraphProps;
	progressBarProps: ProgressBarProps;
	loading: boolean;
	selectedNodeProps: Array<Omit<SelectedNodeProps, "name">>;
	updateNodeRankProps: UpdateNodeRankButtonProps;
	subNetworksProps: Array<SubNetworkGraphProps>;
	appendNetworkButtonProps: AppendNetworkButtonProps;
	chartProps: CenteralityChartProps;
};
export const PlayGroundComponent = (props: PlayGroundComponentProps) => {
	props.chartProps?.xaxis;
	return (
		<>
			<ProgressBar {...props.progressBarProps} />
			<FormComps.Wrap {...props.formProps}>
				<FormComps.Col>
					<Label htmlFor="from">{"開始日"}</Label>
					<DateInput id="from" name="from" defaultValue={"2023-10-01"} />
				</FormComps.Col>
				<FormComps.Col>
					<Label htmlFor="until">{"終了日"}</Label>
					<DateInput id="until" name="until" defaultValue={"2023-12-31"} />
				</FormComps.Col>
				<FormComps.Col>
					<Label htmlFor="keyword">{"キーワード"}</Label>
					<KeywordInput id="keyword" name="keyword" defaultValue={"デジタル"} />
				</FormComps.Col>
				<AccordionPanel.Head>{"詳細設定"}</AccordionPanel.Head>
				<AccordionPanel.Content>
					<FormComps.Col>
						<Label htmlFor="">{"出力品詞"}</Label>
						<ChooseBox>
							<CheckboxLabel htmlFor="noun">
								<PoSpeechCheckbox
									id="noun"
									name="noun"
									value={101}
									defaultChecked
								/>
								{"普通名詞"}
							</CheckboxLabel>
							<CheckboxLabel htmlFor="personName">
								<PoSpeechCheckbox
									id="personName"
									name="personName"
									value={111}
								/>
								{"人名"}
							</CheckboxLabel>
							<CheckboxLabel htmlFor="placeName">
								<PoSpeechCheckbox id="placeName" name="placeName" value={121} />
								{"地名"}
							</CheckboxLabel>
							<CheckboxLabel htmlFor="number">
								<PoSpeechCheckbox id="number" name="number" value={121} />
								{"数"}
							</CheckboxLabel>
							<CheckboxLabel htmlFor="adjective">
								<PoSpeechCheckbox id="adjective" name="adjective" value={201} />
								{"形容詞"}
							</CheckboxLabel>
							<CheckboxLabel htmlFor="adjectiveVerb">
								<PoSpeechCheckbox
									id="adjectiveVerb"
									name="adjectiveVerb"
									value={301}
								/>
								{"形容動詞"}
							</CheckboxLabel>
							<CheckboxLabel htmlFor="verb">
								<PoSpeechCheckbox id="verb" name="verb" value={401} />
								{"動詞"}
							</CheckboxLabel>
						</ChooseBox>
					</FormComps.Col>
					<FormComps.Col>
						<Label htmlFor="stopwords">{"除外ワード"}</Label>
						<StopwordsTextarea
							id="stopwords"
							name="stopwords"
							defaultValue={"総理 必要 尋ね"}
							placeholder={"スペース区切りで複数入力"}
						/>
					</FormComps.Col>
					<FormComps.Col>
						<CategorySelectColumn />
					</FormComps.Col>
					<FormComps.Col>
						<APISelectColumn />
					</FormComps.Col>
				</AccordionPanel.Content>
				{props.loading ? <LoadingButton /> : <StartButton />}
			</FormComps.Wrap>
			<NetworkGraph {...props.networkProps} />
			<SelectedNodesWrap>
				<UpdateNodeRankButton {...props.updateNodeRankProps} />
				{props.selectedNodeProps.map((px, i) => {
					return <SelectedNode key={i} name={i} {...px} />;
				})}
			</SelectedNodesWrap>
			<CenteralityChart {...props.chartProps} />
			{props.subNetworksProps.map((subProps, i) => {
				return <SubNetworkGraph key={i} {...subProps} />;
			})}
			<AppendNetworkButton {...props.appendNetworkButtonProps} />
		</>
	);
};

const CategorySelectColumn = () => {
	const [checked, setChecked] = useState(1);
	return (
		<>
			<Label htmlFor="">{"カテゴライズ"}</Label>
			<ChooseBox>
				<CheckboxLabel htmlFor="issue">
					<Radio
						id="issue"
						name="pick"
						value={1}
						checked={checked === 1}
						onChange={() => setChecked(1)}
					/>
					{"会議ごと"}
				</CheckboxLabel>
				<CheckboxLabel htmlFor="month">
					<Radio
						id="month"
						name="pick"
						value={2}
						checked={checked === 2}
						onChange={() => setChecked(2)}
					/>
					{"月ごと"}
				</CheckboxLabel>
			</ChooseBox>
		</>
	);
};

const APISelectColumn = () => {
	const [checked, setChecked] = useState(1);
	return (
		<>
			<Label htmlFor="">{"使用API"}</Label>
			<ChooseBox>
				<CheckboxLabel htmlFor="speechapi">
					<Radio
						id="speechapi"
						name="api"
						value={1}
						checked={checked === 1}
						onChange={() => setChecked(1)}
					/>
					{"発言単位"}
				</CheckboxLabel>
				<CheckboxLabel htmlFor="meetingapi">
					<Radio
						id="meetingapi"
						name="api"
						value={2}
						checked={checked === 2}
						onChange={() => setChecked(2)}
					/>
					{"会議単位"}
				</CheckboxLabel>
			</ChooseBox>
		</>
	);
};

type AppendNetworkButtonProps = NonNullablePick<
	React.ComponentProps<"div">,
	"onClick"
>;
const AppendNetworkButton = (props: AppendNetworkButtonProps) => {
	return (
		<div className="w-full p-4 mt-10">
			<div
				{...props}
				className={`flex items-center justify-center 
                border-gray-600 border-4 border-dashed 
                rounded-md text-lg font-bold text-gray-600 w-full h-60
                cursor-pointer
                hover:bg-blue-50 transition`}
			>
				{"+ Add Network Graph"}
			</div>
		</div>
	);
};

type ProgressBarProps = {
	progress: number;
};
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

const SelectedNodesWrap = ({ children }: React.ComponentProps<"div">) => (
	<div className="px-2 py-4 flex flex-wrap gap-2 overflow-auto max-h-48">
		{children}
	</div>
);

type UpdateNodeRankButtonProps = {
	buttonProps: NonNullablePick<
		React.ComponentProps<"button">,
		"onClick" | "disabled"
	>;
	spinner: {
		animate: boolean;
	};
};
const UpdateNodeRankButton = (props: UpdateNodeRankButtonProps) => {
	const cn = props.spinner.animate ? "animate-spin" : "";
	return (
		<button
			className="focus:outline-none focus:shadow-outline-blue disabled:opacity-50 disabled:cursor-not-allowed"
			{...props.buttonProps}
		>
			<div
				className={`h-6 w-6 m-auto ${cn} rounded-full border-b-2 border-t-2 border-blue-400`}
			/>
		</button>
	);
};

type SelectedNodeProps = {
	name: number;
	labelProps: NonNullablePick<React.ComponentProps<"div">, "children">;
	checkboxProps: NonNullablePick<
		React.ComponentProps<"input">,
		"checked" | "onChange"
	>;
};

const SelectedNode = (props: SelectedNodeProps) => {
	const id = props.name + "_selected_node";
	const labelProps: React.ComponentProps<"label"> = {
		htmlFor: id,
		...props.labelProps,
	};
	const checkboxProps: React.ComponentProps<"input"> = {
		type: "checkbox",
		id: id,
		...props.checkboxProps,
	};
	return (
		<div className="">
			<input {...checkboxProps} className="peer hidden" />
			<label
				{...labelProps}
				className="py-1 px-2 rounded-full cursor-pointer border-gray-100 border hover:shadow-md peer-checked:text-white peer-checked:bg-blue-400"
			></label>
		</div>
	);
};
