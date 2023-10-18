import React from "react";
import Dropdown from "../../../common/Dropdown";
import { twMerge } from "tailwind-merge";

export default function Hero(props) {
  const { refs, gitRef, setGitRef } = props;

  return (
    <section className="flex flex-col gap-y-[10vh] pt-[10vh] justify-center items-center h-[70vh]">
      <h2 className="text-primary font-medium text-6xl cursor-default">
        Compare Macrobenchmarks
      </h2>
      {refs && refs.length > 0 && (
        <div className="flex gap-x-24">
          <Dropdown.Container
            className="w-[20vw] py-2 border border-primary rounded-md mb-[1px] text-lg shadow-xl"
            defaultIndex={refs.map((r) => r.Name).indexOf(gitRef.left)}
            onChange={(event) => {
              setGitRef((p) => {
                return { ...p, left: event.value };
              });
            }}
          >
            {refs.map((ref, key) => (
              <Dropdown.Option
                key={key}
                className={twMerge(
                  "w-[20vw] relative border-front border border-t-transparent border-opacity-60 bg-background py-2 after:duration-150 after:absolute-cover after:bg-foreground after:bg-opacity-0 hover:after:bg-opacity-10 font-medium",
                  key === 0 && "rounded-t border-t-front",
                  key === refs.length - 1 && "rounded-b"
                )}
              >
                {ref.Name}
              </Dropdown.Option>
            ))}
          </Dropdown.Container>

          <Dropdown.Container
            className="w-[20vw] py-2 border border-primary rounded-md mb-[1px] text-lg shadow-xl"
            defaultIndex={refs.map((r) => r.Name).indexOf(gitRef.right)}
            onChange={(event) => {
              setGitRef((p) => {
                return { ...p, right: event.value };
              });
            }}
          >
            {refs.map((ref, key) => (
              <Dropdown.Option
                key={key}
                className={twMerge(
                  "w-[20vw] relative border-front border border-t-transparent border-opacity-60 bg-background py-2 after:duration-150 after:absolute-cover after:bg-foreground after:bg-opacity-0 hover:after:bg-opacity-10 font-medium",
                  key === 0 && "rounded-t border-t-front",
                  key === refs.length - 1 && "rounded-b"
                )}
              >
                {ref.Name}
              </Dropdown.Option>
            ))}
          </Dropdown.Container>
        </div>
      )}
    </section>
  );
}
