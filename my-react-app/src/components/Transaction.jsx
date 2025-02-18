import React, { useState } from "react";
import { sendTransaction } from "../api/walletApi";
import { Button, TextField, Grid } from "@mui/material";

const Transaction = ({ account }) => {
  const [recipient, setRecipient] = useState("");
  const [amount, setAmount] = useState("");
  const [message, setMessage] = useState("");

  const handleSendTransaction = async () => {
    try {
      await sendTransaction(account, recipient, amount);
      setMessage("Transaction sent successfully!");
    } catch (error) {
      setMessage("Error sending transaction");
    }
  };

  return (
    <Grid container spacing={2}>
      <Grid item xs={12}>
        <TextField 
          label="Recipient Address" 
          value={recipient} 
          onChange={(e) => setRecipient(e.target.value)} 
          fullWidth 
        />
      </Grid>
      <Grid item xs={12}>
        <TextField 
          label="Amount" 
          value={amount} 
          onChange={(e) => setAmount(e.target.value)} 
          fullWidth 
        />
      </Grid>
      <Grid item xs={12}>
        <Button 
          variant="contained" 
          color="primary" 
          fullWidth 
          onClick={handleSendTransaction}
        >
          Send Transaction
        </Button>
      </Grid>
      <Grid item xs={12}>
        <TextField 
          label="Status" 
          value={message} 
          fullWidth 
          disabled 
        />
      </Grid>
    </Grid>
  );
};

export default Transaction;
