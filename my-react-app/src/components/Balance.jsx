import React, { useState } from "react";
import { queryBalance } from "../api/walletApi";
import { Button, TextField, Grid } from "@mui/material";

const Balance = ({ account }) => {
  const [balance, setBalance] = useState("");

  const handleQueryBalance = async () => {
    try {
      const response = await queryBalance(account);
      setBalance(response.data.balance[0].amount + " " + response.data.balance[0].denom);
    } catch (error) {
      setBalance("Error fetching balance");
    }
  };

  return (
    <Grid container spacing={2}>
      <Grid item xs={12}>
        <Button variant="contained" color="primary" fullWidth onClick={handleQueryBalance}>Query Balance</Button>
      </Grid>
      <Grid item xs={12}>
        <TextField label="Balance" value={balance} fullWidth disabled />
      </Grid>
    </Grid>
  );
};

export default Balance;
