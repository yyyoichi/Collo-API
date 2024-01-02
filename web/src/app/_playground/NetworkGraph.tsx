"use client";
import React from "react";
import { LayoutForceAtlas2Control } from "@react-sigma/layout-forceatlas2";
import { SigmaContainer, ControlsContainer, ZoomControl, FullScreenControl, SearchControl } from "@react-sigma/core";
import '@react-sigma/core/lib/react-sigma.min.css';

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
    loaderProps: NetworkGraphLoaderProps;
    selectProps: GroupSelectProps;
}
export const SubNetworkGraph = (props: SubNetworkGraphProps) => {
    return (
        <>
            <div className="relative h-1 bg-gray-600 my-10">
                <div className="absolute flex top-1 z-10 w-full overflow-hidden max-w-md mx-auto p-3 bg-white shadow-md rounded-md">
                    <div {...props.deleteButtonProps} className="cursor-pointer px-2">{"[x]"}</div>
                    <GroupSelect {...props.selectProps} />
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
    return <select {...props.groupSelectProps} className="w-full">
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
    return <option {...props} className="px-2 py-1">
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


