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

import logo from "../../assets/logo.png";

import "./home.css";
import { AppContext } from "../../AppContext";
import executionPipeline from "../../assets/images/execution-pipeline.png"
import executionPipelineDark from "../../assets/images/execution-pipeline-dark.png"

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

const microMacroContent = [
  {
    title: "Gaining Functional Insights",
    points: [
      {
        title: "Focused Evaluation",
        content:
          "Micro benchmarks dissect specific functional units within Vitess, allowing precise assessment of individual components.",
      },
      {
        title : "Golang Advantage",
        content : "Vitess leverages the Go standard library's testing framework and micro-benchmarking tools, ensuring accurate measurements and consistent results. Relying on native go features produces better accuracy."
      },
      {
        title : "Execution Efficiency",
        content : "Micro-benchmarks are effortlessly executed using the default go test runner and arewefastyet's microbench command, facilitating streamlined testing processes."
      },
      {
        title : "Critical Metrics",
        content : "Key performance indicators, such as iteration time measured in nanoseconds and memory allocation in bytes, are derived from micro-benchmark results. Multiple metrics allow better understanding for various audiences."
      },
      {
        title : "Structured Analysis",
        content : "Extracted metrics are meticulously analyzed and stored in a MySQL database, forming a valuable repository for future reference and comparison. Basically analogous to unit tests"
      },
      {
        title : "Granular Performance Assessment",
        content : "Micro benchmarks provide an unparalleled level of granularity, enabling a meticulous examination of individual code units, ensuring that even the smallest performance nuances are captured."
      }
    ],
  },
  {
    title: "Real-World Performance Insights",
    points: [
      {
        title: "Comprehensive Overview",
        content:
          "Macro benchmarks provide a comprehensive view of Vitess' performance, simulating real-world production conditions for accurate evaluations.",
      },
      {
        title : "Cluster Configuration",
        content : "Benchmark Vitess clusters are thoughtfully assembled, encompassing vtgates, vttablets, etcd clusters, and vtctld servers, creating an environment closely resembling real deployments."
      },
      {
        title : "Multi-Step Process",
        content : "Macro benchmarks comprise three sequential stages: preparation, warm-up, and the actual run, systematically capturing the performance trajectory."
      },
      {
        title : "Custom Benchmarking",
        content : "Sysbench, tailored to benchmark various data stores, is utilized for the main benchmarking process, capturing critical metrics such as latency, transactions per second (TPS), and queries per second (QPS)."
      },
      {
        title : "Incorporating Insights",
        content : "Internal cluster metrics and operating system metrics are integrated and processed through a Prometheus backend, enriching the assessment with deeper performance context."
      },
      {
        title : "Informed Optimization",
        content : "The combined insights from macro and micro benchmarks empower users to optimize Vitess effectively for their specific production environments, making well-informed decisions based on both granular functional details and broader performance trends."
      }
    ],
  },
];

const Home = () => {
  const { isColorChanged } = useContext(AppContext);

  return (
    <div className="home">
      <article className="home__top">
        <div className="home_top_gradient" />
        <div className="home_top_content">
          <div className="home_topcontent_text">
            <div className="home_topcontenttext_heading">
              <h3>Vitess Introduces</h3>
              <h1>arewefastyet</h1>
            </div>
            <p>
              A Cutting-Edge Benchmarking Approach for Unparalleled Database
              Speed
            </p>
          </div>
          <div className="home_topcontent_button">
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
        <img src={logo} alt="logo" className="home_top_logo" />
      </article>

      <article className="home__body">
        <section className="home_bodyhowit_works">
          <h1>How it works</h1>
          <div className="how_itworkscards_container">
            {howItWorksItems.map((item, key) => (
              <div className="how_itworks_card" key={key}>
                <h3>{item.title}</h3>
                <p>{item.content}</p>
              </div>
            ))}
          </div>
        </section>

        <section className="micro_and_macro">
          <h1 className="home_bodysection_title">
            Micro and Macro Benchmarks
          </h1>

          <div className="micro_andmacro_cards">
            {microMacroContent.map((item, key) => (
              <div key={key} className="micro_andmacro_card">
                <h3>{item.title}</h3>
                <ul>
                  {microMacroContent[key].points.map((point, i) => (
                    <li><h5>{point.title}</h5> <p>{point.content}</p></li>
                  ))}
                </ul>
              </div>
            ))}
          </div>
        </section>

        <section className="micro_and_macro">
          <h1 className="home_bodysection_title">
            Diagramatic overview
          </h1>
          <img src={isColorChanged ? executionPipeline : executionPipelineDark} alt="execution pipeline" className="execution__pipeline__image" />
        </section>
      </article>
    </div>
  );
};

export default Home;