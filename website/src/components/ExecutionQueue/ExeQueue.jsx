import React from 'react';

import '../ExecutionQueue/exeQueue.css'
const ExeQueue = ({data}) => {
    return (
        <tr>
            <td><a target='_blank' href={`https://github.com/vitessio/vitess/commit/${data.git_ref}`}>{data.git_ref.slice(0, 8)}</a></td>
            <td>{data.source}</td>
            <td>{data.type_of.slice(0, 4)}</td>
            <td>{data.pull_nb}</td>
        </tr>
    );
};

export default ExeQueue;