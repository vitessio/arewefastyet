import React from 'react';

import '../PreviousExecutions/previousexe.css'

const PreviousExe = ({data}) => {

    const getStatusClass = (status) => {
        if (status === 'finished') {
          return 'finished';
        } else if (status === 'failed') {
          return 'failed';
        } else if (status ==='started') {
            return 'started';
        } else {
            return 'default';
        }
    }
    return (
    <>
        <tr>
            <td>{data.uuid.slice(0, 8)}</td>
            <td><a target='_blank' href={`https://github.com/vitessio/vitess/commit/${data.git_ref}`}>{data.git_ref.slice(0,6)}</a></td>
            <td>{data.source}</td>
            <td>{data.started_at}</td>
            <td>{data.finished_at}</td>
            <td>{data.type_of.slice(0,4)}</td>
            <td className='tdPR'>{data.pull_nb}</td>
            <td>{data.golang_version}</td>
            <td ><span className={`data ${getStatusClass(data.status)} tdStatus`}>{data.status}</span></td>
        </tr>
        
        
    </>
    );
};

export default PreviousExe;