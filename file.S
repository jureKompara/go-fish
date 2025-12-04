TEXT go-fish/internal/engine.(*Position).Perft(SB) /home/jk/usb/chessbruh/go-fish/internal/engine/perft.go
  perft.go:31		0x4a1680		493b6610		CMPQ SP, 0x10(R14)					
  perft.go:31		0x4a1684		0f86b2010000		JBE 0x4a183c						
  perft.go:31		0x4a168a		55			PUSHQ BP						
  perft.go:31		0x4a168b		4889e5			MOVQ SP, BP						
  perft.go:31		0x4a168e		4883ec50		SUBQ $0x50, SP						
  perft.go:34		0x4a1692		4889442460		MOVQ AX, 0x60(SP)					
  perft.go:34		0x4a1697		4883fb01		CMPQ BX, $0x1						
  perft.go:34		0x4a169b		7433			JE 0x4a16d0						
  perft.go:34		0x4a169d		0f1f00			NOPL 0(AX)						
  perft.go:45		0x4a16a0		4885db			TESTQ BX, BX						
  perft.go:45		0x4a16a3		7420			JE 0x4a16c5						
  perft.go:34		0x4a16a5		48895c2468		MOVQ BX, 0x68(SP)					
  perft.go:49		0x4a16aa		e871ecffff		CALL go-fish/internal/engine.(*Position).pseudoAll(SB)	
  perft.go:49		0x4a16af		4889442448		MOVQ AX, 0x48(SP)					
  perft.go:49		0x4a16b4		48895c2440		MOVQ BX, 0x40(SP)					
  perft.go:49		0x4a16b9		31c9			XORL CX, CX						
  perft.go:49		0x4a16bb		31d2			XORL DX, DX						
  perft.go:49		0x4a16bd		0f1f00			NOPL 0(AX)						
  perft.go:49		0x4a16c0		e9ed000000		JMP 0x4a17b2						
  perft.go:46		0x4a16c5		b801000000		MOVL $0x1, AX						
  perft.go:46		0x4a16ca		4883c450		ADDQ $0x50, SP						
  perft.go:46		0x4a16ce		5d			POPQ BP							
  perft.go:46		0x4a16cf		c3			RET							
  perft.go:36		0x4a16d0		e84becffff		CALL go-fish/internal/engine.(*Position).pseudoAll(SB)	
  perft.go:36		0x4a16d5		4889442448		MOVQ AX, 0x48(SP)					
  perft.go:36		0x4a16da		48895c2440		MOVQ BX, 0x40(SP)					
  perft.go:36		0x4a16df		31c9			XORL CX, CX						
  perft.go:36		0x4a16e1		31d2			XORL DX, DX						
  perft.go:36		0x4a16e3		eb4e			JMP 0x4a1733						
  perft.go:38		0x4a16e5		488b9cd0e0000000	MOVQ 0xe0(AX)(DX*8), BX					
  perft.go:38		0x4a16ed		e86e010000		CALL go-fish/internal/engine.(*Position).isAttacked(SB)	
  perft.go:38		0x4a16f2		8844241f		MOVB AL, 0x1f(SP)					
  perft.go:41		0x4a16f6		488b442460		MOVQ 0x60(SP), AX					
  perft.go:41		0x4a16fb		8b5c2424		MOVL 0x24(SP), BX					
  perft.go:41		0x4a16ff		90			NOPL							
  perft.go:41		0x4a1700		e8db0a0000		CALL go-fish/internal/engine.(*Position).Unmake(SB)	
  perft.go:39		0x4a1705		488b542430		MOVQ 0x30(SP), DX					
  perft.go:39		0x4a170a		488d7201		LEAQ 0x1(DX), SI					
  perft.go:39		0x4a170e		0fb67c241f		MOVZX 0x1f(SP), DI					
  perft.go:39		0x4a1713		400fb6ff		MOVZX DI, DI						
  perft.go:36		0x4a1717		488b4c2438		MOVQ 0x38(SP), CX					
  perft.go:36		0x4a171c		48ffc1			INCQ CX							
  perft.go:39		0x4a171f		4885ff			TESTQ DI, DI						
  perft.go:39		0x4a1722		480f45f2		CMOVNE DX, SI						
  perft.go:36		0x4a1726		488b442448		MOVQ 0x48(SP), AX					
  perft.go:36		0x4a172b		488b5c2440		MOVQ 0x40(SP), BX					
  perft.go:39		0x4a1730		4889f2			MOVQ SI, DX						
  perft.go:36		0x4a1733		4839cb			CMPQ BX, CX						
  perft.go:36		0x4a1736		7e3a			JLE 0x4a1772						
  perft.go:36		0x4a1738		48894c2438		MOVQ CX, 0x38(SP)					
  perft.go:39		0x4a173d		4889542430		MOVQ DX, 0x30(SP)					
  perft.go:36		0x4a1742		8b1c88			MOVL 0(AX)(CX*4), BX					
  perft.go:36		0x4a1745		895c2424		MOVL BX, 0x24(SP)					
  perft.go:37		0x4a1749		488b442460		MOVQ 0x60(SP), AX					
  perft.go:37		0x4a174e		e8ad030000		CALL go-fish/internal/engine.(*Position).Make(SB)	
  perft.go:38		0x4a1753		488b442460		MOVQ 0x60(SP), AX					
  perft.go:38		0x4a1758		488b88b8000000		MOVQ 0xb8(AX), CX					
  perft.go:38		0x4a175f		488d51ff		LEAQ -0x1(CX), DX					
  perft.go:38		0x4a1763		48f7da			NEGQ DX							
  perft.go:38		0x4a1766		4883fa02		CMPQ DX, $0x2						
  perft.go:38		0x4a176a		0f8275ffffff		JB 0x4a16e5						
  perft.go:38		0x4a1770		eb09			JMP 0x4a177b						
  perft.go:43		0x4a1772		4889d0			MOVQ DX, AX						
  perft.go:43		0x4a1775		4883c450		ADDQ $0x50, SP						
  perft.go:43		0x4a1779		5d			POPQ BP							
  perft.go:43		0x4a177a		c3			RET							
  perft.go:38		0x4a177b		4889d0			MOVQ DX, AX						
  perft.go:38		0x4a177e		b902000000		MOVL $0x2, CX						
  perft.go:38		0x4a1783		e8783efdff		CALL runtime.panicIndex(SB)				
  perft.go:52		0x4a1788		4889542428		MOVQ DX, 0x28(SP)					
  perft.go:54		0x4a178d		488b442460		MOVQ 0x60(SP), AX					
  perft.go:54		0x4a1792		8b5c2420		MOVL 0x20(SP), BX					
  perft.go:54		0x4a1796		e8450a0000		CALL go-fish/internal/engine.(*Position).Unmake(SB)	
  perft.go:49		0x4a179b		488b4c2438		MOVQ 0x38(SP), CX					
  perft.go:49		0x4a17a0		48ffc1			INCQ CX							
  perft.go:49		0x4a17a3		488b442448		MOVQ 0x48(SP), AX					
  perft.go:49		0x4a17a8		488b5c2440		MOVQ 0x40(SP), BX					
  perft.go:52		0x4a17ad		488b542428		MOVQ 0x28(SP), DX					
  perft.go:49		0x4a17b2		4839cb			CMPQ BX, CX						
  perft.go:49		0x4a17b5		7e6e			JLE 0x4a1825						
  perft.go:49		0x4a17b7		48894c2438		MOVQ CX, 0x38(SP)					
  perft.go:52		0x4a17bc		4889542428		MOVQ DX, 0x28(SP)					
  perft.go:49		0x4a17c1		8b1c88			MOVL 0(AX)(CX*4), BX					
  perft.go:49		0x4a17c4		895c2420		MOVL BX, 0x20(SP)					
  perft.go:50		0x4a17c8		488b442460		MOVQ 0x60(SP), AX					
  perft.go:50		0x4a17cd		e82e030000		CALL go-fish/internal/engine.(*Position).Make(SB)	
  perft.go:51		0x4a17d2		488b442460		MOVQ 0x60(SP), AX					
  perft.go:51		0x4a17d7		488b88b8000000		MOVQ 0xb8(AX), CX					
  perft.go:51		0x4a17de		488d51ff		LEAQ -0x1(CX), DX					
  perft.go:51		0x4a17e2		48f7da			NEGQ DX							
  perft.go:51		0x4a17e5		4883fa02		CMPQ DX, $0x2						
  perft.go:51		0x4a17e9		7343			JAE 0x4a182e						
  perft.go:51		0x4a17eb		488b9cd0e0000000	MOVQ 0xe0(AX)(DX*8), BX					
  perft.go:51		0x4a17f3		e868000000		CALL go-fish/internal/engine.(*Position).isAttacked(SB)	
  perft.go:51		0x4a17f8		84c0			TESTL AL, AL						
  perft.go:51		0x4a17fa		7407			JE 0x4a1803						
  perft.go:52		0x4a17fc		488b542428		MOVQ 0x28(SP), DX					
  perft.go:51		0x4a1801		eb85			JMP 0x4a1788						
  perft.go:52		0x4a1803		488b5c2468		MOVQ 0x68(SP), BX					
  perft.go:52		0x4a1808		48ffcb			DECQ BX							
  perft.go:52		0x4a180b		488b442460		MOVQ 0x60(SP), AX					
  perft.go:52		0x4a1810		e86bfeffff		CALL go-fish/internal/engine.(*Position).Perft(SB)	
  perft.go:52		0x4a1815		488b4c2428		MOVQ 0x28(SP), CX					
  perft.go:52		0x4a181a		488d1401		LEAQ 0(CX)(AX*1), DX					
  perft.go:52		0x4a181e		6690			NOPW							
  perft.go:52		0x4a1820		e963ffffff		JMP 0x4a1788						
  perft.go:56		0x4a1825		4889d0			MOVQ DX, AX						
  perft.go:56		0x4a1828		4883c450		ADDQ $0x50, SP						
  perft.go:56		0x4a182c		5d			POPQ BP							
  perft.go:56		0x4a182d		c3			RET							
  perft.go:51		0x4a182e		4889d0			MOVQ DX, AX						
  perft.go:51		0x4a1831		b902000000		MOVL $0x2, CX						
  perft.go:51		0x4a1836		e8c53dfdff		CALL runtime.panicIndex(SB)				
  perft.go:51		0x4a183b		90			NOPL							
  perft.go:31		0x4a183c		4889442408		MOVQ AX, 0x8(SP)					
  perft.go:31		0x4a1841		48895c2410		MOVQ BX, 0x10(SP)					
  perft.go:31		0x4a1846		e8b51cfdff		CALL runtime.morestack_noctxt.abi0(SB)			
  perft.go:31		0x4a184b		488b442408		MOVQ 0x8(SP), AX					
  perft.go:31		0x4a1850		488b5c2410		MOVQ 0x10(SP), BX					
  perft.go:31		0x4a1855		e926feffff		JMP go-fish/internal/engine.(*Position).Perft(SB)	
