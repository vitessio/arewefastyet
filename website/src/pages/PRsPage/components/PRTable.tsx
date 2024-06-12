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

import React, { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { formatDate } from "../../../utils/Utils";
import DisplayList from "../../../common/DisplayList";
import Icon from "../../../common/Icon";
import { PRTableType, prDataTypes } from "@/types";

/**
 * Displays a table of pull requests (PRs).
 * @param {Object} props - The props for the PRTable component.
 * @param {prDataTypes[]} props.data - The array of pull request data.
 * @returns {JSX.Element} - The rendered JSX element.
 */

export default function PRTable({ data }: { data: prDataTypes[] }): JSX.Element {
  const [PRList, setPRList] = useState<PRTableType[]>([]);

  useEffect(() => {
    const updatedPRList = data.map((entry) => {
      const newData: PRTableType = {};

      newData["#"] = (
        <Link
          target="__blank"
          className="text-primary text"
          to={`https://github.com/vitessio/vitess/pull/${entry.ID}`}
        >
          {entry.ID}
        </Link>
      );

      newData["Title"] = entry.Title;

      newData["Author"] = (
        <Link
          target="__blank"
          className="text-primary text font-bold"
          to={`https://github.com/${entry.Author}`}
        >
          {entry.Author}
        </Link>
      );

      newData["Opened At"] = formatDate(entry.CreatedAt);

      newData["Details"] = (
        <Link
          to={`/pr/${entry.ID}`}
          className="flex justify-center text-lg text-primary duration-300 hover:scale-105"
        >
          <Icon icon="open_in_new" />
        </Link>
      );

      return newData;
    });

    setPRList(updatedPRList);
  }, [data]);

  return (
    <section className="p-page mt-20 flex flex-col overflow-y-scroll">
      {PRList.length > 0 && <DisplayList data={PRList} />}
    </section>
  );
}
