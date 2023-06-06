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

import './status.css'

import PreviousExe from '../../components/PreviousExecutions/PreviousExe';
import ExeQueue from '../../components/ExecutionQueue/ExeQueue';
import RingLoader from "react-spinners/RingLoader";

const Status = () => {

  const [isLoading, setIsLoading] = useState(true)
  const [dataQueue, setDataQueue] = useState([]);
  const [dataPreviousExe, setDataPreviousExe] = useState([]);
  
  useEffect(() => {
    const fetchData = async () => {
      try {
        const responseQueue = await fetch('http://localhost:9090/api/queue');
        const responsePreviousExe = await fetch('http://localhost:9090/api/recent');
  
        const jsonDataQueue = await responseQueue.json();
        const jsonDataPreviousExe = await responsePreviousExe.json();
        
        setDataQueue(jsonDataQueue);
        setDataPreviousExe(jsonDataPreviousExe);
        setIsLoading(false)
      } catch (error) {
        console.log('Erreur lors de la récupération des données de l\'API', error);
      }
    };
  
    fetchData();
  }, []);
  



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
                  {/* EXECUTION QUEUE */}

                <article className='queue'>
                    <h3>Executions Queue</h3>
                    <table>
                      <thead>
                        <tr>
                          <th>SHA</th>
                          <th>Source</th>
                          <th>Type</th>
                          <th>Pull Request</th>
                        </tr>
                      </thead>
                      <tbody>
                          {dataQueue.map((queue,index) => {
                            return (
                              <ExeQueue data={queue} key={index}/>
                            )
                          })}
                      </tbody>
                    </table>
                </article>
                <figure className='line'></figure>

                  {/*PREVIOUS EXECUTIONS*/}

                <article className='previousExe'>
                    <h3>Previous Execution</h3>
                    <table>
                        <thead className='previousExe__thead'>
                            <tr className='previousExe__thead__tr '>
                                <th>UUID</th>
                                <th>SHA</th>
                                <th>Source</th>
                                <th>Started</th>
                                <th>Finished</th>
                                <th>Type</th>
                                <th>PR</th>
                                <th>Go Version</th>
                                <th>Status</th>
                            </tr>
                        </thead>
                        <tbody>
                            {dataPreviousExe.map((previousExe,index) => {
                              return (
                                <PreviousExe data={previousExe} key={index}/>
                              )
                            })}
                        </tbody>
                    </table>
                </article>
             </>

             )}
        </div>
        
    );
};

export default Status;