import {useState, useEffect} from 'react'
import moment from 'moment'
import {Table} from 'semantic-ui-react'
import 'semantic-ui-css/semantic.min.css';
import './App.css' 

// handle state and every five seconds update it with external API data 
const useFetch = url => {
  const [state, setState] = useState([])

  useEffect(() => {
    const interval = setInterval(() => {
      fetch(url)
        .then(data => data.json())
        .then(data =>  setState(data))
        .catch(function(error) {
           console.log(error)
        })
    }, 5000)

    return () => clearInterval(interval);
  }, [url, state])

  return state
}

// create table columns
const headerRow = ['Time', 'Door', 'Temperature']

// create table rows
// this implementation contains a bug which makes all rows rendered with warninings
// if a row data contains a warning flag
// haven't found how to resolve it in a reasonable time
const renderBodyRow = ({ TimeStamp, IsOpen, IsOpenLong, Temp, IsExtremeTemp }, i) => ({
  key: TimeStamp || `row-${i}`,
  warning: !!(IsOpen && IsOpen.match('Requires Action')),
  cells: [
    moment(TimeStamp).format('dd:h:mm:ss') || 'No time specified',
    // this makes door notification handle two cases at once. Should be decoupled
    !!IsOpen || IsOpenLong ? { key: 'open', icon: 'attention', content: 'open'} : "closed",
    !!IsExtremeTemp ? {  icon: 'attention', content: Temp, warning: true } : Temp
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
