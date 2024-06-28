/*
Copyright 2023 The Vitess Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import React from "react";
import PropTypes from "prop-types";
import { ResponsiveLine } from "@nivo/line";

type DataPoint = {
  x: string;
  y: string | number;
};

type DataSeries = {
  id: string;
  data: DataPoint[];
};

type ChartData = DataSeries[];

type Title = string;

type Colors = string[];

type ResponsiveChartProps = {
  data: ChartData;
  title: Title;
  colors: Colors;
  isFirstChart?: boolean;
};

/**
 * A responsive line chart component.
 * @param {ResponsiveChartProps} props - The props required are data, title and colors.
 * @returns {JSX.Element} - Render Responsive chart component.
 */

const ResponsiveChart = ({
  data,
  title,
  colors,
}: ResponsiveChartProps): JSX.Element | null => {
  if (data.length === 0 || data[0].data.length === 0) {
    return null;
  }

  return (
    <>
      <h3 className="my-10 text-xl font-medium text-primary">{title}</h3>
      <ResponsiveLine
        data={data}
        colors={colors}
        theme={{
          background: "transparent",
          axis: {
            ticks: {
              text: {
                fontSize: "13px",
                fill: "hsl(var(--foreground))",
              },
            },
          },
          legends: {
            text: {
              fontSize: "14px",
              fill: "hsl(var(--foreground))",
            },
          },
          grid: {
            line: {
              stroke: "hsl(var(--foreground))",
            },
          },
        }}
        tooltip={({ point }) => (
          <div className="flex gap-x-2 bg-background bg-opacity-80 pl-1 pr-5 py-1 rounded">
            <figure
              className="w-1"
              style={{ backgroundColor: point.serieColor }}
            />
            <div className="flex flex-col gap-y-2">
              <span>commit : {String(point.data.x)}</span>
              <span>
                {title.split("(")[0]} : {String(point.data.y)}
              </span>
            </div>
          </div>
        )}
        areaBaselineValue={50}
        margin={{ top: 50, right: 110, bottom: 50, left: 60 }}
        xScale={{ type: "point" }}
        yScale={{
          type: "linear",
          min: 0,
          max: "auto",
          reverse: false,
        }}
        yFormat=" >-.2f"
        axisTop={null}
        axisRight={null}
        pointSize={10}
        isInteractive={true}
        pointBorderWidth={2}
        pointBorderColor={{ from: "serieColor" }}
        pointLabelYOffset={-12}
        areaOpacity={0.1}
        useMesh={true}
        legends={[
          {
            anchor: "bottom-right",
            direction: "column",
            justify: false,
            translateX: 100,
            translateY: 0,
            itemsSpacing: 0,
            itemDirection: "left-to-right",
            itemWidth: 80,
            itemHeight: 20,
            itemOpacity: 0.75,
            symbolSize: 12,
            symbolShape: "circle",
            symbolBorderColor: "rgba(0, 0, 0, .5)",
            effects: [
              {
                on: "hover",
                style: {
                  itemBackground: "rgba(0, 0, 0, .03)",
                  itemOpacity: 1,
                },
              },
            ],
          },
        ]}
      />
    </>
  );
};

ResponsiveChart.propTypes = {
  data: PropTypes.arrayOf(
    PropTypes.shape({
      id: PropTypes.string.isRequired,
      data: PropTypes.arrayOf(
        PropTypes.shape({
          x: PropTypes.string.isRequired,
          y: PropTypes.oneOfType([PropTypes.string, PropTypes.number])
        })
      ).isRequired,
    })
  ).isRequired,
  title: PropTypes.string.isRequired,
  colors: PropTypes.arrayOf(PropTypes.string).isRequired,
  isFirstChart: PropTypes.bool.isRequired,
};

export default ResponsiveChart;
