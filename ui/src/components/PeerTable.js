
import { useState } from 'react';
import { FileSelection, showFileSize, checkForZipType } from './FileSelection';
import CircularProgress from '@mui/material/CircularProgress';

const LearningStatus = {
  NOT_STARTED: 0,
  SPLITTING: 1,
  LEARNING: 2,
  COMBINING: 3,
  DONE: 4
};

const activeHelper = (is_active, is_transmitting) => {
    
    if (is_active && is_transmitting) {
        return (
            <td style={{"color": "orange"}}>BUSY</td>
        )
    }

    if (is_active)  {
        return (
            <td style={{"color": "green"}}>ONLINE</td>
        )
    }

    return (
        <td style={{"color": "red"}}>OFFLINE</td>
    )
}

export const StatusHandler = (props) => {
  const s = props.status

  if (s === LearningStatus.NOT_STARTED) {
    return (<></>)
  }

  let backgroundColor1, backgroundColor2, backgroundColor3 = "white"

  if (s === LearningStatus.SPLITTING) {
    backgroundColor1 = "orange"
    backgroundColor2 = "white"
  } 
  else if (s === LearningStatus.LEARNING) {
    backgroundColor1 = "green"
    backgroundColor2 = "orange"
  }
  else if (s === LearningStatus.COMBINING) {
    backgroundColor1 = "green"
    backgroundColor2 = "green"
    backgroundColor3 = "orange"
  }
  else if ( s === LearningStatus.DONE) {
    backgroundColor1 = "green"
    backgroundColor2 = "green"
    backgroundColor3 = "green"
    console.log(props.model)
  }

  function handleModelDownload(data) {
    const blob = new Blob([data]);
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'model.pth';
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
  }

  if (s !== LearningStatus.DONE) {
    return (
    
      <div className="container">
        
        <div className="box" style={{backgroundColor: backgroundColor1}}>Splitting the training set <span className="arrow">&rarr;</span></div>
        <div className="box" style={{backgroundColor: backgroundColor2}}>Peer Learning in progress <span className="arrow">&rarr;</span></div>
        <div className="box" style={{backgroundColor: backgroundColor3}}>Constructing the model</div>
        <CircularProgress />
      </div>
    )
  }

  return (
    <div className="container">
        <h2>model.pth</h2>
          {/* You can handle the received model data here */}
          <button className="download" onClick={() => handleModelDownload(props.model)}>Download</button>
    </div>
  )
  
}

const togglePeerFromList = (peer, list) => {
    const cp = list;
    const index = cp.indexOf(peer);
    if (index !== -1) {
        // Peer exists, so remove it
        cp.splice(index, 1);
    } else {
        // Peer doesn't exist, so add it
        cp.push(peer);
    }
    console.log(cp)
    return cp
}

async function sendFileAndWaitForResponse(ws, zipFile) {
  try {
      // Send message to server
      ws.send(zipFile)

      // Wait for response message from server
      return new Promise((resolve, reject) => {
        ws.onmessage = (event) => {
              console.log("Received message from server:", event.data);
              resolve(event.data); // Resolve the promise with the received message
          };
      });
  } catch (error) {
      console.error("Error:", error);
  }
}

async function startPeerLearningProcess(ws, msg) {
  try {
      // Send message to server
      ws.send(JSON.stringify(msg))
      // Wait for response message from server
      return new Promise((resolve, reject) => {
        ws.onmessage = (event) => {
              console.log("Received message from server:", event.data);
              resolve(event.data); // Resolve the promise with the received message
          };
      });
  } catch (error) {
      console.error("Error:", error);
  }
}

async function waitForModelParameters(ws) {
  try {
      // Send message to server
      ws.binaryType = "blob";
      // Wait for response message from server
      return new Promise((resolve, reject) => {
        ws.onmessage = (event) => {
              console.log("Received message from server:", event.data);
              resolve(event.data); // Resolve the promise with the received message
          };
      });
  } catch (error) {
      console.error("Error:", error);
  }
}

