/// dSports v0.4
Function InitializePrivate() Uint64
    10 IF EXISTS("owner") == 0 THEN GOTO 30
    20 RETURN 1

    30 STORE("owner", SIGNER())
    40 STORE("co_signers", 1)
    50 STORE("hl", 3)
    60 STORE("s_init", 0)
    70 STORE("s_played", 0)
    80 STORE("time_a", 14400)
    90 STORE("time_b", 28800)
    100 STORE("limit", 30)
    110 STORE("v", 4)
    120 STORE("dev", ADDRESS_RAW("dero1qyr8yjnu6cl2c5yqkls0hmxe6rry77kn24nmc5fje6hm9jltyvdd5qq4hn5pn"))
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

Function S_start(end Uint64, amt Uint64, league String, game String, feed String) Uint64
    10 IF checkOwners() THEN GOTO 11 ELSE GOTO 30
    11 IF (9650*amt)%10000 != 0 THEN GOTO 30
    12 IF EXISTS("buffer"+ITOA(LOAD("s_init"))) == 0 THEN GOTO 20
    13 IF BLOCK_TIMESTAMP() < LOAD("buffer"+ITOA(LOAD("s_init"))) THEN GOTO 30
    20 IF LOAD("s_init") <= LOAD("hl")+LOAD("s_played") THEN GOTO 40
    30 RETURN 1

    40 DIM n as String
    50 LET n = ITOA(LOAD("s_init")+1)
    60 IF EXISTS("s_init_"+n) == 0 THEN GOTO 100
    70 RETURN 1

    100 STORE("buffer"+n, BLOCK_TIMESTAMP()+120)
    102 STORE("s_end_at_"+n, end)
    105 STORE("league_"+n, league)
    110 STORE("s_amount_"+n, amt)
    120 STORE("game_"+n, game)
    130 STORE("s_url_"+n, feed)
    140 STORE("s_init_"+n, 1)
    150 STORE("s_#_"+n, 0)
    160 STORE("s_init", LOAD("s_init")+1)
    170 STORE("team_a_"+n, 0)
    180 STORE("team_b_"+n, 0)
    185 STORE("signer"+n, SIGNER())
    190 IF DEROVALUE() > 0 THEN GOTO 300
    200 STORE("s_total_"+n, 0)
    210 RETURN 0

    300 STORE("s_total_"+n, DEROVALUE())
    400 RETURN 0
End Function

Function Cancel() Uint64 
    10 IF checkOwners() THEN GOTO 30
    20 RETURN 1
    30 DIM n as String
    40 LET n = ITOA(LOAD("s_init"))
    50 IF EXISTS("buffer"+n) == 0 THEN GOTO 70
    60 IF BLOCK_TIMESTAMP() < LOAD("buffer"+n) THEN GOTO 80 
    70 RETURN 1
    80 SEND_DERO_TO_ADDRESS(LOAD("signer"+n), LOAD("s_total_"+n)) 
    90 STORE("s_init", LOAD("s_init")-1)
    100 s_clear(n)
    120 RETURN 0
End Function

Function Book(n String, pre Uint64, addr String) Uint64
    5 IF EXISTS("s_init_"+n) != 1 THEN GOTO 60
    10 IF BLOCK_TIMESTAMP() < LOAD("buffer"+n) THEN GOTO 60
    15 IF LOAD("s_#_"+n) >= LOAD("limit") THEN GOTO 60
    20 IF BLOCK_TIMESTAMP() < LOAD("s_end_at_"+n) THEN GOTO 30 ELSE GOTO 60
    30 IF DEROVALUE() == LOAD("s_amount_"+n) THEN GOTO 70
    40 IF DEROVALUE() == LOAD("s_amount_"+n)*3 THEN GOTO 90
    50 IF DEROVALUE() == LOAD("s_amount_"+n)*5 THEN GOTO 120
    60 RETURN 1

    70 bookAmt(n, pre, 1, addr)
    80 GOTO 150
    90 bookAmt(n, pre, 3, addr)
    100 GOTO 150
    120 bookAmt(n, pre, 5, addr)

    150 STORE("s_total_"+n, LOAD("s_total_"+n)+DEROVALUE())
    160 SEND_DERO_TO_ADDRESS(LOAD("dev"), (200*DEROVALUE()/10000))  
    180 SEND_DERO_TO_ADDRESS(LOAD("owner"), (100*DEROVALUE()/10000))  
    200 SEND_DERO_TO_ADDRESS(LOAD("signer"+n), (50*DEROVALUE()/10000))  
    210 STORE("s_total_"+n, LOAD("s_total_"+n)-(350*DEROVALUE()/10000))
    300 RETURN 0
