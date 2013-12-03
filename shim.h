#ifndef UNITYSCOPE_SHIM_H
#define UNITYSCOPE_SHIM_H

#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

/* A typedef that can be used to represent a std::shared_ptr */

typedef uintptr_t SharedPtrData[2];


void init_reply_ptr(SharedPtrData dest, SharedPtrData src);
void destroy_reply_ptr(SharedPtrData data);


#ifdef __cplusplus
}
#endif

#endif
