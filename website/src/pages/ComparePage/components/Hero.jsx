import React from "react";
import { twMerge } from "tailwind-merge";

export default function Hero(props) {
  const { gitRef, setGitRef } = props;

  return (
    <div className="flex flex-col h-[35vh] justify-center items-center">
      <h1 className="mb-3 text-front text-opacity-70">Enter SHAs to compare commits</h1>
      <div className="flex overflow-hidden bg-gradient-to-br from-primary to-accent p-[2px] rounded-full">
        <ComparisonInput
          name="left"
          className="rounded-l-full"
          setGitRef={setGitRef}
        />
        <ComparisonInput
          name="right"
          className="rounded-r-full "
          setGitRef={setGitRef}
        />
      </div>
    </div>
  );
}

function ComparisonInput(props) {
  const { className, setGitRef, name } = props;

  return (
    <input
      type="text"
      name={name}
      className={twMerge(
        className,
        "relative text-xl px-6 py-2 bg-background focus:border-none focus:outline-none border border-primary"
      )}
      placeholder={`${name} commit SHA`}
      onChange={(event) =>
        setGitRef((p) => {
          return { ...p, [name]: event.target.value };
        })
      }
    />
  );
}
