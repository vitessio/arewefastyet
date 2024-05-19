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

export default function DisplayList(props) {
  const { data } = props;

  return (
    <div className="">
      <div className="widescreen:hidden">
        {data.map((val, key) => (
          <div
            key={key}
            className="cursor-default mb-4 p-4 rounded-lg border border-front border-[1px]"
          >
            {Object.keys(val).map((item, itemKey) => (
              <div
                key={itemKey}
                className="flex justify-between border-b py-2 last:border-b-0 hover:bg-accent "
              >
                <span className="font-semibold">{item}</span>
                <span>{val[item]}</span>
              </div>
            ))}
          </div>
        ))}
      </div>
      <div className="w-full overflow-x-auto">
        <table className="w-full">
          <thead>
            <tr>
              {Object.keys(data[0]).map((item, key) => (
                <th
                  className="border-b min-w-[150px] text-left font-semibold py-2 border-front"
                  key={key}
                >
                  {item}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {data.map((val, key) => {
              return (
                <tr key={key} className="hover:bg-accent cursor-default">
                  {Object.keys(val).map((item, key) => (
                    <td className="border-b py-3 text-left" key={key}>
                      {val[item]}
                    </td>
                  ))}
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>
    </div>
  );
}
