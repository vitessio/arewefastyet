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

import React from 'react';
import { useState, useEffect } from 'react';
import RingLoader from "react-spinners/RingLoader";
import { v4 as uuidv4 } from 'uuid';


import PreviousExe from '../../components/PreviousExecutions/PreviousExe';
import ExeQueue from '../../components/ExecutionQueue/ExeQueue';
import PreviousExeRes from '../../components/PreviousExeResponsive/PreviousExeRes';
import './status.css';





const Status = () => {

  const [isLoading, setIsLoading] = useState(true)
  const [dataQueue, setDataQueue] = useState([]);
  const [dataPreviousExe, setDataPreviousExe] = useState([]);
  const [error, setError] = useState(null);
  const [selectedButtonIndex, setSelectedButtonIndex] = useState(null);

  
  useEffect(() => {
    const fetchData = async () => {
      try {
        const responseQueue = await fetch(`${import.meta.env.VITE_API_URL}queue`);
        const responsePreviousExe = await fetch(`${import.meta.env.VITE_API_URL}recent`);
  
        const jsonDataQueue = await responseQueue.json();
        const jsonDataPreviousExe = await responsePreviousExe.json();
        
        setDataQueue(jsonDataQueue);
        setDataPreviousExe(jsonDataPreviousExe);
        setIsLoading(false)
      } catch (error) {
        console.log('Error while retrieving data from the API', error);
        setError('An error occurred while retrieving data from the API. Please try again.');
        setIsLoading(false);
      }
    };
  
    fetchData();
  }, []);
  

  const handleClick = (index) => {
    setSelectedButtonIndex((prevIndex) => (prevIndex === index ? null : index));
  };
  

    return (
        <div className='status'>

            <article className='status__top justify--content'>
                <div className='status__top__text'>
                    <h2>Status</h2>
                    <span>
                        Lorem ipsum dolor sit amet, consectetur adipiscing elit. Integer a augue mi.
                        Etiam sed imperdiet ligula, vel elementum velit.
                        Phasellus sodales felis eu condimentum convallis.
                        Suspendisse sodales malesuada iaculis. Mauris molestie placerat ex non malesuada.
                        Curabitur eget sagittis eros. Aliquam aliquam sem non tincidunt volutpat. 
                    </span>
                </div>

                <figure className='statusStats'></figure>
            </article>
            <figure className='line'></figure>
            
             {isLoading ? (
              <div className='loadingSpinner'>
                <RingLoader loading={isLoading} color='#E77002' size={300}/>
                </div>
            ): ( 
                  <>
                      {/* EXECUTION QUEUE  */}

                
                  {dataQueue.length > 0 ?(
                    <>
                      <article className='queue'>
                        <h3>Executions Queue</h3>
                        <div className='queue__top flex'>
                            <span className='width--6em'>SHA</span>
                            <span className='width--11em'> Source</span>
                            <span className='width--11em'>Type</span>
                            <span className='width--5em'>Pull Request</span>
                        </div>
                        <figure className='queue__top__line'></figure>
                              {dataQueue.map((queue,index) => {
                                return (
                                  <ExeQueue data={queue} key={index}/>
                                )
                              })}
                          
                      </article>
                      <figure className='line'></figure>
                  </>
                    ) : null }
                  
                  

                    {/* PREVIOUS EXECUTIONS */}
                  
                  {dataPreviousExe.length > 0 ?(
                      <article className='previousExe'>
                      <h3> Previous Executions</h3>
                      <div className='previousExe__top flex'>
                        <span className='width--6em'>UUID</span>
                        <span className='width--6em'>SHA</span>
                        <span className='width--11em'>Source</span>
                        <span className='width--11em'>Started</span>
                        <span className='width--11em'>Finished</span>
                        <span className='width--11em'>Type</span>
                        <span className='width--5em'>PR</span>
                        <span className='width--6em'>Go version</span>
                        <span className='width--6em'>Status</span>
                      </div>
                      <figure className='previousExe__top__line'></figure>
                      
                            {dataPreviousExe.map((previousExe, index) => {
                                  const isEvenIndex = index % 2 === 0;
                                  const backgroundGrey = isEvenIndex ? 'grey--background' : '';
  
                                  return ( 
                                    <PreviousExe data={previousExe} key={index} className={backgroundGrey}/>
                                  )})}
  
                      </article>
                  ) : null }
                  
                  </>
              )}
                 {error ? <div className='apiError'>{error}</div> : null}
              
        </div>
        
    );
};

export default Status;