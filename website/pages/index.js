import Head from 'next/head'
import styles from '../styles/Home.module.css'
import Header from "../components/header";
import Footer from "../components/footer";

export default function Home() {
  return (
    <div>
      <Header />
      <div className={styles.container}>
        <Head>
          <title>arewefastyet</title>
          <meta name="description" content="Vitess' arewefastyet benchmarking tool" />
          <link rel="icon" href="/favicon.ico" />
        </Head>

        <main className={styles.main}>
          <h1 className={styles.title}>
            Welcome to arewefastyet!
          </h1>

          <p className={styles.description}>
            An automated benchmarking system for <a href="https://vitess.io">Vitess</a>.<br/>
            Providing adopters and maintainers a clear vision of how Vitess is performing throughout different releases.
          </p>

          <div className={styles.grid}>
            <a href="https://vitess.io/blog/2021-07-08-announcing-vitess-arewefastyet/" className={styles.card}  target="_blank">
              <h2>Learn &rarr;</h2>
              <p>Read our blog post to learn how arewefastyet works.</p>
            </a>

            <a href="https://vitess.io" className={styles.card} target="_blank">
              <h2>Vitess &rarr;</h2>
              <p>Find out more on the Vitess project.</p>
            </a>
          </div>
        </main>
      </div>
      <Footer />
    </div>
  )
}
