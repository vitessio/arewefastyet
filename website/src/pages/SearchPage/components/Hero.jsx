import React from "react";

export default function Hero({setGitRef}) {
  return (
    <section className="h-[40vh] flex justify-center items-center">
      <div className="p-[3px] bg-gradient-to-br from-primary to-accent rounded-full w-1/2 duration-300 focus-within:p-[1px]">
        <input
          type="text"
          onChange={(e) => setGitRef(e.target.value)}
          className="px-4 py-2 relative text-lg bg-background rounded-inherit w-full text-center focus:outline-none duration-inherit focus:border-none"
          placeholder="Enter commit SHA"
        />
      </div>
    </section>
  );
}
