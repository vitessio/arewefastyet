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
import { ResponsiveLine } from "@nivo/line";

const ResponsiveChart = ({ data, title, colors, isFirstChart }) => {
  return (
    
    data.length > 0 && (
      <div className="chart">
        {isFirstChart ? ( // Conditionally render h3 based on isFirstChart prop
          <h3>{title}</h3>
        ) : (
          <h3 className="chart__title">{title}</h3>
        )}
        <ResponsiveLine
          data={data}
          height={400}
          colors={colors}
          theme={{
            axis: {
              ticks: {
                text: {
                  fontSize: "13px",
                  fill: "white",
                },
              },
            },
            legends: {
              text: {
                fontSize: "14px",
                fill: "white",
              },
            },
          }}
          tooltip={({ point }) => (
            <div className="tooltip flex">
              <figure style={{ backgroundColor: point.serieColor }}></figure>
              <div>commit : {point.data.x}</div>
              <div>y : {point.data.y}</div>
            </div>
          )}
          areaBaselineValue={50}
          margin={{ top: 50, right: 110, bottom: 50, left: 60 }}
          xScale={{ type: "point" }}
          yScale={{
            type: "linear",
            min: "0",
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
      </div>
    )
  );
};

export default ResponsiveChart;
