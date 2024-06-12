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

import bytes from "bytes";

// Types for the status of background
type BackgroundStatus = "finished" | "failed" | "started" | "default";

export const getStatusClass = (status: BackgroundStatus): BackgroundStatus => {
  if (status !== "finished" && status !== "failed" && status !== "started") {
    return "default";
  }
  return status;
};
export const formatDate = (date: string | undefined): string | null => {
  if (!date) return null;

  const parsedDate = new Date(date);
  const year = parsedDate.getFullYear();
  const month = String(parsedDate.getMonth() + 1).padStart(2, "0");
  const day = String(parsedDate.getDate()).padStart(2, "0");
  const hours = String(parsedDate.getHours()).padStart(2, "0");
  const minutes = String(parsedDate.getMinutes()).padStart(2, "0");

  return `${month}/${day}/${year} ${hours}:${minutes}`;
};

// Formatting BYTES TO Bytes
export const formatByte = (byte: number): string => {
  const byteValue = bytes(byte);
  if (byteValue === null) {
    return "0";
  }
  return byteValue.toString();
};

export const fixed = (value: number | null, f: number): string => {
  if (value === null || typeof value === "undefined") {
    return "0";
  }
  return value.toFixed(f);
};

export const secondToMicrosecond = (value: number): string => {
  return fixed(value * 1000000, 2) + "μs";
};

// Error API message
export const errorApi: string =
  "An error occurred while retrieving data from the API. Please try again.";

// Number of pixels to open and close the drop-down
export const openDropDownValue: number = 1000;
export const closeDropDownValue: number = 58;

// Open Drop Open
export const openDropDown = (
  currentValue: number,
  setOpenDropDown: React.Dispatch<React.SetStateAction<number>>
): void => {
  if (currentValue === closeDropDownValue) {
    setOpenDropDown(openDropDownValue);
  } else {
    setOpenDropDown(closeDropDownValue);
  }
};

// Change Value Dropdown
export const valueDropDown = (
  ref: { Name: string; CommitHash: string },
  setDropDown: React.Dispatch<React.SetStateAction<string>>,
  setCommitHash: React.Dispatch<React.SetStateAction<string | null>>,
  setOpenDropDown: React.Dispatch<React.SetStateAction<number>>,
  setChangeUrl?: React.Dispatch<React.SetStateAction<boolean>>
): void => {
  setDropDown(ref.Name);
  setCommitHash(ref.CommitHash);
  setOpenDropDown(closeDropDownValue);
};

// updateCommitHash: This function updates the value of CommitHash based on the provided Git reference and JSON data.
export const updateCommitHash = (
  gitRef: string,
  setCommitHash: React.Dispatch<React.SetStateAction<string | null>>,
  jsonDataRefs: { Name: string; CommitHash: string }[]
): void => {
  const obj = jsonDataRefs.find((item) => item.Name === gitRef);
  setCommitHash(obj ? obj.CommitHash : null);
};

// Number of pixels to open and close the previous executions and microbench tables
export const openTables: number = 400;
export const closeTables: number = 70;
