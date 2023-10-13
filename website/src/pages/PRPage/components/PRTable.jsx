
import React, { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { formatDate } from "../../../utils/Utils";
import DisplayList from "../../../common/DisplayList";

export default function PRTable(props) {
  const { data } = props;

  const [PRList, setPRList] = useState([]);

  useEffect(() => {
    console.log(data)
    for (const entry of data) {
      const newData = {};

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

      newData["Opened At"] = formatDate(entry.CreatedAt)

      setPRList((p) => [...p, newData]);
    }
  }, []);

  return (
    <section className="p-page mt-20">
      {PRList.length > 0 && <DisplayList data={PRList} />}
    </section>
  );
}
