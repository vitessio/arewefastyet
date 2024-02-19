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
import { Link } from "react-router-dom";
import Icon from "../../../common/Icon";

export default function Hero() {
  return (
    <section className="hidden md:flex items-center h-screen relative p-page">
      <div className="absolute-cover overflow-hidden -z-1">
        <div
          className="absolute-cover bg-gradient-to-br from-primary to-theme scale-150"
          style={{ clipPath: "polygon(0% 0%, 35% 0%, 62% 100%, 0% 100%)" }}
        />
      </div>
      <div className="flex flex-1">
        <div className="flex flex-col gap-y-8 flex-1 ml-8">
          <div className="flex flex-col bg-background flex-1 p-10 rounded-3xl gap-y-2">
            <h1 className="text-6xl font-bold text-primary">arewefastyet</h1>
            <p className="text-lg font-normal mt-5 whitespace-nowrap">
              A Benchmarking System for Vitess
            </p>
          </div>
          <div className="flex gap-x-8">
            <Link
              className="bg-black text-white rounded-2xl p-5 flex items-center gap-x-2 duration-300 hover:scale-105 hover:-translate-y-1"
              to="https://vitess.io/blog/2021-07-08-announcing-vitess-arewefastyet"
              target="__blank"
            >
              Blog post
              <Icon className="text-2xl" icon="bookmark" />
            </Link>

            <Link
              className="bg-black text-white rounded-2xl p-5 flex items-center gap-x-2 duration-300 hover:scale-105 hover:-translate-y-1"
              to="https://github.com/vitessio/arewefastyet"
              target="__blank"
            >
              GitHub
              <Icon className="text-2xl" icon="github" />
            </Link>

            <Link
              className="bg-black text-white rounded-2xl p-5 flex items-center gap-x-2 duration-300 hover:scale-105 hover:-translate-y-1"
              to="https://www.vitess.io"
              target="__blank"
            >
              Vitess
              <Icon className="text-2xl" icon="vitess" />
            </Link>
          </div>
        </div>
      </div>
      <div className="basis-1/3 flex justify-start items-center">
        <img
          src="/logo.png"
          alt="logo"
          className="w-11/12"
          style={{
            filter: "drop-shadow(10px 10px 20px rgb(var(--color-primary))",
          }}
        />
      </div>
    </section>
  );
}
