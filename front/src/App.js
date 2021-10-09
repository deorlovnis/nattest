// import logo from './logo.svg';
import {useState, useEffect} from 'react'
import './App.css';

const useFetch = url => {
  const [state, setState] = useState()

  useEffect(() => {
    const interval = setInterval(() => {
      fetch(url)
        .then(data => data.json())
        .then(data => setState(data))
        .catch(function(error) {
          console.log(error)
        })
    }, 5000)

    return () => clearInterval(interval);
  }, [url])

  return state
}

function App() {
  const state = useFetch("http://localhost:5050/api/v1/readings/")
  

  console.log(state)
  return (
    <div className="App">
      <header className="App-header">
      </header>
    </div>
  );
}

export default App;
