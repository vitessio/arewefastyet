import React, {useState} from 'react';
import moment from 'moment';

import './previousExeResMobile.css'

const PreviousExeRes = ({data, className}) => {

    const getStatusClass = (status) => {
        switch (status) {
            case 'finished':
              return 'finished';
            case 'failed':
              return 'failed';
            case 'started':
              return 'started';
            default:
              return 'default';
          }
    }


    // Create a Moment object from the given date string
    const startedDate = moment(data.started_at)
    const finishedDate = moment(data.finished_at)
    
    // Format the date using the format method of Moment.js
    // Here, we use 'DD/MM/YYYY HH:mm:ss' as the desired format
    const formattedStartedDate = startedDate.format('MM/DD/YYYY HH:mm')
    const formattedFinishedDate = finishedDate.format('MM/DD/YYYY HH:mm')


    const [maxHeight, setMaxHeight] = useState(70);

    const handleClick = () => {
        if (maxHeight === 70) {
            setMaxHeight(400);
          } else {
            setMaxHeight(70);
          }
        };

    return (
        <div className={`previousExe__data__mobile ${className}`} style={{ maxHeight: `${maxHeight}px` }}>
            <div className='previousExe__data__mobile__top flex'>
                <span className='width--11em'>{data.source}</span>
                <span  className={`data ${getStatusClass(data.status)} spanStatus width--6em`}>{data.status}</span>
                <span className='width--3em'><i className="fa-solid fa-circle-info" onClick={handleClick}></i></span>
            </div>
            <div className='previousExe__data__mobile__bottom flex'>
                <div className='previousExe__data__mobile__bottom__more flex'>
                    <span>UUID</span>
                    <i className="fa-solid fa-arrow-right"></i>
                    <span >{data.uuid.slice(0, 8)}</span>
                </div>
                <div className='previousExe__data__mobile__bottom__more flex'>
                    <span>SHA</span>
                    <i className="fa-solid fa-arrow-right"></i>
                    <span ><a target='_blank' href={`https://github.com/vitessio/vitess/commit/${data.git_ref}`}>{data.git_ref.slice(0,6)}</a></span>
                </div>
                <div className='previousExe__data__mobile__bottom__more flex'>
                    <span>Type</span>
                    <i className="fa-solid fa-arrow-right"></i>
                    <span>{data.type_of}</span>
                </div>
                <div className='previousExe__data__mobile__bottom__more flex'>
                    <span>Started</span>
                    <i className="fa-solid fa-arrow-right"></i>
                    <span>{formattedStartedDate}</span>
                </div>
                <div className='previousExe__data__mobile__bottom__more flex'>
                    <span>Finished</span>
                    <i className="fa-solid fa-arrow-right"></i>
                    <span>{formattedFinishedDate}</span>
                </div>
                <div className='previousExe__data__mobile__bottom__more flex'>
                    <span>PR</span>
                    <i className="fa-solid fa-arrow-right"></i>
                    <span>{data.pull_nb}</span>
                </div>
                <div className='previousExe__data__mobile__bottom__more flex'>
                    <span>Go version</span>
                    <i className="fa-solid fa-arrow-right"></i>
                    <span>{data.golang_version}</span>
                </div>
            </div>
         
        </div>
    );
};

export default PreviousExeRes;