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

import React, { useContext } from "react";
import { AppContext } from "../../AppContext";

import logo from "../../assets/logo.png";

import "./home.css";

const Home = () => {

  return (
    <div className="home">
      <article className="home__top justify--content">
        <div className="home__top__text">
          <h1 className="header--title">arewefastyet</h1>
          <span>
            Arewefastyet is all about precise performance measurement.
            We test Vitess in various scenarios, assessing query latency, transaction speed, and CPU/Memory usage.
            These insights drive our continuous improvement efforts to ensure Vitess remains at the forefront of performance.
          </span>
        </div>
        <img
          src={logo}
          alt="logo"
          className="home__top__logo"
        />
      </article>
    </div>
  );
};

export default Home;
