// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in licenses/BSD-golang.txt.

// This code is based on compare_arm64.s from Go 1.12.5.

TEXT ·CommonPrefix(SB),$0-56
    MOVD    a_base+0(FP), R0
    MOVD    a_len+8(FP), R1
    MOVD    b_base+24(FP), R2
    MOVD    b_len+32(FP), R3

    CMP R1, R3
	CSEL    LT, R3, R1, R6    // R6 = min(alen, blen)
	// Throughout this function, R7 remembers the original min(alen, blen) and
	// R6 is the number of bytes we still need to compare (with bytes 0 to R7-R6
	// known to match).
	MOVD    R6, R7

	CBZ	R6, samebytes
	CMP $16, R6
	BLT small

	// length >= 16
chunk16_loop:
	LDP.P	16(R0), (R4, R8)
	LDP.P	16(R2), (R5, R9)
	CMP	R4, R5
	BNE	cmp
	CMP	R8, R9
	BNE	cmpnext
	SUB $16, R6
	CMP	$16, R6
	BGE	chunk16_loop
	CBZ	R6, samebytes
	CMP	$8, R6
	BLE	tail
	// the length of tail > 8 bytes
	MOVD.P	8(R0), R4
	MOVD.P	8(R2), R5
	CMP	R4, R5
	BNE	cmp
	SUB	$8, R6
	// compare last 8 bytes
tail:
    SUB $8, R6
	MOVD	(R0)(R6), R4
	MOVD	(R2)(R6), R5
	CMP	R4, R5
	BEQ	samebytes
	MOVD    $8, R6
cmp:
	REV	R4, R4
	REV	R5, R5
cmprev:
	EOR R4, R5, R5
	CLZ R5, R5 // R5 = the number of bits that match
	LSR $3, R5, R5 // R5 = the number of bytes that match
	SUBS R5, R6, R6
	BLT samebytes
ret:
    SUB R6, R7
 	MOVD    R7, ret+48(FP)
	RET
small:
	TBZ	$3, R6, lt_8
	MOVD	(R0), R4
	MOVD	(R2), R5
	CMP	R4, R5
	BNE	cmp
	SUBS    $8, R6, R6
	BEQ	samebytes
	ADD	$8, R0
	ADD	$8, R2
	B	tail
lt_8:
	TBZ	$2, R6, lt_4
	MOVWU	(R0), R4
	MOVWU	(R2), R5
	CMPW	R4, R5
	BNE	cmp
	SUBS	$4, R6
	BEQ	samebytes
	ADD	$4, R0
	ADD	$4, R2
lt_4:
	TBZ	$1, R6, lt_2
	MOVHU	(R0), R4
	MOVHU	(R2), R5
	CMPW	R4, R5
	BNE	cmp
	ADD	$2, R0
	ADD	$2, R2
	SUB $2, R6
lt_2:
	TBZ	$0, R6, samebytes
one:
	MOVBU	(R0), R4
	MOVBU	(R2), R5
	CMPW	R4, R5
	BNE	ret
samebytes:
 	MOVD    R7, ret+48(FP)
	RET
cmpnext:
    SUB $8, R6
	REV	R8, R4
	REV	R9, R5
	B cmprev
