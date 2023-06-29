/// dPrediction  v0.3
Function InitializePrivate() Uint64
    10 IF EXISTS("owner") == 0 THEN GOTO 30
    20 RETURN 1

    30 STORE("owner", SIGNER())
    33 STORE("co_signers", 1)
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
    160 STORE("limit", 30)
    170 STORE("v", 3)
    180 STORE("dev", ADDRESS_RAW("dero1qyr8yjnu6cl2c5yqkls0hmxe6rry77kn24nmc5fje6hm9jltyvdd5qq4hn5pn"))
    200 RETURN 0
End Function

Function checkOwners() Uint64
    10 IF SIGNER() == LOAD("owner") THEN GOTO 200
    20 IF LOAD("co_signers") < 2 THEN GOTO 100
    30 DIM i as Uint64
    40 LET i = 2
    45 IF EXISTS("co_signer"+ITOA(i)) == 0 THEN GOTO 60
    50 IF SIGNER() == LOAD("co_signer"+ITOA(i)) THEN GOTO 200
    60 LET i = i+1
    70 IF i <= 9 THEN GOTO 45
    100 RETURN 0
    200 RETURN 1
End Function

Function AddSigner(new String) Uint64 
    10 IF SIGNER() == LOAD("owner") && LOAD("co_signers") < 9 THEN GOTO 30
    20 RETURN 1
    30 DIM i as Uint64
    40 LET i = 1
    50 LET i = i+1
    60 IF i == 10 THEN GOTO 20
    70 IF EXISTS("co_signer"+ITOA(i)) THEN GOTO 50
    80 STORE("co_signers", LOAD("co_signers")+1)
    90 STORE("co_signer"+ITOA(i), ADDRESS_RAW(new))
    100 RETURN 0
End Function

Function RemoveSigner(remove Uint64) Uint64 
    10 IF SIGNER() == LOAD("owner") THEN GOTO 30
    20 RETURN 1
    30 IF EXISTS("co_signer"+ITOA(remove)) == 0 THEN GOTO 60
    40 STORE("co_signers", LOAD("co_signers")-1)
    50 DELETE("co_signer"+ITOA(remove))
    60 RETURN 0
End Function

Function P_start(end Uint64, amt Uint64, predict String, feed String, mark Uint64) Uint64
    10 IF checkOwners() THEN GOTO 11 ELSE GOTO 30
    11 IF (9650*amt)%10000 != 0 THEN GOTO 30
    20 IF LOAD("p_init") == 0 THEN GOTO 40
    30 RETURN 1

    40 STORE("buffer", BLOCK_TIMESTAMP()+300)
    41 STORE("p_end_at", end)
    50 STORE("p_amount", amt)
    60 STORE("predicting", predict)
    70 STORE("p_url", feed)
    80 STORE("p_init", 1)
    85 STORE("p_#", 0)
    86 IF mark == 0 THEN GOTO 90
    87 STORE("mark", mark)
    88 GOTO 100
    90 DELETE("mark")
    100 DELETE("p_final")
    110 DELETE("p_final_txid")
    115 STORE("signer", SIGNER())
    120 IF DEROVALUE() > 0 THEN GOTO 300
    200 RETURN 0

    300 STORE("p_total", LOAD("p_total")+DEROVALUE())
    400 RETURN 0
End Function

Function Cancel() Uint64 
    10 IF checkOwners() && BLOCK_TIMESTAMP() < LOAD("buffer") THEN GOTO 30 
    20 RETURN 1
    30 SEND_DERO_TO_ADDRESS(LOAD("signer"), LOAD("p_total")) 
    40 p_clear()
    50 RETURN 0
End Function

Function Predict(pre Uint64, addr String) Uint64
    4 IF LOAD("p_init") != 1 THEN GOTO 40
    5 IF BLOCK_TIMESTAMP() < LOAD("buffer") THEN GOTO 40
    10 IF LOAD("p_#") >= LOAD("limit") THEN GOTO 40
    20 IF DEROVALUE() == LOAD("p_amount") THEN GOTO 30 ELSE GOTO 40
    30 IF BLOCK_TIMESTAMP() < LOAD("p_end_at") THEN GOTO 60
    40 RETURN 1

    60 IF pre == 0 THEN GOTO 100
    70 IF pre != 1 THEN GOTO 40 
    80 STORE("p_up", LOAD("p_up")+1)
    90 GOTO 110
    100 STORE("p_down", LOAD("p_down")+1)
    110 IF IS_ADDRESS_VALID(ADDRESS_RAW(addr)) && checkOwners() THEN GOTO 130
    115 STORE("p-"+ITOA(LOAD("p_#")+1)+"-"+ITOA(pre), SIGNER())

    120 GOTO 140
    130 STORE("p-"+ITOA(LOAD("p_#")+1)+"-"+ITOA(pre), ADDRESS_RAW(addr))
    
    140 STORE("p_total", LOAD("p_total")+DEROVALUE())
    150 STORE("p_#", LOAD("p_#")+1)
    160 SEND_DERO_TO_ADDRESS(LOAD("dev"), (200*DEROVALUE()/10000))  
    170 SEND_DERO_TO_ADDRESS(LOAD("owner"), (100*DEROVALUE()/10000))  
    180 SEND_DERO_TO_ADDRESS(LOAD("signer"), (50*DEROVALUE()/10000))   
    190 STORE("p_total", LOAD("p_total")-(350*DEROVALUE()/10000))
    200 RETURN 0
