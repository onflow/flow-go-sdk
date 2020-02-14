// +build relic

#ifndef _REL_MISC_INCLUDE_H
#define _REL_MISC_INCLUDE_H

#include "relic.h"

typedef uint8_t byte;

#define BITS_TO_BYTES(x) ((x+7)>>3)
#define BITS_TO_DIGITS(x) ((x+63)>>6)
#define MIN(a,b) ((a)>(b)?(b):(a))

// Most of the functions are written for ALLOC=AUTO not ALLOC=DYNAMIC

// Debug related functions
void    _bytes_print(char*, byte*, int);
void    _fp_print(char*, fp_t);
void    _bn_print(char*, bn_st*);
void    _ep_print(char*, ep_st*);
void    _ep2_print(char*, ep2_st*);
void    fp_read_raw(fp_t a, const dig_t *raw, int len);
void    _seed_relic(byte*, int);
void    _bn_randZr(bn_t);
void    _G1scalarGenMult(ep_st*, const bn_st*);
void    opswu_test(uint8_t *, const uint8_t *, int);

#endif