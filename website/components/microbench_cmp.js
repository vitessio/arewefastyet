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
        </div>
    )
}
