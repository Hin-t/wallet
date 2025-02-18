import React from "react";
import { createAccount } from "../api/walletApi";
import { Button, TextField, Grid } from "@mui/material";

const Account = ({ setAccount, account }) => {
  const handleCreateAccount = async () => {
    try {
      const response = await createAccount();
      setAccount(response.data.address);
    } catch (error) {
      console.log("Error creating account");
    }
  };

  return (
    <Grid container spacing={2}>
      <Grid item xs={12}>
        <Button variant="contained" color="primary" fullWidth onClick={handleCreateAccount}>Create Account</Button>
      </Grid>
      <Grid item xs={12}>
        <TextField label="Account Address" value={account} fullWidth disabled />
      </Grid>
    </Grid>
  );
};

export default Account;
