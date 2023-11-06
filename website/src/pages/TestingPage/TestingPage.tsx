import React, { useEffect, useState } from "react";
import api from "../../services/api";

export default function DynamicFunctionCaller() {
  const functions = api;

  const [response, setResponse] = useState<any>("");
  const [loading, setLoading] = useState<boolean>(false);
  const [inputValues, setInputValues] = useState<Record<string, string>>({});

  const functionNames = Object.keys(functions);

  async function handleExecute(action: keyof Omit<typeof functions, "client">) {
    if (functions[action]) {
      const parameters = Object.values(inputValues);
      const result = await (functions[action] as Function)(...parameters);
      setResponse(result);
      setLoading(loading);
    }
  }

  const handleInputChange = (paramName: string, value: string) => {
    setInputValues({ ...inputValues, [paramName]: value });
  };

  return (
    <div className="p-page h-[98vh] py-5 flex flex-col gap-y-5">
      <h2 className="text-2xl font-semibold">Dynamic API Testing</h2>

      <div className="flex flex-col gap-y-3">
        {functionNames.map((action) => (
          <div key={action} className="flex gap-x-2">
            <button
              className="bg-teal-600 px-6 py-2 rounded font-medium text-white"
              onClick={() => handleExecute(action as any)}
            >
              {action}
            </button>
            {getParamNames(
              functions[action as keyof Omit<typeof functions, "client">]
            ).map((param: string, index: number) => {
              return (
                <input
                  key={index}
                  type="text"
                  placeholder={`Enter ${param} for ${action}`}
                  onChange={(e) =>
                    handleInputChange(`${action}_${index}`, e.target.value)
                  }
                  className="bg-transparent border border-white border-opacity-20 px-2 rounded"
                />
              );
            })}
          </div>
        ))}
      </div>
      <div className="flex-1 relative">
        <h3>Response:</h3>
        <textarea
          className="h-full w-full bg-black border border-white border-opacity-30 rounded-lg p-3"
          readOnly
          disabled
          value={formatObject(response)}
        />

        {loading && (
          <div className="border-teal-500 border-[5px] animate-spin border-dashed aspect-square w-10 absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 rounded-full" />
        )}
      </div>
    </div>
  );
}

function formatObject(obj: object | Array<any>) {
  return JSON.stringify(obj, null, 4);
}

function getParamNames(func: any) {
  const STRIP_COMMENTS = /((\/\/.*$)|(\/\*[\s\S]*?\*\/))/gm;
  const ARGUMENT_NAMES = /([^\s,]+)/g;
  const fnStr = func.toString().replace(STRIP_COMMENTS, "");
  var result: any = fnStr
    .slice(fnStr.indexOf("(") + 1, fnStr.indexOf(")"))
    .match(ARGUMENT_NAMES);
  if (result === null) result = [];
  return result;
}
