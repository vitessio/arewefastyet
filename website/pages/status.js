import styles from '../styles/Home.module.css'
import stylesStatus from '../styles/Status.module.css'
import Header from "../components/header";
import Footer from "../components/footer";
import Table from 'react-bootstrap/Table';
import Badge from 'react-bootstrap/Badge';
import moment from "moment";
import { useState, useEffect } from 'react'
import Waiter from "./waiter";

export default function Status() {
    const [recent, setRecent] = useState(null);
    const [isRecentLoading, setRecentLoading] = useState(true)
    const [queue, setQueue] = useState(null);
    const [isQueueLoading, setQueueLoading] = useState(true)


    useEffect(() => {
        fetch('http://localhost:9090/api/recent')
            .then((res) => res.json())
            .then((data) => {
                setRecent(data)
                setRecentLoading(false)
            })

        fetch('http://localhost:9090/api/queue')
            .then((res) => res.json())
            .then((data) => {
                setQueue(data)
                setQueueLoading(false)
            })
    }, [])

    if (isRecentLoading || isQueueLoading || !recent || !queue) {
        return <Waiter />
    }

    return (
        <div>
            <Header />
            <div className={styles.container}>
                <div>
                    <div className={stylesStatus.card}>
                        <h4 className={stylesStatus.h4}>Executions Queue</h4>
                        <Table className={stylesStatus.table} striped bordered hover>
                            <thead>
                            <tr>
                                <th>SHA</th>
                                <th>Source</th>
                                <th>Type</th>
                                <th>PR #</th>
                            </tr>
                            </thead>
                            <tbody>
                            {queue.map((item,index)=>{
                                let short_git_ref = item.git_ref.slice(0, 8)
                                let commit_github_link = "https://github.com/vitessio/vitess/commit/"+item.git_ref

                                let pr_item = <></>
                                if (item.pull_nb > 0) {
                                    let pr_link = "https://github.com/vitessio/vitess/pull/" + item.pull_nb
                                    pr_item = <a href={pr_link} target="_blank">{item.pull_nb}</a>
                                }

                                return <tr>
                                    <td><a href={commit_github_link} target="_blank">{short_git_ref}</a></td>
                                    <td>{item.source}</td>
                                    <td>{item.type_of}</td>
                                    <td>{pr_item}</td>
                                </tr>
                            })}
                            </tbody>
                        </Table>
                    </div>

                    <div className={stylesStatus.card}>
                        <h4 className={stylesStatus.h4}>Recent Executions</h4>
                        <Table className={stylesStatus.table} striped bordered hover>
                            <thead>
                                <tr>
                                    <th>ID</th>
                                    <th>SHA</th>
                                    <th>Source</th>
                                    <th>Started</th>
                                    <th>Finished</th>
                                    <th>Type</th>
                                    <th>PR #</th>
                                    <th>Go</th>
                                    <th>Status</th>
                                </tr>
                            </thead>
                            <tbody>
                                {recent.map((item,index)=>{
                                    let uuid = item.uuid.slice(0, 6)
                                    let short_git_ref = item.git_ref.slice(0, 8)
                                    let commit_github_link = "https://github.com/vitessio/vitess/commit/"+item.git_ref

                                    let started_at = ""
                                    if (item.started_at) {
                                        started_at = moment(item.started_at, "YYYY-MM-DDTHH:mm:ssZ").format("MMM Do, hh:mm:ss a");
                                    }

                                    let finished_at = ""
                                    if (item.finished_at) {
                                        finished_at = moment(item.finished_at, "YYYY-MM-DDTHH:mm:ssZ").format("MMM Do, hh:mm:ss a");
                                    }

                                    let pr_item = <></>
                                    if (item.pull_nb > 0) {
                                        let pr_link = "https://github.com/vitessio/vitess/pull/" + item.pull_nb
                                        pr_item = <a href={pr_link} target="_blank">{item.pull_nb}</a>
                                    }

                                    let status_item
                                    switch (item.status) {
                                        case "finished":
                                            status_item = <Badge pill bg="success">{item.status}</Badge>
                                            break
                                        case "failed":
                                            status_item = <Badge pill bg="danger">{item.status}</Badge>
                                            break
                                        case "started":
                                            status_item = <Badge pill bg="primary">{item.status}</Badge>
                                            break
                                        default:
                                            status_item = <Badge pill bg="secondary">{item.status}</Badge>
                                    }

                                    return <tr>
                                        <td>{uuid}</td>
                                        <td><a href={commit_github_link} target="_blank">{short_git_ref}</a></td>
                                        <td>{item.source}</td>
                                        <td>{started_at}</td>
                                        <td>{finished_at}</td>
                                        <td>{item.type_of}</td>
                                        <td>{pr_item}</td>
                                        <td>{item.golang_version}</td>
                                        <td>{status_item}</td>
                                    </tr>
                                })}
                            </tbody>
                        </Table>
                    </div>
                </div>
            </div>
            <Footer />
        </div>
    )
}
