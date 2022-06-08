package prg2p

import "strings"

var rules = `
# =======PREAMBLE========
# LETTERS
# a, ą, b, c, ć, d, e, ę, f, g, h, i, j, k, l, ł, m, n, ń, o, ó, p, q, r, s, ś, t, u, v, w, x, y, z, ź, ż, é, ü, ö, š, ë

# DIGRAPHS AND CONSONANT SOFTENING
# ch, cz, dz, dź, dż, rz, sz, ci, ni, si, zi, dzi

# PHONEMES
# a, a_, b, c, ci, cz, d, dz, dzi, drz, e, e_, f, g, h, i, j, k, l, l_, m, n, ni, o, p, r, s, si, sz, t, u, w, y, z, zi, rz
# + ng  + ni_ + h_
# ======================


# =====DECLARATION======
ALL = a, ą, b, c, ć, d, e, ę, f, g, h, i, j, k, l, ł, m, n, ń, o, ó, p, q, r, s, ś, t, u, v, w, x, y, z, ź, ż, ch, cz, dz, dź, dż, rz, sz, ci, ni, si, zi, dzi, é, ü, ö, š, ë, $

# EMPTY - ANY CHARACTER, THE BEGINNING AND END OF THE WORD INCLUDED
EMPTY = *

# SA - VOWELS
SA = a, e, y, ą, ę, u, o, i, ó

# SA1 - VOWELS
SA1 = a, e, y, ą, ę, u, o, ó

# SP - CONSONANTS
SP = b, c, ć, ch, cz, d, dz, dż, f, g, h, j, k, l, ł, m, n, ń, p, r, rz, s, ś, sz, t, w, z, ż, q, v, x, š

# SB - VOICELESS CONSONANTS
SB = p, k, t, s, c, h, f, ć, ś, ch, cz, sz, ci

# SD - VOICED CONSONANTS
SD = b, d, dz, dź, dż, g, rz, w, z, ź, ż

# SP1 - CONSONANTS
SP1 = p, b, m, w, f, ch, h

# SP2 - CONSONANTS
SP2 = l, r, t, d, k, g

# FR - FRICATIVES
FR = f, w, s, z, ź, ś, ż, h, sz

# SPZ - CONSONANTS FOR zi TRANSCRIPTS
SPZ = d, g, l, n, r, sz, ś

# END - END OF A WORD
END = $
# ======================


# =======RULES==========
# TAB-SEPARATED IN THE FOLLOWING FORMAT:
# BEFORE	CHARACTER	AFTER	 PHONEME

# ALWAYS
EMPTY	a	EMPTY	a
EMPTY	c	EMPTY	c
EMPTY	e	EMPTY	e
EMPTY	j	EMPTY	j
EMPTY	l	EMPTY	l
EMPTY	ł	EMPTY	l_
EMPTY	m	EMPTY	m
EMPTY	o	EMPTY	o
EMPTY	ó	EMPTY	u
EMPTY	p	EMPTY	p
EMPTY	r	EMPTY	r
END	y	EMPTY	j
-END	y	EMPTY	y

# DEVOICING
EMPTY	b	END	p, b
EMPTY	b	SB	p
EMPTY	b	-SB-END	b
EMPTY	dz	END	c, dz
EMPTY	dz	SB	c
EMPTY	dz	-SB-END	dz
EMPTY	dź	END	ci, dzi
EMPTY	dź	SB	ci
EMPTY	dź	-SB-END	dzi
EMPTY	d	END	t, d
EMPTY	d	SB	t
EMPTY	d	-SB-END	d
EMPTY	g	END	k, g
EMPTY	g	SB	k
EMPTY	g	-SB-END	g
EMPTY	rz	END	sz, rz
EMPTY	rz	SB	sz
SB	rz	EMPTY	sz
-SB	rz	-SB-END	rz
EMPTY	w	END	f, w
EMPTY	w	SB	f
SB	w	EMPTY	f
-SB	w	-SB-END	w
EMPTY	z	END	s,z
EMPTY	z	SB	s
EMPTY	z	-SB-END	z
EMPTY	ź	END	si, zi
EMPTY	ź	-END	zi
EMPTY	ż	END	sz, rz
EMPTY	ż	SB	sz
SB	ż	EMPTY	sz
-SB	ż	-SB-END	rz
EMPTY	dż	END	cz, drz
EMPTY	dż	-END	drz

# VOICING
EMPTY	ć	(b)	ci, dzi
EMPTY	ć	-(b)	ci
EMPTY	ś	(b)	si, zi
EMPTY	ś	-(b)	si
EMPTY	s	(b, g)	s, z
EMPTY	s	-(b, g)	s
EMPTY	t	(b, g)	t, d
EMPTY	t	-(b, g)	t
EMPTY	k	(b, ż)	k, g
EMPTY	k	-(b, ż)	k
EMPTY	cz	(b, d)	cz, drz
EMPTY	cz	-(b, d)	cz
EMPTY	f	(g)	f, w
EMPTY	f	-(g)	f

# MISCELLANEA
EMPTY	h	SD	h
EMPTY	h	-SD	h
SP1	i	SA	j
SP2	i	SA	i, j
SP1	i	-SA	i
SP2	i	-SA	i
-SP1-SP2	i	EMPTY	i
EMPTY	n	(k, g)	n
EMPTY	n	(ni, ci, dzi)	n, ni
EMPTY	n	-(k, g, ni, ci, dzi)	n
EMPTY	ń	FR	ni
EMPTY	ń	-FR	ni
(a, e)	u	EMPTY	l_
-(a, e)	u	EMPTY	u

# DIGRAPHS
EMPTY	sz	EMPTY	sz
EMPTY	ch	SD	h
EMPTY	ch	-SD	h
EMPTY	ci	SA1	ci
EMPTY	ci	SP	ci i
EMPTY	ci	END	ci i
EMPTY	dzi	SA1	dzi
EMPTY	dzi	SP	dzi i
EMPTY	dzi	END	dzi i
EMPTY	ni	SA1	ni
EMPTY	ni	(i)	ni j, ni
EMPTY	ni	SP	ni i
EMPTY	ni	END	ni i
EMPTY	si	SA1	si
EMPTY	si	SP	si i
EMPTY	si	END	si i
EMPTY	zi	SA1	zi
EMPTY	zi	-SA1-SPZ	zi i
EMPTY	zi	SPZ	z i
EMPTY	zi	END	zi i

# ĄĘ
EMPTY	ę	END	e, e_
EMPTY	ą	END	o l_, a_, o m
EMPTY	ę	(l, ł)	e
EMPTY	ą	(ł)	o
EMPTY	ę	(b, p)	e m
EMPTY	ą	(b, p)	o m
EMPTY	ę	(m, n)	e
EMPTY	ą	(m)	o
EMPTY	ą	(n, j)	a_, o l_
EMPTY	ę	(w, f, s, z, ż, rz, sz, ch, h)	e_, e l_
EMPTY	ą	(w, f, s, z, ż, rz, sz, ch, h)	a_, o l_
EMPTY	ę	(d, t, c, dz, dż, cz)	e_, e n
EMPTY	ą	(d, t, c, dz, dż, cz)	a_, o n
EMPTY	ę	(dź, dzi, ć, ci)	e_, e ni
EMPTY	ą	(dź, dzi, ć, ci)	a_, o ni
EMPTY	ę	(g, k)	e n, e_
EMPTY	ą	(g, k)	o n, a_
EMPTY	ę	(ź, zi, ś, si)	e l_, e ni
EMPTY	ą	(ź, zi, ś, si)	o l_, o ni

# VARIANTS, CLUSTERS, EXCEPTIONS
EMPTY	trz	EMPTY	t sz, cz
EMPTY	drz	EMPTY	d rz, drz
EMPTY	zsz	EMPTY	s sz, sz
EMPTY	nadz	-(i)	n a d z
EMPTY	nadż	EMPTY	n a d rz
EMPTY	podz	-(i)	p o d z
EMPTY	podż	EMPTY	p o d rz
END	odz	-(i)	o d z
EMPTY	odż	EMPTY	o d rz
EMPTY	budż	EMPTY	b u d rz
EMPTY	śćdzi	EMPTY	si dzi, zi dzi, si ci dzi
EMPTY	śćs	EMPTY	si s, j s, si ci s
EMPTY	śródzi	EMPTY	si r u d zi
EMPTY	sji	EMPTY	s j i, s i
EMPTY	cji	EMPTY	c j i, c i
EMPTY	izm	END	i z m, i s m
EMPTY	izm	(i)	i z m, i zi m
EMPTY	on	(s)	o n, a_
-END	en	(t, k, ci, s)	e n, e_
-(n)	ii	END	i i, i, j i
EMPTY	dźm	(y)	ci m
EMPTY	marz	(ł, n, l)	m a r z
EMPTY	żć	END	si ci, zi ci
EMPTY	klie	EMPTY	k l i j e

# TO REPAIR
EMPTY	czw	EMPTY	cz f
EMPTY	szw	EMPTY	sz f

# FOREIGN
EMPTY	v	EMPTY	w
EMPTY	q	EMPTY	k
EMPTY	x	EMPTY	k s
EMPTY	é	EMPTY	e
EMPTY	ü	EMPTY	u, i
EMPTY	ö	EMPTY	y
EMPTY	š	EMPTY	s
EMPTY	ë	EMPTY	e
# ======================
`

// Rules returns a default set of g2p rules.
func Rules() *strings.Reader {
	r := strings.NewReader(rules)
	return r
}
