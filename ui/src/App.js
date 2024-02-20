import './App.css';
import { useState, useEffect} from 'react';
import { PeerTable } from './components/PeerTable';

function App() {
  
  const [isLoading, setIsLoading] = useState(true);
  const [peers, setPeers] = useState([])
  const [selectedFiles, setSelectedFiles] = useState([]);
  const [connect, setConnect] = useState(false);
  const [ws, setWs] = useState(null)

  const handleFileChange = (event) => {
    const files = event.target.files;
    setSelectedFiles([...selectedFiles, ...files]);
  };

  

  useEffect(() => {
    const ws = new WebSocket("ws://localhost:7070/ws")
    ws.onopen = () => {
      setWs(ws);
    }
  }, [])

  useEffect(() => {
    const fetchData = async () => {
      console.log("moro")
      try {
        const response = await fetch('https://peer-service-qfobv32vvq-lz.a.run.app/api/v1/peers');
        console.log(response)
        if (!response.ok) {
          throw new Error('Failed to fetch data');
        }
        const jsonData = await response.json();
        console.log(jsonData)
        setPeers(jsonData);
        setIsLoading(false);
      } catch (error) {
        console.log("e")
        console.log(error)
        setIsLoading(false);
      }
    };

    fetchData();
  }, []);

  
  return (
    <div className="App">
      {connect && (
        <div className="overlay">
          <div className="window">
            <button style={{}} onClick={() => setConnect(false)}>x</button>
            <PeerTable peers={peers} select={true} ws={ws}/>
            
          </div>
        </div>
      )}
      <div>
        <p>Connected Peers</p>
        <PeerTable peers={peers} select={false} ws={ws}/>
      </div>
      <div>
        <p> File drop </p>
        <input type="file" onChange={handleFileChange} multiple/>
        {selectedFiles.length > 0 && (
        <table>
          <thead>
            <tr>
              <th>File Name</th>
              <th>File Size</th>
            </tr>
          </thead>
          <tbody>
            {selectedFiles.map((file, index) => (
              <tr key={index}>
                <td>{file.name}</td>
                <td>{file.size} bytes</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
      
      </div>
      <div>
        <p>Establish connection</p>
        <button style={{color: "green"}} onClick={() => setConnect(true)}>Connect</button>
      </div>
      
    </div>
  );
}

export default App;