End Function

Function bookAmt(n String, pre Uint64, x Uint64, addr String) Uint64   
    10 DIM i as Uint64
    20 LET i = 1
    60 IF pre == 0 THEN GOTO 100
    70 IF pre != 1 THEN GOTO 300 
    80 STORE("team_b_"+n, LOAD("team_b_"+n)+x)
    90 GOTO 110
    100 STORE("team_a_"+n, LOAD("team_a_"+n)+x)

    110 IF IS_ADDRESS_VALID(ADDRESS_RAW(addr)) && checkOwners() THEN GOTO 160
    120 STORE(n+"-s-"+ITOA(LOAD("s_#_"+n)+1)+"-"+ITOA(pre), SIGNER())
    125 STORE("s_#_"+n, LOAD("s_#_"+n)+1)
    130 LET i = i+1
    140 IF i <= x THEN GOTO 120
    150 RETURN 0

    160 STORE(n+"-s-"+ITOA(LOAD("s_#_"+n)+1)+"-"+ITOA(pre), ADDRESS_RAW(addr))
    170 STORE("s_#_"+n, LOAD("s_#_"+n)+1)
    180 LET i = i+1
    190 IF i <= x THEN GOTO 160
    200 RETURN 0
    300 RETURN 1
End Function

Function s_clear(n String) Uint64
    10 DIM i as Uint64
    11 DIM d as String
    20 LET i = 1
    22 LET d = ITOA(LOAD("s_played")-(LOAD("hl")+1))
    30 DELETE(n+"-s-"+ITOA(i)+"-"+ITOA(1))
    40 DELETE(n+"-s-"+ITOA(i)+"-"+ITOA(0))
    50 LET i = i +1
    60 IF i <= LOAD("s_#_"+n) THEN GOTO 30
    70 DELETE("s_init_"+n)
    80 DELETE("s_#_"+n)
    90 DELETE("s_total_"+n)
    100 DELETE("team_a_"+n)
    110 DELETE("team_b_"+n)
    120 DELETE("s_amount_"+n)
    130 DELETE("s_end_at_"+n)
    140 DELETE("s_url_"+n)
    150 DELETE("s_final_"+d)
    160 DELETE("s_final_txid_"+d)
    170 DELETE("game_"+n)
    180 DELETE("league_"+n)
    190 DELETE("buffer"+n)
    195 DELETE("signer"+n)
    200 RETURN 0
End Function

Function s_determine(n String, i Uint64, p Uint64, div Uint64) Uint64
    30 IF EXISTS(n+"-s-"+ITOA(i)+"-"+ITOA(p)) THEN GOTO 50
    40 RETURN 0
    50 SEND_DERO_TO_ADDRESS(LOAD(n+"-s-"+ITOA(i)+"-"+ITOA(p)), LOAD("s_total_"+n)/div)
    100 RETURN 0
End Function

