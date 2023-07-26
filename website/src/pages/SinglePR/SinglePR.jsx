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

import "../SinglePR/singlePR.css";

const SinglePR = () => {
  return (
    <div className="singlePR flex--column">
      <div className="singlePR__top flex">
        <div>
          <h2>
			[#46546]
            Enhancing VTGate buffering for MoveTables and Shard by Shard
            Migration
          </h2>
          <span>By Frouioui at 07/15/2023 00:19 </span>
        </div>

        <div className="singlePR__link justify--content">
          <a>Compare with base commit</a>
        </div>
      </div>
	  <div>
		
	  </div>
    </div>
  );
};

export default SinglePR;
