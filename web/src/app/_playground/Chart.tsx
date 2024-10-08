import React from "react";
import ReactApexChart from "react-apexcharts";
import type { ApexOptions } from "apexcharts";

export type CenteralityChartProps = {
	series: ApexAxisChartSeries;
	xaxis: {
		categories: string[];
		title: string;
	};
};
const CenteralityChart = ({ series, ...props }: CenteralityChartProps) => {
	const options: ApexOptions = {
		chart: {
			// height: 350,
			type: "line",
			// dropShadow: {
			//     enabled: true,
			//     color: '#000',
			//     top: 18,
			//     left: 7,
			//     blur: 10,
			//     opacity: 0.2
			// },
			toolbar: {
				show: false,
			},
		},
		// colors: ['#77B6EA', '#545454'],
		dataLabels: {
			enabled: true,
		},
		stroke: {
			curve: "straight",
		},
		title: {
			text: "出現単語の中心性推移",
			align: "left",
		},
		// grid: {
		//     borderColor: '#e7e7e7',
		//     row: {
		//         colors: ['#f3f3f3', 'transparent'], // takes an array which will be repeated on columns
		//         opacity: 0.5
		//     },
		// },
		markers: {
			size: 1, // pointer size
		},
		xaxis: {
			categories: props.xaxis.categories,
			title: {
				text: props.xaxis.title,
			},
		},
		yaxis: {
			title: {
				text: "中心性",
			},
			// min: 5,
			// max: 40
		},
		legend: {
			// label position
			position: "top",
			horizontalAlign: "right",
			floating: true,
			offsetY: -25,
			offsetX: -5,
		},
	};

	//@ts-ignore
	return <ReactApexChart options={options} series={series} />;
};
export default CenteralityChart;
