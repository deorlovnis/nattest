import {useState, useEffect} from 'react'
import {Table} from 'semantic-ui-react'
import './App.css';

const useFetch = url => {
  const [state, setState] = useState([])

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
  }, [url, state])

  return state
}

const headerRow = ['Time', 'Door', 'Temperature']

const renderBodyRow = ({ TimeStamp, status, Temp }, i) => ({
  key: TimeStamp || `row-${i}`,
  warning: !!(status && status.match('Requires Action')),
  cells: [
    TimeStamp.T || 'No name specified',
    status ? { key: 'status', icon: 'attention', content: status } : 'Unknown',
    Temp
      ? { key: 'notes', icon: 'attention', content: Temp, warning: true }
      : 'None',
  ],
})


function App() {
  const readings = useFetch("http://localhost:5050/api/v1/readings/")

  return (
    <div >
      <header className="App-header">
        <Table
          celled
          headerRow={headerRow}
          renderBodyRow={renderBodyRow}
          tableData={readings}
        />
      </header>
    </div>
  );
}

export default App;
