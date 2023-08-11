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

import logo from "../../assets/logo.png";

import "./home.css";

const Home = () => {
  return (
    <div className="home">
      <article className="home__top">
        <div className="home__top__gradient" />
        <div className="home__top__content">
          <div className="home__top__content__text">
            <div className="home__top__content__text__heading">
              <h3>Vitess Introduces</h3>
              <h1>arewefastyet</h1>
            </div>
            <p>
              A Cutting-Edge Benchmarking Approach for Unparalleled Database Speed
            </p>
          </div>
          <div className="home__top__content__button">
            <button onClick={() =>window.open("https://vitess.io/blog/2021-07-08-announcing-vitess-arewefastyet","__blank")}>
              Read our blog post <i class="fa-solid fa-bookmark"></i>
            </button >
            <button onClick={() =>window.open("https://github.com/vitessio/arewefastyet","__blank")}>
              Contribute on GitHub
              <i className="fa-brands fa-github"></i>
            </button>
          </div>
        </div>
        <img src={logo} alt="logo" className="home__top__logo" />
      </article>
    </div>
  );
};

export default Home;
