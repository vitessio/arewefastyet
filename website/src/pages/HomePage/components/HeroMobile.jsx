import React from "react";
import { Link } from "react-router-dom";
import Icon from "../../../common/Icon";

const HeroMobile = () => {
  return (
    <section className="md:hidden h-[90vh] w-full flex flex-col justify-center items-center gap-5">
      <Link
        className="bg-front bg-opacity-30 mb-4 no-underline group cursor-pointer relative shadow-2xl shadow-zinc-900 rounded-full p-px text-xs font-semibold leading-6  text-white inline-block"
        to="https://vitess.io/blog/2021-07-08-announcing-vitess-arewefastyet"
        target="__blank"
      >
        <div className="relative flex space-x-2 items-center z-10 rounded-full bg-zinc-950 py-0.5 px-4 ring-1 ring-white/10 ">
          <span>{`Read the Announcement Blog`}</span>
          <svg
            width="16"
            height="16"
            viewBox="0 0 24 24"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path
              stroke="currentColor"
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth="1.5"
              d="M10.75 8.75L14.25 12L10.75 15.25"
            ></path>
          </svg>
        </div>
        <span className="absolute -bottom-0 left-[1.125rem] h-px w-[calc(100%-2.25rem)] bg-gradient-to-r from-[rgba(236,132,3,0)] via-[rgba(236,132,3)] to-rgba(236,132,3,0) transition-opacity duration-500 group-hover:opacity-40"></span>
      </Link>
      <img
        src="/logo.png"
        alt="logo"
        className="w-[70vw]"
        style={{
          filter: "drop-shadow(10px 10px 20px rgb(var(--color-primary))",
        }}
      />
      <div className="flex flex-col justify-center mt-4 gap-y-2">
        <h1 className="text-3xl text-center font-bold text-primary">
          arewefastyet
        </h1>
        <p className="text-lg font-normal mt-5 whitespace-nowrap">
          A Benchmarking System for Vitess
        </p>
      </div>
      <div className="flex gap-x-8">
        <Link
          className="bg-primary text-black rounded-2xl p-3 flex items-center gap-x-2 duration-300 hover:scale-105 hover:-translate-y-1"
          to="https://github.com/vitessio/arewefastyet"
          target="__blank"
        >
          GitHub
          <Icon className="text-2xl" icon="github" />
        </Link>

        <Link
          className="bg-primary text-black rounded-2xl p-3 flex items-center gap-x-2 duration-300 hover:scale-105 hover:-translate-y-1"
          to="https://www.vitess.io"
          target="__blank"
        >
          Vitess
          <Icon className="text-2xl" icon="vitess" />
        </Link>
      </div>
    </section>
  );
};

export default HeroMobile;
