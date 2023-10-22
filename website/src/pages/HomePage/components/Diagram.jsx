import { twMerge } from "tailwind-merge";
import useTheme from "../../../hooks/useTheme";
import React from "react";

export default function Diagram() {
  const theme = useTheme();

  return (
    <section className="p-page flex flex-col items-center my-8">
      <h1 className="bg-primary bg-opacity-20 text-primary text-lg font-semibold px-4 py-2 rounded-full">
        Diagramatic Overview
      </h1>

      <div className="relative duration-1000">
        <img
          className="duration-inherit"
          src="/images/execution-pipeline-dark.png"
          alt="execution pipeline"
        />

        <img
          className={twMerge(
            "absolute-cover duration-inherit",
            theme.current === "dark" && "opacity-0"
          )}
          src="/images/execution-pipeline.png"
          alt="execution pipeline"
        />
      </div>
    </section>
  );
}
