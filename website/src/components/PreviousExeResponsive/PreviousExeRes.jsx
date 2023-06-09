import React from 'react';
import moment from 'moment';

import '../PreviousExeResponsive/previousExeRes.css'

const PreviousExeRes = ({data}) => {

    // Create a Moment object from the given date string
    const startedDate = moment(data.started_at)
    const finishedDate = moment(data.finished_at)
    
    // Format the date using the format method of Moment.js
    // Here, we use 'DD/MM/YYYY HH:mm:ss' as the desired format
    const formattedStartedDate = startedDate.format('MM/DD/YYYY HH:mm')
    const formattedFinishedDate = finishedDate.format('MM/DD/YYYY HH:mm')

    return (
        <>
        <tr className='previousExeRes' >
            <td colSpan="5" >
                <span className='previousExeRes__containter'>
                    <span className='previousExeRes__span'>UUID <i className="fa-solid fa-arrow-right"></i> {data.uuid.slice(0, 8)} <br/><br/></span>
                    <span className='previousExeRes__span'>SHA <i className="fa-solid fa-arrow-right"></i> <a target='_blank' href={`https://github.com/vitessio/vitess/commit/${data.git_ref}`}>{data.git_ref.slice(0,6)}</a> <br/><br/></span>
                    <span className='previousExeRes__span'>Strated <i className="fa-solid fa-arrow-right"></i> {formattedStartedDate} <br/><br/></span>
                    <span className='previousExeRes__span'>Finished <i className="fa-solid fa-arrow-right"></i> {formattedFinishedDate} <br/><br/></span>
                    <span className='previousExeRes__span'>Type <i className="fa-solid fa-arrow-right"></i> {data.type_of} <br/><br/></span>
                    <span className='previousExeRes__span'>PR <i className="fa-solid fa-arrow-right"></i> {data.pull_nb} <br/><br/></span>
                    <span className='previousExeRes__span'>Go version <i className="fa-solid fa-arrow-right"></i> {data.golang_version} <br/><br/></span>

                </span>
            </td>
        </tr>
        
        </>
    );
};

export default PreviousExeRes;