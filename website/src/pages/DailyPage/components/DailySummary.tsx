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
import { twMerge } from "tailwind-merge";

type DailySummaryDataType = {
  data: {
    total_qps: {
      center: number;
      confidence: number;
      range: {
        infinite: boolean;
        unknown: boolean;
        value: number;
      };
    };
  }[];
  name: string;
};

enum BenchmarkType {
  TPCC_UNSHARDED = "TPCC_UNSHARDED",
  OLTP_SET = "OLTP-SET",
  OLTP = "OLTP",
  "OLTP-READONLY" = "OLTP-READONLY",
  TPCC = "TPCC",
  TPCC_FK = "TPCC_FK",
  TPCC_FK_UNMANAGED = "TPCC_FK_UNMANAGED",
}

interface HeroDailySummaryPropsType {
  data: DailySummaryDataType;
  benchmarkType: BenchmarkType;
  setBenchmarktype: React.Dispatch<React.SetStateAction<BenchmarkType>>;
}

/**
 * A component that renders a daily summary chart.
 * @param {HeroDailySummaryPropsType} props - The props required are data, benchmarkType and setter function for benchmarkType .
 * @returns {JSX.Element} - Render hero component for daily summary
 */

const DailySummary = ({
  data,
  setBenchmarktype,
  benchmarkType,
}: HeroDailySummaryPropsType): JSX.Element => {
  const orange = "#E77002";

  const transformedData = [];

  if (data !== null) {
    transformedData.push({
      id: "QPSTotal",
      data: data.data.map((item) => ({
        // x: item.CreatedAt,
        y: item.total_qps.center,
      })),
    });
  }

  const getBenchmarkType = () => {
    setBenchmarktype(data.name as BenchmarkType);
  };

  return (
    <div
      className={twMerge(
        "flex flex-col border border-front rounded-xl w-[20vw] h-[15vh] cursor-pointer hover:bg-accent overflow-hidden duration-300 hover:scale-105",
        benchmarkType === data.name && "bg-accent brightness-125 border-2"
      )}
      onClick={() => {
        getBenchmarkType();
      }}
    >
      {data ? (
        <div className="w-full h-1/2 p-2 pt-4 flex flex-nowrap justify-center items-center">
          <ResponsiveLine
            data={transformedData}
            margin={{ top: 20, right: 30, bottom: 50, left: 10 }}
            enableGridX={false}
            enableGridY={false}
            colors={orange}
            axisBottom={null}
            axisLeft={null}
          />
        </div>
      ) : null}
      <figure className="w-full border-b border-front"></figure>
      <div className="flex justify-between flex-1 p-3 items-center">
        <h3 className="font-extralight">{data.name}</h3>
        <i className="fa-solid fa-arrow-right daily--fa-arrow-right"></i>
      </div>
    </div>
  );
};

DailySummary.propTypes = {
  data: PropTypes.shape({
    name: PropTypes.string.isRequired,
    data: PropTypes.array,
  }),
  setBenchmarktype: PropTypes.func.isRequired,
  isSelected: PropTypes.bool,
};

export default DailySummary;
