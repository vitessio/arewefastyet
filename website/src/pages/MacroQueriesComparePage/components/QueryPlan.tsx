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
import ReactJson from "react-json-pretty";

import "react-json-pretty/themes/monikai.css";
import { twMerge } from "tailwind-merge";
import useModal from "../../../hooks/useModal";
import Icon from "../../../common/Icon";
import { Button } from "@/components/ui/button";

export default function QueryPlan({ data }) {
  const modal = useModal();

  return (
    <div
      className="border border-front p-5 rounded-xl bg-accent shadow-md cursor-pointer duration-300"
      onClick={() => modal.show(<DetailsModal data={data} />)}
    >
      <span className="text-lg border-l-primary border-l-4 py-2 pr-8 rounded-r-xl pl-4 bg-accent tracking-wider">
        {data.Key}
      </span>
      <div className="flex gap-x-3 mt-5">
        {data.Left && !data.Right && <Badge type="warning" content="Only A" />}
        {!data.Left && data.Right && <Badge type="warning" content="Only B" />}

        {data.Left && data.Left.Value.Errors > 0 && (
          <Badge type="error" content="A has Errors" />
        )}
        {data.Right && data.Right.Value.Errors > 0 && (
          <Badge type="error" content="B has Errors" />
        )}

        {data.ExecTimeDiff > 5 && <Badge type="info" content="A fastest" />}
        {data.ExecTimeDiff < -5 && <Badge type="info" content="B fastest" />}
      </div>
    </div>
  );
}

function Badge({ type, content }) {
  return (
    <p
      className={twMerge(
        "px-5 py-1 rounded-full w-max text-sm text-white font-medium",
        type === "error" && "bg-red-600",
        type === "warning" && "bg-purple-700",
        type === "info" && "bg-sky-600"
      )}
    >
      {content}
    </p>
  );
}

function DetailsModal({ data }) {
  const bothPlansExist = data.Left && data.Right;
  const arePlansDifferent =
    bothPlansExist &&
    data.Left.Value.Instructions != data.Right.Value.Instructions;

  const modal = useModal();

  return (
    <div className="w-[80vw] h-[80vh] overflow-y-auto overflow-x-hidden p-5 rounded-xl bg-background border border-front flex flex-col overscroll-contain relative pb-10">
      <Button
        variant="outline"
        className="sticky self-end top-2 right-2 w-max duration-300 hover:scale-125 hover:text-primary text-3xl"
        onClick={modal.hide}
      >
        <Icon icon="close" />
      </Button>

      <div className="flex flex-col items-center gap-y-6">
        <h3 className="text-primary text-3xl my-2 font-medium">Statistics</h3>

        <div className="flex gap-x-4 items-center">
          <span className="font-bold">Query</span>
          <span className="font-light bg-accent px-4 py-1 rounded-md">
            {data.Key}
          </span>
        </div>

        <table className="w-1/2 text-center my-5 border border-front">
          <thead className="py-2 border-b border-front">
            <tr>
              <th />
              <th>A</th>
              <th>B</th>
              <th>Diff</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td>
                <b>ExecutionsCount</b>
              </td>
              {<td>{data.Left && data.Left.Value.ExecCount}</td>}
              {<td>{data.Right && data.Right.Value.ExecCount}</td>}
              <td>{data.ExecCountDiff}</td>
            </tr>
            <tr>
              <td>
                <b>Execution Time</b>
              </td>
              {<td>{data.Left && data.Left.Value.ExecTime}</td>}
              {<td>{data.Right && data.Right.Value.ExecTime}</td>}
              <td>{data.ExecTimeDiff}</td>
            </tr>
            <tr>
              <td>
                <b>Rows Returned</b>
              </td>
              {<td>{data.Left && data.Left.Value.RowsReturned}</td>}
              {<td>{data.Right && data.Right.Value.RowsReturned}</td>}
              <td>{data.RowsReturnedDiff}</td>
            </tr>
            <tr>
              <td>
                <b>Errors</b>
              </td>
              {<td>{data.Left && data.Left.Value.Errors}</td>}
              {<td>{data.Right && data.Right.Value.Errors}</td>}
              <td>{data.ErrorsDiff}</td>
            </tr>
          </tbody>
        </table>
      </div>

      {bothPlansExist && (
        <div
          className={twMerge(
            "px-10 self-center text-center py-2 rounded-full w-max my-5",
            arePlansDifferent ? "bg-red-700" : "bg-green-600"
          )}
        >
          <span>
            Plans are {arePlansDifferent ? "Different" : "Equivalent"}
          </span>
        </div>
      )}

      <div className="flex flex-col items-center gap-y-6 mt-10">
        {data.Left && (
          <>
            <h3 className="text-primary text-3xl my-2 font-medium">
              Query Plan
              {bothPlansExist && arePlansDifferent && " A"}
              {!bothPlansExist && " A"}
            </h3>
            <ReactJson
              data={data.Left.Value.Instructions}
              className="json w-11/12 overflow-auto"
            />
          </>
        )}

        {data.Right && (
          <>
            <h3 className="text-primary text-3xl my-2 font-medium">
              Query Plan
              {bothPlansExist && arePlansDifferent && " B"}
              {!bothPlansExist && " B"}
            </h3>
            <ReactJson
              data={data.Right.Value.Instructions}
              className="json w-11/12 overflow-x-auto"
            />
          </>
        )}
      </div>
    </div>
  );
}
