import express from 'express';
import cors from 'cors';
import { Request, Response } from 'express';
import { body, validationResult } from 'express-validator';
import { getPeers, registerPeer, changePeerState } from './db';
import { validateFields } from './validation';


const app = express();
app.use(express.json()); // json body parsing
app.use(cors( { origin: "*" }));

app.get("/api/v1/peers", async (req: Request, res: Response) => {
    const ip = req.headers['x-forwarded-for'] || req.connection.remoteAddress;
    const results = await getPeers(ip);
    res.status(200).json(results);
});

app.post("/api/v1/register", [ body().custom(validateFields) ] , async (req: Request, res: Response) => {
    const errors = validationResult(req);
    if (!errors.isEmpty()) {
      return res.status(400).json({ errors: errors.array() });
    }
    try {
        const msg = await registerPeer(req.body);
        res.status(200).send(msg)
    } catch (err) {
        res.status(400).send("Error registering peer")
    }
});

app.post("/api/v1/status", async (req: Request, res: Response) => {
    try {
        const msg = await changePeerState(req.body);
        res.status(200).send(msg)
    } catch (err) {
        res.status(400).send("Error changing status of peer")
    }
});

app.listen(3000, () => {console.log("app running port 3000")})
