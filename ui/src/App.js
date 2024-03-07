import './App.css';
import { useState, useEffect} from 'react';
import { PeerTable } from './components/PeerTable';
import CircularProgress from '@mui/material/CircularProgress';

function App() {
  
  const [isLoading, setIsLoading] = useState(true);
  const [peers, setPeers] = useState([]);
  const [yourIp, setYourIp] = useState("");
  
  const [connect, setConnect] = useState(false);
  const [ws, setWs] = useState(null)
  const [fileTransferWs, setFileTransferWs] = useState(null)
  
  useEffect(() => {
    const ws = new WebSocket("ws://localhost:7070/ws")
    ws.onopen = () => {
      setWs(ws);
    };

    const ws2 = new WebSocket("ws://localhost:7070/filetransfer")
    ws2.onopen = () => {
      setFileTransferWs(ws2);
    };
  }, [])

  useEffect(() => {
    const fetchData = async () => {
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

  if (isLoading) {
    return (
      <div className="App">
        {connect && (
          <div className="overlay">
            <div className="window">
              <button className='close-btn' onClick={() => setConnect(false)}>X</button>
              <PeerTable peers={peers} select={true} ws={ws} filetransferWs={fileTransferWs}/>
            </div>
          </div>
        )}
        <div className="container">
          <CircularProgress />
        </div>
        <div>
          <p>Establish connection</p>
          <button style={{color: "green"}} onClick={() => setConnect(true)}>Connect</button>
        </div>
        
      </div>
    );
  }
  
  return (
    <div className="App">
      {connect && (
        <div className="overlay">
          <div className="window">
            <button className='close-btn' onClick={() => setConnect(false)}>X</button>
            <PeerTable peers={peers} select={true} ws={ws} filetransferWs={fileTransferWs}/>
          </div>
        </div>
      )}
      <h1>Peer Learner App</h1>
      <p>Peers</p>
      <div className="container">
        <PeerTable peers={peers} select={false}/>
      </div>
      <div>
        <p>Establish connection</p>
        <button style={{color: "green"}} onClick={() => setConnect(true)}>Connect</button>
      </div>
      
    </div>
  );
}

export default App;
