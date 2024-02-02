import express from 'express';
import { Request, Response } from 'express';
import { body, validationResult } from 'express-validator';
import { getPeers, registerPeer } from './db';
import { validateFields } from './validation';


const app = express();
app.use(express.json()); // json body parsing

app.get("/get-peers", async (req: Request, res: Response) => {
    const results = await getPeers();
    res.status(200).json(results);
});

app.post("/register", [ body().custom(validateFields) ] , async (req: Request, res: Response) => {
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

app.listen(3000, () => {console.log("app running port 3000")})
