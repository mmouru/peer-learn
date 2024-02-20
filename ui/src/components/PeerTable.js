
import { useState } from 'react';

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



const connectPeersAndStartLearning = (event, peers, selected, webSocketConnection) => {
  event.preventDefault();
  // First get ids and ports of peers
  let msg = [];
  
  peers.forEach(peer => {
    if (selected.includes(peer.id)) {
      msg.push(peer)
    }
  });
  webSocketConnection.send(JSON.stringify(msg));
}


export function PeerTable(props) {
    const [selectedPeers, setSelectedPeers] = useState([]);
    const [peers, setPeers] = useState(props.peers);
    const [selectedCol, setSelectedCol] = useState(new Array(props.peers.length).fill(false))
    const [ws, setWs] = useState(props.ws)

  

    if (props.select) {
        return (
          <form onSubmit={(event) => connectPeersAndStartLearning(event, peers, selectedPeers, ws)}>
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
                    <th>
                      select
                    </th>

                </tr>
                {peers.map((peer, idx) => (
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
                      {peer.last_online.slice(0,10)}
                    </td>
                    <td>
                        <input type="checkbox" onClick={() => setSelectedPeers(togglePeerFromList(peer.id, selectedPeers))}/>
                    </td>
                  </tr>
                ))}
            </tbody>
        </table>
        <button type="submit">Connect</button>
        </form>
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
