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

const howItWorksItems = [
  {
    title: "The Execution Engine",
    content:
      "At the heart of AreWeFastYet is the Execution Engine. It orchestrates the entire benchmarking process, ensuring accuracy and reproducibility on a large scale. Each benchmark run, referred to as an 'execution', is initiated from various sources like CLI triggers, scheduled tasks, or events like pull requests or releases. This kicks off a pipeline creation, configured via a YAML file. This file encompasses infrastructure provisioning, results storage, notifications, and more, setting the stage for meticulous benchmarking.",
  },
  {
    title: "Dedicated Benchmarking Servers",
    content:
      "For our production deployment, AreWeFastYet relies on dedicated hardware provided by Equinix Metal. Our benchmarking infrastructure uses m2.xlarge.x86 bare-metal servers, boosting benchmark reliability and accuracy. Terraform handles the server provisioning, ensuring consistent configurations. Dynamic adjustments are applied using Ansible roles, tailored to each benchmark's requirements. Whether it's a macro-benchmark involving Vitess clusters or a micro-benchmark, the server configurations are tailored accordingly.",
  },
  {
    title: "Customized Benchmark Settings",
    content:
      "Different benchmarks demand distinct configurations. For instance, a macro-benchmark necessitates the setup of a Vitess cluster, while a micro-benchmark might not. Server configuration includes package installations, binary setups, network adjustments, and the deployment of both Vitess and AreWeFastYet codebases. The Vitess cluster's settings align with the initial trigger configuration. The default setup examines Vitess performance in a sharded keyspace with six vtgates and two vttablets.",
  },
  {
    title: "Starting Benchmark Runs",
    content:
      "Once the server is primed, the final step is initiating the benchmark run. Ansible triggers AreWeFastYet's CLI to set the benchmark in motion. This comprehensive process, from YAML-based pipeline configuration to dynamic server setup, ensures that every benchmark run is accurate, reproducible, and adaptable to the unique demands of each benchmark type. AreWeFastYet streamlines the complexities of executing benchmarks against Vitess, offering a robust and precise benchmarking solution at scale.",
  },
];

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
              A Cutting-Edge Benchmarking Approach for Unparalleled Database
              Speed
            </p>
          </div>
          <div className="home__top__content__button">
            <button
              onClick={() =>
                window.open(
                  "https://vitess.io/blog/2021-07-08-announcing-vitess-arewefastyet",
                  "__blank"
                )
              }
            >
              Read our blog post <i class="fa-solid fa-bookmark"></i>
            </button>
            <button
              onClick={() =>
                window.open(
                  "https://github.com/vitessio/arewefastyet",
                  "__blank"
                )
              }
            >
              Contribute on GitHub
              <i className="fa-brands fa-github"></i>
            </button>
          </div>
        </div>
        <img src={logo} alt="logo" className="home__top__logo" />
      </article>

      <article className="home__body">
        <section className="home__body__how__it__works">
          <h1>How it works</h1>
          <div className="how__it__works__cards__container">
            {howItWorksItems.map((item, key) => (
              <div className="how__it__works__card" key={key}>
                <h3>{item.title}</h3>
                <p>{item.content}</p>
              </div>
            ))}
          </div>
        </section>
      </article>
    </div>
  );
};

export default Home;