const startLearningProcess = async (event, zipFile, peers, selected, webSocketConnection, fileTransferProtocol, statusChanger, setCompleteModel) => {
  // prevent default form submission, handle all manually
  event.preventDefault();

  statusChanger(LearningStatus.SPLITTING)
  let res = await sendFileAndWaitForResponse(fileTransferProtocol, zipFile)
  
  if (res != "File transfer completed successfully") {
    throw Error("moroo")
  }
  statusChanger(LearningStatus.LEARNING)
  
  // first connect to application and send selected zip file

  // First get ids and ports of peers
  let msg = [];
  
  peers.forEach(peer => {
    if (selected.includes(peer.id)) {
      msg.push(peer)
    }
  });
  console.log(webSocketConnection, fileTransferProtocol)
  res = await startPeerLearningProcess(webSocketConnection, msg);

  if (res != "Peer learning completed!") {
    throw Error("Error in peer training")
  }

  statusChanger(LearningStatus.COMBINING)

  //webSocketConnection.onmessage(setCompleteModel(event.data))
  res = await waitForModelParameters(webSocketConnection);
  setCompleteModel(res)
  statusChanger(LearningStatus.DONE)
}


export function PeerTable(props) {
    const peers = props.peers;
    const ws = props.ws;
    const filetransferWs = props.filetransferWs;

    const [selectedPeers, setSelectedPeers] = useState([]);
    const [selectedFile, setSelectedFiles] = useState({name: "", size: 0, type: ""});
    const [status, setStatus] = useState(LearningStatus.NOT_STARTED);
    const [completeModel, setCompleteModel] = useState(null);

    const handleFileChange = (event) => {
      const files = event.target.files;
      console.log(files)
      setSelectedFiles(files[0]);
    };
    

    if (props.select) {
      
        return (
          <>
          <StatusHandler status={status} model={completeModel} />
          
          <form onSubmit={(event) => startLearningProcess(event, selectedFile, peers, selectedPeers, ws, filetransferWs, setStatus, setCompleteModel)}>
            
          <div className='peercon'>
            <p>Peers</p>
            <table>
            <tbody>
                <tr>
                    <th>
                      ip
                    </th>
                    <th>
                      country
                    </th>
                    <th>
                      cuda
                    </th>
                    <th>
                      online
                    </th>
                    <th>
                      select
                    </th>

                </tr>
                {peers.filter(peer => peer.is_active == 1).map(peer => (
                    <tr key={peer.id} className='choose-tr'>
                    <td>
                      {peer.ip}
                    </td>
                    <td>
                      FI
                    </td>
                    <td>
                      {peer.has_cuda === 1 ? "✅": "❌"}
                    </td>
                    {activeHelper(peer.is_active, peer.is_transmitting)}
                    <td>
                        <input type="checkbox" onClick={() => setSelectedPeers(togglePeerFromList(peer.id, selectedPeers))}/>
                    </td>
                  </tr>
                ))}
            </tbody>
        </table>
        </div>
        <div className='peercon' id="fileselect">
        <div>
        <p> Training set </p>
        <input type="file" onChange={handleFileChange}/>
          <table id="filetable">
            <thead>
              <tr>
                <th>File Name</th>
                <th>File Size</th>
              </tr>
            </thead>
              <tbody>
                <tr>
                  <td>
                    {selectedFile.name}
                  </td>
                  <td>
                    {showFileSize(selectedFile.size)}
                  </td>
                </tr>
              </tbody>
            </table>
          {checkForZipType(selectedFile) ? <p style={{color: "red"}}>Requires zip file!</p> : ""}
        </div>
        </div>
        <button type="submit">Start Learning</button>
        </form>
        </>
        )
    }
    return (
        <table>
            <tbody>
                <tr>
                    <th>
                      ip
                    </th>
                    <th>
                      country
                    </th>
                    <th>
                      cuda
                    </th>
                    <th>
                      online
                    </th>
                    <th>
                      previously
                    </th>

                </tr>
                {props.peers.map(peer => (
                    <tr key={peer.id}>
                    <td>
                      {peer.ip}
                    </td>
                    <td>
                      FI
                    </td>
                    <td>
                      {peer.has_cuda === 0 ? "❌": "✅"}
                    </td>
                    {activeHelper(peer.is_active, peer.is_transmitting)}
                    <td>
                      {peer.last_online.slice(0,10)}
                    </td>
                  </tr>
                ))}
            </tbody>
        </table>
    )
};
