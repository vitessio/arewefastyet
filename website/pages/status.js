import styles from '../styles/Home.module.css'
import stylesStatus from '../styles/Status.module.css'
import Header from "../components/header";
import Footer from "../components/footer";
import useSWR from 'swr'

const fetcher = (...args) => fetch(...args).then((res) => res.json())

export default function Status() {
    const { data, error } = useSWR('https://jsonplaceholder.typicode.com/todos/1', fetcher)

    if (error) return <div>Failed to load</div>
    if (!data) return <div>Loading...</div>
    return (
        <div>
            <Header />
            <div className={styles.container}>
                <p>{data.title}</p>
            </div>
            <Footer />
        </div>
    )
}
