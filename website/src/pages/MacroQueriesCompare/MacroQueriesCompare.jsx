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

import React, { useEffect, useState } from 'react';
import { useLocation } from 'react-router-dom';

import '../MacroQueriesCompare/macroQueriesCompare.css'

import { errorApi } from '../../utils/Utils';
import QueryPlan from '../../components/QueryPlan/QueryPlan';
const MacroQueriesCompare = () => {

    const { state } = useLocation()
    const commitA = state.another
    const commitB = state.more
    const type = state.some.type

    const [dataQueryPlan, setDataQueryPlan] = useState([])

    useEffect(() => {
        const fetchData = async () => {
            try {
                const responseQueryPlan = await fetch(`${import.meta.env.VITE_API_URL}macrobench/compare/queries?ltag=${commitA}&rtag=${commitB}&type=${type}`);

                const jsonDataQueryPlan = await responseQueryPlan.json()
                console.log(jsonDataQueryPlan)
                setDataQueryPlan(jsonDataQueryPlan)
            } catch (error){
                console.log('Error while retrieving data from the API', error);
                setError(errorApi);
            }
        }
        fetchData();
    }, [])
    
    
    
    return (
        <div className='query'>
            <div className='query__top flex--column'>
                <h2>Compare Query Plans</h2>
                <span>Comparing the query plans of two OLTP benchmarks: A and B.</span>
                <span><b>A</b> benchmarked commit <a href={`https://github.com/vitessio/vitess/commit/${commitA}`}>{commitA.slice(0, 8)}</a> using the Gen4 query planner.</span>
                <span><b>B</b> benchmarked commit <a href={`https://github.com/vitessio/vitess/commit/${commitB}`}>{commitB.slice(0, 8)}</a> using the Gen4 query planner.</span>
                <span>Queries are ordered from the worst regression in execution time to the best. All executed queries are shown below.</span>
            </div>
            <figure className='line'></figure>
            {dataQueryPlan.map((queryPlan, index) => {
                return (
                    <QueryPlan key={index} data={queryPlan}/>
                )
            })}

        </div>
    );
};

export default MacroQueriesCompare;