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

export function formatDate(_date: string | number) {
  if (!_date || _date === 0) return null;

  const date = new Date(_date);
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, "0");
  const day = String(date.getDate()).padStart(2, "0");
  const hours = String(date.getHours()).padStart(2, "0");
  const minutes = String(date.getMinutes()).padStart(2, "0");

  return `${month}/${day}/${year} ${hours}:${minutes}`;
}

//FORMATTING BYTES TO GB
export function formatByteForGB(byte: number) {
  const byteValue = bytes(byte);
  if (byteValue === null) {
    return "0";
  }
  return byteValue.toString();
}

export function fixed(value: number, f: number) {
  if (value === null || typeof value === "undefined") {
    return "0";
  }
  return value.toFixed(f);
}

//ERROR API MESSAGE ERROR
export const errorApi =
  "An error occurred while retrieving data from the API. Please try again.";
