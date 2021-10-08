// import logo from './logo.svg';
import {useState, useEffect} from 'react'
import './App.css';

function App() {
  const [state, setState] = useState([])
    useEffect(() => {
        fetch("http://localhost:5050/api/v1/readings/")
        .then(res => res.json())
        .then(data => setState(data))
        .catch(err => console.log(err))
    })

  console.log(state)
  return (
    <div className="App">
      <header className="App-header">
      </header>
    </div>
  );
}

export default App;
