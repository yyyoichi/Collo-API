"use client";
import { PlayGroundComponent } from "./Component";
import { useComponentProps } from "./useComponentProps";

export default function Home() {
	const props = useComponentProps();
	return <PlayGroundComponent {...props} />;
}
