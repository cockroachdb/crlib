// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in licenses/BSD-golang.txt.

// This code is based on compare_amd64.s from Go 1.12.5.

TEXT ·CommonPrefix(SB),$0-56
    MOVQ    a_base+0(FP), SI
    MOVQ    a_len+8(FP), BX
    MOVQ    b_base+24(FP), DI
    MOVQ    b_len+32(FP), DX

	CMPQ	BX, DX
	MOVQ	DX, R8
	CMOVQLT	BX, R8 // R8 = min(alen, blen) = # of bytes to compare
	// Throughout this function, DX remembers the original min(alen, blen) and
	// R8 is the number of bytes we still need to compare (with bytes 0 to
	// DX-R8 known to match).
	MOVQ    R8, DX
	CMPQ	R8, $8
	JB  small

	CMPQ	R8, $63
	JBE	loop
	JMP	big_loop
	RET

// loop checks 16 bytes at a time.
loop:
	CMPQ	R8, $16
	JB  _0through15
	MOVOU	(SI), X0
	MOVOU	(DI), X1
	PCMPEQB X0, X1
	PMOVMSKB X1, AX
	XORQ	$0xffff, AX	// convert EQ to NE
	JNE	diff16	// branch if at least one byte is not equal
	ADDQ	$16, SI
	ADDQ	$16, DI
	SUBQ	$16, R8
	JMP	loop

diff64:
	SUBQ	$48, R8
	JMP	diff16
diff48:
	SUBQ	$32, R8
	JMP	diff16
diff32:
	SUBQ	$16, R8
	// AX = bit mask of differences
diff16:
	BSFQ	AX, BX	// index of first byte that differs
	SUBQ    BX, R8

	SUBQ    R8, DX
	MOVQ    DX, ret+48(FP)
	RET

_0through15: // R8 < 16, DX >= 8
	CMPQ	R8, $8
	JBE	_0through8
	MOVQ	(SI), AX
	MOVQ	(DI), CX
	CMPQ	AX, CX
	JNE	diff8
_0through8:
    // Load last 8 bytes of both.
	MOVQ	-8(SI)(R8*1), AX
	MOVQ	-8(DI)(R8*1), CX
	CMPQ	AX, CX
	JEQ	allsame
	MOVQ    $8, R8

	// AX and CX contain parts of a and b that differ.
diff8:
	BSWAPQ	AX	// reverse order of bytes
	BSWAPQ	CX
	XORQ	AX, CX
	BSRQ	CX, CX	// index of highest bit difference
	SHRQ    $3, CX  // index of highest byte difference

	SUBQ    R8, DX
	ADDQ    $7, DX
	SUBQ    CX, DX
	MOVQ    DX, ret+48(FP)
	RET

	// DX < 8
small:
	LEAQ	(R8*8), CX	// bytes left -> bits left
	NEGQ	CX		//  - bits lift (== 64 - bits left mod 64)
	JEQ	allsame

	// load bytes of a into high bytes of AX
	CMPB	SI, $0xf8
	JA	si_high
	MOVQ	(SI), SI
	JMP	si_finish
si_high:
	MOVQ	-8(SI)(R8*1), SI
	SHRQ	CX, SI
si_finish:
	SHLQ	CX, SI

	// load bytes of b into high bytes of BX
	CMPB	DI, $0xf8
	JA	di_high
	MOVQ	(DI), DI
	JMP	di_finish
di_high:
	MOVQ	-8(DI)(R8*1), DI
	SHRQ	CX, DI
di_finish:
	SHLQ	CX, DI

	BSWAPQ	SI	// reverse order of bytes
	BSWAPQ	DI
	XORQ	SI, DI	// find bit differences
	JEQ	allsame
	BSRQ	DI, CX	// index of highest bit difference
	SHRQ    $3, CX  // index of highest byte difference
	DECQ    DX
	SUBQ    CX, DX
	MOVQ    DX, ret+48(FP)
	RET

allsame:
	MOVQ    DX, ret+48(FP)
	RET

big_loop:
	MOVOU	(SI), X0
	MOVOU	(DI), X1
	PCMPEQB X0, X1
	PMOVMSKB X1, AX
	XORQ	$0xffff, AX
	JNE	diff16

	MOVOU	16(SI), X0
	MOVOU	16(DI), X1
	PCMPEQB X0, X1
	PMOVMSKB X1, AX
	XORQ	$0xffff, AX
	JNE	diff32

	MOVOU	32(SI), X0
	MOVOU	32(DI), X1
	PCMPEQB X0, X1
	PMOVMSKB X1, AX
	XORQ	$0xffff, AX
	JNE	diff48

	MOVOU	48(SI), X0
	MOVOU	48(DI), X1
	PCMPEQB X0, X1
	PMOVMSKB X1, AX
	XORQ	$0xffff, AX
	JNE	diff64

	ADDQ	$64, SI
	ADDQ	$64, DI
	SUBQ	$64, R8
	CMPQ	R8, $64
	JBE	loop
	JMP	big_loop
