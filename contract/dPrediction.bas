///dPrediction v0.1 ♣♣♣♣♣♣

Function InitializePrivate() Uint64
    10 IF EXISTS("owner") == 0 THEN GOTO 30
    20 RETURN 1

    30 STORE("owner", SIGNER())
    40 STORE("p_init", 0)
    50 STORE("p_#", 0)
    60 STORE("p_amount", 0)
    70 STORE("p_total", 0)
    80 STORE("p_up", 0)
    100 STORE("p_down", 0)
    120 STORE("p_played", 0)
    130 STORE("time_a", 1800)
    140 STORE("time_b", 3600)
    150 STORE("time_c", 86400)
    160 STORE("v", 1)
    1000 RETURN 0
End Function


Function P_start(end Uint64, amt Uint64, predict String, feed String) Uint64
    10 IF LOAD("owner") == SIGNER() THEN GOTO 20 ELSE GOTO 30
    20 IF LOAD("p_init") == 0 THEN GOTO 40
    30 RETURN 1

    40 STORE("p_end_at", end)
    50 STORE("p_amount", amt)
    60 STORE("predicting", predict)
    70 STORE("p_url", feed)
    80 STORE("p_init", 1)
    85 STORE("p_#", 0)
    90 DELETE("mark")
    100 DELETE("p_final")
    110 DELETE("p_final_txid")
    120 IF DEROVALUE() > 0 THEN GOTO 300
    200 RETURN 0

    300 STORE("p_total", LOAD("p_total")+DEROVALUE())
    400 RETURN 0
    1000 RETURN 0
End Function


Function Predict(pre Uint64, name String) Uint64
    5 IF name == "" THEN GOTO 40
    10 IF LOAD("p_init") == 1 THEN GOTO 20 ELSE GOTO 40
    20 IF DEROVALUE() == LOAD("p_amount") THEN GOTO 30 ELSE GOTO 40
    30 IF BLOCK_TIMESTAMP() < LOAD("p_end_at") THEN GOTO 60
    40 RETURN 1

    60 IF pre == 0 THEN GOTO 100
    70 IF pre != 1 THEN GOTO 40 
    80 STORE("p_up", LOAD("p_up")+1)
    90 GOTO 110
    100 STORE("p_down", LOAD("p_down")+1)
    110 STORE("p-"+ITOA(LOAD("p_#")+1)+"-"+ITOA(pre), SIGNER())
    120 IF EXISTS(HEX(SIGNER())) THEN GOTO 140
    125 IF EXISTS("u_"+name) THEN GOTO 40
    130 STORE(HEX(SIGNER()), "u_"+name)
    135 STORE("u_"+name, 0)
    140 STORE("p_total", LOAD("p_total")+DEROVALUE())
    150 STORE("p_#", LOAD("p_#")+1)
    160 SEND_DERO_TO_ADDRESS(LOAD("owner"), (300*DEROVALUE()/10000))  
    170 STORE("p_total", LOAD("p_total")-(300*DEROVALUE()/10000))
    1000 RETURN 0
End Function



Function Post(price Uint64) Uint64
    10 IF EXISTS("mark") == 0 THEN GOTO 20 ELSE GOTO 30
    20 IF LOAD("owner") == SIGNER() THEN GOTO 40
    30 RETURN 1

    40 IF BLOCK_TIMESTAMP() >= LOAD("p_end_at") && BLOCK_TIMESTAMP() <= LOAD("p_end_at")+LOAD("time_a") THEN GOTO 60
    50 RETURN 1
    60 STORE("mark", price)
    100 RETURN 0
End Function


Function NameChange(name String) Uint64
    10 IF EXISTS("u_"+name) THEN GOTO 50
    20 IF name == "" THEN GOTO 50
    30 IF EXISTS(HEX(SIGNER())) == 0 THEN GOTO 50
    40 IF DEROVALUE() == 10000 THEN GOTO 60
    50 RETURN 1

    60 STORE(HEX(SIGNER()), "u_"+name)
    70 STORE("u_"+name, 0)
    80 SEND_DERO_TO_ADDRESS(LOAD("owner"), 10000)
    1000 RETURN 0
End Function


Function Remove() Uint64
    10 IF EXISTS(HEX(SIGNER())) == 0 THEN GOTO 30
    20 IF DEROVALUE() == 5000 THEN GOTO 40
    30 RETURN 1

    40 DELETE(HEX(SIGNER()))
    50 SEND_DERO_TO_ADDRESS(LOAD("owner"), 5000)
    1000 RETURN 0
End Function


