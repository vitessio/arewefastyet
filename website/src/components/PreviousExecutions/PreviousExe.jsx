import React from 'react';
import moment from 'moment';

import '../PreviousExecutions/previousexe.css'

const PreviousExe = ({data}) => {
        /*BACKGROUND STATUS*/ 

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

        //Create a Moment object from the given date string
        
        const startedDate = moment(data.started_at)
        const finishedDate = moment(data.finished_at)
        // Format the date using the format method of Moment.js
        // Here, we use 'DD/MM/YYYY HH:mm:ss' as the desired format

        const formattedStartedDate = startedDate.format('MM/DD/YYYY HH:mm')
        const formattedFinishedDate = finishedDate.format('MM/DD/YYYY HH:mm')
    return (
    <>
        <tr>
            <td>{data.uuid.slice(0, 8)}</td>
            <td><a target='_blank' href={`https://github.com/vitessio/vitess/commit/${data.git_ref}`}>{data.git_ref.slice(0,6)}</a></td>
            <td>{data.source}</td>
            <td>{formattedStartedDate}</td>
            <td>{formattedFinishedDate}</td>
            <td>{data.type_of.slice(0,4)}</td>
            <td className='tdPR'>{data.pull_nb}</td>
            <td>{data.golang_version}</td>
            <td ><span className={`data ${getStatusClass(data.status)} tdStatus`}>{data.status}</span></td>
        </tr>
        
        
    </>
    );
};

export default PreviousExe;