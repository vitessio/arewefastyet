import stylesMacrobench from '../styles/Macrobench.module.css'
import stylesWaiter from "../styles/Waiter.module.css";
import Table from 'react-bootstrap/Table';
import Spinner from "react-bootstrap/Spinner";
import { useState, useEffect } from 'react'
import prettyBytes from 'pretty-bytes';

export default function MacrobenchCmp(props) {
    const [benchmarks, setBenchmarks] = useState([]);
    const [benchmarksLoading, setBenchmarksLoading] = useState(true);

    useEffect(() => {
        setBenchmarksLoading(true)
        fetch('http://localhost:9090/api/macrobench/compare?rtag='+props.to.commit_hash+'&ltag='+props.from.commit_hash)
            .then((res) => res.json())
            .then((data) => {
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
            {benchmarks.map((item,index) => {
                return <div key={index} className={stylesMacrobench.card}>
                    <h5>{item.type}</h5>
                    <Table className={stylesMacrobench.table} striped bordered hover>
                        <thead>
                            <tr>
                                <th scope="col" className={stylesMacrobench.thsm}></th>
                                <th scope="col" className={stylesMacrobench.thmd}><a target="_blank" href={"https://github.com/vitessio/vitess/commit/"+props.from.commit_hash}>{props.from.name}</a></th>
                                <th scope="col" className={stylesMacrobench.thmd}><a target="_blank" href={"https://github.com/vitessio/vitess/commit/"+props.to.commit_hash}>{props.to.name}</a></th>
                                <th scope="col" className={stylesMacrobench.thsm}>Improved by %</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr>
                                <th scope="row">QPS Total</th>
                                <td>{item.diff.Compare.Result.qps.total}</td>
                                <td>{item.diff.Reference.Result.qps.total}</td>
                                <td>{item.diff.Diff.qps.total}</td>
                            </tr>
                            <tr>
                                <th scope="row">QPS Reads</th>
                                <td>{item.diff.Compare.Result.qps.reads}</td>
                                <td>{item.diff.Reference.Result.qps.reads}</td>
                                <td>{item.diff.Diff.qps.reads}</td>
                            </tr>
                            <tr>
                                <th scope="row">QPS Writes</th>
                                <td>{item.diff.Compare.Result.qps.writes}</td>
                                <td>{item.diff.Reference.Result.qps.writes}</td>
                                <td>{item.diff.Diff.qps.writes}</td>
                            </tr>
                            <tr>
                                <th scope="row">QPS Other</th>
                                <td>{item.diff.Compare.Result.qps.other}</td>
                                <td>{item.diff.Reference.Result.qps.other}</td>
                                <td>{item.diff.Diff.qps.other}</td>
                            </tr>
                            <tr>
                                <th scope="row">TPS</th>
                                <td>{item.diff.Compare.Result.tps}</td>
                                <td>{item.diff.Reference.Result.tps}</td>
                                <td>{item.diff.Diff.tps}</td>
                            </tr>
                            <tr>
                                <th scope="row">Latency</th>
                                <td>{item.diff.Compare.Result.latency}</td>
                                <td>{item.diff.Reference.Result.latency}</td>
                                <td>{item.diff.Diff.latency}</td>
                            </tr>
                            <tr>
                                <th scope="row">Errors</th>
                                <td>{item.diff.Compare.Result.errors}</td>
                                <td>{item.diff.Reference.Result.errors}</td>
                                <td>{item.diff.Diff.errors}</td>
                            </tr>
                            <tr>
                                <th scope="row">Reconnects</th>
                                <td>{item.diff.Compare.Result.reconnects}</td>
                                <td>{item.diff.Reference.Result.reconnects}</td>
                                <td>{item.diff.Diff.reconnects}</td>
                            </tr>
                            <tr>
                                <th scope="row">Time</th>
                                <td>{item.diff.Compare.Result.time}</td>
                                <td>{item.diff.Reference.Result.time}</td>
                                <td>{item.diff.Diff.time}</td>
                            </tr>
                            <tr>
                                <th scope="row">Threads</th>
                                <td>{item.diff.Compare.Result.threads}</td>
                                <td>{item.diff.Reference.Result.threads}</td>
                                <td>{item.diff.Diff.threads}</td>
                            </tr>
                            <tr>
                                <th scope="row">Total CPU Time</th>
                                <td>{item.diff.Compare.Metrics.TotalComponentsCPUTime}</td>
                                <td>{item.diff.Reference.Metrics.TotalComponentsCPUTime}</td>
                                <td>{item.diff.DiffMetrics.TotalComponentsCPUTime}</td>
                            </tr>
                            <tr>
                                <th scope="row">Total vtgate CPU Time</th>
                                <td>{item.diff.Compare.Metrics.ComponentsCPUTime ? item.diff.Compare.Metrics.ComponentsCPUTime.vtgate : 0}</td>
                                <td>{item.diff.Reference.Metrics.ComponentsCPUTime ? item.diff.Reference.Metrics.ComponentsCPUTime.vtgate : 0}</td>
                                <td>{item.diff.DiffMetrics.ComponentsCPUTime ? item.diff.DiffMetrics.ComponentsCPUTime.vtgate : 0}</td>
                            </tr>
                            <tr>
                                <th scope="row">Total vttablet CPU Time</th>
                                <td>{item.diff.Compare.Metrics.ComponentsCPUTime ? item.diff.Compare.Metrics.ComponentsCPUTime.vttablet : 0}</td>
                                <td>{item.diff.Reference.Metrics.ComponentsCPUTime ? item.diff.Reference.Metrics.ComponentsCPUTime.vttablet : 0}</td>
                                <td>{item.diff.DiffMetrics.ComponentsCPUTime ? item.diff.DiffMetrics.ComponentsCPUTime.vttablet : 0}</td>
                            </tr>
                            <tr>
                                <th scope="row">Total Allocs Bytes</th>
                                <td>{prettyBytes(item.diff.Compare.Metrics.TotalComponentsMemStatsAllocBytes)}</td>
                                <td>{prettyBytes(item.diff.Reference.Metrics.TotalComponentsMemStatsAllocBytes)}</td>
                                <td>{item.diff.DiffMetrics.TotalComponentsMemStatsAllocBytes}</td>
                            </tr>
                            <tr>
                                <th scope="row">Total vtgate CPU Time</th>
                                <td>{item.diff.Compare.Metrics.ComponentsMemStatsAllocBytes ? prettyBytes(item.diff.Compare.Metrics.ComponentsMemStatsAllocBytes.vtgate) : 0}</td>
                                <td>{item.diff.Reference.Metrics.ComponentsMemStatsAllocBytes ? prettyBytes(item.diff.Reference.Metrics.ComponentsMemStatsAllocBytes.vtgate) : 0}</td>
                                <td>{item.diff.DiffMetrics.ComponentsMemStatsAllocBytes ? item.diff.DiffMetrics.ComponentsMemStatsAllocBytes.vtgate : 0}</td>
                            </tr>
                            <tr>
                                <th scope="row">Total vttablet CPU Time</th>
                                <td>{item.diff.Compare.Metrics.ComponentsMemStatsAllocBytes ? prettyBytes(item.diff.Compare.Metrics.ComponentsMemStatsAllocBytes.vttablet) : 0}</td>
                                <td>{item.diff.Reference.Metrics.ComponentsMemStatsAllocBytes ? prettyBytes(item.diff.Reference.Metrics.ComponentsMemStatsAllocBytes.vttablet) : 0}</td>
                                <td>{item.diff.DiffMetrics.ComponentsMemStatsAllocBytes ? item.diff.DiffMetrics.ComponentsMemStatsAllocBytes.vttablet : 0}</td>
                            </tr>
                        </tbody>
                    </Table>
                </div>
            })}
        </div>
    )
}
