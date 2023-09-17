import React, { useEffect, useState } from "react";
import { twMerge } from "tailwind-merge";

export default function Diagram() {
  const [dark, setDark] = useState(false);

  useEffect(() => {
    if (document.documentElement.getAttribute("data-theme") == "dark") {
      setDark(true);
    }
  }, []);

  return (
    <section className="p-page relative">
      <img
        className={twMerge("duration-500", dark ? "opacity-0" : "opacity-100")}
        src="/images/execution-pipeline.png"
        alt="execution pipeline"
      />
      <img
        className={twMerge("absolute-cover duration-500 z-1", !dark ? "opacity-0" : "opacity-100")}
        src="/images/execution-pipeline-dark.png"
        alt="execution pipeline"
      />
    </section>
  );
}
