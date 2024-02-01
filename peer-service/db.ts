import * as mysql from 'mysql2/promise';
import { FieldPacket } from 'mysql2/promise';

interface RegisterBody {
    ip: string,
    port: string,
    peerId: string,
}

const pool = mysql.createPool({
    socketPath: "",
    user: "",
    password: "",
    database: ""
});

// get all peers info
export const getPeers = async () => {
    const query = `SELECT * from Peers ORDER by last_online DESC`;
    try {
        const connection = await pool.getConnection();
        const res = await connection.query(query);
        return res[0];
    } catch (err) {
        console.error("caught error fetching all peers", err);
        throw err;
    }
};

// for registering peer
export const registerPeer = async (body: RegisterBody) => {
    const peer = body;
    console.log(peer.ip, "peer")
    // current date + time for last_online column
    const formattedDate = new Date().toISOString().slice(0, 19).replace('T', ' ');

    const checkIfHostAlreadyExists = `SELECT * FROM Peers WHERE ip = "${peer.ip}";`
    let hostCheckResults, updateResults, registerResults : [any, FieldPacket[]];
    const connection = await pool.getConnection();
    try {
        hostCheckResults = await connection.query(checkIfHostAlreadyExists);
    } catch (err) {
        console.error('Error executing query:', err);
        throw err;
    }
    //console.log(hostCheckResults)
    if (hostCheckResults[0]) {
        console.log("update existing information for found IP address");
        const updateStateQuery = `UPDATE Peers SET port = "${peer.port}", peer_id = "${peer.peerId}", is_active=1, last_online = "${formattedDate}"
                                    WHERE ip = "${peer.ip}";`;
        
        try {
            updateResults = await connection.query(updateStateQuery);
            return updateResults;
        } catch (err) {
            console.error('Error updating peer info:', err);
            throw err;
        }
    }

    const registerQuery = `INSERT INTO Peers (ip, port, peer_id, is_active, is_transmitting, last_online)
                VALUES ("${peer.ip}", "${peer.port}", "${peer.peerId}", 1, 0, "${formattedDate}");`

    try {
        const registerResults = await connection.query(registerQuery)
        return registerResults;
    } catch (err) {
        console.error('Error updating peer info:', err);
        throw err;
    }
};