Function p_clear() Uint64
    10 DIM i as Uint64
    20 LET i = 0
    30 DELETE("p-"+ITOA(i+1)+"-"+ITOA(1))
    40 DELETE("p-"+ITOA(i+1)+"-"+ITOA(0))
    50 LET i = i +1
    60 IF i <= LOAD("p_#") THEN GOTO 30
    70 STORE("p_init", 0)
    80 STORE("p_#", 0)
    90 STORE("p_total", 0)
    100 STORE("p_up", 0)
    110 STORE("p_down", 0)
    120 STORE("p_amount", 0)
    130 DELETE("p_end_at")
    140 DELETE("p_url")
    200 RETURN 0
End Function


Function p_determine(i Uint64, p Uint64, div Uint64) Uint64
    10 IF LOAD("owner") == SIGNER() THEN GOTO 30
    20 RETURN 1
    
    30 IF EXISTS("p-"+ITOA(i+1)+"-"+ITOA(p)) THEN GOTO 50
    40 RETURN 0
    50 SEND_DERO_TO_ADDRESS(LOAD("p-"+ITOA(i+1)+"-"+ITOA(p)), LOAD("p_total")/div)
    1000 RETURN 0
End Function


Function addWins(i Uint64, p Uint64) Uint64
    10 IF LOAD("owner") == SIGNER() THEN GOTO 30
    20 RETURN 1
    
    30 IF EXISTS("p-"+ITOA(i+1)+"-"+ITOA(p)) THEN GOTO 50
    40 RETURN 0
    50 STORE(LOAD(HEX(LOAD("p-"+ITOA(i+1)+"-"+ITOA(p)))), LOAD(LOAD(HEX(LOAD("p-"+ITOA(i+1)+"-"+ITOA(p)))))+1)
    1000 RETURN 0
End Function


Function P_end(price Uint64) Uint64
    10 IF LOAD("owner") == SIGNER() THEN GOTO 20 ELSE GOTO 30
    20 IF BLOCK_TIMESTAMP() >= LOAD("p_end_at")+LOAD("time_c") && BLOCK_TIMESTAMP() <= LOAD("p_end_at")+LOAD("time_c")+LOAD("time_b") THEN GOTO 40 ELSE GOTO 600
    30 RETURN 1
    40 DIM i, p as Uint64
    50 IF price < LOAD("mark") THEN GOTO 100
    60 IF price == LOAD("mark") THEN GOTO 30
    70 LET p = 1
    80 MAPSTORE("winners", LOAD("p_up"))
    90 IF MAPGET("winners") == 0 THEN GOTO 800 ELSE GOTO 120

    100 LET p = 0
    110 MAPSTORE("winners", LOAD("p_down"))
    115 IF MAPGET("winners") == 0 THEN GOTO 800
    120 SEND_DERO_TO_ADDRESS(LOAD("owner"), LOAD("p_total")%MAPGET("winners")) 

    130 LET i = 0
    140 p_determine(i, p, MAPGET("winners"))
    150 addWins(i, p)
    160 LET i = i +1
    170 IF i <= LOAD("p_#") THEN GOTO 140
    180 endStore(price)
    190 p_clear()
    500 RETURN 0
 
    600 IF LOAD("p_#") == 0 THEN GOTO 800
    610 DIM i as Uint64
    620 LET i = 0
    630 p_determine(i, 0, LOAD("p_#"))
    640 p_determine(i, 1, LOAD("p_#"))
    650 LET i = i +1
    660 IF i <= LOAD("p_#") THEN GOTO 630
    690 p_clear()
    700 endStore(price)
    750 RETURN 0
    800 SEND_DERO_TO_ADDRESS(LOAD("owner"), LOAD("p_total")) 
    810 endStore(price)
    820 p_clear()
    1000 RETURN 0
End Function


Function endStore(price Uint64) Uint64 
    10 STORE("p_played", LOAD("p_played")+1)
    20 STORE("p_final", LOAD("predicting")+"_"+ITOA(price))
    30 STORE("p_final_txid", TXID())
    40 DELETE("predicting")
    100 RETURN 0
End Function


Function UpdateCode(code String) Uint64  
    10 IF LOAD("owner") == SIGNER() THEN GOTO 30
    20 RETURN 1

    30 IF code == "" THEN GOTO 100 
    40 IF LOAD("p_init") == 1 THEN GOTO 100
    50 UPDATE_SC_CODE(code)
    60 STORE("v", LOAD("v")+1)
    100 RETURN 0
End Function


Function VarUpdate(ta Uint64, tb Uint64, tc Uint64) Uint64  
    10 IF LOAD("owner") == SIGNER() THEN GOTO 30
    20 RETURN 1

    30 IF LOAD("p_init") == 1 THEN GOTO 100
    40 STORE("time_a", ta)
    50 STORE("time_b", tb)
    60 STORE("time_c", tc)
    100 RETURN 0
End Function