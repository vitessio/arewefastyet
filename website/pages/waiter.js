import styles from '../styles/Waiter.module.css'
import Header from "../components/header";
import Footer from "../components/footer";
import Spinner from 'react-bootstrap/Spinner';

export default function Waiter() {
    return (
        <div>
            <Header />
            <div className={styles.spinner}>
                <Spinner animation="border" role="status">
                    <span className="visually-hidden">Loading...</span>
                </Spinner>
            </div>
            <Footer />
        </div>
    )
}
