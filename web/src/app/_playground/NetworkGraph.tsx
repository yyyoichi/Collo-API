"use client";
import React from "react";
import { LayoutForceAtlas2Control } from "@react-sigma/layout-forceatlas2";
import { SigmaContainer, ControlsContainer, ZoomControl, FullScreenControl, SearchControl } from "@react-sigma/core";
import '@react-sigma/core/lib/react-sigma.min.css';
import Link, { LinkProps } from "next/link";
import { LoadingButton, StartButton } from "./Forms";

export type MainNetworkGraphProps = {
    loaderProps: NetworkGraphLoaderProps;
}
export const NetworkGraph = (props: MainNetworkGraphProps) => {
    return (
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
    )
}

export type SubNetworkGraphProps = {
    deleteButtonProps: NonNullablePick<React.ComponentProps<"div">, "onClick">;
    contentsProps: {
        metaProps: Array<NonNullablePick<React.ComponentProps<"a">, "href" | "children">>;
        loading: boolean;
        top3Button: NonNullablePick<React.ComponentProps<"input">, "disabled" | "onClick">;
    },
    loaderProps: NetworkGraphLoaderProps;
    selectProps: GroupSelectProps;
}
export const SubNetworkGraph = (props: SubNetworkGraphProps) => {
    return (
        <>
            <div className="relative h-1 bg-gray-600 my-10">
                {/* absolute panel */}
                <div className="absolute top-1 z-10 w-full overflow-x-hidden max-w-lg mx-auto py-4 px-3 bg-white shadow-md rounded-md resize h-fit max-h-[580px]">
                    {/* headers */}
                    <div className="flex my-5">
                        <div {...props.deleteButtonProps} className="mr-2 px-2 cursor-pointer rounded-sm hover:bg-red-100">
                            {"x"}
                        </div>
                        <GroupSelect {...props.selectProps} />
                    </div>
                    {/* contents */}
                    <div className="text-gray-600 my-2">
                        <h3 className="">{"対象会議録"}</h3>
                        <div className="mx-2">
                            {
                                props.contentsProps.metaProps.map((metaProp, i) => {
                                    return (
                                        <div className="text-sm after:contents" key={i}>
                                            <a className="hover:text-blue-600" href={metaProp.href} target="_blank">
                                                {metaProp.children}
                                            </a>
                                        </div>
                                    )
                                })
                            }
                        </div>
                    </div>
                    {
                        props.contentsProps.loading ? (
                            <LoadingButton />
                        ) : (
                            <StartButton type="button" value={"主要3単語のネットワークを取得する"} {...props.contentsProps.top3Button} />
                        )
                    }
                </div>
            </div>
            <SigmaContainer style={{ height: "600px" }}>
                <NetworkGraphLoader {...props.loaderProps} />
                <ControlsContainer position={"bottom-right"}>
                    <ZoomControl />
                    <FullScreenControl />
                    <LayoutForceAtlas2Control settings={{ settings: { slowDown: 10 } }} />
                </ControlsContainer>
            </SigmaContainer>
        </>
    )
}
type GroupSelectProps = {
    groupSelectProps: NonNullablePick<React.ComponentProps<"select">, "onChange">;
    groupOptionProps: Array<Omit<GroupOptionProps, "key">>;
}
const GroupSelect = (props: GroupSelectProps) => {
    return <select {...props.groupSelectProps} className="w-full rounded border focus:outline-none focus:border-blue-500">
        <option value={""}>{"選択してください"}</option>
        {
            props.groupOptionProps.map((p, i) => {
                return <GroupOption key={i} {...p} />
            })
        }
    </select>
}
type GroupOptionProps = NonNullablePick<React.ComponentProps<"option">, "key" | "value" | "children">
const GroupOption = ({ children, ...props }: GroupOptionProps) => {
    return <option {...props} className="px-2 py-1 w-full text-sm">
        {children}
    </option>
}

type NetworkGraphLoaderProps = {
    useLoadingGraphEffect: () => void
}
const NetworkGraphLoader = (props: NetworkGraphLoaderProps) => {
    props.useLoadingGraphEffect();
    return null
}


