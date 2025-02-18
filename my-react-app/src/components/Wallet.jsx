import React, { useState } from "react";
import Account from "./Account";
import Balance from "./Balance";
import Transaction from "./Transaction";
import { Button, Container, Paper, Typography } from "@mui/material";

const Wallet = ({ setIsLoggedIn }) => {
  const [account, setAccount] = useState("");

  return (
    <Container maxWidth="sm">
      <Paper style={{ padding: 20, marginTop: 20 }}>
        <Typography variant="h4" align="center">Blockchain Wallet</Typography>
        <Account setAccount={setAccount} account={account} />
        <Balance account={account} />
        <Transaction account={account} />
        <Button variant="contained" color="secondary" fullWidth onClick={() => setIsLoggedIn(false)}>Logout</Button>
      </Paper>
    </Container>
  );
};

export default Wallet;
