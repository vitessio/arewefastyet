import React from "react";
import { Link } from "react-router-dom";

export default function Hero() {
  return (
    <section className="flex items-center relative h-screen p-page">
      <div className="absolute-cover overflow-hidden -z-1">
        <div
          className="absolute-cover bg-gradient-to-br from-primary to-accent scale-150"
          style={{ clipPath: "polygon(0% 0%, 35% 0%, 62% 100%, 0% 100%)" }}
        />
      </div>
      <div className="flex flex-1">
        <div className="flex flex-col gap-y-8 flex-1 ml-8">
          <div className="flex flex-col bg-background flex-1 p-10 rounded-3xl gap-y-2">
            <h3 className="text-4xl font-semibold">Vitess Introduces</h3>
            <h1 className="text-6xl font-bold text-primary">arewefastyet</h1>
            <p className="text-lg font-normal mt-5 whitespace-nowrap">
              A Benchmarking System for Vitess
            </p>
          </div>
          <div className="flex gap-x-8">
            <Link
              className="bg-black text-white rounded-2xl p-5 flex items-center gap-x-2"
              to="https://vitess.io/blog/2021-07-08-announcing-vitess-arewefastyet"
              target="__blank"
            >
              Read our blog post <i className="fa-solid fa-bookmark"></i>
            </Link>
            <Link
              className="bg-black text-white rounded-2xl p-5 flex items-center gap-x-2"
              to="https://github.com/vitessio/arewefastyet"
              target="__blank"
            >
              Contribute on GitHub
              <i className="fa-brands fa-github"></i>
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
