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

import React, { useState, useEffect } from 'react';
import RingLoader from "react-spinners/RingLoader";

import '../CRON/cron.css'

import { errorApi} from '../../utils/Utils';
import CronSummary from '../../components/CRONComponents/CRONSummary/CronSummary';

const CRON = () => {

    const [dataCronSummary, setDataCronSummary] = useState([])
    const [error, setError] = useState(null);
    const [isLoading, setIsLoading] = useState(true);

    useEffect(() => {
        const fetchData = async () => {
          try {
            const responseCron = await fetch(`${import.meta.env.VITE_API_URL}cron`);
      
            const jsonDataCron = await responseCron.json();
            console.log(jsonDataCron)
            setDataCronSummary(jsonDataCron);
            setIsLoading(false)
            

          } catch (error) {
            console.log('Error while retrieving data from the API', error);
            setError(errorApi);
            setIsLoading(false);
          }
        };
      
        fetchData();
      }, []);

    return (
        <div className='cron'>
            <div className='cron__top'>
            <h2>CRON</h2>
                    <span>
                            Lorem ipsum dolor sit amet, consectetur adipiscing elit. Integer a augue mi.
                            Etiam sed imperdiet ligula, vel elementum velit.
                            Phasellus sodales felis eu condimentum convallis.
                            Suspendisse sodales malesuada iaculis. Mauris molestie placerat ex non malesuada.
                            Curabitur eget sagittis eros. Aliquam aliquam sem non tincidunt volutpat. 
                    </span>
            </div>
            <figure className='line'></figure>
            {error ? (
                    <div className='macrobench__apiError'>{error}</div> 
                ) : (
                    isLoading ? (
                        <div className='loadingSpinner'>
                            <RingLoader loading={isLoading} color='#E77002' size={300}/>
                            </div>
                        ): ( 
                            <>
                                <div className='cron__summary__container justify--content'>
                                    {dataCronSummary.map((cronSummary, index) => {
                                        return (
                                            <CronSummary key={index} data={cronSummary}/>
                                        )
                                    })}
                                </div>
                                <figure className='line'></figure>
                            </>
                    )
                )}
        </div>
    );
};

export default CRON;