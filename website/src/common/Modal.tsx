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
import useModal from "../hooks/useModal";
import { twMerge } from "tailwind-merge";

export default function Modal(): JSX.Element {
  const modal = useModal();

  return (
    <article
      className={twMerge(
        "z-[1000] bg-black backdrop-blur-sm bg-opacity-10 flex justify-center items-center fixed top-0 left-0 w-full h-full duration-300",
        modal.element ? "opacity-100" : "opacity-0 pointer-events-none"
      )}
    >
      <div
        className={twMerge(
          "duration-inherit ease-out",
          !modal.element && " scale-150 translate-y-full opacity-25 blur-md"
        )}
      >
        {modal.element}
      </div>
    </article>
  );
}