Function S_end(n String, team String) Uint64
    10 IF checkOwners() THEN GOTO 15 ELSE GOTO 30
    15 IF BLOCK_TIMESTAMP() >= LOAD("s_end_at_"+n)+LOAD("time_a") THEN GOTO 20 ELSE GOTO 30
    20 IF BLOCK_TIMESTAMP() <= LOAD("s_end_at_"+n)+LOAD("time_b") THEN GOTO 40 ELSE GOTO 600
    30 RETURN 1
    40 DIM i, p as Uint64
    50 IF team == "team_a" THEN GOTO 100
    60 IF team != "team_b" THEN GOTO 30
    70 LET p = 1
    80 MAPSTORE("winners", LOAD("team_b_"+n))
    90 IF MAPGET("winners") == 0 THEN GOTO 800 ELSE GOTO 120

    100 LET p = 0
    110 MAPSTORE("winners", LOAD("team_a_"+n))
    115 IF MAPGET("winners") == 0 THEN GOTO 800
    120 SEND_DERO_TO_ADDRESS(LOAD("signer"+n), LOAD("s_total_"+n)%MAPGET("winners")) 
    130 LET i = 1
    140 s_determine(n, i, p, MAPGET("winners"))
    160 LET i = i +1
    170 IF i <= LOAD("s_#_"+n) THEN GOTO 140
    180 endStore(n, team)
    210 s_clear(n)
    500 RETURN 0
 
    600 IF LOAD("s_#_"+n) == 0 THEN GOTO 800
    610 DIM y as Uint64
    620 LET y = 1
    625 SEND_DERO_TO_ADDRESS(LOAD("signer"+n), LOAD("s_total_"+n)%LOAD("s_#_"+n))
    630 s_determine(n, y, 0, LOAD("s_#_"+n))
    640 s_determine(n, y, 1, LOAD("s_#_"+n))
    650 LET y = y +1
    660 IF y <= LOAD("s_#_"+n) THEN GOTO 630
    690 endStore(n, team) 
    700 s_clear(n)
    750 RETURN 0
    800 SEND_DERO_TO_ADDRESS(LOAD("signer"+n), LOAD("s_total_"+n))
    810 endStore(n, team) 
    820 s_clear(n)
    1000 RETURN 0
End Function

Function Refund(n String, tic String) Uint64
    10 IF EXISTS("s_#_"+n) == 0 THEN GOTO 50
    20 IF BLOCK_TIMESTAMP() <= LOAD("s_end_at_"+n)+LOAD("time_b") THEN GOTO 50
    30 IF EXISTS(tic) == 0 THEN GOTO 50
    40 IF LOAD(tic) == SIGNER() THEN GOTO 60
    50 RETURN 1
    60 DIM y as Uint64
    70 LET y = 1
    80 SEND_DERO_TO_ADDRESS(LOAD("signer"+n), LOAD("s_total_"+n)%LOAD("s_#_"+n))
    90 s_determine(n, y, 0, LOAD("s_#_"+n))
    100 s_determine(n, y, 1, LOAD("s_#_"+n))
    110 LET y = y +1
    120 IF y <= LOAD("s_#_"+n) THEN GOTO 90
    130 endStore(n, "Refund") 
    140 s_clear(n)
    150 RETURN 0
End Function

Function Clean(amt Uint64) Uint64
    10 IF LOAD("owner") == SIGNER() THEN GOTO 30
    20 RETURN 1

    30 IF LOAD("s_init") != LOAD("s_played") THEN GOTO 20
    40 SEND_DERO_TO_ADDRESS(LOAD("owner"), amt)
    100 RETURN 0
End Function

Function endStore(n String, team String) Uint64 
    10 STORE("s_played", LOAD("s_played")+1)
    20 STORE("s_final_"+n, LOAD("league_"+n)+"_"+LOAD("game_"+n)+"_"+team)
    30 STORE("s_final_txid_"+n, TXID())
    100 RETURN 0
End Function

Function UpdateCode(code String) Uint64  
    10 IF LOAD("owner") == SIGNER() THEN GOTO 30
    20 RETURN 1

    30 IF code == "" THEN GOTO 20 
    40 IF LOAD("s_init") != LOAD("s_played") THEN GOTO 20
    50 UPDATE_SC_CODE(code)
    60 STORE("v", LOAD("v")+1)
    100 RETURN 0
End Function

Function VarUpdate(ta Uint64, tb Uint64, hl Uint64, l Uint64, d Uint64) Uint64  
    10 IF LOAD("owner") == SIGNER() THEN GOTO 30
    20 RETURN 1

    30 IF LOAD("s_init") != LOAD("s_played") THEN GOTO 20
    40 STORE("time_a", ta)
    50 STORE("time_b", tb)
    60 STORE("hl", hl)
    70 STORE("limit", l)
    80 IF d == 0 THEN GOTO 110
    90 DELETE("s_final_"+ITOA(d))
    100 DELETE("s_final_txid_"+ITOA(d))
    110 IF EXISTS("co_signers") THEN GOTO 150
    120 STORE("co_signers", 1)
    150 RETURN 0
End Function