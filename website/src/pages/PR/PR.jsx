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

import React, { useState, useEffect } from "react";

import "../PR/PR.css";

import { errorApi } from "../../utils/Utils";
import { Link } from "react-router-dom";

const PR = () => {
  const [dataPRList, setDataPRList] = useState([]);
  const [error, setError] = useState(null);
  const [prNumber, setPrNumber] = useState("");
  const [dataPRInfo, setDataPRInfo] = useState([]);
  const [destination, setDestination] = useState("");

  useEffect(() => {
    const fetchData = async () => {
      try {
        const responsePRList = await fetch(
          `${import.meta.env.VITE_API_URL}pr/list`
        );

        const jsonDataPRList = await responsePRList.json();
        console.log(jsonDataPRList);
        setDataPRList(jsonDataPRList);
      } catch (error) {
        console.log("Error while retrieving data from the API", error);
        setError(errorApi);
      }
    };

    fetchData();
  }, []);

  useEffect(() => {
    const fetchData = async () => {
      if (prNumber === "") {
        return;
      }
      try {
        const responsePRInfo = await fetch(
          `${import.meta.env.VITE_API_URL}pr/info/${prNumber}`
        );

        const jsonDataPRInfo = await responsePRInfo.json();
        console.log(jsonDataPRInfo.Main);
        setDataPRInfo(jsonDataPRInfo);

        if (jsonDataPRInfo.Main) {
          window.location.href = `/compare?ltag=${jsonDataPRInfo.Main}&rtag=${jsonDataPRInfo.PR}`;
        }
      } catch (error) {
        console.log("Error while retrieving data from the API", error);
        setError(errorApi);
      }
    };

    fetchData();
  }, [prNumber]);

  const handlePrInfo = (e) => {
    const number = e.toString();
    setPrNumber(number);
  };

  return (
    <div className="pr">
      <div className="pr__top justify--content">
        <h1>Pull Request</h1>
        <span>
          Lorem ipsum dolor sit amet, consectetur adipiscing elit. Integer a
          augue mi. Etiam sed imperdiet ligula, vel elementum velit. Phasellus
          sodales felis eu condimentum convallis. Suspendisse sodales malesuada
          iaculis. Mauris molestie placerat ex non malesuada. Curabitur eget
          sagittis eros. Aliquam aliquam sem non tincidunt volutpat.
        </span>
      </div>
      <figure className="line"></figure>
      <div className="pr__container justify--content">
        {dataPRList.map((PRList, index) => {
          return (
            <span key={index} onClick={() => handlePrInfo(PRList)}>
              {PRList}
            </span>
          );
        })}
      </div>
    </div>
  );
};

export default PR;
