import React, {useState, useEffect} from 'react';

import '../PRcomponents/PRGitInfo'

const PRGitInfo = ({pull_nb}) => {
    const [dataPRGit, setDataPRGit] = useState([]);

    useEffect(() => {
        const fetchData = async () => {
          try {
            const responsePRGit = await fetch(
                `https://api.github.com/repos/vitessio/vitess/pulls/${pull_nb}`
            );
    
            const jsonDataPRGit = await responsePRGit.json();
            console.log(jsonDataPRGit);
            setDataPRGit(jsonDataPRGit);
          } catch (error) {
            console.log("Error while retrieving data from the API", error);
            setError(errorApi);
          }
        };
    
        fetchData();
      }, []);
      
    return (
        <div className='prGit'>
            <h2>{pull_nb}</h2>
        </div>
    );
};

export default PRGitInfo;