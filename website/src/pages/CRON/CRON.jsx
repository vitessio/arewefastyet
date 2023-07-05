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
import { useNavigate } from 'react-router-dom';
import { ResponsiveLine } from '@nivo/line'


import '../CRON/cron.css'

import { errorApi, formatByteForGB} from '../../utils/Utils';
import CronSummary from '../../components/CRONComponents/CRONSummary/CronSummary';

const CRON = () => {

    const urlParams = new URLSearchParams(window.location.search);
    const [dataCronSummary, setDataCronSummary] = useState([])
    const [dataCron, setDataCron] = useState([])
    const [error, setError] = useState(null);
    const [isLoading, setIsLoading] = useState(true);
    const [isLoading2, setIsLoading2] = useState(true);
    const [benchmarkType, setBenchmarktype] = useState(urlParams.get('type') == null ? '' : urlParams.get('type'))

    useEffect(() => {
        const fetchData = async () => {
          try {
            const responseCronSummary = await fetch(`${import.meta.env.VITE_API_URL}cronsummary`);
      
            const jsonDataCronSummary = await responseCronSummary.json();
            
            setDataCronSummary(jsonDataCronSummary);
            setIsLoading(false)

          } catch (error) {
            console.log('Error while retrieving data from the API', error);
            setError(errorApi);
            setIsLoading(false);
          }
        };
      
        fetchData();
      }, []);

      useEffect(() => {
        const fetchData = async () => {
          try {
            const responseCron = await fetch(`${import.meta.env.VITE_API_URL}cron?type=${benchmarkType}`);
      
            const jsonDataCron = await responseCron.json();
            
            setDataCron(jsonDataCron);
            setIsLoading2(false)
            
          } catch (error) {
            console.log('Error while retrieving data from the API', error);
            setError(errorApi);
            setIsLoading2(false);
          }
        };
      
        fetchData();
      }, [benchmarkType]);

    // Changing the URL relative to the reference of a selected benchmark.
    // Storing the carousel position as a URL parameter.
    const navigate = useNavigate();

      useEffect(() => {
        navigate(`?type=${benchmarkType}`)
    }, [benchmarkType]) 

    
        const TPSData = [{
            id: 'TPS',
            data: [],
        }];

        const QPSData = [{
            id: 'Reads',
            data: [],
        },
        {
            id: 'Total',
            data: [],
        },
        
        {
            id: 'Writes',
            data: [],
        }];              

        const latencyData = [{
            id: 'Latency',
            data: [],
        }];

        const CPUTimeData = [
            {
            id:'Total',
            data:[]
        },
        {
            id:'vtgate',
            data:[]
        },
        {
            id:'vttablet',
            data:[]
        }
    ];

        const MemBytesData = [
            {
            id:'Total',
            data:[]
        },
        {
            id:'vtgate',
            data:[]
        },
        {
            id:'vttablet',
            data:[]
        }
    ]

        for (const item of dataCron) {
            const xValue = item.GitRef.slice(0, 8);
                
            // TPS Data
            
            TPSData[0].data.push({
                x: xValue,
                y: item.Result.tps
            });

            // QPS Data
            
            QPSData[0].data.push({
                x: xValue,
                y: item.Result.qps.reads
            });

            
            QPSData[1].data.push({ 
                x: xValue,
                y: item.Result.qps.total
            });

            
            QPSData[2].data.push({
                x: xValue,
                y: item.Result.qps.writes
            });

            // Latency Data
            
            latencyData[0].data.push({
                x: xValue,
                y: item.Result.latency
            });

            // CPUTime Data

            CPUTimeData[0].data.push({
                x: xValue,
                y: item.Metrics.TotalComponentsCPUTime
            })

            CPUTimeData[1].data.push({
                x: xValue,
                y: item.Metrics.ComponentsCPUTime.vtgate
            })

            CPUTimeData[2].data.push({
                x: xValue,
                y: item.Metrics.ComponentsCPUTime.vttablet
            })

            //MemStatsAllocBytes Data

            MemBytesData[0].data.push({
                x: xValue,
                y: formatByteForGB(item.Metrics.TotalComponentsMemStatsAllocBytes)
            })

            MemBytesData[1].data.push({
                x: xValue,
                y: formatByteForGB(item.Metrics.ComponentsMemStatsAllocBytes.vtgate)
            })

            MemBytesData[2].data.push({
                x: xValue,
                y: formatByteForGB(item.Metrics.ComponentsMemStatsAllocBytes.vttablet)
            })
        }

        
      
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
                                            <CronSummary key={index} data={cronSummary} setBenchmarktype={setBenchmarktype}/>
                                        )
                                    })}
                                </div>
                                <figure className='line'></figure>
                                {isLoading2 ? (
                                        <div className='loadingSpinner'>
                                            <RingLoader loading={isLoading2} color='#E77002' size={300}/>
                                        </div>
                                    ) : (
                                            <div className='cron__container'>
                                                <h3 >QPS (Queries per second)</h3>
                                                <div className='chart'>
                                                    <ResponsiveLine
                                                        data={QPSData}

                                                        height={450}
                                                        colors={['Yellow', 'orange', 'brown']}
                                                        theme={{
                                                            axis:{
                                                                ticks:{
                                                                    text:{
                                                                        fontSize:'13px',
                                                                        fill:'white'
                                                                    }
                                                                }
                                                            },
                                                            legends:{
                                                                text:{
                                                                    fontSize:'14px',
                                                                    fill:'white'
                                                                }
                                                            }
                                                        }}
                                                        tooltip={({ point }) => (
                                                            <div className='tooltip flex'>
                                                                <figure style={{ backgroundColor: point.serieColor }}></figure>
                                                              <div>x : {point.data.x}</div>
                                                              <div>y : {point.data.y}</div>
                                                            </div>
                                                          )}
                                                        areaBaselineValue={50}
                                                        margin={{ top: 50, right: 110, bottom: 50, left: 60 }}
                                                        xScale={{ type: 'point' }}
                                                        yScale={{
                                                            type: 'linear',
                                                            min: '0',
                                                            max: 'auto',
                                                            reverse: false
                                                        }}
                                                        yFormat=" >-.2f"
                                                        axisTop={null}
                                                        axisRight={null}
                                                        pointSize={10}
                                                        isInteractive={true}
                                                        pointBorderWidth={2}
                                                        pointBorderColor={{ from: 'serieColor' }}
                                                        pointLabelYOffset={-12}
                                                        areaOpacity={0.1}
                                                        useMesh={true}
                                                        legends={[
                                                            {
                                                                anchor: 'bottom-right',
                                                                direction: 'column',
                                                                justify: false,
                                                                translateX: 100,
                                                                translateY: 0,
                                                                itemsSpacing: 0,
                                                                itemDirection: 'left-to-right',
                                                                itemWidth: 80,
                                                                itemHeight: 20,
                                                                itemOpacity: 0.75,
                                                                symbolSize: 12,
                                                                symbolShape: 'circle',
                                                                symbolBorderColor: 'rgba(0, 0, 0, .5)',
                                                                effects: [
                                                                    {
                                                                        on: 'hover',
                                                                        style: {
                                                                            itemBackground: 'rgba(0, 0, 0, .03)',
                                                                            itemOpacity: 1
                                                                        }
                                                                    }
                                                                ]
                                                            }
                                                        ]}
                                                        />
                                                </div>
                                                <h3 className='chart__title'>TPS (Transactions per second)</h3>
                                                 <div className='chart'>
                                                    <ResponsiveLine
                                                        data={TPSData}
                                                        height={400}
                                                        colors={['#E77002']}
                                                        theme={{
                                                            axis:{
                                                                ticks:{
                                                                    text:{
                                                                        fontSize:'13px',
                                                                        fill:'white'
                                                                    }
                                                                }
                                                            },
                                                            legends:{
                                                                    text:{
                                                                        fontSize:'14px',
                                                                        fill:'white'
                                                                    }
                                                            }
                                                        }}
                                                        tooltip={({ point }) => (
                                                            <div className='tooltip flex'>
                                                                <figure style={{ backgroundColor: point.serieColor }}></figure>
                                                              <div>x : {point.data.x}</div>
                                                              <div>y : {point.data.y}</div>
                                                            </div>
                                                          )}
                                                        areaBaselineValue={50}
                                                        margin={{ top: 50, right: 110, bottom: 50, left: 60 }}
                                                        xScale={{ type: 'point' }}
                                                        yScale={{
                                                            type: 'linear',
                                                            min: '0',
                                                            max: 'auto',
                                                            stacked: true,
                                                            reverse: false
                                                        }}
                                                        yFormat=" >-.2f"
                                                        axisTop={null}
                                                        axisRight={null}
                                                        pointSize={10}
                                                        isInteractive={true}
                                                        pointBorderWidth={2}
                                                        pointBorderColor={{ from: 'serieColor' }}
                                                        pointLabelYOffset={-12}
                                                        areaOpacity={0.1}
                                                        useMesh={true}
                                                        legends={[
                                                            {
                                                                anchor: 'bottom-right',
                                                                direction: 'column',
                                                                justify: false,
                                                                translateX: 100,
                                                                translateY: 0,
                                                                itemsSpacing: 0,
                                                                itemDirection: 'left-to-right',
                                                                itemWidth: 80,
                                                                itemHeight: 20,
                                                                itemOpacity: 0.75,
                                                                symbolSize: 12,
                                                                symbolShape: 'circle',
                                                                symbolBorderColor: 'rgba(0, 0, 0, .5)',
                                                                effects: [
                                                                    {
                                                                        on: 'hover',
                                                                        style: {
                                                                            itemBackground: 'rgba(0, 0, 0, .03)',
                                                                            itemOpacity: 1
                                                                        }
                                                                    }
                                                                ]
                                                            }
                                                        ]}
                                                        />
                                                </div>                                              
                                                <h3 className='chart__title'>Latency (Milliseconds)</h3>
                                                <div className='chart'>
                                                    <ResponsiveLine
                                                        data={latencyData}
                                                        height={400}
                                                        theme={{
                                                            axis:{
                                                                ticks:{
                                                                    text:{
                                                                        fontSize:'13px',
                                                                        fill:'white'
                                                                    }
                                                                }
                                                            },
                                                            legends:{
                                                                text:{
                                                                    fontSize:'14px',
                                                                    fill:'white'
                                                                }
                                                            }
                                                        }}
                                                        tooltip={({ point }) => (
                                                            <div className='tooltip flex'>
                                                                <figure style={{ backgroundColor: point.serieColor }}></figure>
                                                              <div>x : {point.data.x}</div>
                                                              <div>y : {point.data.y}</div>
                                                            </div>
                                                          )}
                                                        colors={['#E77002']}
                                                        areaBaselineValue={50}
                                                        margin={{ top: 50, right: 110, bottom: 50, left: 60 }}
                                                        xScale={{ type: 'point' }}
                                                        yScale={{
                                                            type: 'linear',
                                                            min: '0',
                                                            max: 'auto',
                                                            stacked: true,
                                                            reverse: false
                                                        }}
                                                        yFormat=" >-.2f"
                                                        axisTop={null}
                                                        axisRight={null}
                                                        pointSize={10}
                                                        isInteractive={true}
                                                        pointBorderWidth={2}
                                                        pointBorderColor={{ from: 'serieColor' }}
                                                        pointLabelYOffset={-12}
                                                        areaOpacity={0.1}
                                                        useMesh={true}
                                                        legends={[
                                                            {
                                                                anchor: 'bottom-right',
                                                                direction: 'column',
                                                                justify: false,
                                                                translateX: 100,
                                                                translateY: 0,
                                                                itemsSpacing: 0,
                                                                itemDirection: 'left-to-right',
                                                                itemWidth: 80,
                                                                itemHeight: 20,
                                                                itemOpacity: 0.75,
                                                                symbolSize: 12,
                                                                symbolShape: 'circle',
                                                                symbolBorderColor: 'rgba(0, 0, 0, .5)',
                                                                effects: [
                                                                    {
                                                                        on: 'hover',
                                                                        style: {
                                                                            itemBackground: 'rgba(0, 0, 0, .03)',
                                                                            itemOpacity: 1
                                                                        }
                                                                    }
                                                                ]
                                                            }
                                                        ]}
                                                        
                                                        />
                                                </div>
                                                <h3 className='chart__title'>CPUTime</h3>
                                                <div className='chart'>
                                                    <ResponsiveLine
                                                        data={CPUTimeData}
                                                        
                                                        height={400}
                                                        colors={['Yellow', 'orange', 'brown']}
                                                        theme={{
                                                            axis:{
                                                                ticks:{
                                                                    text:{
                                                                        fontSize:'13px',
                                                                        fill:'white'
                                                                    }
                                                                }
                                                            },
                                                            legends:{
                                                                    text:{
                                                                        fontSize:'14px',
                                                                        fill:'white'
                                                                    }
                                                            }
                                                        }}
                                                        tooltip={({ point }) => (
                                                            <div className='tooltip flex'>
                                                                <figure style={{ backgroundColor: point.serieColor }}></figure>
                                                              <div>x : {point.data.x}</div>
                                                              <div>y : {point.data.y}</div>
                                                            </div>
                                                          )}
                                                        areaBaselineValue={50}
                                                        margin={{ top: 50, right: 110, bottom: 50, left: 60 }}
                                                        xScale={{ type: 'point' }}
                                                        yScale={{
                                                            type: 'linear',
                                                            min: '0',
                                                            max: 'auto',
                                                            reverse: false
                                                        }}
                                                        yFormat=" >-.2f"
                                                        axisTop={null}
                                                        axisRight={null}
                                                        pointSize={10}
                                                        isInteractive={true}
                                                        pointBorderWidth={2}
                                                        pointBorderColor={{ from: 'serieColor' }}
                                                        pointLabelYOffset={-12}
                                                        areaOpacity={0.1}
                                                        useMesh={true}
                                                        legends={[
                                                            {
                                                                anchor: 'bottom-right',
                                                                direction: 'column',
                                                                justify: false,
                                                                translateX: 100,
                                                                translateY: 0,
                                                                itemsSpacing: 0,
                                                                itemDirection: 'left-to-right',
                                                                itemWidth: 80,
                                                                itemHeight: 20,
                                                                itemOpacity: 0.75,
                                                                symbolSize: 12,
                                                                symbolShape: 'circle',
                                                                symbolBorderColor: 'rgba(0, 0, 0, .5)',
                                                                effects: [
                                                                    {
                                                                        on: 'hover',
                                                                        style: {
                                                                            itemBackground: 'rgba(0, 0, 0, .03)',
                                                                            itemOpacity: 1
                                                                        }
                                                                    }
                                                                ]
                                                            }
                                                        ]}
                                                        />
                                                </div>  
                                                <h3 className='chart__title'>MemStatsAllocBytes</h3>
                                                <div className='chart'>
                                                    <ResponsiveLine
                                                        data={MemBytesData}
                                                        
                                                        height={400}
                                                        colors={['Yellow', 'orange', 'brown']}
                                                        theme={{
                                                            axis:{
                                                                ticks:{
                                                                    text:{
                                                                        fontSize:'13px',
                                                                        fill:'white'
                                                                    }
                                                                }
                                                            },
                                                            legends:{
                                                                    text:{
                                                                        fontSize:'14px',
                                                                        fill:'white'
                                                                    }
                                                            }
                                                        }}
                                                        tooltip={({ point }) => (
                                                            <div className='tooltip flex'>
                                                                <figure style={{ backgroundColor: point.serieColor }}></figure>
                                                              <div>x : {point.data.x}</div>
                                                              <div>y : {point.data.y}</div>
                                                            </div>
                                                          )}
                                                        areaBaselineValue={50}
                                                        margin={{ top: 50, right: 110, bottom: 50, left: 60 }}
                                                        xScale={{ type: 'point' }}
                                                        yScale={{
                                                            type: 'linear',
                                                            min: '0',
                                                            max: 'auto',
                                                            reverse: false
                                                        }}
                                                        yFormat=" >-.2f"
                                                        axisTop={null}
                                                        axisRight={null}
                                                        pointSize={10}
                                                        isInteractive={true}
                                                        pointBorderWidth={2}
                                                        pointBorderColor={{ from: 'serieColor' }}
                                                        pointLabelYOffset={-12}
                                                        areaOpacity={0.1}
                                                        useMesh={true}
                                                        legends={[
                                                            {
                                                                anchor: 'bottom-right',
                                                                direction: 'column',
                                                                justify: false,
                                                                translateX: 100,
                                                                translateY: 0,
                                                                itemsSpacing: 0,
                                                                itemDirection: 'left-to-right',
                                                                itemWidth: 80,
                                                                itemHeight: 20,
                                                                itemOpacity: 0.75,
                                                                symbolSize: 12,
                                                                symbolShape: 'circle',
                                                                symbolBorderColor: 'rgba(0, 0, 0, .5)',
                                                                effects: [
                                                                    {
                                                                        on: 'hover',
                                                                        style: {
                                                                            itemBackground: 'rgba(0, 0, 0, .03)',
                                                                            itemOpacity: 1
                                                                        }
                                                                    }
                                                                ]
                                                            }
                                                        ]}
                                                        />
                                                </div> 
                                            </div>
                                    )}
                            </>
                    )
                    
                )}
        </div>
    );
};

export default CRON;