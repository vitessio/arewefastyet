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

import useApiCall from "../../utils/Hook";
import RingLoader from "react-spinners/RingLoader";

import PRTable from "./components/PRTable";
import { prDataTypes } from "@/types";
import Hero, { HeroProps } from "@/common/Hero";

export default function PRsPage() {
  const {
    data: dataPRList,
    isLoading: isPRListLoading,
    error: PRListError,
  } = useApiCall<prDataTypes>(`${import.meta.env.VITE_API_URL}pr/list`);

  const heroProps: HeroProps = {
    title: "Pull Request",
    description: (
      <p>
      If a given Pull Request on vitessio/vitess is labelled with the
      <span className="bg-red-800 text-white px-2 py-1 rounded-2xl ml-2">Benchmark Me</span> label
      the Pull Request will be handled and benchmarked by arewefastyet. For each commit on the Pull Request there will be two benchmarks: one on the Pull Request's HEAD and another on the base of the Pull Request.
      <br />
      <br />
      On this page you can find all benchmarked Pull Requests.
    </p>
    ),
  };

  return (
    <>
      <Hero title={heroProps.title} description={heroProps.description}/>

      {isPRListLoading && (
        <div className="flex justify-center w-full my-16">
          <RingLoader loading={isPRListLoading} color="#E77002" size={300} />
        </div>
      )}

      {PRListError ? <div className="text-red-500 text-center my-2">{PRListError}</div> : null}

      {!isPRListLoading && dataPRList && <PRTable data={dataPRList} />}
    </>
  );
}
