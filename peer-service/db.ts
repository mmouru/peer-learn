import * as mysql from 'mysql2/promise';
import { FieldPacket } from 'mysql2/promise';

interface RegisterBody {
    ip: string,
    port: string,
    peerId: string,
    hasCuda: boolean
}

interface StatusBody {
    ip: string,
    peerId: string,
    active?: string,
    isTransmitting?: string
}

const pool = mysql.createPool({
    socketPath: process.env.SOCKETPATH,
    user: process.env.DBUSER,
    password: process.env.DBPASS,
    database: process.env.DATABASE,
    connectionLimit : 50,
});

// get all peers info
export const getPeers = async (ip: any) => {
    const query = `SELECT * from Peers WHERE ip != "${ip}" ORDER by last_online DESC`;
    try {
        const connection = await pool.getConnection();
        const res = await connection.query(query);
        connection.release();
        return res[0];
    } catch (err) {
        console.error("caught error fetching all peers", err);
        throw err;
    }
};

// for registering peer
export const registerPeer = async (body: RegisterBody) => {
    const peer = body;
    console.log(peer, "peer")
    const hasCuda = peer.hasCuda ? 1 : 0;
    // current date + time for last_online column
    const formattedDate = new Date().toISOString().slice(0, 19).replace('T', ' ');

    const checkIfHostAlreadyExists = `SELECT * FROM Peers WHERE ip = "${peer.ip}";`
    let hostCheckResults : [any, FieldPacket[]];
    const connection = await pool.getConnection();
    try {
        hostCheckResults = await connection.query(checkIfHostAlreadyExists);
    } catch (err) {
        console.error('Error executing query:', err);
        throw err;
    }
    
    console.log(hostCheckResults)
    console.log(hostCheckResults[0])

    if (hostCheckResults[0].length > 0) {
        console.log("update existing information for found IP address");
        const updateStateQuery = `UPDATE Peers SET port = "${peer.port}", peer_id = "${peer.peerId}", has_cuda=${hasCuda}, is_active=1, is_transmitting=0, last_online = "${formattedDate}"
                                    WHERE ip = "${peer.ip}";`;
        
        try {
            await connection.query(updateStateQuery);
            connection.release();
            return "Updated existing peer information";
        } catch (err) {
            console.error('Error updating peer info:', err);
            throw err;
        }
    }

    const registerQuery = `INSERT INTO Peers (ip, port, peer_id, has_cuda, is_active, is_transmitting, last_online)
                VALUES ("${peer.ip}", "${peer.port}", "${peer.peerId}", has_cuda=${hasCuda}, 1, 0, "${formattedDate}");`

    try {
        await connection.query(registerQuery)
        connection.release();
        return "Added new peer to db";
    } catch (err) {
        console.error('Error updating peer info:', err);
        throw err;
    }
};

export const changePeerState = async (body: StatusBody) => {
    const peer = body;
    let updatePeerQuery = "";
    const connection = await pool.getConnection();
    if (peer.active) {
        const active = peer.active;
        updatePeerQuery = `UPDATE Peers SET is_active = "${active}" WHERE ip = "${peer.ip}" AND peer_id = "${peer.peerId}";`
    }

    else if (peer.isTransmitting) {
        const transmitting = peer.isTransmitting;
        updatePeerQuery = `UPDATE Peers SET is_transmitting = "${transmitting}" WHERE ip = "${peer.ip}" AND peer_id = "${peer.peerId}";`
    }

    try {
        await connection.query(updatePeerQuery)
        connection.release();
        return "Updated peer status";
    } catch (err) {
        console.error('Error updating peer info:', err);
        throw err;
    }
}
