import stylesCard from '../styles/Card.module.css'
import stylesTable from '../styles/Table.module.css'
import stylesWaiter from "../styles/Waiter.module.css";
import Table from 'react-bootstrap/Table';
import Spinner from "react-bootstrap/Spinner";
import { useState, useEffect } from 'react'
import prettyBytes from 'pretty-bytes';

export default function MicrobenchCmp(props) {
    const [benchmarks, setBenchmarks] = useState([]);
    const [benchmarksLoading, setBenchmarksLoading] = useState(true);

    useEffect(() => {
        setBenchmarksLoading(true)
        fetch('http://localhost:9090/api/microbench/compare?rtag='+props.to.commit_hash+'&ltag='+props.from.commit_hash)
            .then((res) => res.json())
            .then((data) => {
                if (data === null) {
                    data = []
                }
                setBenchmarks(data)
                setBenchmarksLoading(false)
            })
    }, [props])

    if (benchmarksLoading) {
        return <div className={stylesWaiter.spinner}>
            <Spinner animation="border" role="status">
                <span className="visually-hidden">Loading...</span>
            </Spinner>
        </div>
    }

    return (
        <div>
            <div className={stylesCard.card}>
                <Table className={stylesTable.table} striped bordered hover>
                    <thead>
                    <tr>
                        <th scope="col">Package</th>
                        <th scope="col">Benchmark Name</th>
                        <th scope="col" colSpan="3">Number of Iterations</th>
                        <th scope="col" colSpan="3">Time</th>
                        <th scope="col" colSpan="3">Bytes</th>
                        <th scope="col" colSpan="3">Megabytes</th>
                        <th scope="col" colSpan="3">Allocations/op</th>
                    </tr>
                    </thead>
                    <tbody>
                    <tr>
                        <th scope="col"></th>
                        <th scope="col"></th>
                        <th scope="col">{props.from.name}</th>
                        <th scope="col">{props.to.name}</th>
                        <th scope="col">Diff %</th>
                        <th scope="col">{props.from.name}</th>
                        <th scope="col">{props.to.name}</th>
                        <th scope="col">Diff %</th>
                        <th scope="col">{props.from.name}</th>
                        <th scope="col">{props.to.name}</th>
                        <th scope="col">Diff %</th>
                        <th scope="col">{props.from.name}</th>
                        <th scope="col">{props.to.name}</th>
                        <th scope="col">Diff %</th>
                        <th scope="col">{props.from.name}</th>
                        <th scope="col">{props.to.name}</th>
                        <th scope="col">Diff %</th>
                    </tr>
                    {benchmarks.map((item,index) => {
                        return <tr>

                        </tr>
                    })}
                    </tbody>
                </Table>
            </div>
        </div>
    )
}
