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
import PropTypes from 'prop-types';
import { ResponsiveLine } from "@nivo/line";

import "./dailySummary.css";

const DailySummary = ({ data, setBenchmarktype, isSelected, handleClick  }) => {
  const orange = "#E77002";
  
  const transformedData = [];

  if (data.Data !== null) {
    transformedData.push({
      id: "QPSTotal",
      data: data.Data.map((item) => ({
        x: item.CreatedAt,
        y: item.QPSTotal,
      })),
    });
  }

  const getBenchmarkType = () => {
    setBenchmarktype(data.Name);
  };

  return (
    <div className={`dailySummary flex--column ${isSelected ? "dailySummary--selected" : ""}`}  onClick={() => {
      handleClick();
      getBenchmarkType();
    }}>
      {data ? (
        <div className="dailySummary__chart">
          <ResponsiveLine
            data={transformedData}
            margin={{ top: 20, right: 30, bottom: 50, left: 10 }}
            height={300}
            enableGridX={false}
            enableGridY={false}
            colors={orange}
            axisBottom={null}
            axisLeft={null}
          />
        </div>
      ) : null}
      <figure className="dailySummary__line"></figure>
      <div className="dailySummary__text">
        <h3>{data.Name}</h3>
        <i className="fa-solid fa-arrow-right daily--fa-arrow-right"></i>
      </div>
    </div>
  );
};

CronSummary.propTypes = {
  data: PropTypes.shape({
    Name: PropTypes.string.isRequired,
    Data: PropTypes.array,
  }),
  setBenchmarktype: PropTypes.func.isRequired,
  isSelected: PropTypes.bool,
  handleClick: PropTypes.func.isRequired,
};

export default DailySummary;
