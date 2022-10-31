import styles from '../styles/Home.module.css'
import stylesFooter from '../styles/Footer.module.css'
import Container from 'react-bootstrap/Container';
import Navbar from 'react-bootstrap/Navbar';
import Nav from 'react-bootstrap/Nav';

export default function Footer() {
    return (
        <div className={styles.container}>
            <footer className={stylesFooter.footer}>
                <p>Copyright &copy; 2022</p>
            </footer>
        </div>
    )
}
