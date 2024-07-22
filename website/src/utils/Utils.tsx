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

// BACKGROUND STATUS
export const getStatusClass = (status) => {
  if (status != "finished" && status != "failed" && status != "started") {
    return "default";
  }
  return status;
};

// FORMATDATE
export const formatDate = (date) => {
  if (!date || date === 0) return null;

  date = new Date(date);
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, "0");
  const day = String(date.getDate()).padStart(2, "0");
  const hours = String(date.getHours()).padStart(2, "0");
  const minutes = String(date.getMinutes()).padStart(2, "0");

  return `${month}/${day}/${year} ${hours}:${minutes}`;
};

//FORMATTING BYTES TO Bytes
export const formatByte = (byte: number) => {
  const byteValue = bytes(byte);
  if (byteValue === null) {
    return "0";
  }
  return byteValue;
};

export function fixed(value: number, fractionDigits: number): string {
  if (value === null || typeof value === "undefined" || isNaN(value) || value === 0) {
    return "0";
  }
  return value.toFixed(fractionDigits);
}

export function secondToMicrosecond(value: number): string {
  return fixed(value * 1000000, 2) + "μs";
}

//ERROR API MESSAGE ERROR

export const errorApi: string =
  "An error occurred while retrieving data from the API. Please try again.";

//NUMBER OF PIXELS TO OPEN AND CLOSE THE DROP-DOWN
export const openDropDownValue = 1000;
export const closeDropDownValue = 58;

// OPEN DROP DOWN

export const openDropDown = (currentValue, setOpenDropDown) => {
  if (currentValue === closeDropDownValue) {
    setOpenDropDown(openDropDownValue);
  } else {
    setOpenDropDown(closeDropDownValue);
  }
};

// CHANGE VALUE DROPDOWN

export const valueDropDown = (
  ref,
  setDropDown,
  setCommitHash,
  setOpenDropDown,
  setChangeUrl
) => {
  setDropDown(ref.Name);
  setCommitHash(ref.CommitHash);
  setOpenDropDown(closeDropDownValue);
};

// updateCommitHash: This function updates the value of CommitHash based on the provided Git reference and JSON data.
export const updateCommitHash = (gitRef, setCommitHash, jsonDataRefs) => {
  const obj = jsonDataRefs.find((item) => item.Name === gitRef);
  setCommitHash(obj ? obj.CommitHash : null);
};

////THE NUMBER OF PIXELS THAT ARE USED TO OPEN AND CLOSE THE PREVIOUS EXECUTIONS AND MICROBENCH TABLES
export const openTables = 400;
export const closeTables = 70;

export function formatGitRef(gitRef: string): string {
  return gitRef.slice(0, 8);
}