End Function

Function Post(price Uint64) Uint64
    10 IF EXISTS("mark") == 0 THEN GOTO 20 ELSE GOTO 30
    20 IF checkOwners() THEN GOTO 40
    30 RETURN 1

    40 IF BLOCK_TIMESTAMP() >= LOAD("p_end_at") && BLOCK_TIMESTAMP() <= LOAD("p_end_at")+LOAD("time_a") THEN GOTO 60
    50 RETURN 1
    60 STORE("mark", price)
    100 RETURN 0
End Function

Function p_clear() Uint64
    10 DIM i as Uint64
    20 LET i = 1
    30 DELETE("p-"+ITOA(i)+"-"+ITOA(1))
    40 DELETE("p-"+ITOA(i)+"-"+ITOA(0))
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
    150 DELETE("buffer")
    160 DELETE("signer")
    200 RETURN 0
End Function

Function p_determine(i Uint64, p Uint64, div Uint64) Uint64
    30 IF EXISTS("p-"+ITOA(i)+"-"+ITOA(p)) THEN GOTO 50
    40 RETURN 0
    50 SEND_DERO_TO_ADDRESS(LOAD("p-"+ITOA(i)+"-"+ITOA(p)), LOAD("p_total")/div)
    100 RETURN 0
End Function

Function P_end(price Uint64) Uint64
    10 IF checkOwners() THEN GOTO 15 ELSE GOTO 30
    15 IF BLOCK_TIMESTAMP() >= LOAD("p_end_at")+LOAD("time_c") THEN GOTO 20 ELSE GOTO 30
    20 IF BLOCK_TIMESTAMP() <= LOAD("p_end_at")+LOAD("time_c")+LOAD("time_b") THEN GOTO 40 ELSE GOTO 600
    30 RETURN 1
    40 IF EXISTS("mark") == 0 THEN GOTO 30
    45 DIM i, p as Uint64
    50 IF price < LOAD("mark") THEN GOTO 100
    60 IF price == LOAD("mark") THEN GOTO 600
    70 LET p = 1
    80 MAPSTORE("winners", LOAD("p_up"))
    90 IF MAPGET("winners") == 0 THEN GOTO 800 ELSE GOTO 120

    100 LET p = 0
    110 MAPSTORE("winners", LOAD("p_down"))
    115 IF MAPGET("winners") == 0 THEN GOTO 800
    120 SEND_DERO_TO_ADDRESS(LOAD("signer"), LOAD("p_total")%MAPGET("winners")) 

    130 LET i = 1
    140 p_determine(i, p, MAPGET("winners"))
    160 LET i = i +1
    170 IF i <= LOAD("p_#") THEN GOTO 140
    180 endStore(price)
    190 p_clear()
    500 RETURN 0
 
    600 IF LOAD("p_#") == 0 THEN GOTO 800
    610 DIM y as Uint64
    620 LET y = 1
    625 SEND_DERO_TO_ADDRESS(LOAD("signer"), LOAD("p_total")%LOAD("p_#"))
    630 p_determine(y, 0, LOAD("p_#"))
    640 p_determine(y, 1, LOAD("p_#"))
    650 LET y = y +1
    660 IF y <= LOAD("p_#") THEN GOTO 630
    690 p_clear()
    700 endStore(price)
    750 RETURN 0
    800 SEND_DERO_TO_ADDRESS(LOAD("signer"), LOAD("p_total")) 
    810 endStore(price)
    820 p_clear()
    1000 RETURN 0
End Function

Function Refund(tic String) Uint64
    10 IF LOAD("p_#") == 0 THEN GOTO 50
    20 IF BLOCK_TIMESTAMP() <= LOAD("p_end_at")+LOAD("time_c")+LOAD("time_b") THEN GOTO 50
    30 IF EXISTS(tic) == 0 THEN GOTO 50
    40 IF LOAD(tic) == SIGNER() THEN GOTO 60
    50 RETURN 1
    60 DIM y as Uint64
    70 LET y = 1
    80 SEND_DERO_TO_ADDRESS(LOAD("signer"), LOAD("p_total")%LOAD("p_#"))
    90 p_determine(y, 0, LOAD("p_#"))
    100 p_determine(y, 1, LOAD("p_#"))
    110 LET y = y +1
    120 IF y <= LOAD("p_#") THEN GOTO 90
    130 p_clear()
    140 endStore(0)
    150 RETURN 0
End Function

Function Clean(amt Uint64) Uint64
    10 IF LOAD("owner") == SIGNER() THEN GOTO 30
    20 RETURN 1

    30 IF LOAD("p_init") == 1 THEN GOTO 20 
    40 SEND_DERO_TO_ADDRESS(LOAD("owner"), amt)
    100 RETURN 0
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

Function VarUpdate(ta Uint64, tb Uint64, tc Uint64, l Uint64) Uint64  
    10 IF LOAD("owner") == SIGNER() THEN GOTO 30
    20 RETURN 1

    30 IF LOAD("p_init") == 1 THEN GOTO 100
    40 STORE("time_a", ta)
    50 STORE("time_b", tb)
    60 STORE("time_c", tc)
    70 STORE("limit", l)
    80 IF EXISTS("co_signers") THEN GOTO 100
    90 STORE("co_signers", 1)
    100 RETURN 0
End Function