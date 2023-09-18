import React from "react";

export default function DisplayTable(props) {
  const { data } = props;

  return (
    <table className="flex-1">
      <thead>
        <tr>
          {Object.keys(data[0]).map((item, key) => (
            <th className="border-b font-semibold py-2 border-front" key={key}>
              {item}
            </th>
          ))}
        </tr>
      </thead>
      <tbody>
        {data.map((val, key) => {
          return (
            <tr
              key={key}
              className="hover:bg-foreground hover:bg-opacity-10 cursor-default"
            >
              {Object.keys(val).map((item, key) => (
                <td className="border-b py-3 text-center" key={key}>
                  {val[item]}
                </td>
              ))}
            </tr>
          );
        })}
      </tbody>
    </table>
  );
}
